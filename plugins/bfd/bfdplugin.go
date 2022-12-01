// SPDX-License-Identifier: Apache-2.0

// Copyright 2021 PANTHEON.tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate descriptor-adapter --descriptor-name Bfd --value-type *bfd.BFD --import "go.pantheon.tech/stonework/proto/bfd" --output-dir "descriptor"

package bfdplugin

import (
	"context"
	"errors"
	"sync"

	govppapi "go.fd.io/govpp/api"

	"go.ligato.io/cn-infra/v2/infra"
	"go.ligato.io/cn-infra/v2/logging"
	grpc_plugin "go.ligato.io/cn-infra/v2/rpc/grpc"

	"go.ligato.io/vpp-agent/v3/plugins/govppmux"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"

	"go.pantheon.tech/stonework/plugins/bfd/descriptor"
	"go.pantheon.tech/stonework/plugins/bfd/vppcalls"
	"go.pantheon.tech/stonework/proto/bfd"

	_ "go.pantheon.tech/stonework/plugins/bfd/vppcalls/vpp2106"
	_ "go.pantheon.tech/stonework/plugins/bfd/vppcalls/vpp2202"
	_ "go.pantheon.tech/stonework/plugins/bfd/vppcalls/vpp2210"
)

// BfdPlugin groups required BFD dependencies and descriptors
type BfdPlugin struct {
	sync.Mutex
	Deps

	// VPP API handler
	bfdHandler vppcalls.BfdVppAPI
	vppCh      govppapi.Channel

	// BFD events
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	bfdEventChan chan *bfd.BFDEvent
	bfdEventSubs []bfdEventSub
	grpcSrv      *grpcService
}

// Deps is a set of BFD plugin dependencies
type Deps struct {
	infra.PluginDeps
	KVScheduler kvs.KVScheduler
	GoVpp       govppmux.API
	IfPlugin    ifplugin.API
	GRPC        grpc_plugin.Server
}

type bfdEventSub struct {
	name      string
	ctx       context.Context
	eventChan chan<- *bfd.BFDEvent
}

type grpcService struct {
	bfd.UnimplementedBFDWatcherServer

	// Deps:
	Log       logging.Logger
	BfdPlugin API
}

// WatchBfdEvents subscribes for BFD state change notifications.
func (p *BfdPlugin) WatchBFDEvents(ctx context.Context, subName string, eventChan chan<- *bfd.BFDEvent) error {
	p.Lock()
	defer p.Unlock()
	p.bfdEventSubs = append(p.bfdEventSubs, bfdEventSub{
		ctx:       ctx,
		name:      subName,
		eventChan: eventChan,
	})
	return nil
}

// Init the VPP handler and register descriptors
func (p *BfdPlugin) Init() error {
	p.bfdEventChan = make(chan *bfd.BFDEvent, 1000)
	p.ctx, p.cancel = context.WithCancel(context.Background())

	p.bfdHandler = vppcalls.CompatibleBfdVppHandler(p.GoVpp, p.IfPlugin.GetInterfaceIndex(), p.Log)
	if p.bfdHandler == nil {
		return errors.New("bfdHandler is not available")
	}

	bfdDescriptor := descriptor.NewBfdDescriptor(p.bfdHandler, p.Log)
	if err := p.KVScheduler.RegisterKVDescriptor(bfdDescriptor); err != nil {
		return err
	}

	// allow to watch BFD events over gRPC
	grpcServer := p.GRPC.GetServer()
	if grpcServer == nil {
		return errors.New("gRPC server is not initialized")
	}
	p.grpcSrv = &grpcService{BfdPlugin: p, Log: p.Log.NewLogger("bfd-grpc-srv")}
	bfd.RegisterBFDWatcherServer(grpcServer, p.grpcSrv)
	return nil
}

// AfterInit subscribes for watching BFD state events.
func (p *BfdPlugin) AfterInit() error {
	watch := func() error {
		if err := p.bfdHandler.WatchBfdEvents(p.ctx, p.bfdEventChan); err != nil {
			return err
		}
		p.wg.Add(1)
		go p.processBfdEvents()
		return nil
	}
	if err := watch(); err != nil {
		return err
	}

	p.GoVpp.OnReconnect(func() {
		p.cancel()
		p.wg.Wait()
		if err := watch(); err != nil {
			p.Log.Warnf("WatchBFDEvents failed: %v", err)
		}
	})
	return nil
}

func (p *BfdPlugin) processBfdEvents() {
	defer p.wg.Done()
	for {
		select {
		case ev := <-p.bfdEventChan:
			p.Lock()
			for i := 0; i < len(p.bfdEventSubs); {
				sub := p.bfdEventSubs[i]
				select {
				case <-sub.ctx.Done():
					// subscription ended
					p.bfdEventSubs = append(p.bfdEventSubs[:i], p.bfdEventSubs[i+1:]...)
					p.Log.Debugf("subscription '%s' ended", sub.name)
					continue
				case sub.eventChan <- ev:
					// ok
				default:
					p.Log.Warnf("failed to deliver BFD notification to subscriber: %s",
						sub.name)
				}
				i++
			}
			p.Unlock()
		case <-p.ctx.Done():
			return
		}
	}
}

// Close the event channel
func (p *BfdPlugin) Close() error {
	p.cancel()
	p.wg.Wait()
	close(p.bfdEventChan)
	return nil
}

// WatchBFDEvents allows to subscribe for BFD events over gRPC.
func (s *grpcService) WatchBFDEvents(req *bfd.WatchBFDEventsRequest, srv bfd.BFDWatcher_WatchBFDEventsServer) error {
	bfdEventChan := make(chan *bfd.BFDEvent, 1000)
	defer close(bfdEventChan)
	err := s.BfdPlugin.WatchBFDEvents(srv.Context(), req.GetSubscriptionLabel(), bfdEventChan)
	if err != nil {
		return err
	}
	for {
		select {
		case ev := <-bfdEventChan:
			if err := srv.Send(ev); err != nil {
				s.Log.Warnf("gRPC server send error: %v", err)
				return err
			}
		case <-srv.Context().Done():
			return srv.Context().Err()
		}
	}
}

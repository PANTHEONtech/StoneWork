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

//go:generate descriptor-adapter --descriptor-name ISISX --value-type *vpp_isisx.ISISXConnection --import "go.pantheon.tech/stonework/proto/isisx" --output-dir "descriptor"

package isisxplugin

import (
	"github.com/go-errors/errors"
	govppapi "go.fd.io/govpp/api"
	"go.ligato.io/cn-infra/v2/health/statuscheck"
	"go.ligato.io/cn-infra/v2/infra"
	"go.ligato.io/vpp-agent/v3/plugins/govppmux"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"

	"go.pantheon.tech/stonework/plugins/isisx/descriptor"
	"go.pantheon.tech/stonework/plugins/isisx/vppcalls"

	_ "go.pantheon.tech/stonework/plugins/isisx/vppcalls/vpp2106"
	_ "go.pantheon.tech/stonework/plugins/isisx/vppcalls/vpp2202"
)

// ISISXPlugin is a plugin that manages ISIS protocol packet forwarding.
// ISIS packets can't be handled by ABX because ISIS packets won't reach ACL in VPP (they are dropped
// by OSI VPP node). Therefore ISIS crossconnect plugin for VPP was created. Packet are grabbed sooner
// and based on protocol filtering. This ISISXPlugin should handle configuration(control plane) for that
// ISISX VPP plugin.
type ISISXPlugin struct {
	Deps

	// GoVPP channels
	vppCh govppapi.Channel

	// handlers and descriptors
	isisxHandler    vppcalls.ISISXVppAPI
	isisxDescriptor *kvs.KVDescriptor
}

// Deps represents dependencies for the plugin.
type Deps struct {
	infra.PluginDeps
	Scheduler   kvs.KVScheduler
	GoVppmux    govppmux.API
	IfPlugin    ifplugin.API
	StatusCheck statuscheck.PluginStatusWriter // optional
}

// Init initializes ISISX plugin.
func (p *ISISXPlugin) Init() error {
	// init handler
	p.isisxHandler = vppcalls.CompatibleIsisxVppHandler(p.GoVppmux, p.IfPlugin.GetInterfaceIndex(), p.Log)
	if p.isisxHandler == nil {
		return errors.New("isisxHandler is not available")
	}

	// init & register descriptor
	p.isisxDescriptor = descriptor.NewISIXDescriptor(p.isisxHandler, p.Log)
	if err := p.Deps.Scheduler.RegisterKVDescriptor(p.isisxDescriptor); err != nil {
		return err
	}

	return nil
}

// AfterInit registers plugin with StatusCheck.
func (p *ISISXPlugin) AfterInit() error {
	if p.StatusCheck != nil {
		p.StatusCheck.Register(p.PluginName, nil)
	}
	return nil
}

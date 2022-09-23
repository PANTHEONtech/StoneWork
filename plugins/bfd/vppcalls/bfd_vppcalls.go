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

package vppcalls

import (
	"context"

	govppapi "go.fd.io/govpp/api"
	"go.ligato.io/cn-infra/v2/logging"

	"go.ligato.io/vpp-agent/v3/plugins/vpp"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"

	"go.pantheon.tech/stonework/proto/bfd"
)

// BfdVppAPI defines methods to add, delete and watch BFD configuration
type BfdVppAPI interface {
	// AddBfd creates BFD session attached to the defined interface
	// with given configuration ID.
	AddBfd(confID uint32, bfd *bfd.BFD) error

	// DeleteBfd removes existing BFD session.
	DeleteBfd(bfd *bfd.BFD) error

	// DumpBfd returns retrieved BFD data together with BFD state.
	DumpBfd() ([]*BfdDetails, error)

	// WatchBfdEvents starts BFD event watcher.
	WatchBfdEvents(ctx context.Context, eventChan chan<- *bfd.BFDEvent) error
}

// BfdDetails represents retrieved BFD data
type BfdDetails struct {
	Config          *bfd.BFD
	State           bfd.BFDEvent_SessionState
	ConfKey         uint32
	BfdKey          uint8
	IsAuthenticated bool
}

var handler = vpp.RegisterHandler(vpp.HandlerDesc{
	Name:       "bfd",
	HandlerAPI: (*BfdVppAPI)(nil),
})

func AddBfdHandlerVersion(version vpp.Version, msgs []govppapi.Message,
	h func(ch govppapi.Channel, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) BfdVppAPI,
) {
	handler.AddVersion(vpp.HandlerVersion{
		Version: version,
		Check: func(c vpp.Client) error {
			ch, err := c.NewAPIChannel()
			if err != nil {
				return err
			}
			return ch.CheckCompatiblity(msgs...)
		},
		NewHandler: func(c vpp.Client, a ...interface{}) vpp.HandlerAPI {
			ch, err := c.NewAPIChannel()
			if err != nil {
				return err
			}
			return h(ch, a[0].(ifaceidx.IfaceMetadataIndex), a[1].(logging.Logger))
		},
	})
}

func CompatibleBfdVppHandler(c vpp.Client, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) BfdVppAPI {
	if v := handler.FindCompatibleVersion(c); v != nil {
		return v.NewHandler(c, ifIdx, log).(BfdVppAPI)
	}
	return nil
}

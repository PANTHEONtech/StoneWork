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
	govppapi "git.fd.io/govpp.git/api"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/plugins/vpp"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"
	proto_isisx "go.pantheon.tech/stonework/proto/isisx"
)

// ISISXVppAPI provides read/write methods required to handle VPP ISIS-protocol-based forwarding
type ISISXVppAPI interface {
	ISISXVppRead

	// GetISISXVersion retrieves version of the VPP ISISX plugin
	GetISISXVersion() (ver string, err error)
	// AddISISXConnection creates new ISISX unidirectional cross-connection between 2 interfaces
	AddISISXConnection(inputInterface, outputInterface string) error
	// DeleteISISXConnection deletes existing ISISX unidirectional cross-connection between 2 interfaces
	DeleteISISXConnection(inputInterface, outputInterface string) error
}

// ISISXVppRead provides read methods for ISISX plugin
type ISISXVppRead interface {
	// DumpISISXConnections retrieves VPP ISISX configuration.
	DumpISISXConnections() ([]*proto_isisx.ISISXConnection, error)
}

var handler = vpp.RegisterHandler(vpp.HandlerDesc{
	Name:       "isisx",
	HandlerAPI: (*ISISXVppAPI)(nil),
})

func AddIsisxHandlerVersion(version vpp.Version, msgs []govppapi.Message,
	h func(ch govppapi.Channel, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) ISISXVppAPI,
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

func CompatibleIsisxVppHandler(c vpp.Client, ifIdx ifaceidx.IfaceMetadataIndex,
	log logging.Logger) ISISXVppAPI {
	if v := handler.FindCompatibleVersion(c); v != nil {
		return v.NewHandler(c, ifIdx, log).(ISISXVppAPI)
	}
	return nil
}

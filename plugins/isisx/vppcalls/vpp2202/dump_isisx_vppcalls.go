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

package vpp2202

import (
	"github.com/go-errors/errors"

	"go.pantheon.tech/stonework/plugins/binapi/vpp2202/isisx"
	proto_isisx "go.pantheon.tech/stonework/proto/isisx"
)

// DumpISISXConnections retrieves VPP ISISX configuration.
func (h *ISISXVppHandler) DumpISISXConnections() ([]*proto_isisx.ISISXConnection, error) {
	var connections []*proto_isisx.ISISXConnection

	// make multi request
	req := &isisx.IsisxConnectionDump{}
	reqCtx := h.callsChannel.SendMultiRequest(req)
	for {
		reply := &isisx.IsisxConnectionDetails{}
		last, err := reqCtx.ReceiveReply(reply)
		if err != nil {
			return nil, err
		}
		if last {
			break
		}

		// translate interface indexes to names
		inputInterfaceName, _, exists := h.ifIndexes.LookupBySwIfIndex(reply.Connection.RxSwIfIndex)
		if !exists {
			return nil, errors.Errorf("input interface %d not found", reply.Connection.RxSwIfIndex)
		}
		outputInterfaceName, _, exists := h.ifIndexes.LookupBySwIfIndex(reply.Connection.TxSwIfIndex)
		if !exists {
			return nil, errors.Errorf("output interface %d not found", reply.Connection.TxSwIfIndex)
		}
		// remember connection
		connections = append(connections, &proto_isisx.ISISXConnection{
			InputInterface:  inputInterfaceName,
			OutputInterface: outputInterfaceName,
		})
	}

	return connections, nil
}

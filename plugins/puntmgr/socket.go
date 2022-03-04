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

package puntmgr

import (
	"google.golang.org/protobuf/proto"

	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"

	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// socketPunt implements PuntHandler for PuntRequest_PUNT_TO_SOCKET
type socketPunt struct{}

func NewSocketPuntHandler() PuntHandler {
	return &socketPunt{}
}

// PuntSelector is used only by socketPunt to ensure that no two punt requests attempt to configure the same
// punt rule.
func PuntSelector(puntKey string) string {
	return "vpp/punt-to-socket/" + puntKey
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *socketPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if socket := punt.GetPuntToSocket(); socket != nil {
		var socketPath, key string
		switch {
		case socket.GetException() != nil:
			socketPath = socket.GetException().GetSocketPath()
			key = models.Key(socket.GetException())
		case socket.GetToHost() != nil:
			socketPath = socket.GetToHost().GetSocketPath()
			key = models.Key(socket.GetToHost())
		}
		return []InterconnectReq{
			{
				link: &AFUnixLink{
					socketPath: socketPath,
				},
				vppSelector: PuntSelector(key),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *socketPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	return nil
}

// CanMultiplex enables interconnection multiplexing for this punting. It could be enabled in certain cases:
// 1. two or more punts of this type can coexist even if they have the same vpp selector
// 2. one or more punts of this type can coexist with other type of punts on the same (TAP-only)
// interconnection if they all have the same vpp selector and cnf selector.
// The TAP-backed interconnection is shared for multiple multiplexing punts with the same cnf selector
// (same network namespace) and vpp selector.
func (p *socketPunt) CanMultiplex() bool {
	return false
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *socketPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	socket := puntReq.GetPuntToSocket()
	var punt proto.Message
	switch {
	case socket.GetException() != nil:
		punt = socket.GetException()
	case socket.GetToHost() != nil:
		punt = socket.GetToHost()
	}
	if remove {
		txn.Delete(punt)
	} else {
		txn.Update(punt)
	}
	return nil
}

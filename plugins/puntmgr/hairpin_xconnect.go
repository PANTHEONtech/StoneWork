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
	"github.com/golang/protobuf/proto"

	"go.ligato.io/vpp-agent/v3/client"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	vpp_l2 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l2"

	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// hairpinXConnPunt implements PuntHandler for PuntRequest_HAIRPIN_XCONNECT
type hairpinXConnPunt struct{}

func NewHairpinXConnPuntHandler() PuntHandler {
	return &hairpinXConnPunt{}
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *hairpinXConnPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if hairpinXConn := punt.GetHairpinXConnect(); hairpinXConn != nil {
		return []InterconnectReq{
			{
				link: &InterfaceLink{},
				// Selector = interface name, i.e. same as used by Hairpin and ABX, both of which
				// are mutually exclusive with Hairpin-XConnect.
				vppSelector: VppInterfaceSelector(hairpinXConn.VppInterface1),
			},
			{
				link: &InterfaceLink{},
				// Selector = interface name, i.e. same as used by Hairpin and ABX, both of which
				// are mutually exclusive with Hairpin-XConnect.
				vppSelector: VppInterfaceSelector(hairpinXConn.VppInterface2),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *hairpinXConnPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	// L2 VPP interfaces
	if hairpinXConn := punt.GetHairpinXConnect(); hairpinXConn != nil {
		deps = append(deps,
			kvs.Dependency{
				Label: punt.GetLabel() + "-hairpin-xconnect-" + hairpinXConn.GetVppInterface1(),
				Key:   vpp_interfaces.InterfaceKey(hairpinXConn.GetVppInterface1()),
			}, kvs.Dependency{
				Label: punt.GetLabel() + "-hairpin-xconnect-" + hairpinXConn.GetVppInterface2(),
				Key:   vpp_interfaces.InterfaceKey(hairpinXConn.GetVppInterface2()),
			})
	}
	return deps
}

// CanMultiplex enables interconnection multiplexing for this punting. It could be enabled in certain cases:
// 1. two or more punts of this type can coexist even if they have the same vpp selector
// 2. one or more punts of this type can coexist with other type of punts on the same (TAP-only)
// interconnection if they all have the same vpp selector and cnf selector.
// The TAP-backed interconnection is shared for multiple multiplexing punts with the same cnf selector
// (same network namespace) and vpp selector.
func (p *hairpinXConnPunt) CanMultiplex() bool {
	return false
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *hairpinXConnPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	vppIface1 := puntReq.GetHairpinXConnect().VppInterface1
	vppIface2 := puntReq.GetHairpinXConnect().VppInterface2
	var (
		icIface1, icIface2 *pb.PuntMetadata_Interface
	)
	if interconnects[0].Id.VppSelector == VppInterfaceSelector(vppIface1) {
		icIface1 = interconnects[0].VppInterface
		icIface2 = interconnects[1].VppInterface
	} else {
		icIface1 = interconnects[1].VppInterface
		icIface2 = interconnects[0].VppInterface
	}

	config := []proto.Message{
		&vpp_l2.XConnectPair{
			ReceiveInterface:  vppIface1,
			TransmitInterface: icIface1.Name,
		},
		&vpp_l2.XConnectPair{
			ReceiveInterface:  icIface1.Name,
			TransmitInterface: vppIface1,
		},
		&vpp_l2.XConnectPair{
			ReceiveInterface:  vppIface2,
			TransmitInterface: icIface2.Name,
		},
		&vpp_l2.XConnectPair{
			ReceiveInterface:  icIface2.Name,
			TransmitInterface: vppIface2,
		},
	}
	if remove {
		txn.Delete(config...)
	} else {
		txn.Update(config...)
	}
	return nil
}

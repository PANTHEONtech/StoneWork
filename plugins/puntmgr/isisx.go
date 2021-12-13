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
	"fmt"
	"strconv"

	"go.ligato.io/vpp-agent/v3/client"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	vppisisx "go.pantheon.tech/stonework/proto/isisx"
	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// isisxPunt implements PuntHandler for PuntRequest_ISISX
type isisxPunt struct {
}

func NewIsisxPuntHandler() PuntHandler {
	return &isisxPunt{}
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *isisxPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if isisx := punt.GetIsisx(); isisx != nil {
		return []InterconnectReq{
			{
				link: &InterfaceLink{
					vrf:               isisx.Vrf,
					withoutCNFVrf:     isisx.WithoutCnfVrf,
					unnumberedToIface: isisx.VppInterface,
				},
				// Selector = interface name, it is the same as used by Hairpin, Hairpin XConnect and ABX,
				// => mutually exclusive with them if they can't multiplex
				// => mutually exclusive with Hairpin, Hairpin XConnect, but sharing TAP interconnection
				// with ABX (ABX and ISISX can multiplex)
				vppSelector: VppInterfaceSelector(isisx.VppInterface),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *isisxPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	// L3 VPP interface
	if isisx := punt.GetIsisx(); isisx != nil {
		deps = append(deps,
			kvs.Dependency{
				Label: punt.GetLabel() + "-isisx-" + isisx.GetVppInterface(),
				AnyOf: kvs.AnyOfDependency{
					KeyPrefixes: []string{vpp_interfaces.InterfaceAddressPrefix(isisx.GetVppInterface())},
				},
			})
		if isisx.GetVrf() != 0 {
			// interface is inside the VRF (irrelevant whether it is IPv4 or IPv6 VRF)
			deps = append(deps, kvs.Dependency{
				Label: fmt.Sprintf("%s-isisx-vrf-%d", punt.GetLabel(), isisx.GetVrf()),
				AnyOf: kvs.AnyOfDependency{
					KeyPrefixes: []string{
						vpp_interfaces.InterfaceVrfKeyPrefix(isisx.GetVppInterface()) + strconv.Itoa(int(isisx.GetVrf())),
					},
				},
			})
		}
	}
	return
}

// CanMultiplex enables interconnection multiplexing for this punting. It could be enabled in certain cases:
// 1. two or more punts of this type can coexist even if they have the same vpp selector
// 2. one or more punts of this type can coexist with other type of punts on the same (TAP-only)
// interconnection if they all have the same vpp selector and cnf selector.
// The TAP-backed interconnection is shared for multiple multiplexing punts with the same cnf selector
// (same network namespace) and vpp selector.
func (p *isisxPunt) CanMultiplex() bool {
	return true // see vppSelector for explanation
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *isisxPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	// create bidirectional ISIS protocol tunnel
	toCnf := &vppisisx.ISISXConnection{
		InputInterface:  puntReq.GetIsisx().GetVppInterface(),
		OutputInterface: interconnects[0].GetVppInterface().Name,
	}
	fromCnf := &vppisisx.ISISXConnection{
		InputInterface:  toCnf.OutputInterface,
		OutputInterface: toCnf.InputInterface,
	}

	if remove {
		txn.Delete(toCnf, fromCnf)
	} else {
		txn.Update(toCnf, fromCnf)
	}

	return nil
}

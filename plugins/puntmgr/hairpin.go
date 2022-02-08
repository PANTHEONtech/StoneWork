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
	"strings"

	"go.ligato.io/vpp-agent/v3/client"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	vpp_l2 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l2"
	vpp_l3 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"
	"google.golang.org/protobuf/proto"

	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// hairpinPunt implements PuntHandler for PuntRequest_HAIRPIN
type hairpinPunt struct{}

func NewHairpinPuntHandler() PuntHandler {
	return &hairpinPunt{}
}

// HairpinInterfaceSelector is used only by Hairpin to ensure that no two hairpin punt requests would try to create
// hairpin interface of the same name.
func HairpinInterfaceSelector(ifaceName string) string {
	return "vpp/hairpin/interface/" + ifaceName
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *hairpinPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if hairpin := punt.GetHairpin(); hairpin != nil {
		return []InterconnectReq{
			{
				link: &InterfaceLink{},
				// Selector = interface name, i.e. same as used by Hairpin-XConnect and ABX, both of which
				// are mutually exclusive with Hairpin.
				vppSelector: VppInterfaceSelector(hairpin.VppInterface),
			},
			{
				link: &InterfaceLink{
					interfaceName:  hairpin.HairpinInterface.Name,
					physAddress:    hairpin.HairpinInterface.PhysAddress,
					ipAddresses:    hairpin.HairpinInterface.IpAddresses,
					vrf:            hairpin.HairpinInterface.Vrf,
					withDhcpClient: hairpin.HairpinInterface.WithDhcpClient,
					mtu:            hairpin.HairpinInterface.Mtu,
				},
				vppSelector: HairpinInterfaceSelector(hairpin.HairpinInterface.Name),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *hairpinPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	// L2 VPP interfaces
	if hairpin := punt.GetHairpin(); hairpin != nil {
		deps = append(deps,
			kvs.Dependency{
				Label: punt.GetLabel() + "-hairpin-" + hairpin.GetVppInterface(),
				Key:   vpp_interfaces.InterfaceKey(hairpin.GetVppInterface()),
			})
		if vrf := hairpin.GetHairpinInterface().GetVrf(); vrf != 0 {
			hasIpv4, hasIpv6 := getIPAddressVersions(hairpin.GetHairpinInterface().GetIpAddresses())
			if hasIpv4 {
				deps = append(deps, kvs.Dependency{
					Label: fmt.Sprintf("%s-hairpin-vrf-v4-%d", punt.GetLabel(), vrf),
					Key:   vpp_l3.VrfTableKey(vrf, vpp_l3.VrfTable_IPV4),
				})
			}
			if hasIpv6 {
				deps = append(deps, kvs.Dependency{
					Label: fmt.Sprintf("%s-hairpin-vrf-v6-%d", punt.GetLabel(), vrf),
					Key:   vpp_l3.VrfTableKey(vrf, vpp_l3.VrfTable_IPV6),
				})
			}
		}
	}
	return deps
}

// CanMultiplex enables interconnection multiplexing for this punting. It could be enabled in certain cases:
// 1. two or more punts of this type can coexist even if they have the same vpp selector
// 2. one or more punts of this type can coexist with other type of punts on the same (TAP-only)
// interconnection if they all have the same vpp selector and cnf selector.
// The TAP-backed interconnection is shared for multiple multiplexing punts with the same cnf selector
// (same network namespace) and vpp selector.
func (p *hairpinPunt) CanMultiplex() bool {
	return false
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *hairpinPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	vppIface := puntReq.GetHairpin().VppInterface
	var icIface *pb.PuntMetadata_Interface
	if interconnects[0].Id.VppSelector == VppInterfaceSelector(vppIface) {
		icIface = interconnects[0].VppInterface
	} else {
		icIface = interconnects[1].VppInterface
	}

	config := []proto.Message{
		&vpp_l2.XConnectPair{
			ReceiveInterface:  vppIface,
			TransmitInterface: icIface.Name,
		},
		&vpp_l2.XConnectPair{
			ReceiveInterface:  icIface.Name,
			TransmitInterface: vppIface,
		},
	}
	if remove {
		txn.Delete(config...)
	} else {
		txn.Update(config...)
	}
	return nil
}

// getIPAddressVersions returns two flags to tell whether the provided list of addresses
// contains IPv4 and/or IPv6 type addresses
func getIPAddressVersions(ipAddrs []string) (hasIPv4, hasIPv6 bool) {
	for _, ip := range ipAddrs {
		if strings.Contains(ip, ":") {
			hasIPv6 = true
		} else {
			hasIPv4 = true
		}
	}
	return
}

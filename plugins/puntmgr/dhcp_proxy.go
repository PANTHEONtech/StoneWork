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
	"strings"

	"go.ligato.io/vpp-agent/v3/client"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_l3 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"

	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// dhcpProxyPunt implements PuntHandler for PuntRequest_DHCP_PROXY
type dhcpProxyPunt struct{}

func NewDhcpProxyPuntHandler() PuntHandler {
	return &dhcpProxyPunt{}
}

// VrfSelector ensures that there is at most one DHCP proxy configured for a given VRF.
func VrfSelector(vrf uint32) string {
	return "vpp/vrf/" + strconv.Itoa(int(vrf))
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *dhcpProxyPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if dhcpProxy := punt.GetDhcpProxy(); dhcpProxy != nil {
		return []InterconnectReq{
			{
				link: &InterfaceLink{
					vrf:            dhcpProxy.Vrf,
					withoutCNFVrf:  dhcpProxy.WithoutCnfVrf,
					allocateSubnet: true,
				},
				vppSelector: VrfSelector(dhcpProxy.Vrf),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *dhcpProxyPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	if dhcpProxy := punt.GetDhcpProxy(); dhcpProxy != nil {
		if vrf := dhcpProxy.GetVrf(); vrf != 0 {
			deps = append(deps, kvs.Dependency{
				Label: fmt.Sprintf("%s-dhcp-proxy-vrf-%d", punt.GetLabel(), vrf),
				Key:   vpp_l3.VrfTableKey(vrf, vpp_l3.VrfTable_IPV4),
			})
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
func (p *dhcpProxyPunt) CanMultiplex() bool {
	return false
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *dhcpProxyPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	vppIP := strings.SplitN(interconnects[0].VppInterface.IpAddresses[0], "/", 2)[0]
	cnfIP := strings.SplitN(interconnects[0].CnfInterface.IpAddresses[0], "/", 2)[0]
	dhcpProxy := &vpp_l3.DHCPProxy{
		SourceIpAddress: vppIP,
		RxVrfId:         puntReq.GetDhcpProxy().GetVrf(),
		Servers: []*vpp_l3.DHCPProxy_DHCPServer{
			{
				VrfId:     puntReq.GetDhcpProxy().GetVrf(),
				IpAddress: cnfIP,
			},
		},
	}
	if remove {
		txn.Delete(dhcpProxy)
	} else {
		txn.Update(dhcpProxy)
	}
	return nil
}

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
	"hash/fnv"
	"net"
	"strconv"

	"go.ligato.io/vpp-agent/v3/client"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"
	vppacl "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	vppabx "go.pantheon.tech/stonework/proto/abx"

	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

const (
	anyIPv4Addr      = "0.0.0.0"
	anyIPv6Addr      = "::"
	anyAddressPrefix = "/0"
	abxIndexPrefix   = 111 // 8 bits

	anyAddrAlias   = "any"
	localAddrAlias = "local"
)

// abxPunt implements PuntHandler for PuntRequest_ABX
type abxPunt struct {
	ifPlugin ifplugin.API
	// ABX Punt handler needs in-memory cache for the sake of multiplexing
	abxPunts map[string][]abxPuntMeta
}

type abxPuntMeta struct {
	puntId   puntID
	puntReq  *pb.PuntRequest
	icIface  string
	priority uint32
}

func NewAbxPuntHandler(ifPlugin ifplugin.API) PuntHandler {
	return &abxPunt{
		ifPlugin: ifPlugin,
		abxPunts: make(map[string][]abxPuntMeta),
	}
}

// VppInterfaceSelector selects (all/some) packets received or sent through a given VPP interface.
// It can be used by punts that need to reserve entire interface and cannot share it with other
// punt types (that use this same selector). The same interface can be punted multiple times
// only within the same punt type if it supports multiplexing like it is the case with ABX.
func VppInterfaceSelector(ifaceName string) string {
	return "vpp/interface/" + ifaceName
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *abxPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if abx := punt.GetAbx(); abx != nil {
		return []InterconnectReq{
			{
				link: &InterfaceLink{
					vrf:               abx.Vrf,
					withoutCNFVrf:     abx.WithoutCnfVrf,
					unnumberedToIface: abx.VppInterface,
				},
				// Selector = interface name, i.e. same as used by Hairpin and Hairpin XConnect, both of which
				// are mutually exclusive with ABX.
				vppSelector: VppInterfaceSelector(abx.VppInterface),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *abxPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	// L3 VPP interface
	if abx := punt.GetAbx(); abx != nil {
		deps = append(deps,
			kvs.Dependency{
				Label: punt.GetLabel() + "-abx-" + abx.GetVppInterface(),
				AnyOf: kvs.AnyOfDependency{
					KeyPrefixes: []string{vpp_interfaces.InterfaceAddressPrefix(abx.GetVppInterface())},
				},
			})
		if abx.GetVrf() != 0 {
			// interface is inside the VRF (irrelevant whether it is IPv4 or IPv6 VRF)
			deps = append(deps, kvs.Dependency{
				Label: fmt.Sprintf("%s-abx-vrf-%d", punt.GetLabel(), abx.GetVrf()),
				AnyOf: kvs.AnyOfDependency{
					KeyPrefixes: []string{
						vpp_interfaces.InterfaceVrfKeyPrefix(abx.GetVppInterface()) + strconv.Itoa(int(abx.GetVrf())),
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
func (p *abxPunt) CanMultiplex() bool {
	return true
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *abxPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	interconnect := interconnects[0]
	vppInterface := puntReq.GetAbx().GetVppInterface()

	// obtain priority for this ABX
	var abxPrio, maxPrio uint32
	for _, abx := range p.abxPunts[vppInterface] {
		if abx.priority > maxPrio {
			maxPrio = abx.priority
		}
		if abx.icIface == interconnect.VppInterface.Name {
			abxPrio = abx.priority
		}
	}
	if abxPrio == 0 {
		abxPrio = maxPrio + 1
	}

	// update the in-memory cache used for multiplexing
	if remove {
		var filtered []abxPuntMeta
		for _, abx := range p.abxPunts[vppInterface] {
			if abx.puntId != puntId {
				filtered = append(filtered, abx)
			}
		}
		if len(filtered) == 0 {
			delete(p.abxPunts, vppInterface)
		} else {
			p.abxPunts[vppInterface] = filtered
		}
	} else {
		p.abxPunts[vppInterface] = append(p.abxPunts[vppInterface],
			abxPuntMeta{
				icIface:  interconnect.VppInterface.Name,
				puntId:   puntId,
				puntReq:  puntReq,
				priority: abxPrio,
			})
	}

	ingressAcl := &vppacl.ACL{
		Name: p.getACLName(vppInterface, interconnect.VppInterface.Name),
	}
	for _, abx := range p.abxPunts[vppInterface] {
		if abx.icIface != interconnect.VppInterface.Name {
			continue
		}
		ingressRules, err := p.processAclRules(vppInterface, abx.puntReq.GetAbx().GetIngressAclRules())
		if err != nil {
			return err
		}
		for _, ipRule := range ingressRules {
			ingressAcl.Rules = append(ingressAcl.Rules, &vppacl.ACL_Rule{
				Action: vppacl.ACL_Rule_PERMIT,
				IpRule: ipRule,
			})
		}
	}

	egressAcl := &vppacl.ACL{
		Name: p.getACLName(interconnect.VppInterface.Name, vppInterface),
	}
	for _, abx := range p.abxPunts[vppInterface] {
		if abx.icIface != interconnect.VppInterface.Name {
			continue
		}
		egressRules, err := p.processAclRules(vppInterface, abx.puntReq.GetAbx().GetEgressAclRules())
		if err != nil {
			return err
		}
		for _, ipRule := range egressRules {
			egressAcl.Rules = append(egressAcl.Rules, &vppacl.ACL_Rule{
				Action: vppacl.ACL_Rule_PERMIT,
				IpRule: ipRule,
			})
		}
	}

	ingressAbx := &vppabx.ABX{
		Index:           p.getABXIndex(vppInterface, interconnect.VppInterface.Name),
		AclName:         ingressAcl.Name,
		OutputInterface: interconnect.VppInterface.Name,
		DstMac:          interconnect.CnfInterface.PhysAddress,
		AttachedInterfaces: []*vppabx.ABX_AttachedInterface{
			{
				InputInterface: vppInterface,
				Priority:       abxPrio,
			},
		},
	}

	egressAbx := &vppabx.ABX{
		Index:           p.getABXIndex(interconnect.VppInterface.Name, vppInterface),
		AclName:         egressAcl.Name,
		OutputInterface: vppInterface,
		AttachedInterfaces: []*vppabx.ABX_AttachedInterface{
			{
				InputInterface: interconnect.VppInterface.Name,
			},
		},
	}

	if remove {
		if len(ingressAcl.Rules) > 0 {
			txn.Update(ingressAcl, ingressAbx)
		} else {
			txn.Delete(ingressAcl, ingressAbx)
		}
		if len(egressAcl.Rules) > 0 {
			txn.Update(egressAcl, egressAbx)
		} else {
			txn.Delete(egressAcl, egressAbx)
		}
	} else {
		txn.Update(ingressAcl, ingressAbx)
		if len(egressAcl.Rules) > 0 {
			txn.Update(egressAcl, egressAbx)
		}
	}
	return nil
}

// processAclRules translates special address constants like "local" and "any" to actual IP addresses.
func (p *abxPunt) processAclRules(vppInterface string, in []*vppacl.ACL_Rule_IpRule) (out []*vppacl.ACL_Rule_IpRule, err error) {
	// obtain the list of IP addresses assigned to the interface
	ifMeta, exists := p.ifPlugin.GetInterfaceIndex().LookupByName(vppInterface)
	if !exists || ifMeta == nil {
		return nil, fmt.Errorf("required VPP interface %s was not found", vppInterface)
	}
	if len(ifMeta.IPAddresses) == 0 {
		return nil, fmt.Errorf("VPP interface %s does not have any IP address assigned", vppInterface)
	}

	// parse IP addresses assigned to the interface
	var (
		hasIPv4, hasIPv6               bool
		localIPv4Addrs, localIPv6Addrs []string
	)
	allOnesIPv4Mask := net.CIDRMask(8*net.IPv4len, 8*net.IPv4len)
	allOnesIPv6Mask := net.CIDRMask(8*net.IPv6len, 8*net.IPv6len)
	for _, ipAddr := range ifMeta.IPAddresses {
		ip, _, err := net.ParseCIDR(ipAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse IP address %s assigned to VPP interface %s: %v",
				ipAddr, vppInterface, err)
		}
		if ip.To4() != nil {
			hasIPv4 = true
			ipNet := &net.IPNet{
				IP:   ip.To4(),
				Mask: allOnesIPv4Mask,
			}
			localIPv4Addrs = append(localIPv4Addrs, ipNet.String())
		} else {
			hasIPv6 = true
			ipNet := &net.IPNet{
				IP:   ip.To16(),
				Mask: allOnesIPv6Mask,
			}
			localIPv6Addrs = append(localIPv6Addrs, ipNet.String())
		}
	}

	for _, aclRule := range in {
		if aclRule.Ip == nil {
			out = append(out, aclRule)
			continue
		}

		// check for IP address mismatch
		src := aclRule.Ip.SourceNetwork
		dst := aclRule.Ip.DestinationNetwork
		var forIPv4, forIPv6 bool
		forAny := true
		isAddrAlias := func(addr string) bool {
			return addr == "" || addr == anyAddrAlias || addr == localAddrAlias
		}
		parseRealAddr := func(addr string) error {
			if !isAddrAlias(addr) {
				ip, _, err := net.ParseCIDR(addr)
				if err != nil {
					return fmt.Errorf("failed to parse IP network %s used with ABX ACL: %v",
						addr, err)
				}
				if ip.To4() != nil {
					forIPv4 = true
					forAny = false
				} else {
					forIPv6 = true
					forAny = false
				}
			}
			return nil
		}
		if err = parseRealAddr(src); err != nil {
			return nil, err
		}
		if err = parseRealAddr(dst); err != nil {
			return nil, err
		}
		if forIPv4 && !hasIPv4 {
			return nil, fmt.Errorf("ABX ACL IP rule %v is for IPv4 but interface %s has no IPv4 address assigned",
				aclRule.Ip, vppInterface)
		}
		if forIPv6 && !hasIPv6 {
			return nil, fmt.Errorf("ABX ACL IP rule %v is for IPv6 but interface %s has no IPv6 address assigned",
				aclRule.Ip, vppInterface)
		}

		// translate address aliases to actual IP addresses
		if !isAddrAlias(src) && !isAddrAlias(dst) {
			// nothing to do
			out = append(out, aclRule)
			continue
		}
		combineAddrs := func(srcAddrs, dstAddrs []string) {
			for _, src := range srcAddrs {
				for _, dst := range dstAddrs {
					out = append(out, &vppacl.ACL_Rule_IpRule{
						Ip: &vppacl.ACL_Rule_IpRule_Ip{
							SourceNetwork:      src,
							DestinationNetwork: dst,
							Protocol:           aclRule.Ip.Protocol,
						},
						Icmp: aclRule.Icmp,
						Tcp:  aclRule.Tcp,
						Udp:  aclRule.Udp,
					})
				}
			}
		}
		// --> IPv4
		if forIPv4 || (forAny && hasIPv4) {
			translateAlias := func(alias string) []string {
				switch {
				case !isAddrAlias(alias):
					return []string{alias}
				case alias == "":
					fallthrough
				case alias == anyAddrAlias:
					return []string{anyIPv4Addr + anyAddressPrefix}
				case alias == localAddrAlias:
					return localIPv4Addrs
				}
				return nil // unreachable
			}
			combineAddrs(translateAlias(src), translateAlias(dst))
		}
		// --> IPv6
		if forIPv6 || (forAny && hasIPv6) {
			translateAlias := func(alias string) []string {
				switch {
				case !isAddrAlias(alias):
					return []string{alias}
				case alias == "":
					fallthrough
				case alias == anyAddrAlias:
					return []string{anyIPv6Addr + anyAddressPrefix}
				case alias == localAddrAlias:
					return localIPv6Addrs
				}
				return nil // unreachable
			}
			combineAddrs(translateAlias(src), translateAlias(dst))
		}
	}
	return out, nil
}

func (p *abxPunt) getAbxLabel(inputIface, outputIface string) string {
	return fmt.Sprintf("abx/from/%s/to/%s", inputIface, outputIface)
}

// GetACLName returns name of the ACL used to match traffic for ABX-based packet punting.
func (p *abxPunt) getACLName(inputIface, outputIface string) string {
	return p.getAbxLabel(inputIface, outputIface)
}

// getABXIndex returns the identifier of ABX rule used for ABX-based packet punting.
func (p *abxPunt) getABXIndex(inputIface, outputIface string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(p.getAbxLabel(inputIface, outputIface)))
	hash := h.Sum32()
	return (uint32(abxIndexPrefix) << 24) | (hash & 0xffffff)
}

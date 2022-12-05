// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 PANTHEON.tech
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

package vpp2210

import (
	"encoding/binary"
	"fmt"
	"net"

	natba "go.pantheon.tech/stonework/plugins/binapi/vpp2210/nat64"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2210/nat_types"
	nat "go.pantheon.tech/stonework/proto/nat64"
)

// Num protocol representation
const (
	ICMP uint8 = 1
	TCP  uint8 = 6
	UDP  uint8 = 17
)

// Nat64IPv6PrefixDump dumps all IPv6 prefixes configured for NAT64.
func (h *Nat64VppHandler) Nat64IPv6PrefixDump() (prefixes []*nat.Nat64IPv6Prefix, err error) {
	req := &natba.Nat64PrefixDump{}
	reqContext := h.callsChannel.SendMultiRequest(req)

	for {
		msg := &natba.Nat64PrefixDetails{}
		stop, err := reqContext.ReceiveReply(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to dump NAT64 IPv6 prefixes: %v", err)
		}
		if stop {
			break
		}
		prefixes = append(prefixes, &nat.Nat64IPv6Prefix{
			VrfId:  msg.VrfID,
			Prefix: msg.Prefix.ToIPNet().String(),
		})
	}

	return prefixes, nil
}

// Nat64InterfacesDump dumps NAT64 config of all NAT64-enabled interfaces.
func (h *Nat64VppHandler) Nat64InterfacesDump() (interfaces []*nat.Nat64Interface, err error) {
	req := &natba.Nat64InterfaceDump{}
	reqContext := h.callsChannel.SendMultiRequest(req)

	for {
		msg := &natba.Nat64InterfaceDetails{}
		stop, err := reqContext.ReceiveReply(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to dump NAT64 interfaces: %v", err)
		}
		if stop {
			break
		}
		ifName, _, found := h.ifIndexes.LookupBySwIfIndex(uint32(msg.SwIfIndex))
		if !found {
			h.log.Warnf("Interface with index %d not found in the mapping", msg.SwIfIndex)
			continue
		}
		ifType := nat.Nat64Interface_IPV6_INSIDE
		if msg.Flags&nat_types.NAT_IS_OUTSIDE != 0 {
			ifType = nat.Nat64Interface_IPV4_OUTSIDE
		}
		interfaces = append(interfaces, &nat.Nat64Interface{
			Name: ifName,
			Type: ifType,
		})
	}
	return interfaces, nil
}

// Nat64AddressPoolsDump dumps all configured NAT64 address pools.
// Note that VPP returns configured addresses one-by-one, loosing information about address pools
// configured with multiple addresses through IP ranges. Provide expected configuration to group
// addresses from the same range.
func (h *Nat64VppHandler) Nat64AddressPoolsDump(correlateWith []*nat.Nat64AddressPool) (pools []*nat.Nat64AddressPool, err error) {
	req := &natba.Nat64PoolAddrDump{}
	reqContext := h.callsChannel.SendMultiRequest(req)

	for {
		msg := &natba.Nat64PoolAddrDetails{}
		stop, err := reqContext.ReceiveReply(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to dump NAT64 address pools: %v", err)
		}
		if stop {
			break
		}
		address := net.IP(msg.Address[:]).String()
		pools = append(pools, &nat.Nat64AddressPool{
			VrfId:   msg.VrfID,
			FirstIp: address,
			LastIp:  address,
		})
	}
	return correlateAddressPools(pools, correlateWith), nil
}

// Nat64StaticBIBsDump dumps NAT64 static bindings.
func (h *Nat64VppHandler) Nat64StaticBIBsDump() (bibs []*nat.Nat64StaticBIB, err error) {
	req := &natba.Nat64BibDump{
		Proto: ^uint8(0), // ALL
	}
	reqContext := h.callsChannel.SendMultiRequest(req)

	for {
		msg := &natba.Nat64BibDetails{}
		stop, err := reqContext.ReceiveReply(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to dump NAT64 static BIBs: %v", err)
		}
		if stop {
			break
		}
		if msg.Flags&nat_types.NAT_IS_STATIC == 0 {
			// dynamic entry
			continue
		}
		var proto nat.Nat64StaticBIB_Protocol
		switch msg.Proto {
		case TCP:
			proto = nat.Nat64StaticBIB_TCP
		case UDP:
			proto = nat.Nat64StaticBIB_UDP
		case ICMP:
			proto = nat.Nat64StaticBIB_ICMP
		}
		bibs = append(bibs, &nat.Nat64StaticBIB{
			VrfId:              msg.VrfID,
			Protocol:           proto,
			InsideIpv6Address:  net.IP(msg.IAddr[:]).String(),
			InsidePort:         uint32(msg.IPort),
			OutsideIpv4Address: net.IP(msg.OAddr[:]).String(),
			OutsidePort:        uint32(msg.OPort),
		})
	}
	return bibs, nil
}

func correlateAddressPools(dumped, correlateWith []*nat.Nat64AddressPool) (correlated []*nat.Nat64AddressPool) {
	if len(correlateWith) == 0 {
		return dumped
	}
	dumpedMap := make(map[uint32]map[uint32]*nat.Nat64AddressPool) // VRF -> address as int -> pool
	for _, pool := range dumped {
		if _, initialized := dumpedMap[pool.VrfId]; !initialized {
			dumpedMap[pool.VrfId] = make(map[uint32]*nat.Nat64AddressPool)
		}
		dumpedMap[pool.VrfId][ip2int(net.ParseIP(pool.FirstIp))] = pool
	}
	for _, correlate := range correlateWith {
		var firstIP, lastIP net.IP
		firstIP = net.ParseIP(correlate.FirstIp)
		if correlate.LastIp != "" {
			lastIP = net.ParseIP(correlate.LastIp)
		} else {
			lastIP = firstIP
		}
		if firstIP == nil || lastIP == nil {
			continue
		}
		first := ip2int(firstIP)
		last := ip2int(lastIP)
		if first >= last {
			continue
		}
		allDumped := true
		for i := first; i <= last; i++ {
			if _, isDumped := dumpedMap[correlate.VrfId][i]; !isDumped {
				allDumped = false
				break
			}
		}
		if allDumped {
			for i := first; i <= last; i++ {
				delete(dumpedMap[correlate.VrfId], i)
			}
			correlated = append(correlated, correlate)
		}
	}
	for vrf := range dumpedMap {
		for _, pool := range dumpedMap[vrf] {
			correlated = append(correlated, pool)
		}
	}
	return correlated
}

func ip2int(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

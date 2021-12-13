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

package vpp2106

import (
	"fmt"
	"net"

	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/interface_types"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/ip_types"
	natba "go.pantheon.tech/stonework/plugins/binapi/vpp2106/nat64"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/nat_types"

	"go.pantheon.tech/stonework/proto/nat64"
)

// AddNat64IPv6Prefix adds IPv6 prefix for NAT64 (used to embed IPv4 address).
func (h *Nat64VppHandler) AddNat64IPv6Prefix(vrf uint32, prefix string) error {
	return h.handleNat64IPv6Prefix(vrf, prefix, true)
}

// DelNat64IPv6Prefix removes existing IPv6 prefix previously configured for NAT64.
func (h *Nat64VppHandler) DelNat64IPv6Prefix(vrf uint32, prefix string) error {
	return h.handleNat64IPv6Prefix(vrf, prefix, false)
}

// EnableNat64Interface enables NAT64 for provided interface.
func (h *Nat64VppHandler) EnableNat64Interface(iface string, natIfaceType nat64.Nat64Interface_Type) error {
	return h.handleNat64Interface(iface, natIfaceType, true)
}

// DisableNat64Interface disables NAT64 for provided interface.
func (h *Nat64VppHandler) DisableNat64Interface(iface string, natIfaceType nat64.Nat64Interface_Type) error {
	return h.handleNat64Interface(iface, natIfaceType, false)
}

// AddNat64AddressPool adds new IPV4 address pool into the NAT64 pools.
func (h *Nat64VppHandler) AddNat64AddressPool(vrf uint32, firstIP, lastIP string) error {
	return h.handleNat64AddressPool(vrf, firstIP, lastIP, true)
}

// DelNat64AddressPool removes existing IPv4 address pool from the NAT64 pools.
func (h *Nat64VppHandler) DelNat64AddressPool(vrf uint32, firstIP, lastIP string) error {
	return h.handleNat64AddressPool(vrf, firstIP, lastIP, false)
}

// AddNat64StaticBIB creates new NAT64 static binding.
func (h *Nat64VppHandler) AddNat64StaticBIB(bib *nat64.Nat64StaticBIB) error {
	return h.handleNat64StaticBIB(bib, true)
}

// DelNat64StaticBIB removes existing NAT64 static binding.
func (h *Nat64VppHandler) DelNat64StaticBIB(bib *nat64.Nat64StaticBIB) error {
	return h.handleNat64StaticBIB(bib, false)
}

// Calls VPP binary API to set/unset NAT64 IPv6 prefix.
func (h *Nat64VppHandler) handleNat64IPv6Prefix(vrf uint32, prefix string, isAdd bool) error {
	ipv6Prefix, err := ipTo6Prefix(prefix)
	if err != nil {
		return fmt.Errorf("unable to parse IPv6 prefix %s: %v", prefix, err)
	}
	req := &natba.Nat64AddDelPrefix{
		VrfID:  vrf,
		Prefix: ipv6Prefix,
		IsAdd:  isAdd,
	}
	reply := &natba.Nat64AddDelPrefixReply{}
	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return err
	}
	return nil
}

// Calls VPP binary API to set/unset interface NAT64 feature.
func (h *Nat64VppHandler) handleNat64Interface(iface string, natIfaceType nat64.Nat64Interface_Type, isAdd bool) error {
	// get interface metadata
	ifaceMeta, found := h.ifIndexes.LookupByName(iface)
	if !found {
		return fmt.Errorf("failed to get metadata for interface: %s", iface)
	}
	var flags nat_types.NatConfigFlags
	switch natIfaceType {
	case nat64.Nat64Interface_IPV6_INSIDE:
		flags = nat_types.NAT_IS_INSIDE
	case nat64.Nat64Interface_IPV4_OUTSIDE:
		flags = nat_types.NAT_IS_OUTSIDE
	}
	req := &natba.Nat64AddDelInterface{
		SwIfIndex: interface_types.InterfaceIndex(ifaceMeta.SwIfIndex),
		Flags:     flags,
		IsAdd:     isAdd,
	}
	reply := &natba.Nat64AddDelInterfaceReply{}
	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return err
	}
	return nil
}

// Calls VPP binary API to add/del NAT64 address pool.
func (h *Nat64VppHandler) handleNat64AddressPool(vrf uint32, firstIP, lastIP string, isAdd bool) error {
	startAddr, err := ip_types.ParseIP4Address(firstIP)
	if err != nil {
		return fmt.Errorf("unable to parse address %s from the NAT64 pool: %v", firstIP, err)
	}
	endAddr := startAddr
	if lastIP != "" {
		endAddr, err = ip_types.ParseIP4Address(lastIP)
		if err != nil {
			return fmt.Errorf("unable to parse address %s from the NAT64 pool: %v", lastIP, err)
		}
	}
	req := &natba.Nat64AddDelPoolAddrRange{
		StartAddr: startAddr,
		EndAddr:   endAddr,
		VrfID:     vrf,
		IsAdd:     isAdd,
	}
	reply := &natba.Nat64AddDelPoolAddrRangeReply{}
	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return err
	}
	return nil
}

// Calls VPP binary API to add/del NAT64 static binding.
func (h *Nat64VppHandler) handleNat64StaticBIB(bib *nat64.Nat64StaticBIB, isAdd bool) error {
	inAddr, err := ip_types.ParseIP6Address(bib.GetInsideIpv6Address())
	if err != nil {
		return fmt.Errorf("unable to parse inside IPv6 address (%s): %v", bib.GetInsideIpv6Address(), err)
	}
	outAddr, err := ip_types.ParseIP4Address(bib.GetOutsideIpv4Address())
	if err != nil {
		return fmt.Errorf("unable to parse outside IPv4 address (%s): %v", bib.GetOutsideIpv4Address(), err)
	}
	var proto uint8
	switch bib.GetProtocol() {
	case nat64.Nat64StaticBIB_TCP:
		proto = TCP
	case nat64.Nat64StaticBIB_UDP:
		proto = UDP
	case nat64.Nat64StaticBIB_ICMP:
		proto = ICMP
	default:
		h.log.Warnf("Unknown protocol %v, defaulting to TCP", bib.GetProtocol())
		proto = TCP
	}
	req := &natba.Nat64AddDelStaticBib{
		IAddr: inAddr,
		OAddr: outAddr,
		IPort: uint16(bib.GetInsidePort()),
		OPort: uint16(bib.GetOutsidePort()),
		Proto: proto,
		VrfID: bib.GetVrfId(),
		IsAdd: isAdd,
	}
	reply := &natba.Nat64AddDelStaticBibReply{}
	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return err
	}
	return nil
}

func ipTo6Prefix(ipStr string) (prefix ip_types.IP6Prefix, err error) {
	_, netIP, err := net.ParseCIDR(ipStr)
	if err != nil {
		return prefix, fmt.Errorf("invalid IP (%q): %v", ipStr, err)
	}
	if ip4 := netIP.IP.To4(); ip4 != nil {
		return prefix, fmt.Errorf("required IPv6, provided IPv4 prefix: %q", ipStr)
	}
	copy(prefix.Address[:], netIP.IP.To16())
	prefixLen, _ := netIP.Mask.Size()
	prefix.Len = uint8(prefixLen)
	return
}

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

package vpp2106_test

import (
	"net"
	"sort"
	"testing"

	. "github.com/onsi/gomega"

	"go.ligato.io/cn-infra/v2/logging/logrus"

	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/ip_types"
	natba "go.pantheon.tech/stonework/plugins/binapi/vpp2106/nat64"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/nat_types"
	vpp_vpe "go.pantheon.tech/stonework/plugins/binapi/vpp2106/vpe"

	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/vppmock"

	"go.pantheon.tech/stonework/plugins/nat64/vppcalls"
	"go.pantheon.tech/stonework/plugins/nat64/vppcalls/vpp2106"
	"go.pantheon.tech/stonework/proto/nat64"
)

func TestNat64IPv6PrefixDump(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(
		&natba.Nat64PrefixDetails{
			VrfID:  5,
			Prefix: ipTo6Prefix("64:ff9b::/96"),
		},
		&natba.Nat64PrefixDetails{
			VrfID:  3,
			Prefix: ipTo6Prefix("2004::/32"),
		})
	ctx.MockVpp.MockReply(&vpp_vpe.ControlPingReply{})

	prefixes, err := natHandler.Nat64IPv6PrefixDump()
	Expect(err).To(Succeed())

	Expect(prefixes).To(HaveLen(2))

	Expect(prefixes[0].VrfId).To(BeEquivalentTo(5))
	Expect(prefixes[0].Prefix).To(Equal("64:ff9b::/96"))

	Expect(prefixes[1].VrfId).To(BeEquivalentTo(3))
	Expect(prefixes[1].Prefix).To(Equal("2004::/32"))
}

func TestNat64InterfacesDump(t *testing.T) {
	ctx, natHandler, swIfIndexes := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(
		&natba.Nat64InterfaceDetails{
			SwIfIndex: 1,
			Flags:     nat_types.NAT_IS_OUTSIDE,
		},
		&natba.Nat64InterfaceDetails{
			SwIfIndex: 2,
			Flags:     nat_types.NAT_IS_INSIDE,
		})
	ctx.MockVpp.MockReply(&vpp_vpe.ControlPingReply{})

	swIfIndexes.Put("if0", &ifaceidx.IfaceMetadata{SwIfIndex: 1})
	swIfIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 2})

	interfaces, err := natHandler.Nat64InterfacesDump()
	Expect(err).To(Succeed())

	Expect(interfaces).To(HaveLen(2))

	Expect(interfaces[0].Name).To(Equal("if0"))
	Expect(interfaces[0].Type).To(Equal(nat64.Nat64Interface_IPV4_OUTSIDE))

	Expect(interfaces[1].Name).To(Equal("if1"))
	Expect(interfaces[1].Type).To(Equal(nat64.Nat64Interface_IPV6_INSIDE))
}

func TestNat64AddressPoolsDump(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	// address pool
	ctx.MockVpp.MockReply(
		&natba.Nat64PoolAddrDetails{
			Address: ipTo4Address("10.10.1.1"),
			VrfID:   1,
		},
		&natba.Nat64PoolAddrDetails{
			Address: ipTo4Address("80.80.80.2"),
			VrfID:   2,
		},
		&natba.Nat64PoolAddrDetails{
			Address: ipTo4Address("192.168.10.3"),
			VrfID:   2,
		},
		&natba.Nat64PoolAddrDetails{
			Address: ipTo4Address("192.168.10.4"),
			VrfID:   2,
		},
		&natba.Nat64PoolAddrDetails{
			Address: ipTo4Address("192.168.10.5"),
			VrfID:   2,
		})
	ctx.MockVpp.MockReply(&vpp_vpe.ControlPingReply{})

	pools, err := natHandler.Nat64AddressPoolsDump([]*nat64.Nat64AddressPool{
		{
			VrfId:   2,
			FirstIp: "192.168.10.3",
			LastIp:  "192.168.10.5",
		},
	})
	Expect(err).To(Succeed())

	Expect(pools).To(HaveLen(3))

	sort.Slice(pools, func(i, j int) bool {
		return pools[i].FirstIp < pools[j].FirstIp
	})

	Expect(pools[0].FirstIp).To(Equal("10.10.1.1"))
	Expect(pools[0].LastIp).To(Equal("10.10.1.1"))
	Expect(pools[0].VrfId).To(BeEquivalentTo(1))

	Expect(pools[1].FirstIp).To(Equal("192.168.10.3"))
	Expect(pools[1].LastIp).To(Equal("192.168.10.5"))
	Expect(pools[1].VrfId).To(BeEquivalentTo(2))

	Expect(pools[2].FirstIp).To(Equal("80.80.80.2"))
	Expect(pools[2].LastIp).To(Equal("80.80.80.2"))
	Expect(pools[2].VrfId).To(BeEquivalentTo(2))
}

func TestNat64StaticBIBsDump(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(
		&natba.Nat64BibDetails{
			IAddr: ipTo6Address("2000::3"),
			IPort: 8080,
			OAddr: ipTo4Address("172.16.2.3"),
			OPort: 80,
			VrfID: 5,
			Proto: 6, // TCP
			Flags: nat_types.NAT_IS_STATIC,
		},
		&natba.Nat64BibDetails{
			IAddr: ipTo6Address("2000::8"),
			IPort: 9090,
			OAddr: ipTo4Address("10.10.5.5"),
			OPort: 90,
			VrfID: 0,
			Proto: 17, // UDP
			Flags: nat_types.NAT_IS_STATIC,
		},
		&natba.Nat64BibDetails{
			IAddr: ipTo6Address("2000::10"),
			OAddr: ipTo4Address("80.80.80.1"),
			VrfID: 4,
			Proto: 1, // ICMP
			Flags: nat_types.NAT_IS_STATIC,
		})
	ctx.MockVpp.MockReply(&vpp_vpe.ControlPingReply{})

	bibs, err := natHandler.Nat64StaticBIBsDump()
	Expect(err).To(Succeed())

	Expect(bibs).To(HaveLen(3))

	Expect(bibs[0].InsideIpv6Address).To(Equal("2000::3"))
	Expect(bibs[0].InsidePort).To(BeEquivalentTo(8080))
	Expect(bibs[0].OutsideIpv4Address).To(Equal("172.16.2.3"))
	Expect(bibs[0].OutsidePort).To(BeEquivalentTo(80))
	Expect(bibs[0].VrfId).To(BeEquivalentTo(5))
	Expect(bibs[0].Protocol).To(Equal(nat64.Nat64StaticBIB_TCP))

	Expect(bibs[1].InsideIpv6Address).To(Equal("2000::8"))
	Expect(bibs[1].InsidePort).To(BeEquivalentTo(9090))
	Expect(bibs[1].OutsideIpv4Address).To(Equal("10.10.5.5"))
	Expect(bibs[1].OutsidePort).To(BeEquivalentTo(90))
	Expect(bibs[1].VrfId).To(BeEquivalentTo(0))
	Expect(bibs[1].Protocol).To(Equal(nat64.Nat64StaticBIB_UDP))

	Expect(bibs[2].InsideIpv6Address).To(Equal("2000::10"))
	Expect(bibs[2].InsidePort).To(BeEquivalentTo(0))
	Expect(bibs[2].OutsideIpv4Address).To(Equal("80.80.80.1"))
	Expect(bibs[2].OutsidePort).To(BeEquivalentTo(0))
	Expect(bibs[2].VrfId).To(BeEquivalentTo(4))
	Expect(bibs[2].Protocol).To(Equal(nat64.Nat64StaticBIB_ICMP))
}

func ipTo6Address(ipStr string) (addr ip_types.IP6Address) {
	netIP := net.ParseIP(ipStr)
	Expect(netIP).ToNot(BeNil())
	ip4 := netIP.To4()
	Expect(ip4).To(BeNil())
	copy(addr[:], netIP.To16())
	return
}

func ipTo6Prefix(ipStr string) (prefix ip_types.IP6Prefix) {
	_, netIP, err := net.ParseCIDR(ipStr)
	Expect(err).To(Succeed())
	ip4 := netIP.IP.To4()
	Expect(ip4).To(BeNil())
	copy(prefix.Address[:], netIP.IP.To16())
	prefixLen, _ := netIP.Mask.Size()
	prefix.Len = uint8(prefixLen)
	return
}

func ipTo4Address(ipStr string) (addr ip_types.IP4Address) {
	netIP := net.ParseIP(ipStr)
	if ip4 := netIP.To4(); ip4 != nil {
		copy(addr[:], ip4)
	}
	return
}

func natTestSetup(t *testing.T) (*vppmock.TestCtx, vppcalls.Nat64VppAPI, ifaceidx.IfaceMetadataIndexRW) {
	ctx := vppmock.SetupTestCtx(t)
	log := logrus.NewLogger("test-log")
	swIfIndexes := ifaceidx.NewIfaceIndex(logrus.DefaultLogger(), "test-sw_if_indexes")
	natHandler := vpp2106.NewNat64VppHandler(ctx.MockChannel, swIfIndexes, log)
	return ctx, natHandler, swIfIndexes
}

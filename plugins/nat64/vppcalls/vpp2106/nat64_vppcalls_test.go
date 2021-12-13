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
	"testing"

	. "github.com/onsi/gomega"

	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"

	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/ip_types"
	natba "go.pantheon.tech/stonework/plugins/binapi/vpp2106/nat64"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/nat_types"
	"go.pantheon.tech/stonework/proto/nat64"
)

func TestAddNat64IPv6Prefix(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	_, prefix1, _ := net.ParseCIDR("64:ff9b::/96")
	_, prefix2, _ := net.ParseCIDR("2004::/32")

	ctx.MockVpp.MockReply(&natba.Nat64AddDelPrefixReply{})
	err := natHandler.AddNat64IPv6Prefix(0, prefix1.String())
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelPrefix)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(ip6PrefixToIPNet(msg.Prefix)).To(BeEquivalentTo(prefix1.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(0))

	ctx.MockVpp.MockReply(&natba.Nat64AddDelPrefixReply{})
	err = natHandler.AddNat64IPv6Prefix(5, prefix2.String())
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelPrefix)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(ip6PrefixToIPNet(msg.Prefix)).To(BeEquivalentTo(prefix2.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(5))

	err = natHandler.AddNat64IPv6Prefix(1, "invalid prefix")
	Expect(err).Should(HaveOccurred())
	err = natHandler.AddNat64IPv6Prefix(1, "192.168.1.0/24")
	Expect(err).Should(HaveOccurred())
}

func TestDelNat64IPv6Prefix(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	_, prefix1, _ := net.ParseCIDR("64:ff9b::/96")
	_, prefix2, _ := net.ParseCIDR("2004::/32")

	ctx.MockVpp.MockReply(&natba.Nat64AddDelPrefixReply{})
	err := natHandler.DelNat64IPv6Prefix(0, prefix1.String())
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelPrefix)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(ip6PrefixToIPNet(msg.Prefix)).To(BeEquivalentTo(prefix1.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(0))

	ctx.MockVpp.MockReply(&natba.Nat64AddDelPrefixReply{})
	err = natHandler.DelNat64IPv6Prefix(5, prefix2.String())
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelPrefix)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(ip6PrefixToIPNet(msg.Prefix)).To(BeEquivalentTo(prefix2.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(5))

	err = natHandler.DelNat64IPv6Prefix(1, "invalid prefix")
	Expect(err).Should(HaveOccurred())
	err = natHandler.DelNat64IPv6Prefix(1, "192.168.1.0/24")
	Expect(err).Should(HaveOccurred())
}

func TestEnableNat64Interface(t *testing.T) {
	ctx, natHandler, swIfIndexes := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	swIfIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 2})

	ctx.MockVpp.MockReply(&natba.Nat64AddDelInterfaceReply{})
	err := natHandler.EnableNat64Interface("if1", nat64.Nat64Interface_IPV4_OUTSIDE)
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelInterface)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(msg.SwIfIndex).To(BeEquivalentTo(2))
	Expect(msg.Flags).To(BeEquivalentTo(nat_types.NAT_IS_OUTSIDE))

	ctx.MockVpp.MockReply(&natba.Nat64AddDelInterfaceReply{})
	err = natHandler.EnableNat64Interface("if1", nat64.Nat64Interface_IPV6_INSIDE)
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelInterface)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(msg.SwIfIndex).To(BeEquivalentTo(2))
	Expect(msg.Flags).To(BeEquivalentTo(nat_types.NAT_IS_INSIDE))

	err = natHandler.EnableNat64Interface("non-existent interface", nat64.Nat64Interface_IPV4_OUTSIDE)
	Expect(err).Should(HaveOccurred())
}

func TestDisableNat64Interface(t *testing.T) {
	ctx, natHandler, swIfIndexes := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	swIfIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 2})

	ctx.MockVpp.MockReply(&natba.Nat64AddDelInterfaceReply{})
	err := natHandler.DisableNat64Interface("if1", nat64.Nat64Interface_IPV4_OUTSIDE)
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelInterface)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(msg.SwIfIndex).To(BeEquivalentTo(2))
	Expect(msg.Flags).To(BeEquivalentTo(nat_types.NAT_IS_OUTSIDE))

	ctx.MockVpp.MockReply(&natba.Nat64AddDelInterfaceReply{})
	err = natHandler.DisableNat64Interface("if1", nat64.Nat64Interface_IPV6_INSIDE)
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelInterface)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(msg.SwIfIndex).To(BeEquivalentTo(2))
	Expect(msg.Flags).To(BeEquivalentTo(nat_types.NAT_IS_INSIDE))

	err = natHandler.DisableNat64Interface("non-existent interface", nat64.Nat64Interface_IPV4_OUTSIDE)
	Expect(err).Should(HaveOccurred())
}

func TestAddNat64AddressPool(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	addr1 := net.ParseIP("10.0.0.1").To4()
	addr2 := net.ParseIP("10.0.0.10").To4()

	// first IP only
	ctx.MockVpp.MockReply(&natba.Nat64AddDelPoolAddrRangeReply{})
	err := natHandler.AddNat64AddressPool(0, addr1.String(), "")
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelPoolAddrRange)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(msg.StartAddr.String()).To(BeEquivalentTo(addr1.String()))
	Expect(msg.EndAddr.String()).To(BeEquivalentTo(addr1.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(0))

	// first IP + last IP
	ctx.MockVpp.MockReply(&natba.Nat64AddDelPoolAddrRangeReply{})
	err = natHandler.AddNat64AddressPool(5, addr1.String(), addr2.String())
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelPoolAddrRange)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(msg.StartAddr.String()).To(BeEquivalentTo(addr1.String()))
	Expect(msg.EndAddr.String()).To(BeEquivalentTo(addr2.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(5))

	err = natHandler.AddNat64AddressPool(0, "invalid address", "")
	Expect(err).Should(HaveOccurred())
	err = natHandler.AddNat64AddressPool(0, "invalid address", addr2.String())
	Expect(err).Should(HaveOccurred())
	err = natHandler.AddNat64AddressPool(0, addr1.String(), "invalid address")
	Expect(err).Should(HaveOccurred())
}

func TestDelNat64AddressPool(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	addr1 := net.ParseIP("10.0.0.1").To4()
	addr2 := net.ParseIP("10.0.0.10").To4()

	// first IP only
	ctx.MockVpp.MockReply(&natba.Nat64AddDelPoolAddrRangeReply{})
	err := natHandler.DelNat64AddressPool(0, addr1.String(), "")
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelPoolAddrRange)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(msg.StartAddr.String()).To(BeEquivalentTo(addr1.String()))
	Expect(msg.EndAddr.String()).To(BeEquivalentTo(addr1.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(0))

	// first IP + last IP
	ctx.MockVpp.MockReply(&natba.Nat64AddDelPoolAddrRangeReply{})
	err = natHandler.DelNat64AddressPool(5, addr1.String(), addr2.String())
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelPoolAddrRange)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(msg.StartAddr.String()).To(BeEquivalentTo(addr1.String()))
	Expect(msg.EndAddr.String()).To(BeEquivalentTo(addr2.String()))
	Expect(msg.VrfID).To(BeEquivalentTo(5))

	err = natHandler.DelNat64AddressPool(0, "invalid address", "")
	Expect(err).Should(HaveOccurred())
	err = natHandler.DelNat64AddressPool(0, "invalid address", addr2.String())
	Expect(err).Should(HaveOccurred())
	err = natHandler.DelNat64AddressPool(0, addr1.String(), "invalid address")
	Expect(err).Should(HaveOccurred())
}

func TestAddNat64StaticBIB(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	inAddr := net.ParseIP("2000::3").To16()
	outAddr := net.ParseIP("172.16.2.3").To4()

	// TCP
	ctx.MockVpp.MockReply(&natba.Nat64AddDelStaticBibReply{})
	err := natHandler.AddNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              10,
		InsideIpv6Address:  inAddr.String(),
		InsidePort:         8000,
		OutsideIpv4Address: outAddr.String(),
		OutsidePort:        80,
		Protocol:           nat64.Nat64StaticBIB_TCP,
	})
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelStaticBib)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(msg.IAddr.String()).To(BeEquivalentTo(inAddr.String()))
	Expect(msg.IPort).To(BeEquivalentTo(8000))
	Expect(msg.OAddr.String()).To(BeEquivalentTo(outAddr.String()))
	Expect(msg.OPort).To(BeEquivalentTo(80))
	Expect(msg.VrfID).To(BeEquivalentTo(10))
	Expect(msg.Proto).To(BeEquivalentTo(6))

	// UDP
	ctx.MockVpp.MockReply(&natba.Nat64AddDelStaticBibReply{})
	err = natHandler.AddNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              0,
		InsideIpv6Address:  inAddr.String(),
		InsidePort:         9000,
		OutsideIpv4Address: outAddr.String(),
		OutsidePort:        90,
		Protocol:           nat64.Nat64StaticBIB_UDP,
	})
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelStaticBib)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeTrue())
	Expect(msg.IAddr.String()).To(BeEquivalentTo(inAddr.String()))
	Expect(msg.IPort).To(BeEquivalentTo(9000))
	Expect(msg.OAddr.String()).To(BeEquivalentTo(outAddr.String()))
	Expect(msg.OPort).To(BeEquivalentTo(90))
	Expect(msg.VrfID).To(BeEquivalentTo(0))
	Expect(msg.Proto).To(BeEquivalentTo(17))

	// Invalid data
	err = natHandler.AddNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              0,
		InsideIpv6Address:  "192.168.1.1", // expecting IPv6
		InsidePort:         9000,
		OutsideIpv4Address: outAddr.String(),
		OutsidePort:        90,
		Protocol:           nat64.Nat64StaticBIB_UDP,
	})
	Expect(err).Should(HaveOccurred())
	err = natHandler.AddNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              0,
		InsideIpv6Address:  inAddr.String(),
		InsidePort:         9000,
		OutsideIpv4Address: "invalid IP address",
		OutsidePort:        90,
		Protocol:           nat64.Nat64StaticBIB_UDP,
	})
	Expect(err).Should(HaveOccurred())
}

func TestDelNat64StaticBIB(t *testing.T) {
	ctx, natHandler, _ := natTestSetup(t)
	defer ctx.TeardownTestCtx()

	inAddr := net.ParseIP("2000::3").To16()
	outAddr := net.ParseIP("172.16.2.3").To4()

	// TCP
	ctx.MockVpp.MockReply(&natba.Nat64AddDelStaticBibReply{})
	err := natHandler.DelNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              10,
		InsideIpv6Address:  inAddr.String(),
		InsidePort:         8000,
		OutsideIpv4Address: outAddr.String(),
		OutsidePort:        80,
		Protocol:           nat64.Nat64StaticBIB_TCP,
	})
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok := ctx.MockChannel.Msg.(*natba.Nat64AddDelStaticBib)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(msg.IAddr.String()).To(BeEquivalentTo(inAddr.String()))
	Expect(msg.IPort).To(BeEquivalentTo(8000))
	Expect(msg.OAddr.String()).To(BeEquivalentTo(outAddr.String()))
	Expect(msg.OPort).To(BeEquivalentTo(80))
	Expect(msg.VrfID).To(BeEquivalentTo(10))
	Expect(msg.Proto).To(BeEquivalentTo(6))

	// UDP
	ctx.MockVpp.MockReply(&natba.Nat64AddDelStaticBibReply{})
	err = natHandler.DelNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              0,
		InsideIpv6Address:  inAddr.String(),
		InsidePort:         9000,
		OutsideIpv4Address: outAddr.String(),
		OutsidePort:        90,
		Protocol:           nat64.Nat64StaticBIB_UDP,
	})
	Expect(err).ShouldNot(HaveOccurred())

	msg, ok = ctx.MockChannel.Msg.(*natba.Nat64AddDelStaticBib)
	Expect(ok).To(BeTrue())
	Expect(msg.IsAdd).To(BeFalse())
	Expect(msg.IAddr.String()).To(BeEquivalentTo(inAddr.String()))
	Expect(msg.IPort).To(BeEquivalentTo(9000))
	Expect(msg.OAddr.String()).To(BeEquivalentTo(outAddr.String()))
	Expect(msg.OPort).To(BeEquivalentTo(90))
	Expect(msg.VrfID).To(BeEquivalentTo(0))
	Expect(msg.Proto).To(BeEquivalentTo(17))

	// Invalid data
	err = natHandler.DelNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              0,
		InsideIpv6Address:  "192.168.1.1", // expecting IPv6
		InsidePort:         9000,
		OutsideIpv4Address: outAddr.String(),
		OutsidePort:        90,
		Protocol:           nat64.Nat64StaticBIB_UDP,
	})
	Expect(err).Should(HaveOccurred())
	err = natHandler.DelNat64StaticBIB(&nat64.Nat64StaticBIB{
		VrfId:              0,
		InsideIpv6Address:  inAddr.String(),
		InsidePort:         9000,
		OutsideIpv4Address: "invalid IP address",
		OutsidePort:        90,
		Protocol:           nat64.Nat64StaticBIB_UDP,
	})
	Expect(err).Should(HaveOccurred())
}

func ip6PrefixToIPNet(prefix ip_types.IP6Prefix) string {
	ipNet := &net.IPNet{}
	ipNet.IP = make(net.IP, net.IPv6len)
	copy(ipNet.IP[:], prefix.Address[:])
	ipNet.Mask = net.CIDRMask(int(prefix.Len), 8*net.IPv6len)
	return ipNet.String()
}

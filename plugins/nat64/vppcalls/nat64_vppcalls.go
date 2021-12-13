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

package vppcalls

import (
	govppapi "git.fd.io/govpp.git/api"
	"go.ligato.io/cn-infra/v2/logging"

	"go.ligato.io/vpp-agent/v3/plugins/vpp"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"
	"go.pantheon.tech/stonework/proto/nat64"
)

// Nat64VppAPI provides methods for managing VPP NAT-64 configuration.
type Nat64VppAPI interface {
	Nat64VppRead

	// AddNat64IPv6Prefix adds IPv6 prefix for NAT64 (used to embed IPv4 address).
	AddNat64IPv6Prefix(vrf uint32, prefix string) error
	// DelNat64IPv6Prefix removes existing IPv6 prefix previously configured for NAT64.
	DelNat64IPv6Prefix(vrf uint32, prefix string) error
	// EnableNat64Interface enables NAT64 for provided interface.
	EnableNat64Interface(iface string, natIfaceType nat64.Nat64Interface_Type) error
	// DisableNat64Interface disables NAT64 for provided interface.
	DisableNat64Interface(iface string, natIfaceType nat64.Nat64Interface_Type) error
	// AddNat64AddressPool adds new IPV4 address pool into the NAT64 pools.
	AddNat64AddressPool(vrf uint32, firstIP, lastIP string) error
	// DelNat64AddressPool removes existing IPv4 address pool from the NAT64 pools.
	DelNat64AddressPool(vrf uint32, firstIP, lastIP string) error
	// AddNat64StaticBIB creates new NAT64 static binding.
	AddNat64StaticBIB(bib *nat64.Nat64StaticBIB) error
	// DelNat64StaticBIB removes existing NAT64 static binding.
	DelNat64StaticBIB(bib *nat64.Nat64StaticBIB) error
}

// Nat64VppRead provides read methods for VPP NAT-64 configuration.
type Nat64VppRead interface {
	// Nat64IPv6PrefixDump dumps all IPv6 prefixes configured for NAT64.
	Nat64IPv6PrefixDump() ([]*nat64.Nat64IPv6Prefix, error)
	// Nat64InterfacesDump dumps NAT64 config of all NAT64-enabled interfaces.
	Nat64InterfacesDump() ([]*nat64.Nat64Interface, error)
	// Nat64AddressPoolsDump dumps all configured NAT64 address pools.
	// Note that VPP returns configured addresses one-by-one, loosing information about address pools
	// configured with multiple addresses through IP ranges. Provide expected configuration to group
	// addresses from the same range.
	Nat64AddressPoolsDump(correlateWith []*nat64.Nat64AddressPool) ([]*nat64.Nat64AddressPool, error)
	// Nat64StaticBIBsDump dumps NAT64 static bindings.
	Nat64StaticBIBsDump() ([]*nat64.Nat64StaticBIB, error)
}

var handler = vpp.RegisterHandler(vpp.HandlerDesc{
	Name:       "nat64",
	HandlerAPI: (*Nat64VppAPI)(nil),
})

func AddNat64HandlerVersion(version vpp.Version, msgs []govppapi.Message,
	h func(ch govppapi.Channel, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) Nat64VppAPI,
) {
	handler.AddVersion(vpp.HandlerVersion{
		Version: version,
		Check: func(c vpp.Client) error {
			ch, err := c.NewAPIChannel()
			if err != nil {
				return err
			}
			return ch.CheckCompatiblity(msgs...)
		},
		NewHandler: func(c vpp.Client, a ...interface{}) vpp.HandlerAPI {
			ch, err := c.NewAPIChannel()
			if err != nil {
				return err
			}
			return h(ch, a[0].(ifaceidx.IfaceMetadataIndex), a[1].(logging.Logger))
		},
	})
}

func CompatibleNat64VppHandler(c vpp.Client, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) Nat64VppAPI {
	if v := handler.FindCompatibleVersion(c); v != nil {
		return v.NewHandler(c, ifIdx, log).(Nat64VppAPI)
	}
	return nil
}

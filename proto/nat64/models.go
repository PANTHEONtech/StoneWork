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

package nat64

import (
	"go.ligato.io/vpp-agent/v3/pkg/models"
)

const ModuleName = "vpp.nat"

var (
	ModelNat64IPv6Prefix = models.Register(&Nat64IPv6Prefix{}, models.Spec{
		Module:  ModuleName,
		Type:    "nat64-prefix",
		Version: "v1",
	}, models.WithNameTemplate("vrf/{{.VrfId}}"))

	ModelNat64Interface = models.Register(&Nat64Interface{}, models.Spec{
		Module:  ModuleName,
		Type:    "nat64-interface",
		Version: "v1",
	}, models.WithNameTemplate("{{.Name}}"))

	ModelNat64AddressPool = models.Register(&Nat64AddressPool{}, models.Spec{
		Module:  ModuleName,
		Type:    "nat64-pool",
		Version: "v1",
	}, models.WithNameTemplate(
		"vrf/{{.VrfId}}"+
			"/address/{{.FirstIp}}"+
			"{{if and .LastIp (ne .FirstIp .LastIp)}}-{{.LastIp}}{{end}}",
	))

	ModelNat64StaticBIB = models.Register(&Nat64StaticBIB{}, models.Spec{
		Module:  ModuleName,
		Type:    "nat64-static-bib",
		Version: "v1",
	}, models.WithNameTemplate(
		"vrf/{{.VrfId}}/proto/{{.Protocol}}"+
			"/inaddr/{{.InsideIpv6Address}}/inport/{{.InsidePort}}"+
			"/outaddr/{{.OutsideIpv4Address}}/outport/{{.OutsidePort}}"))
)

// Nat64IPv6PrefixKey returns the key used in NB DB to store the configuration of a NAT64 IPv6 prefix
// inside a given VRF.
func Nat64IPv6PrefixKey(vrf uint32) string {
	return models.Key(&Nat64IPv6Prefix{
		VrfId: vrf,
	})
}

// Nat64InterfaceKey returns the key used in NB DB to store the configuration of the
// given NAT64 interface.
func Nat64InterfaceKey(name string) string {
	return models.Key(&Nat64Interface{
		Name: name,
	})
}

// Nat64AddressPoolKey returns the key used in NB DB to store the configuration of the
// given NAT64 address pool.
func Nat64AddressPoolKey(vrf uint32, firstIP, lastIP string) string {
	return models.Key(&Nat64AddressPool{
		VrfId:   vrf,
		FirstIp: firstIP,
		LastIp:  lastIP,
	})
}

// Nat64StaticBIBKey returns the key used in NB DB to store the configuration of the
// given NAT64 static BIB.
func Nat64StaticBIBKey(bib *Nat64StaticBIB) string {
	return models.Key(bib)
}

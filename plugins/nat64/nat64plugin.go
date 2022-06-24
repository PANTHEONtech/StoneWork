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

//go:generate descriptor-adapter --descriptor-name NAT64IPv6Prefix --value-type *nat64.Nat64IPv6Prefix --import "go.pantheon.tech/stonework/proto/nat64" --output-dir "descriptor"
//go:generate descriptor-adapter --descriptor-name NAT64Interface --value-type *nat64.Nat64Interface --import "go.pantheon.tech/stonework/proto/nat64" --output-dir "descriptor"
//go:generate descriptor-adapter --descriptor-name NAT64AddressPool --value-type *nat64.Nat64AddressPool --import "go.pantheon.tech/stonework/proto/nat64" --output-dir "descriptor"
//go:generate descriptor-adapter --descriptor-name NAT64StaticBIB --value-type *nat64.Nat64StaticBIB --import "go.pantheon.tech/stonework/proto/nat64" --output-dir "descriptor"

package nat64plugin

import (
	"github.com/pkg/errors"

	"go.ligato.io/cn-infra/v2/infra"

	"go.ligato.io/vpp-agent/v3/plugins/govppmux"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"

	"go.pantheon.tech/stonework/plugins/nat64/descriptor"
	"go.pantheon.tech/stonework/plugins/nat64/vppcalls"

	_ "go.pantheon.tech/stonework/plugins/nat64/vppcalls/vpp2106"
)

// NAT64Plugin configures VPP NAT.
type NAT64Plugin struct {
	Deps

	// handlers
	nat64Handler vppcalls.Nat64VppAPI
}

// Deps lists dependencies of the NAT plugin.
type Deps struct {
	infra.PluginDeps
	KVScheduler kvs.KVScheduler
	VPP         govppmux.API
	IfPlugin    ifplugin.API
}

// Init registers NAT64-related descriptors.
func (p *NAT64Plugin) Init() (err error) {
	if !p.VPP.IsPluginLoaded("nat") {
		p.Log.Warnf("VPP plugin NAT64 was disabled by VPP")
		return nil
	}

	// init handlers
	p.nat64Handler = vppcalls.CompatibleNat64VppHandler(p.VPP, p.IfPlugin.GetInterfaceIndex(), p.Log)
	if p.nat64Handler == nil {
		return errors.New("nat64Handler is not available")
	}

	// init and register descriptors
	nat64IPv6PrefixDescriptor := descriptor.NewNAT64IPv6PrefixDescriptor(p.nat64Handler, p.Log)
	nat64InterfaceDescriptor := descriptor.NewNAT64InterfaceDescriptor(p.nat64Handler, p.Log)
	nat64AddressPoolDescriptor := descriptor.NewNAT64AddressPoolDescriptor(p.nat64Handler, p.Log)
	nat64StaticBIBDescriptor := descriptor.NewNAT64StaticBIBDescriptor(p.nat64Handler, p.Log)

	err = p.KVScheduler.RegisterKVDescriptor(
		nat64IPv6PrefixDescriptor,
		nat64InterfaceDescriptor,
		nat64AddressPoolDescriptor,
		nat64StaticBIBDescriptor,
	)
	if err != nil {
		return err
	}

	return nil
}

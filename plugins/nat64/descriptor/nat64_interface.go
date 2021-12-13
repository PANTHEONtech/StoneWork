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

package descriptor

import (
	"go.ligato.io/cn-infra/v2/logging"

	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_ifdescriptor "go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/descriptor"
	interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"

	"go.pantheon.tech/stonework/plugins/nat64/descriptor/adapter"
	"go.pantheon.tech/stonework/plugins/nat64/vppcalls"
	"go.pantheon.tech/stonework/proto/nat64"
)

const (
	// NAT64InterfaceDescriptorName is the name of the descriptor for enabling/disabling NAT64 for VPP interfaces.
	NAT64InterfaceDescriptorName = "vpp-nat64-interface"

	// dependency labels
	natInterfaceDep = "interface-exists"
)

// NAT64InterfaceDescriptor teaches KVScheduler how to enable NAT64 for VPP interfaces.
type NAT64InterfaceDescriptor struct {
	log        logging.Logger
	natHandler vppcalls.Nat64VppAPI
}

// NewNAT64InterfaceDescriptor creates a new instance of the NAT64Interface descriptor.
func NewNAT64InterfaceDescriptor(natHandler vppcalls.Nat64VppAPI, log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &NAT64InterfaceDescriptor{
		natHandler: natHandler,
		log:        log.NewLogger("nat64-iface-descriptor"),
	}

	typedDescr := &adapter.NAT64InterfaceDescriptor{
		Name:                 NAT64InterfaceDescriptorName,
		NBKeyPrefix:          nat64.ModelNat64Interface.KeyPrefix(),
		ValueTypeName:        nat64.ModelNat64Interface.ProtoName(),
		KeySelector:          nat64.ModelNat64Interface.IsKeyValid,
		KeyLabel:             nat64.ModelNat64Interface.StripKeyPrefix,
		Create:               ctx.Create,
		Delete:               ctx.Delete,
		Retrieve:             ctx.Retrieve,
		Dependencies:         ctx.Dependencies,
		RetrieveDependencies: []string{vpp_ifdescriptor.InterfaceDescriptorName},
	}
	return adapter.NewNAT64InterfaceDescriptor(typedDescr)
}

// Create enables NAT64 on an interface.
func (d *NAT64InterfaceDescriptor) Create(key string, natIface *nat64.Nat64Interface) (metadata interface{}, err error) {
	err = d.natHandler.EnableNat64Interface(natIface.Name, natIface.Type)
	return
}

// Delete disables NAT64 on an interface.
func (d *NAT64InterfaceDescriptor) Delete(key string, natIface *nat64.Nat64Interface, metadata interface{}) (err error) {
	return d.natHandler.DisableNat64Interface(natIface.Name, natIface.Type)
}

// Retrieve returns the current NAT64 interface configuration.
func (d *NAT64InterfaceDescriptor) Retrieve(correlate []adapter.NAT64InterfaceKVWithMetadata) (
	retrieved []adapter.NAT64InterfaceKVWithMetadata, err error) {
	ifaces, err := d.natHandler.Nat64InterfacesDump()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		retrieved = append(retrieved, adapter.NAT64InterfaceKVWithMetadata{
			Key:    nat64.Nat64InterfaceKey(iface.Name),
			Value:  iface,
			Origin: kvs.FromNB,
		})
	}
	return
}

// Dependencies lists the interface as the only dependency.
func (d *NAT64InterfaceDescriptor) Dependencies(key string, natIface *nat64.Nat64Interface) []kvs.Dependency {
	return []kvs.Dependency{
		{
			Label: natInterfaceDep,
			Key:   interfaces.InterfaceKey(natIface.Name),
		},
	}
}

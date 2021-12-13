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
	"errors"
	"net"

	"go.ligato.io/cn-infra/v2/logging"

	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	l3 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"

	"go.pantheon.tech/stonework/plugins/nat64/descriptor/adapter"
	"go.pantheon.tech/stonework/plugins/nat64/vppcalls"
	"go.pantheon.tech/stonework/proto/nat64"
)

const (
	// NAT64StaticBIBDescriptorName is the name of the descriptor for configuring NAT64 Static BIBs.
	NAT64StaticBIBDescriptorName = "vpp-nat64-static-BIB"
)

// A list of non-retriable errors:
var (
	// errUndefinedPort is returned when (TCP or UDP) port is not defined.
	errUndefinedPort = errors.New("undefined port number")
)

// NAT64StaticBIBDescriptor teaches KVScheduler how to configure NAT64 static bindings.
type NAT64StaticBIBDescriptor struct {
	log        logging.Logger
	natHandler vppcalls.Nat64VppAPI
}

// NewNAT64StaticBIBDescriptor creates a new instance of the NAT64StaticBIB descriptor.
func NewNAT64StaticBIBDescriptor(natHandler vppcalls.Nat64VppAPI, log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &NAT64StaticBIBDescriptor{
		natHandler: natHandler,
		log:        log.NewLogger("nat64-static-BIB-descriptor"),
	}

	typedDescr := &adapter.NAT64StaticBIBDescriptor{
		Name:          NAT64StaticBIBDescriptorName,
		NBKeyPrefix:   nat64.ModelNat64StaticBIB.KeyPrefix(),
		ValueTypeName: nat64.ModelNat64StaticBIB.ProtoName(),
		KeySelector:   nat64.ModelNat64StaticBIB.IsKeyValid,
		KeyLabel:      nat64.ModelNat64StaticBIB.StripKeyPrefix,
		Validate:      ctx.Validate,
		Create:        ctx.Create,
		Delete:        ctx.Delete,
		Retrieve:      ctx.Retrieve,
		Dependencies:  ctx.Dependencies,
	}
	return adapter.NewNAT64StaticBIBDescriptor(typedDescr)
}

// Validate validates NAT64 BIB configuration.
func (d *NAT64StaticBIBDescriptor) Validate(key string, bib *nat64.Nat64StaticBIB) error {
	if bib.OutsidePort == 0 {
		return kvs.NewInvalidValueError(errUndefinedPort, "outside_port")
	}
	if bib.InsidePort == 0 {
		return kvs.NewInvalidValueError(errUndefinedPort, "inside_port")
	}
	inAddr := net.ParseIP(bib.InsideIpv6Address)
	if inAddr == nil || inAddr.To4() != nil {
		return kvs.NewInvalidValueError(errInvalidIPAddress, "inside_ipv6_address")
	}
	outAddr := net.ParseIP(bib.OutsideIpv4Address)
	if outAddr == nil || outAddr.To4() == nil {
		return kvs.NewInvalidValueError(errInvalidIPAddress, "outside_ipv4_address")
	}
	return nil
}

// Create adds new NAT44 static binding.
func (d *NAT64StaticBIBDescriptor) Create(key string, bib *nat64.Nat64StaticBIB) (metadata interface{}, err error) {
	err = d.natHandler.AddNat64StaticBIB(bib)
	return
}

// Delete removes existing NAT44 static binding.
func (d *NAT64StaticBIBDescriptor) Delete(key string, bib *nat64.Nat64StaticBIB, metadata interface{}) (err error) {
	return d.natHandler.DelNat64StaticBIB(bib)
}

// Retrieve returns the currently configured static NAT64 bindings.
func (d *NAT64StaticBIBDescriptor) Retrieve(correlate []adapter.NAT64StaticBIBKVWithMetadata) (
	retrieved []adapter.NAT64StaticBIBKVWithMetadata, err error) {
	bibs, err := d.natHandler.Nat64StaticBIBsDump()
	if err != nil {
		return nil, err
	}
	for _, bib := range bibs {
		retrieved = append(retrieved, adapter.NAT64StaticBIBKVWithMetadata{
			Key:    nat64.Nat64StaticBIBKey(bib),
			Value:  bib,
			Origin: kvs.FromNB,
		})
	}
	return
}

// Dependencies lists the VRF (for both IP versions) as the only dependency.
func (d *NAT64StaticBIBDescriptor) Dependencies(key string, bib *nat64.Nat64StaticBIB) (deps []kvs.Dependency) {
	if bib.VrfId != 0 && bib.VrfId != ^uint32(0) {
		deps = append(deps,
			kvs.Dependency{
				Label: mappingVrfDep,
				Key:   l3.VrfTableKey(bib.VrfId, l3.VrfTable_IPV4),
			},
			kvs.Dependency{
				Label: mappingVrfDep,
				Key:   l3.VrfTableKey(bib.VrfId, l3.VrfTable_IPV6),
			})
	}
	return
}

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
	// NAT64IPv6PrefixDescriptorName is the name of the descriptor for configuring NAT64 IPv6 prefix.
	NAT64IPv6PrefixDescriptorName = "vpp-nat64-ipv6-prefix"

	// dependency labels
	mappingVrfDep = "vrf-table-exists"
)

// A list of non-retriable errors:
var (
	// errInvalidIPv6Prefix is returned when IPv6 prefix is not valid for NAT64.
	errInvalidIPv6Prefix = errors.New("invalid IPv6 prefix")
)

// NAT64IPv6PrefixDescriptor teaches KVScheduler how to configure NAT64 IPv6 Prefix.
type NAT64IPv6PrefixDescriptor struct {
	log        logging.Logger
	natHandler vppcalls.Nat64VppAPI
}

// NewNAT64IPv6PrefixDescriptor creates a new instance of the NAT64IPv6Prefix descriptor.
func NewNAT64IPv6PrefixDescriptor(natHandler vppcalls.Nat64VppAPI, log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &NAT64IPv6PrefixDescriptor{
		natHandler: natHandler,
		log:        log.NewLogger("nat64-ipv6-prefix-descriptor"),
	}

	typedDescr := &adapter.NAT64IPv6PrefixDescriptor{
		Name:          NAT64IPv6PrefixDescriptorName,
		NBKeyPrefix:   nat64.ModelNat64IPv6Prefix.KeyPrefix(),
		ValueTypeName: nat64.ModelNat64IPv6Prefix.ProtoName(),
		KeySelector:   nat64.ModelNat64IPv6Prefix.IsKeyValid,
		KeyLabel:      nat64.ModelNat64IPv6Prefix.StripKeyPrefix,
		Create:        ctx.Create,
		Delete:        ctx.Delete,
		Retrieve:      ctx.Retrieve,
		Dependencies:  ctx.Dependencies,
	}
	return adapter.NewNAT64IPv6PrefixDescriptor(typedDescr)
}

// Validate validates NAT64 IPv6 prefix.
func (d *NAT64IPv6PrefixDescriptor) Validate(key string, prefix *nat64.Nat64IPv6Prefix) error {
	_, prefixIPNet, err := net.ParseCIDR(prefix.Prefix)
	if err != nil || prefixIPNet.IP.To4() != nil {
		return kvs.NewInvalidValueError(errInvalidIPv6Prefix, "prefix")
	}
	var validPrefixLen = []int{32, 40, 48, 56, 64, 96}
	var lenIsValid bool
	prefixLen, _ := prefixIPNet.Mask.Size()
	for _, valid := range validPrefixLen {
		if prefixLen == valid {
			lenIsValid = true
			break
		}
	}
	if !lenIsValid {
		return kvs.NewInvalidValueError(errInvalidIPv6Prefix, "prefix")
	}
	return nil
}

// Create configures NAT64 IPv6 prefix for a given VRF.
func (d *NAT64IPv6PrefixDescriptor) Create(key string, prefix *nat64.Nat64IPv6Prefix) (metadata interface{}, err error) {
	err = d.natHandler.AddNat64IPv6Prefix(prefix.VrfId, prefix.Prefix)
	return
}

// Delete removes NAT64 IPv6 prefix from a given VRF.
func (d *NAT64IPv6PrefixDescriptor) Delete(key string, prefix *nat64.Nat64IPv6Prefix, metadata interface{}) (err error) {
	return d.natHandler.DelNat64IPv6Prefix(prefix.VrfId, prefix.Prefix)
}

// Retrieve returns the currently configured NAT64 IPv6 prefixes.
func (d *NAT64IPv6PrefixDescriptor) Retrieve(correlate []adapter.NAT64IPv6PrefixKVWithMetadata) (
	retrieved []adapter.NAT64IPv6PrefixKVWithMetadata, err error) {
	prefixes, err := d.natHandler.Nat64IPv6PrefixDump()
	if err != nil {
		return nil, err
	}
	for _, prefix := range prefixes {
		retrieved = append(retrieved, adapter.NAT64IPv6PrefixKVWithMetadata{
			Key:    nat64.Nat64IPv6PrefixKey(prefix.VrfId),
			Value:  prefix,
			Origin: kvs.FromNB,
		})
	}
	return
}

// Dependencies lists the VRF (for both IP versions) as the only dependency.
func (d *NAT64IPv6PrefixDescriptor) Dependencies(key string, prefix *nat64.Nat64IPv6Prefix) (deps []kvs.Dependency) {
	if prefix.VrfId != 0 && prefix.VrfId != ^uint32(0) {
		deps = append(deps,
			kvs.Dependency{
				Label: mappingVrfDep,
				Key:   l3.VrfTableKey(prefix.VrfId, l3.VrfTable_IPV4),
			},
			kvs.Dependency{
				Label: mappingVrfDep,
				Key:   l3.VrfTableKey(prefix.VrfId, l3.VrfTable_IPV6),
			})
	}
	return
}

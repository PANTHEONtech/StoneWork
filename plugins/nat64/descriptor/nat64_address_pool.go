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
	"bytes"
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
	// NAT64AddressPoolDescriptorName is the name of the descriptor for NAT64 IP address pools.
	NAT64AddressPoolDescriptorName = "vpp-nat64-address-pool"

	// dependency labels
	addressVrfDep = "vrf-table-exists"
)

// A list of non-retriable errors:
var (
	// errInvalidIPAddress is returned when IP address from NAT64 address pool cannot be parsed.
	errInvalidIPAddress = errors.New("invalid IP address")
	// errInvalidLastPoolAddress is returned when last IP of the pool is not higher than first IP of the pool, or empty.
	errInvalidLastPoolAddress = errors.New("last IP should be higher than first IP, or empty")
)

// NAT64AddressPoolDescriptor teaches KVScheduler how to add/remove VPP NAT64 IP address pools.
type NAT64AddressPoolDescriptor struct {
	log        logging.Logger
	natHandler vppcalls.Nat64VppAPI
}

// NewNAT64AddressPoolDescriptor creates a new instance of the NAT64AddressPoolDescriptor.
func NewNAT64AddressPoolDescriptor(natHandler vppcalls.Nat64VppAPI, log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &NAT64AddressPoolDescriptor{
		natHandler: natHandler,
		log:        log.NewLogger("nat64-address-pool-descriptor"),
	}
	typedDescr := &adapter.NAT64AddressPoolDescriptor{
		Name:            NAT64AddressPoolDescriptorName,
		NBKeyPrefix:     nat64.ModelNat64AddressPool.KeyPrefix(),
		ValueTypeName:   nat64.ModelNat64AddressPool.ProtoName(),
		KeySelector:     nat64.ModelNat64AddressPool.IsKeyValid,
		KeyLabel:        nat64.ModelNat64AddressPool.StripKeyPrefix,
		ValueComparator: ctx.EquivalentAddressPools,
		Validate:        ctx.Validate,
		Create:          ctx.Create,
		Delete:          ctx.Delete,
		Retrieve:        ctx.Retrieve,
		Dependencies:    ctx.Dependencies,
	}
	return adapter.NewNAT64AddressPoolDescriptor(typedDescr)
}

// EquivalentAddressPools compares two address pools for equivalency.
func (d *NAT64AddressPoolDescriptor) EquivalentAddressPools(key string, oldPool, newPool *nat64.Nat64AddressPool) bool {
	if oldPool.VrfId != newPool.VrfId || oldPool.FirstIp != newPool.FirstIp {
		return false
	}
	if d.getLastIP(oldPool) != d.getLastIP(newPool) {
		return false
	}
	return true
}

// Validate validates configuration for NAT64 IP addresses pool.
func (d *NAT64AddressPoolDescriptor) Validate(key string, pool *nat64.Nat64AddressPool) error {
	firstIp := net.ParseIP(pool.FirstIp)
	if firstIp == nil {
		return kvs.NewInvalidValueError(errInvalidIPAddress, "first_ip")
	}
	if pool.LastIp != "" {
		lastIp := net.ParseIP(pool.LastIp)
		if lastIp == nil {
			return kvs.NewInvalidValueError(errInvalidIPAddress, "last_ip")
		}
		if bytes.Compare(firstIp, lastIp) > 0 {
			// last IP should be empty or higher than first IP
			return kvs.NewInvalidValueError(errInvalidLastPoolAddress, "last_ip")
		}
	}
	return nil
}

// Create adds IP address pool into VPP NAT64 address pools.
func (d *NAT64AddressPoolDescriptor) Create(key string, pool *nat64.Nat64AddressPool) (metadata interface{}, err error) {
	return nil,
		d.natHandler.AddNat64AddressPool(pool.VrfId, pool.FirstIp, pool.LastIp)
}

// Delete removes IP address pool from VPP NAT64 address pools.
func (d *NAT64AddressPoolDescriptor) Delete(key string, pool *nat64.Nat64AddressPool, metadata interface{}) error {
	return d.natHandler.DelNat64AddressPool(pool.VrfId, pool.FirstIp, pool.LastIp)
}

// Retrieve returns NAT64 IP address pools configured on VPP.
func (d *NAT64AddressPoolDescriptor) Retrieve(correlate []adapter.NAT64AddressPoolKVWithMetadata) (
	retrieved []adapter.NAT64AddressPoolKVWithMetadata, err error) {
	var expected []*nat64.Nat64AddressPool
	for _, pool := range correlate {
		expected = append(expected, pool.Value)
	}
	natPools, err := d.natHandler.Nat64AddressPoolsDump(expected)
	if err != nil {
		return nil, err
	}
	for _, pool := range natPools {
		retrieved = append(retrieved, adapter.NAT64AddressPoolKVWithMetadata{
			Key:    nat64.Nat64AddressPoolKey(pool.VrfId, pool.FirstIp, pool.LastIp),
			Value:  pool,
			Origin: kvs.FromNB,
		})
	}
	return
}

// Dependencies lists non-zero and non-all-ones (IPv4) VRF as the only dependency.
func (d *NAT64AddressPoolDescriptor) Dependencies(key string, pool *nat64.Nat64AddressPool) []kvs.Dependency {
	if pool.VrfId == 0 || pool.VrfId == ^uint32(0) {
		return nil
	}
	return []kvs.Dependency{
		{
			Label: addressVrfDep,
			Key:   l3.VrfTableKey(pool.VrfId, l3.VrfTable_IPV4),
		},
	}
}

func (d *NAT64AddressPoolDescriptor) getLastIP(pool *nat64.Nat64AddressPool) string {
	if pool.LastIp != "" {
		return pool.LastIp
	}
	return pool.FirstIp
}

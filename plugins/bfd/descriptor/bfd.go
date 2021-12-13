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
	"fmt"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
	"go.ligato.io/cn-infra/v2/logging"

	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vppif "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"

	"go.pantheon.tech/stonework/plugins/bfd/descriptor/adapter"
	"go.pantheon.tech/stonework/plugins/bfd/vppcalls"
	"go.pantheon.tech/stonework/proto/bfd"
)

const (
	// name of the descriptor
	bfdDescriptorName = "bfd-descriptor"

	// name of the local interface dependency
	bfdLocalInterfaceDep = "bfd-local-interface-dep"
)

// validation errors
var (
	// ErrBfdIPAddressMissing is returned if the BFD has no local or peer IP address defined
	ErrBfdIPAddressMissing = errors.New("BFD: local or peer IP address is missing")

	// ErrBfdIPAddressInvalid is returned if local or peer IP address is malformed
	ErrBfdIPAddressInvalid = errors.New("BFD: local or peer IP addresses is invalid")

	// ErrBfdDetectMultiplierInvalid is returned if the detect multiplier is a null value
	ErrBfdDetectMultiplierInvalid = errors.New("BFD: detect multiplier must be non-zero value")
)

// BfdDescriptor defines BFD session, model definition and validation
type BfdDescriptor struct {
	log logging.Logger

	// handler manages VPP calls
	handler vppcalls.BfdVppAPI

	// index map cache stores configuration ID for sessions (so they do not need to be
	// provided externally)
	indexCache map[uint32]*bfd.BFD
	mx         sync.Mutex
}

// NewBfdDescriptor initializes BFD descriptor
func NewBfdDescriptor(handler vppcalls.BfdVppAPI, log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &BfdDescriptor{
		handler:    handler,
		indexCache: make(map[uint32]*bfd.BFD),
		log:        log.NewLogger(bfdDescriptorName),
	}
	typed := &adapter.BfdDescriptor{
		Name:          bfdDescriptorName,
		KeySelector:   bfd.ModelBFD.IsKeyValid,
		ValueTypeName: bfd.ModelBFD.ProtoName(),
		KeyLabel:      bfd.ModelBFD.StripKeyPrefix,
		NBKeyPrefix:   bfd.ModelBFD.KeyPrefix(),
		Validate:      ctx.Validate,
		Create:        ctx.Create,
		Delete:        ctx.Delete,
		Retrieve:      ctx.Retrieve,
		Dependencies:  ctx.Dependencies,
	}
	return adapter.NewBfdDescriptor(typed)
}

// Validate BFD session IP addresses and detect multiplier value
func (d *BfdDescriptor) Validate(_ string, bfdEntry *bfd.BFD) error {
	// validate local IP address
	if bfdEntry.GetLocalIp() == "" {
		return kvs.NewInvalidValueError(ErrBfdIPAddressMissing, "local_ip")
	}
	if ip := net.ParseIP(bfdEntry.GetLocalIp()); ip == nil {
		return kvs.NewInvalidValueError(ErrBfdIPAddressInvalid, "local_ip")
	}
	// validate peer IP address
	if bfdEntry.GetPeerIp() == "" {
		return kvs.NewInvalidValueError(ErrBfdIPAddressMissing, "peer_ip")
	}
	if ip := net.ParseIP(bfdEntry.GetPeerIp()); ip == nil {
		return kvs.NewInvalidValueError(ErrBfdIPAddressInvalid, "peer_ip")
	}
	// detect multiplier
	if bfdEntry.GetDetectMultiplier() == 0 {
		return kvs.NewInvalidValueError(ErrBfdDetectMultiplierInvalid, "detect_multiplier")
	}

	return nil
}

// Create add a new BFD session
func (d *BfdDescriptor) Create(_ string, bfdEntry *bfd.BFD) (metadata interface{}, err error) {
	err = d.handler.AddBfd(d.addBfdConfID(bfdEntry), bfdEntry)
	return nil, err
}

// Delete existing BFD session
func (d *BfdDescriptor) Delete(_ string, bfdEntry *bfd.BFD, _ interface{}) error {
	d.delBfdConfID(bfdEntry)
	err := d.handler.DeleteBfd(bfdEntry)
	return err
}

func (d *BfdDescriptor) Retrieve(correlate []adapter.BfdKVWithMetadata) (dump []adapter.BfdKVWithMetadata, err error) {
	bfdList, err := d.handler.DumpBfd()
	if err != nil {
		return nil, fmt.Errorf("failed to dump BFD data: %v", err)
	}
	for _, bfdEntry := range bfdList {
		dump = append(dump, adapter.BfdKVWithMetadata{
			Key:    bfd.BFDKey(bfdEntry.Config.GetInterface(), bfdEntry.Config.GetPeerIp()),
			Value:  bfdEntry.Config,
			Origin: kvs.FromNB,
		})
	}
	return dump, nil
}

// Dependencies define interface where the BFD session is attached on
func (d *BfdDescriptor) Dependencies(_ string, bfdEntry *bfd.BFD) []kvs.Dependency {
	var dependencies []kvs.Dependency

	// the interface must exist
	if bfdEntry.GetInterface() != "" {
		dependencies = append(dependencies, kvs.Dependency{
			Label: bfdLocalInterfaceDep,
			Key:   vppif.InterfaceKey(bfdEntry.GetInterface()),
		})
	}

	return dependencies
}

func (d *BfdDescriptor) addBfdConfID(bfdEntry *bfd.BFD) uint32 {
	d.mx.Lock()
	defer d.mx.Unlock()

	var confID uint32 = 1
	for {
		if _, exists := d.indexCache[confID]; !exists {
			d.indexCache[confID] = bfdEntry
			return confID
		}
		confID++
	}
}

func (d *BfdDescriptor) delBfdConfID(bfdEntry *bfd.BFD) {
	d.mx.Lock()
	defer d.mx.Unlock()

	var confID uint32 = 1
	for {
		if val, ok := d.indexCache[confID]; ok && proto.Equal(val, bfdEntry) {
			delete(d.indexCache, confID)
			return
		}
		confID++
	}
}

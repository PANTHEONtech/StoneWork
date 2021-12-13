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
	"fmt"

	"github.com/go-errors/errors"
	"github.com/golang/protobuf/proto"

	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"

	"go.pantheon.tech/stonework/plugins/abx/abxidx"
	"go.pantheon.tech/stonework/plugins/abx/vppcalls"
	"go.pantheon.tech/stonework/proto/abx"
)

const (
	// ABXToInterfaceDescriptorName is name for descriptor
	ABXToInterfaceDescriptorName = "vpp-abx-to-interface"

	// dependency labels
	interfaceDep = "interface-exists"
)

// ABXToInterfaceDescriptor represents assignment of interface to ABX policy.
type ABXToInterfaceDescriptor struct {
	log        logging.Logger
	abxHandler vppcalls.ABXVppAPI
	abxIndex   abxidx.ABXMetadataIndex
	ifPlugin   ifplugin.API
}

// NewABXToInterfaceDescriptor returns new ABXInterface descriptor
func NewABXToInterfaceDescriptor(abxIndex abxidx.ABXMetadataIndex, abxHandler vppcalls.ABXVppAPI, ifPlugin ifplugin.API, log logging.PluginLogger) *api.KVDescriptor {
	ctx := &ABXToInterfaceDescriptor{
		log:        log,
		abxHandler: abxHandler,
		abxIndex:   abxIndex,
		ifPlugin:   ifPlugin,
	}

	return &api.KVDescriptor{
		Name:         ABXToInterfaceDescriptorName,
		KeySelector:  ctx.IsABXInterfaceKey,
		Create:       ctx.Create,
		Delete:       ctx.Delete,
		Dependencies: ctx.Dependencies,
	}
}

// IsABXInterfaceKey returns true if the key is identifying ABX policy interface (derived value)
func (d *ABXToInterfaceDescriptor) IsABXInterfaceKey(key string) bool {
	_, _, isABXToInterfaceKey := vpp_abx.ParseToInterfaceKey(key)
	return isABXToInterfaceKey
}

// Create binds interface to ABX.
func (d *ABXToInterfaceDescriptor) Create(key string, emptyVal proto.Message) (metadata api.Metadata, err error) {
	// validate and get all required values
	abxIdx, ifIdx, priority, err := d.process(key)
	if err != nil {
		d.log.Error(err)
		return nil, err
	}

	// attach interface to ABX policy
	return nil, d.abxHandler.AbxAttachInterface(abxIdx, ifIdx, priority)
}

// Delete unbinds interface from ABX.
func (d *ABXToInterfaceDescriptor) Delete(key string, emptyVal proto.Message, metadata api.Metadata) (err error) {
	// validate and get all required values
	abxIdx, ifIdx, priority, err := d.process(key)
	if err != nil {
		d.log.Error(err)
		return err
	}

	// detach interface to ABX policy
	return d.abxHandler.AbxDetachInterface(abxIdx, ifIdx, priority)
}

// Dependencies lists the interface as the only dependency for the binding.
func (d *ABXToInterfaceDescriptor) Dependencies(key string, emptyVal proto.Message) []api.Dependency {
	_, ifName, _ := vpp_abx.ParseToInterfaceKey(key)
	return []api.Dependency{
		{
			Label: interfaceDep,
			Key:   vpp_interfaces.InterfaceKey(ifName),
		},
	}
}

// returns a bunch of values needed to attach/detach interface to/from ABX
func (d *ABXToInterfaceDescriptor) process(key string) (abxIdx, ifIdx, priority uint32, err error) {
	// parse ABX and interface name
	abxIndex, ifName, isValid := vpp_abx.ParseToInterfaceKey(key)
	if !isValid {
		err = fmt.Errorf("ABX to interface key %s is not valid", key)
		return
	}
	// obtain ABX index
	abxData, exists := d.abxIndex.LookupByName(abxIndex)
	if !exists {
		err = errors.Errorf("failed to obtain metadata for ABX %s", abxIndex)
		return
	}

	// obtain interface index
	ifData, exists := d.ifPlugin.GetInterfaceIndex().LookupByName(ifName)
	if !exists {
		err = errors.Errorf("failed to obtain metadata for interface %s", ifName)
		return
	}

	// find other interface parameters from metadata
	for _, attachedIf := range abxData.Attached {
		if attachedIf.InputInterface == ifName {
			priority = attachedIf.Priority
		}
	}
	return abxData.Index, ifData.SwIfIndex, priority, nil
}

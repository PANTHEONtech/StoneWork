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
	"github.com/go-errors/errors"
	"go.ligato.io/cn-infra/v2/idxmap"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/aclplugin/aclidx"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/aclplugin/descriptor"
	ifdescriptor "go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/descriptor"
	acl "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"
	interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.pantheon.tech/stonework/plugins/abx/abxidx"
	"go.pantheon.tech/stonework/plugins/abx/descriptor/adapter"
	"go.pantheon.tech/stonework/plugins/abx/vppcalls"
	abx "go.pantheon.tech/stonework/proto/abx"
)

const (
	// ABXDescriptorName is descriptor name
	ABXDescriptorName = "vpp-abx"

	// dependency labels
	aclDep = "acl-exists"

	// dependency labels
	abxInterfaceDep = "interface-exists"
)

// A list of non-retriable errors:
var (
	// ErrABXWithoutACL is returned when ABX configuration does not contain associated access list.
	ErrABXWithoutACL = errors.New("ABX configuration defined without ACL")
)

// ABXDescriptor is descriptor for ABX
type ABXDescriptor struct {
	// dependencies
	log        logging.Logger
	abxHandler vppcalls.ABXVppAPI

	// runtime
	aclIndex aclidx.ACLMetadataIndex
}

// NewABXDescriptor is constructor for ABX descriptor and returns descriptor
// suitable for registration (via adapter) with the KVScheduler.
func NewABXDescriptor(abxHandler vppcalls.ABXVppAPI, aclIndex aclidx.ACLMetadataIndex,
	logger logging.PluginLogger) *api.KVDescriptor {
	ctx := &ABXDescriptor{
		log:        logger.NewLogger("abx-descriptor"),
		aclIndex:   aclIndex,
		abxHandler: abxHandler,
	}
	typedDescr := &adapter.ABXDescriptor{
		Name:          ABXDescriptorName,
		NBKeyPrefix:   abx.ModelABX.KeyPrefix(),
		ValueTypeName: abx.ModelABX.ProtoName(),
		KeySelector:   abx.ModelABX.IsKeyValid,
		KeyLabel:      abx.ModelABX.StripKeyPrefix,
		WithMetadata:  true,
		MetadataMapFactory: func() idxmap.NamedMappingRW {
			return abxidx.NewABXIndex(ctx.log, "vpp-abx-index")
		},
		ValueComparator:      ctx.EquivalentABXs,
		Validate:             ctx.Validate,
		Create:               ctx.Create,
		Delete:               ctx.Delete,
		Retrieve:             ctx.Retrieve,
		DerivedValues:        ctx.DerivedValues,
		Dependencies:         ctx.Dependencies,
		RetrieveDependencies: []string{ifdescriptor.InterfaceDescriptorName, descriptor.ACLDescriptorName},
	}
	return adapter.NewABXDescriptor(typedDescr)
}

// EquivalentABXs compares related ACL name, list of attached interfaces and forwarding paths to
// specify ABS equality.
func (d *ABXDescriptor) EquivalentABXs(key string, oldABX, newABX *abx.ABX) bool {
	// check index and associated ACL
	if oldABX.AclName != newABX.AclName {
		return false
	}
	if oldABX.DstMac != newABX.DstMac {
		return false
	}
	if oldABX.OutputInterface != newABX.OutputInterface {
		return false
	}
	// compare attached interfaces
	return equivalentABXAttachedInterfaces(oldABX.AttachedInterfaces, newABX.AttachedInterfaces)
}

// Validate validates VPP abx configuration.
func (d *ABXDescriptor) Validate(key string, abxData *abx.ABX) error {
	if abxData.AclName == "" {
		return api.NewInvalidValueError(ErrABXWithoutACL, "acl_name")
	}
	return nil
}

// Create validates ABX (mainly index), verifies ACL existence and configures ABX policy. Attached interfaces
// are put to metadata together with the ABX index to make it available for other ABX descriptors.
func (d *ABXDescriptor) Create(key string, abxData *abx.ABX) (*abxidx.ABXMetadata, error) {
	// get ACL index
	aclData, exists := d.aclIndex.LookupByName(abxData.AclName)
	if !exists {
		err := errors.Errorf("failed to obtain metadata for ACL %s", abxData.AclName)
		d.log.Error(err)
		return nil, err
	}

	// add new ABX policy
	if err := d.abxHandler.AddAbxPolicy(abxData.Index, aclData.Index, abxData.OutputInterface, abxData.DstMac); err != nil {
		d.log.Error(err)
		return nil, err
	}

	// fill the metadata
	metadata := &abxidx.ABXMetadata{
		Index:    abxData.Index,
		Attached: abxData.AttachedInterfaces,
	}

	return metadata, nil
}

// Delete removes ABX policy
func (d *ABXDescriptor) Delete(key string, abxData *abx.ABX, metadata *abxidx.ABXMetadata) error {
	// ACL ID is not required
	return d.abxHandler.DeleteAbxPolicy(metadata.Index)
}

// Retrieve returns ABX policies from the VPP.
func (d *ABXDescriptor) Retrieve(correlate []adapter.ABXKVWithMetadata) (abxs []adapter.ABXKVWithMetadata, err error) {
	// Retrieve VPP configuration.
	abxPolicies, err := d.abxHandler.DumpABXPolicy()
	if err != nil {
		return nil, errors.Errorf("failed to dump ABX policy: %v", err)
	}

	for _, abxPolicy := range abxPolicies {
		abxs = append(abxs, adapter.ABXKVWithMetadata{
			Key:   abx.Key(abxPolicy.ABX.Index),
			Value: abxPolicy.ABX,
			Metadata: &abxidx.ABXMetadata{
				Index:    abxPolicy.Meta.PolicyID,
				Attached: abxPolicy.ABX.AttachedInterfaces,
			},
			Origin: api.FromNB,
		})
	}

	return abxs, nil
}

// DerivedValues returns list of derived values for ABX.
func (d *ABXDescriptor) DerivedValues(key string, value *abx.ABX) (derived []api.KeyValuePair) {
	for _, attachedIf := range value.GetAttachedInterfaces() {
		derived = append(derived, api.KeyValuePair{
			Key:   abx.ToInterfaceKey(value.Index, attachedIf.InputInterface),
			Value: &emptypb.Empty{},
		})
	}
	return derived
}

// A list of ABX dependencies (ACL).
func (d *ABXDescriptor) Dependencies(key string, abxData *abx.ABX) (dependencies []api.Dependency) {
	// access list
	dependencies = append(dependencies, api.Dependency{
		Label: aclDep,
		Key:   acl.Key(abxData.AclName),
	})

	if abxData.OutputInterface != "" {
		dependencies = append(dependencies, api.Dependency{
			Label: abxInterfaceDep,
			Key:   interfaces.InterfaceKey(abxData.OutputInterface),
		})
	}

	return dependencies
}

func equivalentABXAttachedInterfaces(oldIfs, newIfs []*abx.ABX_AttachedInterface) bool {
	if len(oldIfs) != len(newIfs) {
		return false
	}
	// compare values in list ignoring order
	for _, oldIf := range oldIfs {
		var found bool
		for _, newIf := range newIfs {
			if proto.Equal(oldIf, newIf) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

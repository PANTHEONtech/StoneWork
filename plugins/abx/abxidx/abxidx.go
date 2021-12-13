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

package abxidx

import (
	"go.ligato.io/cn-infra/v2/idxmap"
	"go.ligato.io/cn-infra/v2/logging"
	abx "go.pantheon.tech/stonework/proto/abx"

	"go.ligato.io/vpp-agent/v3/pkg/idxvpp"
)

// ABXMetadataIndex provides read-only access to mapping between ABX indexes (generated in the ABX plugin)
// and ABX names.
type ABXMetadataIndex interface {
	// LookupIdx looks up previously stored item identified by index in the mapping.
	LookupByName(name string) (metadata *ABXMetadata, exists bool)

	// LookupName looks up previously stored item identified by name in the mapping.
	LookupByIndex(idx uint32) (name string, metadata *ABXMetadata, exists bool)
}

// ABXMetadataIndexRW is mapping between ABX indexes (generated in the ABX plugin) and ABX names.
type ABXMetadataIndexRW interface {
	ABXMetadataIndex
	idxmap.NamedMappingRW
}

// ABXMetadata represents metadata for ABX.
type ABXMetadata struct {
	Index    uint32
	Attached []*abx.ABX_AttachedInterface
}

// Attached is helper struct for metadata (ABX attached interface).
type Attached struct {
	Name     string
	Priority uint32
}

// GetIndex returns index of the ABX.
func (m *ABXMetadata) GetIndex() uint32 {
	return m.Index
}

// ABXMetadataDto represents an item sent through watch channel in abxIndex.
type ABXMetadataDto struct {
	idxmap.NamedMappingEvent
	Metadata *ABXMetadata
}

type abxMetadataIndex struct {
	idxmap.NamedMappingRW

	log         logging.Logger
	nameToIndex idxvpp.NameToIndex
}

// NewABXIndex creates new instance of abxMetadataIndex.
func NewABXIndex(logger logging.Logger, title string) ABXMetadataIndexRW {
	mapping := idxvpp.NewNameToIndex(logger, title, indexAbxMetadata)
	return &abxMetadataIndex{
		NamedMappingRW: mapping,
		log:            logger,
		nameToIndex:    mapping,
	}
}

// LookupByName looks up previously stored item identified by index in mapping.
func (abxIdx *abxMetadataIndex) LookupByName(name string) (metadata *ABXMetadata, exists bool) {
	meta, found := abxIdx.GetValue(name)
	if found {
		if typedMeta, ok := meta.(*ABXMetadata); ok {
			return typedMeta, found
		}
	}
	return nil, false
}

// LookupByIndex looks up previously stored item identified by name in mapping.
func (abxIdx *abxMetadataIndex) LookupByIndex(idx uint32) (name string, metadata *ABXMetadata, exists bool) {
	var item idxvpp.WithIndex
	name, item, exists = abxIdx.nameToIndex.LookupByIndex(idx)
	if exists {
		var isIfaceMeta bool
		metadata, isIfaceMeta = item.(*ABXMetadata)
		if !isIfaceMeta {
			exists = false
		}
	}
	return
}

// indexMetadata is an index function used for ABX metadata.
func indexAbxMetadata(metaData interface{}) map[string][]string {
	indexes := make(map[string][]string)

	ifMeta, ok := metaData.(*ABXMetadata)
	if !ok || ifMeta == nil {
		return indexes
	}

	return indexes
}

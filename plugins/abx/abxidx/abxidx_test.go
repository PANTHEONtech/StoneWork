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

package abxidx_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"go.ligato.io/cn-infra/v2/logging"
	"go.pantheon.tech/stonework/plugins/abx/abxidx"
)

func TestABXIndexLookupByName(t *testing.T) {
	RegisterTestingT(t)
	abxIndex := abxidx.NewABXIndex(logging.DefaultLogger, "abx-index")

	abxIndex.Put("val1", &abxidx.ABXMetadata{Index: 10})
	abxIndex.Put("val2", &abxidx.ABXMetadata{Index: 20})
	abxIndex.Put("val3", 10)

	metadata, exists := abxIndex.LookupByName("val1")
	Expect(exists).To(BeTrue())
	Expect(metadata).ToNot(BeNil())
	Expect(metadata.Index).To(Equal(uint32(10)))

	metadata, exists = abxIndex.LookupByName("val2")
	Expect(exists).To(BeTrue())
	Expect(metadata).ToNot(BeNil())
	Expect(metadata.Index).To(Equal(uint32(20)))

	metadata, exists = abxIndex.LookupByName("val3")
	Expect(exists).To(BeFalse())
	Expect(metadata).To(BeNil())

	metadata, exists = abxIndex.LookupByName("val4")
	Expect(exists).To(BeFalse())
	Expect(metadata).To(BeNil())
}

func TestABXIndexLookupByIndex(t *testing.T) {
	RegisterTestingT(t)
	abxIndex := abxidx.NewABXIndex(logging.DefaultLogger, "abx-index")

	abxIndex.Put("val1", &abxidx.ABXMetadata{Index: 10})
	abxIndex.Put("val2", &abxidx.ABXMetadata{Index: 20})

	name, metadata, exists := abxIndex.LookupByIndex(10)
	Expect(exists).To(BeTrue())
	Expect(name).To(Equal("val1"))
	Expect(metadata).ToNot(BeNil())
	Expect(metadata.Index).To(Equal(uint32(10)))

	name, metadata, exists = abxIndex.LookupByIndex(20)
	Expect(exists).To(BeTrue())
	Expect(name).To(Equal("val2"))
	Expect(metadata).ToNot(BeNil())
	Expect(metadata.Index).To(Equal(uint32(20)))

	name, metadata, exists = abxIndex.LookupByIndex(30)
	Expect(exists).To(BeFalse())
	Expect(name).To(Equal(""))
	Expect(metadata).To(BeNil())
}

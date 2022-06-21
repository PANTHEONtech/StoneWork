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

package vpp2106_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2009/vpe"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/isisx"
)

func TestDumpISISXConnections(t *testing.T) {
	ctx, isisxHandler, ifIndexes := isisxTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(
		&isisx.IsisxConnectionDetails{
			Connection: isisx.IsisxConnection{
				RxSwIfIndex: 11,
				TxSwIfIndex: 12,
			},
		},
		&isisx.IsisxConnectionDetails{
			Connection: isisx.IsisxConnection{
				RxSwIfIndex: 13,
				TxSwIfIndex: 14,
			},
		})
	ctx.MockVpp.MockReply(&vpe.ControlPingReply{})

	ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
	ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
	ifIndexes.Put("interface3", &ifaceidx.IfaceMetadata{SwIfIndex: 13})
	ifIndexes.Put("interface4", &ifaceidx.IfaceMetadata{SwIfIndex: 14})

	isisxConnetions, err := isisxHandler.DumpISISXConnections()
	Expect(err).To(Succeed())
	Expect(isisxConnetions).To(HaveLen(2))
	Expect(isisxConnetions[0].InputInterface).To(Equal("interface1"))
	Expect(isisxConnetions[0].OutputInterface).To(Equal("interface2"))
	Expect(isisxConnetions[1].InputInterface).To(Equal("interface3"))
	Expect(isisxConnetions[1].OutputInterface).To(Equal("interface4"))
}

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

package vpp2202_test

import (
	"testing"

	. "github.com/onsi/gomega"
	govppapi "go.fd.io/govpp/api"
	"go.ligato.io/cn-infra/v2/logging/logrus"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/vppmock"

	"go.pantheon.tech/stonework/plugins/binapi/vpp2202/isisx"
	"go.pantheon.tech/stonework/plugins/isisx/vppcalls"
	"go.pantheon.tech/stonework/plugins/isisx/vppcalls/vpp2202"
)

func isisxTestSetup(t *testing.T) (*vppmock.TestCtx, vppcalls.ISISXVppAPI, ifaceidx.IfaceMetadataIndexRW) {
	ctx := vppmock.SetupTestCtx(t)
	log := logrus.NewLogger("test-log")
	ifIdx := ifaceidx.NewIfaceIndex(log, "if-index")
	isisxHandler := vpp2202.NewISISXVppHandler(ctx.MockChannel, ifIdx, log)
	return ctx, isisxHandler, ifIdx
}

func TestGetISISXVersion(t *testing.T) {
	ctx, isisxHandler, _ := isisxTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&isisx.IsisxPluginGetVersionReply{
		Major: 1,
		Minor: 0,
	})
	version, err := isisxHandler.GetISISXVersion()

	Expect(err).To(BeNil())
	Expect(version).To(Equal("1.0"))
}

func TestAddISISXConnection(t *testing.T) {
	// Prepare different cases
	cases := []struct {
		Name             string
		InputInterface   string
		OutputInterface  string
		ExpectFailure    bool
		MockReply        govppapi.Message
		PrepareIfIndexes func(ifIndexes ifaceidx.IfaceMetadataIndexRW)
	}{
		{
			Name:            "happy path",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   false,
			MockReply:       &isisx.IsisxConnectionAddDelReply{},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
				ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
			},
		},
		{
			Name:            "no index for input interface",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   true,
			MockReply:       &isisx.IsisxConnectionAddDelReply{},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
			},
		},
		{
			Name:            "no index for output interface",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   true,
			MockReply:       &isisx.IsisxConnectionAddDelReply{},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
			},
		},
		{
			Name:            "error from vpp",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   true,
			MockReply:       &isisx.IsisxConnectionAddDelReply{Retval: 1},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
				ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
			},
		},
	}

	// Run all cases
	for _, td := range cases {
		t.Run(td.Name, func(t *testing.T) {
			ctx, isisxHandler, ifIndexes := isisxTestSetup(t)
			defer ctx.TeardownTestCtx()

			// prepare for case
			td.PrepareIfIndexes(ifIndexes)
			ctx.MockVpp.MockReply(td.MockReply)

			// make the call and verify
			err := isisxHandler.AddISISXConnection(td.InputInterface, td.OutputInterface)
			if td.ExpectFailure {
				Expect(err).Should(HaveOccurred())
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	}
}

func TestDeleteISISXConnection(t *testing.T) {
	// Prepare different cases
	cases := []struct {
		Name             string
		InputInterface   string
		OutputInterface  string
		ExpectFailure    bool
		MockReply        govppapi.Message
		PrepareIfIndexes func(ifIndexes ifaceidx.IfaceMetadataIndexRW)
	}{
		{
			Name:            "happy path",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   false,
			MockReply:       &isisx.IsisxConnectionAddDelReply{},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
				ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
			},
		},
		{
			Name:            "no index for input interface",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   true,
			MockReply:       &isisx.IsisxConnectionAddDelReply{},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
			},
		},
		{
			Name:            "no index for output interface",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   true,
			MockReply:       &isisx.IsisxConnectionAddDelReply{},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
			},
		},
		{
			Name:            "error from vpp",
			InputInterface:  "interface1",
			OutputInterface: "interface2",
			ExpectFailure:   true,
			MockReply:       &isisx.IsisxConnectionAddDelReply{Retval: 1},
			PrepareIfIndexes: func(ifIndexes ifaceidx.IfaceMetadataIndexRW) {
				ifIndexes.Put("interface1", &ifaceidx.IfaceMetadata{SwIfIndex: 11})
				ifIndexes.Put("interface2", &ifaceidx.IfaceMetadata{SwIfIndex: 12})
			},
		},
	}

	// Run all cases
	for _, td := range cases {
		t.Run(td.Name, func(t *testing.T) {
			ctx, isisxHandler, ifIndexes := isisxTestSetup(t)
			defer ctx.TeardownTestCtx()

			// prepare for case
			td.PrepareIfIndexes(ifIndexes)
			ctx.MockVpp.MockReply(td.MockReply)

			// make the call and verify
			err := isisxHandler.DeleteISISXConnection(td.InputInterface, td.OutputInterface)
			if td.ExpectFailure {
				Expect(err).Should(HaveOccurred())
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	}
}

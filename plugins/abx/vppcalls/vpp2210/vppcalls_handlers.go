// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 PANTHEON.tech
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

package vpp2210

import (
	govppapi "go.fd.io/govpp/api"
	"go.ligato.io/cn-infra/v2/logging"

	"go.ligato.io/vpp-agent/v3/plugins/vpp/aclplugin/aclidx"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"

	"go.pantheon.tech/stonework/plugins/abx/vppcalls"
	binapi "go.pantheon.tech/stonework/plugins/binapi/vpp2210"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2210/abx"
)

func init() {
	var msgs []govppapi.Message
	msgs = append(msgs, abx.AllMessages()...)

	vppcalls.AddAbxHandlerVersion(binapi.Version, msgs, NewABXVppHandler)
}

// ABXVppHandler is accessor for abx-related vppcalls methods
type ABXVppHandler struct {
	callsChannel govppapi.Channel
	aclIndexes   aclidx.ACLMetadataIndex
	ifIndexes    ifaceidx.IfaceMetadataIndex
	log          logging.Logger
}

// NewABXVppHandler returns new ABXVppHandler.
func NewABXVppHandler(calls govppapi.Channel, aclIdx aclidx.ACLMetadataIndex, ifIdx ifaceidx.IfaceMetadataIndex,
	log logging.Logger) vppcalls.ABXVppAPI {
	return &ABXVppHandler{
		callsChannel: calls,
		aclIndexes:   aclIdx,
		ifIndexes:    ifIdx,
		log:          log,
	}
}

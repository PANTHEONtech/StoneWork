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

package vpp2306

import (
	govppapi "go.fd.io/govpp/api"
	"go.ligato.io/cn-infra/v2/logging"

	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"

	binapi "go.pantheon.tech/stonework/plugins/binapi/vpp2306"
	natba "go.pantheon.tech/stonework/plugins/binapi/vpp2306/nat64"
	"go.pantheon.tech/stonework/plugins/nat64/vppcalls"
)

func init() {
	var msgs []govppapi.Message
	msgs = append(msgs, natba.AllMessages()...)

	vppcalls.AddNat64HandlerVersion(binapi.Version, msgs, NewNat64VppHandler)
}

// Nat64VppHandler is accessor for NAT64-related vppcalls methods.
type Nat64VppHandler struct {
	callsChannel govppapi.Channel
	ifIndexes    ifaceidx.IfaceMetadataIndex
	log          logging.Logger
}

// NewNat64VppHandler creates new instance of NAT64 vppcalls handler.
func NewNat64VppHandler(callsChan govppapi.Channel,
	ifIndexes ifaceidx.IfaceMetadataIndex, log logging.Logger,
) vppcalls.Nat64VppAPI {
	return &Nat64VppHandler{
		callsChannel: callsChan,
		ifIndexes:    ifIndexes,
		log:          log,
	}
}

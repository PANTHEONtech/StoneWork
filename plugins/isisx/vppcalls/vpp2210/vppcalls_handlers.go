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
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"

	binapi "go.pantheon.tech/stonework/plugins/binapi/vpp2210"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2210/isisx"
	"go.pantheon.tech/stonework/plugins/isisx/vppcalls"
)

func init() {
	var msgs []govppapi.Message
	msgs = append(msgs, isisx.AllMessages()...)

	vppcalls.AddIsisxHandlerVersion(binapi.Version, msgs, NewISISXVppHandler)
}

// ISISXVppHandler is accessor for isisx-related vppcalls methods
type ISISXVppHandler struct {
	callsChannel govppapi.Channel
	ifIndexes    ifaceidx.IfaceMetadataIndex
	log          logging.Logger
}

// NewISISXVppHandler returns new ISISXVppHandler.
func NewISISXVppHandler(calls govppapi.Channel, ifIdx ifaceidx.IfaceMetadataIndex,
	log logging.Logger) vppcalls.ISISXVppAPI {
	return &ISISXVppHandler{
		callsChannel: calls,
		ifIndexes:    ifIdx,
		log:          log,
	}
}

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

package bfdplugin

import (
	"go.ligato.io/cn-infra/v2/config"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/cn-infra/v2/rpc/grpc"

	"go.ligato.io/vpp-agent/v3/plugins/govppmux"
	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"
)

const (
	// PluginName is the name of the BFD Plugin.
	PluginName = "bfd"
)

// DefaultPlugin is a default instance of BfdPlugin.
var DefaultPlugin = *NewPlugin()

// NewPlugin creates a new Plugin with provided options
func NewPlugin(opts ...Option) *BfdPlugin {
	p := &BfdPlugin{}
	p.PluginName = PluginName
	p.KVScheduler = &kvscheduler.DefaultPlugin
	p.GoVpp = &govppmux.DefaultPlugin
	p.IfPlugin = &ifplugin.DefaultPlugin
	p.GRPC = &grpc.DefaultPlugin

	for _, o := range opts {
		o(p)
	}

	if p.Log == nil {
		p.Log = logging.ForPlugin(p.String())
	}

	if p.Cfg == nil {
		p.Cfg = config.ForPlugin(p.String())
	}

	return p
}

// Option is a function that can be used in NewPlugin allowing
// plugin customization
type Option func(plugin *BfdPlugin)

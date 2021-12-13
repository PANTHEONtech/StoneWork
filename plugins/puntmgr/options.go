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

package puntmgr

import (
	"go.ligato.io/cn-infra/v2/config"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/cn-infra/v2/rpc/grpc"
	"go.ligato.io/cn-infra/v2/servicelabel"

	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler"
	"go.ligato.io/vpp-agent/v3/plugins/linux/nsplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"

	"go.pantheon.tech/stonework/plugins/cnfreg"
)

const (
	// PluginName is the name of the Punting Manager Plugin.
	// Config file name is `PluginName + ".conf"`
	PluginName = "puntmgr"
)

// DefaultPlugin is a default instance of the Punting Manager.
var DefaultPlugin = *NewPlugin()

// NewPlugin creates a new Plugin with provided options
func NewPlugin(opts ...Option) *Plugin {
	p := &Plugin{}
	p.PluginName = PluginName
	p.ServiceLabel = &servicelabel.DefaultPlugin
	p.GRPCServer = &grpc.DefaultPlugin
	p.CnfRegistry = &cnfreg.DefaultPlugin
	p.IfPlugin = &ifplugin.DefaultPlugin
	p.NsPlugin = &nsplugin.DefaultPlugin
	p.CfgClient = client.LocalClient
	p.KVScheduler = &kvscheduler.DefaultPlugin

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
type Option func(plugin *Plugin)

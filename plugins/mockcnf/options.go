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

package mockcnf

import (
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler"
	cnfreg_plugin "go.pantheon.tech/stonework/plugins/cnfreg"
	puntmgr_plugin "go.pantheon.tech/stonework/plugins/puntmgr"
)

var DefaultPlugin = *NewPlugin()

func NewPlugin(opts ...Option) *Plugin {
	p := &Plugin{}

	p.PluginName = "mockcnf"
	p.KVScheduler = &kvscheduler.DefaultPlugin
	p.PuntManager = &puntmgr_plugin.DefaultPlugin
	p.CnfRegistry = &cnfreg_plugin.DefaultPlugin

	for _, o := range opts {
		o(p)
	}

	if p.Log == nil {
		p.Log = logging.ForPlugin(p.String())
	}

	return p
}

// Option is a function that can be used in NewPlugin to customize Plugin.
type Option func(*Plugin)

// UseDeps returns Option that can inject custom dependencies.
func UseDeps(f func(*Deps)) Option {
	return func(p *Plugin) {
		f(&p.Deps)
	}
}

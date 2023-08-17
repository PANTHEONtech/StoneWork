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

//go:generate descriptor-adapter --descriptor-name ABX --value-type *vpp_abx.ABX --meta-type *abxidx.ABXMetadata --import "go.pantheon.tech/stonework/plugins/abx/abxidx" --import "go.pantheon.tech/stonework/proto/abx" --output-dir "descriptor"

package abxplugin

import (
	"github.com/go-errors/errors"
	govppapi "go.fd.io/govpp/api"
	"go.ligato.io/cn-infra/v2/health/statuscheck"
	"go.ligato.io/cn-infra/v2/infra"
	"go.ligato.io/vpp-agent/v3/plugins/govppmux"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/aclplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"

	"go.pantheon.tech/stonework/plugins/abx/abxidx"
	"go.pantheon.tech/stonework/plugins/abx/descriptor"
	"go.pantheon.tech/stonework/plugins/abx/vppcalls"

	_ "go.pantheon.tech/stonework/plugins/abx/vppcalls/vpp2202"
	_ "go.pantheon.tech/stonework/plugins/abx/vppcalls/vpp2210"
	_ "go.pantheon.tech/stonework/plugins/abx/vppcalls/vpp2306"
)

// ABXPlugin is a plugin that manages ACL-based forwarding.
type ABXPlugin struct {
	Deps

	// GoVPP channels
	vppCh govppapi.Channel

	abxHandler             vppcalls.ABXVppAPI
	abxDescriptor          *descriptor.ABXDescriptor
	abxInterfaceDescriptor *descriptor.ABXToInterfaceDescriptor

	// index maps
	abxIndex abxidx.ABXMetadataIndex
}

// Deps represents dependencies for the plugin.
type Deps struct {
	infra.PluginDeps
	Scheduler   kvs.KVScheduler
	GoVppmux    govppmux.API
	ACLPlugin   aclplugin.API
	IfPlugin    ifplugin.API
	StatusCheck statuscheck.PluginStatusWriter // optional
}

// Init initializes ABX plugin.
func (p *ABXPlugin) Init() error {
	// init handler
	p.abxHandler = vppcalls.CompatibleAbxVppHandler(p.GoVppmux, p.ACLPlugin.GetACLIndex(), p.IfPlugin.GetInterfaceIndex(), p.Log)
	if p.abxHandler == nil {
		return errors.New("abxHandler is not available")
	}

	// init & register descriptor
	abxDescriptor := descriptor.NewABXDescriptor(p.abxHandler, p.ACLPlugin.GetACLIndex(), p.Log)
	if err := p.Deps.Scheduler.RegisterKVDescriptor(abxDescriptor); err != nil {
		return err
	}

	// obtain read-only reference to index map
	var withIndex bool
	metadataMap := p.Scheduler.GetMetadataMap(abxDescriptor.Name)
	p.abxIndex, withIndex = metadataMap.(abxidx.ABXMetadataIndex)
	if !withIndex {
		return errors.New("missing index with ABX metadata")
	}

	// init & register derived value descriptor
	abxInterfaceDescriptor := descriptor.NewABXToInterfaceDescriptor(p.abxIndex, p.abxHandler, p.IfPlugin, p.Log)
	if err := p.Deps.Scheduler.RegisterKVDescriptor(abxInterfaceDescriptor); err != nil {
		return err
	}

	return nil
}

// AfterInit registers plugin with StatusCheck.
func (p *ABXPlugin) AfterInit() error {
	if p.StatusCheck != nil {
		p.StatusCheck.Register(p.PluginName, nil)
	}
	return nil
}

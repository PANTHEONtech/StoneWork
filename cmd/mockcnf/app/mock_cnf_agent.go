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

package app

import (
	"fmt"

	"go.ligato.io/cn-infra/v2/datasync"
	"go.ligato.io/cn-infra/v2/datasync/kvdbsync"
	"go.ligato.io/cn-infra/v2/datasync/kvdbsync/local"
	"go.ligato.io/cn-infra/v2/db/keyval/etcd"
	"go.ligato.io/cn-infra/v2/health/probe"
	"go.ligato.io/cn-infra/v2/health/statuscheck"
	"go.ligato.io/cn-infra/v2/infra"
	"go.ligato.io/cn-infra/v2/logging/logmanager"
	"go.ligato.io/cn-infra/v2/rpc/rest"

	"go.ligato.io/vpp-agent/v3/plugins/configurator"
	linux_ifplugin "go.ligato.io/vpp-agent/v3/plugins/linux/ifplugin"
	linux_iptablesplugin "go.ligato.io/vpp-agent/v3/plugins/linux/iptablesplugin"
	linux_l3plugin "go.ligato.io/vpp-agent/v3/plugins/linux/l3plugin"
	linux_nsplugin "go.ligato.io/vpp-agent/v3/plugins/linux/nsplugin"
	"go.ligato.io/vpp-agent/v3/plugins/netalloc"
	"go.ligato.io/vpp-agent/v3/plugins/orchestrator"
	"go.ligato.io/vpp-agent/v3/plugins/orchestrator/localregistry"
	"go.ligato.io/vpp-agent/v3/plugins/orchestrator/watcher"
	"go.ligato.io/vpp-agent/v3/plugins/restapi"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/l2plugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/l3plugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/puntplugin"

	abx "go.pantheon.tech/stonework/plugins/abx"
	cnfreg_plugin "go.pantheon.tech/stonework/plugins/cnfreg"
	isisx "go.pantheon.tech/stonework/plugins/isisx"
	mockcnf_plugin "go.pantheon.tech/stonework/plugins/mockcnf"
	puntmgr_plugin "go.pantheon.tech/stonework/plugins/puntmgr"
	"go.pantheon.tech/stonework/proto/cnfreg"
	"go.pantheon.tech/stonework/proto/mockcnf"
)

// MockCnfAgent defines plugins which will be loaded and their order.
type MockCnfAgent struct {
	infra.PluginName
	LogManager *logmanager.Plugin

	VPP
	Linux
	Netalloc    *netalloc.Plugin
	MockCnf     *mockcnf_plugin.Plugin
	CnfRegistry *cnfreg_plugin.Plugin
	PuntManager *puntmgr_plugin.Plugin

	Orchestrator *orchestrator.Plugin
	ETCDDataSync *kvdbsync.Plugin

	Configurator *configurator.Plugin
	RESTAPI      *restapi.Plugin
	Probe        *probe.Plugin
	StatusCheck  *statuscheck.Plugin
}

// New creates new MockCnfAgent instance.
func New() *MockCnfAgent {
	etcdDataSync := kvdbsync.NewPlugin(kvdbsync.UseKV(&etcd.DefaultPlugin))

	writers := datasync.KVProtoWriters{
		etcdDataSync,
	}
	statuscheck.DefaultPlugin.Transport = writers

	initFileRegistry := localregistry.NewInitFileRegistryPlugin()
	watchers := watcher.NewPlugin(watcher.UseWatchers(
		local.DefaultRegistry,
		initFileRegistry,
		etcdDataSync,
	))
	orchestrator.DefaultPlugin.Watcher = watchers
	orchestrator.DefaultPlugin.StatusPublisher = writers
	orchestrator.EnabledGrpcMetrics()

	ifplugin.DefaultPlugin.Watcher = etcdDataSync
	puntplugin.DefaultPlugin.PublishState = writers

	cnfreg_plugin.DefaultPlugin.PuntMgr = &puntmgr_plugin.DefaultPlugin
	cnfreg_plugin.DefaultPlugin.HTTPPlugin = &rest.DefaultPlugin
	cnfreg_plugin.DefaultPlugin.CnfIndex = mockcnf.MockCnfIndex()

	var (
		vpp                VPP
		configuratorPlugin *configurator.Plugin
		restPlugin         *restapi.Plugin
	)
	switch cnfreg_plugin.DefaultPlugin.GetCnfMode() {
	case cnfreg.CnfMode_STANDALONE:
		vpp = DefaultVPP()
		configuratorPlugin = &configurator.DefaultPlugin
		restPlugin = &restapi.DefaultPlugin
		linux_ifplugin.DefaultPlugin.VppIfPlugin = &ifplugin.DefaultPlugin
		ifplugin.DefaultPlugin.LinuxIfPlugin = &linux_ifplugin.DefaultPlugin
		ifplugin.DefaultPlugin.NsPlugin = &linux_nsplugin.DefaultPlugin
	case cnfreg.CnfMode_STONEWORK_MODULE:
		vpp = DisabledVPP()
		puntmgr_plugin.DefaultPlugin.IfPlugin = nil
	case cnfreg.CnfMode_STONEWORK:
		panic("invalid CNF mode")
	}
	linux := DefaultLinux()

	return &MockCnfAgent{
		PluginName:   "MockCNF",
		LogManager:   &logmanager.DefaultPlugin,
		Orchestrator: &orchestrator.DefaultPlugin,
		ETCDDataSync: etcdDataSync,
		VPP:          vpp,
		Linux:        linux,
		Netalloc:     &netalloc.DefaultPlugin,
		MockCnf:      &mockcnf_plugin.DefaultPlugin,
		CnfRegistry:  &cnfreg_plugin.DefaultPlugin,
		PuntManager:  &puntmgr_plugin.DefaultPlugin,
		Configurator: configuratorPlugin,
		RESTAPI:      restPlugin,
		Probe:        &probe.DefaultPlugin,
		StatusCheck:  &statuscheck.DefaultPlugin,
	}
}

// Init initializes main plugin.
func (a *MockCnfAgent) Init() error {
	a.StatusCheck.Register(a.PluginName, nil)
	a.StatusCheck.ReportStateChange(a.PluginName, statuscheck.Init, nil)
	return nil
}

// AfterInit executes resync.
func (a *MockCnfAgent) AfterInit() error {
	if err := orchestrator.DefaultPlugin.InitialSync(); err != nil {
		return fmt.Errorf("failure in initial sync: %v", err)
	}
	a.StatusCheck.ReportStateChange(a.PluginName, statuscheck.OK, nil)
	return nil
}

// Close could close used resources.
func (a *MockCnfAgent) Close() error {
	return nil
}

// VPP contains some VPP plugins.
type VPP struct {
	IfPlugin *ifplugin.IfPlugin
	L2Plugin *l2plugin.L2Plugin
	L3Plugin *l3plugin.L3Plugin
	ABX      *abx.ABXPlugin
	ISISX    *isisx.ISISXPlugin
}

func DefaultVPP() VPP {
	return VPP{
		IfPlugin: &ifplugin.DefaultPlugin,
		L2Plugin: &l2plugin.DefaultPlugin,
		L3Plugin: &l3plugin.DefaultPlugin,
		ABX:      &abx.DefaultPlugin,
		ISISX:    &isisx.DefaultPlugin,
	}
}

func DisabledVPP() VPP {
	return VPP{}
}

// Linux contains some Linux plugins.
type Linux struct {
	IfPlugin       *linux_ifplugin.IfPlugin
	L3Plugin       *linux_l3plugin.L3Plugin
	NSPlugin       *linux_nsplugin.NsPlugin
	IPTablesPlugin *linux_iptablesplugin.IPTablesPlugin
}

func DefaultLinux() Linux {
	return Linux{
		IfPlugin:       &linux_ifplugin.DefaultPlugin,
		L3Plugin:       &linux_l3plugin.DefaultPlugin,
		NSPlugin:       &linux_nsplugin.DefaultPlugin,
		IPTablesPlugin: &linux_iptablesplugin.DefaultPlugin,
	}
}

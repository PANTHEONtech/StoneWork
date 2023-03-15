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
	"go.ligato.io/vpp-agent/v3/plugins/telemetry"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/abfplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/aclplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ipsecplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/l2plugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/l3plugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/natplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/puntplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/srplugin"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/stnplugin"

	abx "go.pantheon.tech/stonework/plugins/abx"
	bfd "go.pantheon.tech/stonework/plugins/bfd"
	"go.pantheon.tech/stonework/plugins/cnfreg"
	isisxplugin "go.pantheon.tech/stonework/plugins/isisx"
	"go.pantheon.tech/stonework/plugins/nat64"
	"go.pantheon.tech/stonework/plugins/puntmgr"
)

// StoneWorkAgent defines plugins which will be loaded and their order.
type StoneWorkAgent struct {
	infra.PluginName
	LogManager *logmanager.Plugin

	VPP
	Linux
	Netalloc    *netalloc.Plugin
	CnfRegistry *cnfreg.Plugin
	PuntManager *puntmgr.Plugin

	Orchestrator *orchestrator.Plugin
	ETCDDataSync *kvdbsync.Plugin

	Configurator *configurator.Plugin
	RESTAPI      *restapi.Plugin
	Probe        *probe.Plugin
	StatusCheck  *statuscheck.Plugin
	Telemetry    *telemetry.Plugin
}

// New creates new StoneWorkAgent instance.
func New() *StoneWorkAgent {
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

	linux_ifplugin.DefaultPlugin.VppIfPlugin = &ifplugin.DefaultPlugin
	ifplugin.DefaultPlugin.LinuxIfPlugin = &linux_ifplugin.DefaultPlugin
	ifplugin.DefaultPlugin.NsPlugin = &linux_nsplugin.DefaultPlugin

	cnfreg.DefaultPlugin.PuntMgr = &puntmgr.DefaultPlugin

	vpp := DefaultVPP()
	linux := DefaultLinux()

	defaultConfig := telemetry.DefaultConfig
	telemetry.DefaultConfig = func() *telemetry.Config {
		cfg := defaultConfig()
		cfg.Disabled = true
		return cfg
	}

	return &StoneWorkAgent{
		PluginName:   "StoneWorkAgent",
		LogManager:   &logmanager.DefaultPlugin,
		Orchestrator: &orchestrator.DefaultPlugin,
		ETCDDataSync: etcdDataSync,
		VPP:          vpp,
		Linux:        linux,
		Netalloc:     &netalloc.DefaultPlugin,
		CnfRegistry:  &cnfreg.DefaultPlugin,
		PuntManager:  &puntmgr.DefaultPlugin,
		Configurator: &configurator.DefaultPlugin,
		RESTAPI:      &restapi.DefaultPlugin,
		Probe:        &probe.DefaultPlugin,
		StatusCheck:  &statuscheck.DefaultPlugin,
		Telemetry:    &telemetry.DefaultPlugin,
	}
}

// Init initializes main plugin.
func (a *StoneWorkAgent) Init() error {
	a.StatusCheck.Register(a.PluginName, nil)
	a.StatusCheck.ReportStateChange(a.PluginName, statuscheck.Init, nil)
	return nil
}

// AfterInit executes resync.
func (a *StoneWorkAgent) AfterInit() error {
	if err := orchestrator.DefaultPlugin.InitialSync(); err != nil {
		return fmt.Errorf("failure in initial sync: %v", err)
	}
	a.StatusCheck.ReportStateChange(a.PluginName, statuscheck.OK, nil)
	return nil
}

// Close could close used resources.
func (a *StoneWorkAgent) Close() error {
	return nil
}

// VPP contains all VPP plugins.
type VPP struct {
	ABFPlugin   *abfplugin.ABFPlugin
	ACLPlugin   *aclplugin.ACLPlugin
	IfPlugin    *ifplugin.IfPlugin
	IPSecPlugin *ipsecplugin.IPSecPlugin
	L2Plugin    *l2plugin.L2Plugin
	L3Plugin    *l3plugin.L3Plugin
	NATPlugin   *natplugin.NATPlugin
	NAT64Plugin *nat64plugin.NAT64Plugin
	PuntPlugin  *puntplugin.PuntPlugin
	STNPlugin   *stnplugin.STNPlugin
	SRPlugin    *srplugin.SRPlugin
	ABX         *abx.ABXPlugin
	ISISX       *isisxplugin.ISISXPlugin
	BFD         *bfd.BfdPlugin
}

func DefaultVPP() VPP {
	return VPP{
		ABFPlugin:   &abfplugin.DefaultPlugin,
		ACLPlugin:   &aclplugin.DefaultPlugin,
		IfPlugin:    &ifplugin.DefaultPlugin,
		IPSecPlugin: &ipsecplugin.DefaultPlugin,
		L2Plugin:    &l2plugin.DefaultPlugin,
		L3Plugin:    &l3plugin.DefaultPlugin,
		NATPlugin:   &natplugin.DefaultPlugin,
		NAT64Plugin: &nat64plugin.DefaultPlugin,
		PuntPlugin:  &puntplugin.DefaultPlugin,
		STNPlugin:   &stnplugin.DefaultPlugin,
		SRPlugin:    &srplugin.DefaultPlugin,
		ABX:         &abx.DefaultPlugin,
		ISISX:       &isisxplugin.DefaultPlugin,
		BFD:         &bfd.DefaultPlugin,
	}
}

// Linux contains all Linux plugins.
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

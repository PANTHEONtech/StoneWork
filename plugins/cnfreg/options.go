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

package cnfreg

import (
	"os"

	"go.ligato.io/cn-infra/v2/config"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/cn-infra/v2/rpc/grpc"
	"go.ligato.io/cn-infra/v2/servicelabel"

	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler"

	pb "go.pantheon.tech/stonework/proto/cnfreg"
)

const (
	// PluginName is the name of the CNF Registry Plugin.
	PluginName = "cnfreg"

	// Environment variable that is used to select the CNF mode.
	CnfModeEnvVar = "CNF_MODE"

	// Environment variable that is used to select network interface for CNF management.
	// If not defined then it will be selected automatically.
	CnfMgmtInterfaceEnvVar = "CNF_MGMT_INTERFACE"
	// Environment variable that is used to specify what IP subnet the management
	// IP addresses are allocated from. For example, if CNFs are deployed using docker
	// compose, then this variable should contain the subnet (~ IP pool) used by the bridge.
	// It is an alternative hint to CNF_MGMT_INTERFACE (of lower priority) for
	// correctly selecting the management interface.
	CnfMgmtSubnetEnvVar = "CNF_MGMT_SUBNET"
)

// DefaultPlugin is a default instance of CNF Registry.
var DefaultPlugin = *NewPlugin()

// NewPlugin creates a new Plugin with provided options
func NewPlugin(opts ...Option) *Plugin {
	p := &Plugin{}
	p.PluginName = PluginName
	p.MgmtInterface = os.Getenv(CnfMgmtInterfaceEnvVar)
	p.MgmtSubnet = os.Getenv(CnfMgmtSubnetEnvVar)
	p.ServiceLabel = &servicelabel.DefaultPlugin
	p.KVScheduler = &kvscheduler.DefaultPlugin
	p.GRPCPlugin = &grpc.DefaultPlugin

	// Note: Punt Manager not injected by default due to a cyclical dependency between these two plugins
	//p.PuntMgr = &puntmgr_plugin.DefaultPlugin

	cnfMode := pb.CnfMode_value[os.Getenv(CnfModeEnvVar)]
	p.cnfMode = pb.CnfMode(cnfMode)

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

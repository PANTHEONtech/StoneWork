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
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"go.ligato.io/cn-infra/v2/infra"
	grpc_plugin "go.ligato.io/cn-infra/v2/rpc/grpc"
	"go.ligato.io/cn-infra/v2/rpc/rest"
	"go.ligato.io/cn-infra/v2/servicelabel"

	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/client/remoteclient"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/proto/ligato/generic"

	pb "go.pantheon.tech/stonework/proto/cnfreg"
	"go.pantheon.tech/stonework/proto/puntmgr"
)

// CNF Registry plugin allows to load a CNF module into the StoneWork (all-in-one VPP distribution; SW for short)
// during the Init phase. CNF can be built as another image and run as a separate container.
// This allows to enable/disable CNF without having to rebuild StoneWork docker image or the agent binary.
// Apart from single common VPP it is also possible to share network namespace between all/some CNFs and therefore
// integrate different network functions inside the Linux network stack
// (e.g. to use OSPF-learned routes to connect with a BGP peer).
// The plugin operates in one of the 3 following modes depending on the value of the "CNF_MODE" environment variable:
//  1. STANDALONE (default, i.e. assumed if the variable is not defined):
//     - CNF is used on its own, potentially chained with other CNFs using for example NSM
//     (i.e. each VPP-based CNF runs its own VPP instance)
//     - The CNF Registry plugin is also used by a Standalone CNF, but merely to keep track of CNF Index ID .
//  2. STONEWORK_MODULE:
//     - CNF used as a SW-module
//     - VPP-based CNFs do not run VPP inside their container, instead they connect with the all-in-one VPP of StoneWork
//     - in this mode the Registry acts as a client of the Registry running by the StoneWork agent
//     - internally the plugin uses gRPC to exchange all the information needed between the Registries of CNF and SW
//     to load the CNF and use with the all-in-one VPP
//     - CNF should use only those methods of the plugin which are defined by the CnfAPI interface
//  3. STONEWORK:
//     - CNF Registry plugin is used by StoneWork firstly to discover all the enabled CNFs and then to collect
//     all the information about them to be able to integrate them with the all-in-one VPP
//     - for each CNF, StoneWork needs to learn the NB-facing configuration models, traffic Punting to use (some
//     subset of the traffic typically needs to be diverted from VPP into the Linux network stack for the Linux-based
//     CNF to process)
//     - StoneWork should use only those methods of the plugin which are defined by the StoneWorkAPI interface
type Plugin struct {
	Deps
	pb.UnimplementedCnfDiscoveryServer

	cnfMode   pb.CnfMode
	config    *Config
	ipAddress net.IP // management IP address (discovered during Init)

	// CNF-mode specific attributes
	sw    swAttrs    // STONEWORK
	swMod swModAttrs // STONEWORK_MODULE
}

// Deps is a set of dependencies of the CNF Registry plugin
type Deps struct {
	infra.PluginDeps
	CnfDeps
	MgmtInterface string
	MgmtSubnet    string
	ServiceLabel  servicelabel.ReaderAPI
	KVScheduler   kvs.KVScheduler
	GRPCPlugin    *grpc_plugin.Plugin
	HTTPPlugin    *rest.Plugin // optional
	PuntMgr       PuntManagerAPI
}

// Dependencies to define/provide for CNF running standalone or as a SW-Module.
type CnfDeps struct {
	// CNF index that should be unique for each CNF.
	// Allocate 1 for the first CNF and give +1 for each new CNF.
	CnfIndex int
}

// CnfRegistryAPI encapsulates all methods exposed by CNF Registry.
// Please note, however, that not all the methods may be available in the given CNF mode.
type CnfRegistryAPI interface {
	CnfAPI
	StoneWorkAPI
}

// PuntRequestsClb is used to inform the plugin about how the packet punting should be configured
// for a given configuration item.
type PuntRequestsClb func(configItem proto.Message) *puntmgr.PuntRequests

// ItemDepsClb is used to inform the plugin about all the dependencies of a given configuration
// item (apart from punt dependencies which are determined from the punt requests).
type ItemDepsClb func(configItem proto.Message) []*pb.ConfigItemDependency

// CnfModelCallbacks groups (optional) callbacks that can be assigned to a model registration.
type CnfModelCallbacks struct {
	PuntRequests     PuntRequestsClb
	ItemDependencies ItemDepsClb
}

// API to be used by CNF (standalone or as SW-Module)
type CnfAPI interface {
	// Get mode in which the CNF operates with respect to other CNFs.
	GetCnfMode() pb.CnfMode
	// Returns gRPC port that should be used by this CNF.
	// Not to be used by StoneWork or a standalone CNF (they should respect what is in grpc.conf).
	GetGrpcPort() (port int)
	// Returns HTTP port that should be used by this CNF.
	// Not to be used by StoneWork or a standalone CNF (they should respect what is in http.conf).
	GetHttpPort() (port int)
	// Returns gRPC connection established with StoneWork.
	GetSWGrpcConn() (conn grpc.ClientConnInterface, err error)
	// Returns remote configuration client connected with the StoneWork.
	// Can be used to apply configuration into the VPP
	// (determined by the operation of a CNF, e.g. IP routes received over BGP).
	GetSWCfgClient() (cfgClient client.GenericClient, err error)
	// RegisterCnfModel registers configuration model implemented by the CNF.
	// This method should be used by a SW-Module CNF in the Init phase to convey the definition of a CNF NB API model
	// into StoneWork, which will then act as a proxy for all the operations over that model.
	RegisterCnfModel(model models.KnownModel, descriptor *kvs.KVDescriptor, callbacks *CnfModelCallbacks) error
}

// API to be used by StoneWork.
type StoneWorkAPI interface {
	// Returns gRPC connection established with the given SW-Module CNF.
	GetCnfGrpcConn(cnfMsLabel string) (conn grpc.ClientConnInterface, err error)
	// Returns remote configuration client connected with the given SW-Module CNF.
	GetCnfCfgClient(cnfMsLabel string) (cfgClient client.GenericClient, err error)
}

// APIs of Punt Manager that are used by CNF registry.
type PuntManagerAPI interface {
	// AddPunt is used by StoneWork or standalone CNF to configure punt between VPP and the CNF.
	AddPunt(cnfMsLabel, key string, punt *puntmgr.PuntRequest) error
	// DelPunt is used by StoneWork or standalone CNF to un-configure punt between VPP and the CNF.
	DelPunt(cnfMsLabel, key string, label string) error
	// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
	GetPuntDependencies(cnfMsLabel string, punt *puntmgr.PuntRequest) (deps []kvs.Dependency)
}

// Attributes specific to StoneWork (i.e. not used by CNF).
type swAttrs struct {
	modules map[string]swModule // key = cnf microservice label
}

// CNF used as a StoneWork Module.
type swModule struct {
	cnfMsLabel string
	ipAddress  string
	grpcPort   int
	httpPort   int
	grpcConn   grpc.ClientConnInterface
	cnfClient  pb.CnfDiscoveryClient
	cfgClient  client.GenericClient
	cnfModels  []cnfModel
}

// Attributes specific to StoneWork Module (i.e. not used by standalone CNF or StoneWork itself).
type swModAttrs struct {
	sync.Mutex
	discovered  bool
	swIpAddress string
	swGrpcPort  int
	swHttpPort  int
	swGrpcConn  grpc.ClientConnInterface
	swCfgClient client.GenericClient
	models      []exposedModel
}

type exposedModel struct {
	model      models.KnownModel
	descriptor *kvs.KVDescriptor
	callbacks  *CnfModelCallbacks
}

// Model exposed by a CNF.
type cnfModel struct {
	info         *models.ModelInfo
	withPunt     bool
	withDeps     bool
	withRetrieve bool
}

// Init initializes internal attributes and depending on the mode does the following:
// case STONEWORK_MODULE:
//   - registers gRPC handler for DiscoverCnf and GetPuntRequests
//
// case STONEWORK:
//   - waits few seconds for all CNFs to write pid files
//   - then for each CNF:
//   - creates grpcConnection with the CNF
//   - obtains models using the meta service
//   - register models (exposed by Cnf)
//   - creates CnfDescriptorProxy for each module
func (p *Plugin) Init() (err error) {
	if p.GetCnfMode() != pb.CnfMode_STONEWORK {
		// check CNF dependencies
		if p.CnfIndex == 0 {
			return errors.New("CnfIndex not defined")
		}
	}
	p.config, err = p.loadConfig()
	if err != nil {
		return err
	}

	// discover management IP address
	// TODO: support FQDN (e.g. K8s service name)
	p.ipAddress, err = p.discoverMyIP()
	if err != nil {
		if p.cnfMode == pb.CnfMode_STANDALONE {
			// Standalone CNF does not really need management IP address
			p.Log.Warn(err)
			err = nil
		} else {
			return err
		}
	} else {
		p.Log.Infof("Discovered management IP address: %v", p.ipAddress)
	}

	switch p.cnfMode {
	case pb.CnfMode_STONEWORK_MODULE:
		// inject gRPC and HTTP ports to use by SW-Module
		p.GRPCPlugin.Config.Endpoint = fmt.Sprintf("0.0.0.0:%d", p.GetGrpcPort())
		if p.HTTPPlugin != nil {
			p.HTTPPlugin.Config.Endpoint = fmt.Sprintf("0.0.0.0:%d", p.GetHttpPort())
		}

		// serve CnfDiscovery methods
		grpcServer := p.GRPCPlugin.GetServer()
		if grpcServer == nil {
			return errors.New("gRPC server is not initialized")
		}
		pb.RegisterCnfDiscoveryServer(grpcServer, p)

	case pb.CnfMode_STONEWORK:
		// CNF "discovery"
		go p.cnfDiscovery(make(chan struct{}))
		p.Log.Debugf("Discovered CNFs: %+v", p.sw.modules)

		// setup proxy for each config module exposed by every CNF
		// for _, swMod := range p.sw.modules {
		// 	err := p.initCnfProxy(swMod)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}
	return nil
}

// AfterInit is used by CNF-Module to write pid file under a known directory for StoneWork to discover it.
func (p *Plugin) AfterInit() (err error) {
	if p.cnfMode == pb.CnfMode_STONEWORK_MODULE {
		err = p.writePidFile()
		if err != nil {
			return err
		}
	}
	return nil
}

// Close is NOOP.
func (p *Plugin) Close() error {
	return nil
}

// Get mode in which the CNF operates with respect to other CNFs.
func (p *Plugin) GetCnfMode() pb.CnfMode {
	return p.cnfMode
}

// Returns gRPC port that should be used by this CNF.
// Not to be used by StoneWork or a standalone CNF (they should respect what is in grpc.conf).
func (p *Plugin) GetGrpcPort() (port int) {
	if p.cnfMode != pb.CnfMode_STONEWORK_MODULE {
		panic(fmt.Errorf("method GetGrpcPort is not available in the CNF mode %v", p.cnfMode))
	}
	return p.config.SwModGrpcBasePort + p.CnfIndex
}

// Returns gRPC port that should be used by this CNF.
// Not to be used by StoneWork or a standalone CFN (they should respect what is in http.conf).
func (p *Plugin) GetHttpPort() (port int) {
	if p.cnfMode != pb.CnfMode_STONEWORK_MODULE {
		panic(fmt.Errorf("method GetHttpPort is not available in the CNF mode %v", p.cnfMode))
	}
	return p.config.SwModHttpBasePort + p.CnfIndex
}

// Returns gRPC connection established with StoneWork.
func (p *Plugin) GetSWGrpcConn() (conn grpc.ClientConnInterface, err error) {
	if p.cnfMode != pb.CnfMode_STONEWORK_MODULE {
		panic(fmt.Errorf("method GetSWCfgClient is not available in the CNF mode %v", p.cnfMode))
	}
	if p.swMod.swGrpcConn == nil {
		return nil, errors.New("gRPC connection with StoneWork is not yet established")
	}
	return p.swMod.swGrpcConn, nil
}

// Returns remote configuration client connected with the StoneWork.
// Can be used in combination to apply configuration into the VPP
// (determined by the operation of a CNF, e.g. IP routes received over BGP).
func (p *Plugin) GetSWCfgClient() (cfgClient client.GenericClient, err error) {
	if p.cnfMode != pb.CnfMode_STONEWORK_MODULE {
		panic(fmt.Errorf("method GetSWCfgClient is not available in the CNF mode %v", p.cnfMode))
	}
	if p.swMod.swCfgClient == nil {
		return nil, errors.New("gRPC connection with StoneWork is not yet established")
	}
	return p.swMod.swCfgClient, nil
}

// RegisterCnfModel registers configuration model implemented by the CNF.
// This method should be used by a SW-Module CNF in the Init phase to convey the definition of a CNF NB API model
// into StoneWork, which will then act as a proxy for all the operations over that model.
func (p *Plugin) RegisterCnfModel(model models.KnownModel, descriptor *kvs.KVDescriptor, callbacks *CnfModelCallbacks) error {
	switch p.cnfMode {
	case pb.CnfMode_STONEWORK_MODULE:
		p.swMod.Lock()
		defer p.swMod.Unlock()
		if p.swMod.discovered {
			return errors.New("CNF has been already discovered by StoneWork")
		}
		p.swMod.models = append(p.swMod.models, exposedModel{
			model:      model,
			descriptor: descriptor,
			callbacks:  callbacks,
		})

	case pb.CnfMode_STANDALONE:
		// nothing to do
		return nil
	case pb.CnfMode_STONEWORK:
		panic(fmt.Errorf("method RegisterCnfModel is not available in the CNF mode %v", p.cnfMode))
	}
	return nil
}

// Returns gRPC connection established with the given SW-Module CNF.
func (p *Plugin) GetCnfGrpcConn(cnfMsLabel string) (conn grpc.ClientConnInterface, err error) {
	if p.cnfMode != pb.CnfMode_STONEWORK {
		panic(fmt.Errorf("method GetCnfGrpcConn is not available in the CNF mode %v", p.cnfMode))
	}
	// No need to lock p.sw - it is not changed anymore after Init
	swModule, loaded := p.sw.modules[cnfMsLabel]
	if !loaded {
		return nil, fmt.Errorf("CNF %s is not loaded as StoneWork Module", cnfMsLabel)
	}
	if swModule.grpcConn == nil {
		return nil, fmt.Errorf("gRPC connection with CNF %s is not yet established", cnfMsLabel)
	}
	return swModule.grpcConn, nil
}

// Returns remote configuration client connected with the given SW-Module CNF.
func (p *Plugin) GetCnfCfgClient(cnfMsLabel string) (cfgClient client.GenericClient, err error) {
	if p.cnfMode != pb.CnfMode_STONEWORK {
		panic(fmt.Errorf("method GetCnfGrpcConn is not available in the CNF mode %v", p.cnfMode))
	}
	// No need to lock p.sw - it is not changed anymore after Init
	swModule, loaded := p.sw.modules[cnfMsLabel]
	if !loaded {
		return nil, fmt.Errorf("CNF %s is not loaded as StoneWork Module", cnfMsLabel)
	}
	if swModule.cfgClient == nil {
		return nil, fmt.Errorf("gRPC connection with CNF %s is not yet established", cnfMsLabel)
	}
	return swModule.cfgClient, nil
}

// DiscoverCnf is served by the CNFRegistry of each SW-Module CNF.
// It is called by StoneWork during Init of CNFRegistry.
func (p *Plugin) DiscoverCnf(ctx context.Context, req *pb.DiscoverCnfReq) (resp *pb.DiscoverCnfResp, err error) {
	p.Log.Debugf("Handling DiscoverCnf(%+v)", req)
	resp = &pb.DiscoverCnfResp{CnfMsLabel: p.ServiceLabel.GetAgentLabel()}
	p.swMod.Lock()
	defer p.swMod.Unlock()
	if p.swMod.discovered {
		return resp, errors.New("CNF has been already discovered")
	}
	for _, expModel := range p.swMod.models {
		resp.ConfigModels = append(resp.ConfigModels, &pb.DiscoverCnfResp_ConfigModel{
			ProtoName:    expModel.model.ProtoName(),
			WithPunt:     expModel.callbacks != nil && expModel.callbacks.PuntRequests != nil,
			WithDeps:     expModel.callbacks != nil && expModel.callbacks.ItemDependencies != nil,
			WithRetrieve: expModel.descriptor.Retrieve != nil,
		})
	}
	p.swMod.discovered = true
	p.swMod.swIpAddress = req.GetSwIpAddress()
	p.swMod.swGrpcPort = int(req.GetSwGrpcPort())
	p.swMod.swHttpPort = int(req.GetSwHttpPort())

	// establish connection with StoneWork asynchronously in the background
	// (at this point StoneWork is still in the Init phase and the gRPC server is not listening yet)
	p.swMod.swGrpcConn, err = grpc.Dial(
		fmt.Sprintf("%s:%d", p.swMod.swIpAddress, p.swMod.swGrpcPort), grpc.WithInsecure())
	p.swMod.swCfgClient, err = remoteclient.NewClientGRPC(p.swMod.swGrpcConn)
	return resp, err
}

// GetPuntRequests is served by CNFRegistry of a SW-Module CNF and returns the set of packet punting
// requests corresponding to the given configuration item.
func (p *Plugin) GetPuntRequests(ctx context.Context, item *generic.Item) (puntReqs *puntmgr.PuntRequests, err error) {
	p.Log.Debugf("Handling GetPuntRequests(%+v)", item)
	puntReqs = &puntmgr.PuntRequests{}
	p.swMod.Lock()
	defer p.swMod.Unlock()
	if !p.swMod.discovered {
		return puntReqs, errors.New("CNF has not been yet discovered, execute DiscoverCnf first")
	}

	// find PuntRequests callback corresponding to the model of the item
	model, err := models.GetModelForItem(item)
	if err != nil {
		return puntReqs, fmt.Errorf("failed to get model: %w", err)
	}
	var puntReqClb PuntRequestsClb
	for _, expModel := range p.swMod.models {
		if expModel.model.Name() == model.Name() {
			if expModel.callbacks != nil {
				puntReqClb = expModel.callbacks.PuntRequests
			}
			break
		}
	}
	if puntReqClb == nil {
		return puntReqs, fmt.Errorf("punt not requested for item %v: %w",
			item.GetId(), err)
	}
	value, err := models.UnmarshalItem(item)
	if err != nil {
		return puntReqs, fmt.Errorf("UnmarshalItem failed: %w", err)
	}
	puntReqs = puntReqClb(value)
	return puntReqs, nil
}

// GetItemDependencies is served by CNFRegistry of a SW-Module CNF and returns
// the set of dependencies of the given configuration item (apart from punt deps which are determined
// from punt requests).
func (p *Plugin) GetItemDependencies(ctx context.Context, item *generic.Item) (itemDeps *pb.GetDependenciesResp, err error) {
	p.Log.Debugf("Handling GetItemDependencies(%+v)", item)
	itemDeps = &pb.GetDependenciesResp{}
	p.swMod.Lock()
	defer p.swMod.Unlock()
	if !p.swMod.discovered {
		return itemDeps, errors.New("CNF has not been yet discovered, execute DiscoverCnf first")
	}

	// find ItemDependencies callback corresponding to the model of the item
	model, err := models.GetModelForItem(item)
	if err != nil {
		return itemDeps, fmt.Errorf("failed to get model: %w", err)
	}
	var itemDepsClb ItemDepsClb
	for _, expModel := range p.swMod.models {
		if expModel.model.Name() == model.Name() {
			if expModel.callbacks != nil {
				itemDepsClb = expModel.callbacks.ItemDependencies
			}
			break
		}
	}
	if itemDepsClb == nil {
		return itemDeps, fmt.Errorf("no dependencies required for item %v: %w",
			item.GetId(), err)
	}
	value, err := models.UnmarshalItem(item)
	if err != nil {
		return itemDeps, fmt.Errorf("UnmarshalItem failed: %w", err)
	}
	itemDeps.Dependencies = itemDepsClb(value)
	return itemDeps, nil
}

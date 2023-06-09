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
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"go.ligato.io/cn-infra/v2/datasync/kvdbsync/local"
	"go.ligato.io/cn-infra/v2/infra"
	"go.ligato.io/cn-infra/v2/rpc/grpc"
	"go.ligato.io/cn-infra/v2/servicelabel"
	"google.golang.org/protobuf/proto"

	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/linux/nsplugin"
	"go.ligato.io/vpp-agent/v3/plugins/orchestrator/contextdecorator"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin"
	linux_namespace "go.ligato.io/vpp-agent/v3/proto/ligato/linux/namespace"

	cnfreg_plugin "go.pantheon.tech/stonework/plugins/cnfreg"
	"go.pantheon.tech/stonework/proto/cnfreg"
	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// These constants specify label for Internal StoneWork configuration (that is configuration
// not configured by the user or SW-Modules).
const InternalConfigLabelKey = "io.ligato.from-client"
const InternalConfigLabelValue = "stonework"

// Punt icManager plugins allows for multiple ligato plugins and even distributed agents to request packet punting
// between VPP and the same or distinct Linux network namespace(s). Unless there is a conflict between punt requests,
// the manager will ensure that common configuration items are shared and properly updated (e.g. ABX rules, TAP
// connection, etc.). The manager supports different kinds of packet punting approaches for L2 or L3 source VPP
// interfaces, with memifs, TAPs or AF-UNIX sockets used to deliver packets to the Linux network stack / user-space
// application.
// The plugin can be used by:
//   - STANDALONE CNF (even for a single punt it is a good practise to use the plugin),
//   - StoneWork to orchestrate punt between the all-in-one VPP and every SW-Module,
//   - and by a SW-Module to learn the metadata about a created punt configuration.
type Plugin struct {
	pb.UnimplementedPuntManagerServer
	sync.Mutex

	Deps
	config *Config

	notifDescr   *puntNotifDescriptor
	puntHandlers map[pb.PuntRequest_PuntType]PuntHandler
	icManager    InterconnectManager
	punts        map[puntID]*punt
}

// Deps is a set of dependencies of the Punt Manager plugin
type Deps struct {
	infra.PluginDeps
	ServiceLabel servicelabel.ReaderAPI
	GRPCServer   grpc.Server
	CnfRegistry  cnfreg_plugin.CnfRegistryAPI
	IfPlugin     ifplugin.API
	NsPlugin     nsplugin.API
	CfgClient    client.GenericClient
	KVScheduler  kvs.KVScheduler
}

// Punt Manager API.
type PuntManagerAPI interface {
	// GetPuntMetadata returns metadata about configured packet punt between VPP and the CNF.
	// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
	GetPuntMetadata(cnfMsLabel, key, label string) *pb.PuntMetadata
	// GetAllCNFPunts returns metadata of all punts created for the given CNF.
	// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
	GetAllCNFPunts(cnfMsLabel string) []*pb.PuntMetadata
	// AddPunt is used by StoneWork or standalone CNF to configure punt between VPP and the CNF.
	// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
	AddPunt(cnfMsLabel, key string, punt *pb.PuntRequest) error
	// DelPunt is used by StoneWork or standalone CNF to un-configure punt between VPP and the CNF.
	// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
	DelPunt(cnfMsLabel, key string, label string) error
	// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
	GetPuntDependencies(cnfMsLabel string, punt *pb.PuntRequest) (deps []kvs.Dependency)
}

// API to obtain names of configuration items generated for punts.
// Deprecated: use Punt metadata that can be obtained using GetPuntMetadata().
type PuntManagerNamingAPI interface {
	// GetLinuxVrfName returns the name used for Linux VRF device corresponding to the given VPP VRF.
	// Method is "static" in the sense that it can be called anytime, regardless of the internal state of the plugin.
	// Deprecated: use Punt metadata that can be obtained using GetPuntMetadata().
	GetLinuxVrfName(vrf uint32) string
}

// InterconnectLink is one of the:
//   - AF-UNIX socket
//   - pair of interfaces (memif or TAP)
//
// and each type has type-specific parameters.
type InterconnectLink interface {
	isInterconnectLink()
	equivalent(InterconnectLink) bool
}

// AF-UNIX socket between VPP and a CNF.
type AFUnixLink struct {
	socketPath string
}

func (*AFUnixLink) isInterconnectLink() {}

func (l *AFUnixLink) equivalent(link InterconnectLink) bool {
	l2, isAfUnixLink := link.(*AFUnixLink)
	if !isAfUnixLink {
		return false
	}
	return l.socketPath == l2.socketPath
}

// Interface-based interconnect (either with memif or TAP).
type InterfaceLink struct {
	// Name of the interface on the VPP side of the interconnect.
	// If empty then the interface name will be generated by PuntManager.
	interfaceName string
	// PhysAddress represents physical address (MAC) of the VPP side of the interconnect.
	// If empty then the MAC address will be generated by PuntManager.
	physAddress string
	// IPAddresses define list of IP addresses for the interface and must be
	// defined in the following format: <ipAddress>/<ipPrefix>.
	// Do not combine with unnumberedToIface or allocSubnet.
	ipAddresses []string
	// ID of VRF table that the interface is assigned to.
	vrf uint32
	// If enabled the VRF on the CNF/Linux side will not be created.
	withoutCNFVrf bool
	// Enable DHCP client on the VPP interface.
	withDhcpClient bool
	// Maximum transmission unit.
	mtu uint32
	// Reference to another VPP interface from which this interconnect will "borrow" the IP address.
	unnumberedToIface string
	// Enable if Punt Manager should allocate IP addresses for both ends of the interconnect.
	allocateSubnet bool
}

func (*InterfaceLink) isInterconnectLink() {}

func (l *InterfaceLink) equivalent(link InterconnectLink) bool {
	l2, isIfLink := link.(*InterfaceLink)
	if !isIfLink {
		return false
	}
	return l.interfaceName == l2.interfaceName &&
		l.physAddress == l2.physAddress &&
		l.vrf == l2.vrf &&
		l.withoutCNFVrf == l2.withoutCNFVrf &&
		l.withDhcpClient == l2.withDhcpClient &&
		l.mtu == l2.mtu &&
		l.unnumberedToIface == l2.unnumberedToIface &&
		l.allocateSubnet == l2.allocateSubnet &&
		isSubsetOf(l.ipAddresses, l2.ipAddresses) &&
		isSubsetOf(l2.ipAddresses, l.ipAddresses)
}

// Request to build a VPP<->CNF interconnect.
type InterconnectReq struct {
	link InterconnectLink
	// Select what/where to punt on the vpp side.
	// Different interconnects can still have the same vppSelector if cnfSelectors (determined by PuntManager) are different.
	// (i.e. vppSelector + cnfSelector = interconnect unique ID)
	vppSelector string //  SPAN src-interface + direction, ABX interface, H(x) interface, DHCP vrf ID
}

// PuntHandler should be implemented one for each punt type.
type PuntHandler interface {
	// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
	// for this punt request.
	GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq

	// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
	GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency)

	// CanMultiplex enables interconnection multiplexing for this punting. It could be enabled in certain cases:
	// 1. two or more punts of this type can coexist even if they have the same vpp selector
	// 2. one or more punts of this type can coexist with other type of punts on the same (TAP-only)
	// interconnection if they all have the same vpp selector and cnf selector.
	// The TAP-backed interconnection is shared for multiple multiplexing punts with the same cnf selector
	// (same network namespace) and vpp selector.
	CanMultiplex() bool

	// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
	ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
		interconnects []*pb.PuntMetadata_Interconnect, remove bool) error
}

// Unlike pb.PuntID this can be also used as a map key.
type puntID struct {
	cnfMsLabel string
	key        string
	label      string
}

func (id puntID) String() string {
	return fmt.Sprintf("%s|%s|%s", id.cnfMsLabel, id.key, id.label)
}

// punt groups request with metadata (=response).
type punt struct {
	state    pb.PuntState
	request  *pb.PuntRequest
	metadata *pb.PuntMetadata
}

// Init initializes internal attributes and in the case of STONEWORK_MODULE also starts gRPC server
// for RegisterCreatedPunt and UnregisterDeletedPunt methods.
func (p *Plugin) Init() (err error) {
	p.punts = make(map[puntID]*punt)
	p.puntHandlers = make(map[pb.PuntRequest_PuntType]PuntHandler)

	p.config, err = p.loadConfig()
	if err != nil {
		return err
	}

	cnfMode := p.CnfRegistry.GetCnfMode()
	if cnfMode == cnfreg.CnfMode_STONEWORK_MODULE {
		// serve UpdatePuntState
		grpcServer := p.GRPCServer.GetServer()
		if grpcServer == nil {
			return errors.New("gRPC server is not initialized")
		}
		pb.RegisterPuntManagerServer(grpcServer, p)
	}

	// register descriptor for punt notifications
	var kvDescr *kvs.KVDescriptor
	p.notifDescr, kvDescr = newPuntNotifDescriptor(p.KVScheduler, p.Log.NewLogger(NotifDescriptorName))
	err = p.KVScheduler.RegisterKVDescriptor(kvDescr)
	if err != nil {
		return err
	}

	// register punt handlers
	p.puntHandlers[pb.PuntRequest_HAIRPIN_XCONNECT] = NewHairpinXConnPuntHandler()
	p.puntHandlers[pb.PuntRequest_HAIRPIN] = NewHairpinPuntHandler()
	p.puntHandlers[pb.PuntRequest_SPAN] = NewSpanPuntHandler()
	p.puntHandlers[pb.PuntRequest_ABX] = NewAbxPuntHandler(p.IfPlugin)
	p.puntHandlers[pb.PuntRequest_PUNT_TO_SOCKET] = NewSocketPuntHandler()
	p.puntHandlers[pb.PuntRequest_DHCP_PROXY] = NewDhcpProxyPuntHandler()
	p.puntHandlers[pb.PuntRequest_ISISX] = NewIsisxPuntHandler()

	// prepare interconnect manager
	_, allocCidr, err := net.ParseCIDR(p.config.InterconnectAllocCIDR)
	if err != nil {
		return fmt.Errorf("failed to parse \"interconnect-alloc-cidr\": %w", err)
	}
	p.icManager = NewInterconnectManager(p.Log.NewLogger("icManager"), p.IfPlugin, p.ServiceLabel,
		p.NsPlugin, allocCidr)
	return nil
}

// Close is NOOP.
func (p *Plugin) Close() error {
	return nil
}

// GetPuntMetadata returns metadata about configured packet punting between VPP and the CNF.
// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
func (p *Plugin) GetPuntMetadata(cnfMsLabel, key, label string) *pb.PuntMetadata {
	p.Lock()
	defer p.Unlock()

	if cnfMsLabel == "" {
		cnfMsLabel = p.ServiceLabel.GetAgentLabel()
	}
	id := puntID{
		cnfMsLabel: cnfMsLabel,
		key:        key,
		label:      label,
	}
	punt, exists := p.punts[id]
	if !exists {
		return nil
	}
	return punt.metadata
}

// GetAllCNFPunts returns metadata of all punts created for the given CNF.
// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
func (p *Plugin) GetAllCNFPunts(cnfMsLabel string) (punts []*pb.PuntMetadata) {
	p.Lock()
	defer p.Unlock()

	if cnfMsLabel == "" {
		cnfMsLabel = p.ServiceLabel.GetAgentLabel()
	}
	for id, punt := range p.punts {
		if id.cnfMsLabel == cnfMsLabel {
			punts = append(punts, punt.metadata)
		}
	}
	return punts
}

// AddPunt is used by StoneWork or standalone CNF to configure punt between VPP and the CNF.
// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
func (p *Plugin) AddPunt(cnfMsLabel, key string, puntReq *pb.PuntRequest) error {
	p.Lock()
	defer p.Unlock()

	cnfMode := p.CnfRegistry.GetCnfMode()
	if cnfMode == cnfreg.CnfMode_STONEWORK_MODULE {
		panic(fmt.Errorf("method AddPunt is not available in the CNF mode %v", cnfMode))
	}
	if cnfMode == cnfreg.CnfMode_STANDALONE && puntReq.InterconnectType == pb.PuntRequest_MEMIF {
		return errors.New("it is not supported to punt with memif within a standalone CNF")
	}

	// check for duplicity
	if cnfMsLabel == "" {
		cnfMsLabel = p.ServiceLabel.GetAgentLabel()
	}
	id := puntID{
		cnfMsLabel: cnfMsLabel,
		key:        key,
		label:      puntReq.GetLabel(),
	}
	_, exists := p.punts[id]
	if exists {
		return fmt.Errorf("punt already exists: %v", id)
	}

	// prepare gRPC clients
	var (
		err             error
		remoteCfgClient client.GenericClient
		remoteCnfClient pb.PuntManagerClient
		remoteTxn       client.ChangeRequest
	)
	if cnfMode == cnfreg.CnfMode_STONEWORK {
		remoteCfgClient, err = p.CnfRegistry.GetCnfCfgClient(cnfMsLabel)
		if err != nil {
			p.Log.Error(err)
			return err
		}
		remoteTxn = remoteCfgClient.ChangeRequest()
		cnfConn, err := p.CnfRegistry.GetCnfGrpcConn(cnfMsLabel)
		if err != nil {
			p.Log.Error(err)
			return err
		}
		remoteCnfClient = pb.NewPuntManagerClient(cnfConn)
	}

	// obtain interconnect requirements from the punt handler
	puntHandler, hasHandler := p.puntHandlers[puntReq.GetPuntType()]
	if !hasHandler {
		return fmt.Errorf("punt type %v is not supported", puntReq.GetPuntType())
	}
	withMultiplex := puntHandler.CanMultiplex()
	icReqs := puntHandler.GetInterconnectReqs(puntReq)

	// try to create interconnects
	localTxn := newPuntChangeRequest(map[string]string{InternalConfigLabelKey: InternalConfigLabelValue})
	icType := puntReq.InterconnectType
	enableGso := puntReq.EnableGso
	interconnects, err := p.icManager.AddInterconnects(localTxn, remoteTxn, id, icReqs, icType, enableGso, withMultiplex)
	if err != nil {
		p.Log.Error(err)
		return err
	}

	// try to configure punt
	err = puntHandler.ConfigurePunt(localTxn, id, puntReq, interconnects, false)
	if err != nil {
		p.Log.Error(err)
		// cleanup (of the IC manager internal state)
		_ = p.icManager.DelInterconnects(localTxn, remoteTxn, id)
		return err
	}

	// store punt metadata
	puntMeta := &pb.PuntMetadata{
		Id: &pb.PuntID{
			CnfMsLabel: id.cnfMsLabel,
			Key:        id.key,
			Label:      id.label,
		},
		Interconnects: interconnects,
	}
	puntState := pb.PuntState_INIT
	p.punts[id] = &punt{
		state:    puntState,
		request:  puntReq,
		metadata: puntMeta,
	}

	// send announcement about created packet punting into CNF
	if cnfMode == cnfreg.CnfMode_STONEWORK {
		_, err = remoteCnfClient.UpdatePuntState(context.Background(),
			&pb.UpdatePuntStateReq{
				Metadata: puntMeta,
				State:    puntState,
			})
		if err != nil {
			// ignore any errors at this point
			p.Log.Error(err)
			err = nil
		}
	}

	// configure punt asynchronously
	go func() {
		// ignore errors - they may get fixed by retrying
		if err = localTxn.Send(context.Background()); err != nil {
			p.Log.Error(err)
		}
		if cnfMode == cnfreg.CnfMode_STONEWORK {
			if err = remoteTxn.Send(context.Background()); err != nil {
				p.Log.Error(err)
			}
		}

		puntState = pb.PuntState_CREATED
		if cnfMode == cnfreg.CnfMode_STONEWORK {
			_, err = remoteCnfClient.UpdatePuntState(context.Background(),
				&pb.UpdatePuntStateReq{
					Metadata: puntMeta,
					State:    puntState,
				})
			if err != nil {
				// ignore any errors at this point
				p.Log.Error(err)
				err = nil
			}
		}

		// publish notification about newly configured punt
		p.notifDescr.notify(id, false)

		p.Lock()
		defer p.Unlock()
		punt, exists := p.punts[id]
		if !exists {
			// highly unlikely
			p.Log.Warnf("punt removed before it was fully configured")
			return
		}
		punt.state = puntState
	}()
	return nil
}

// DelPunt is used by StoneWork or standalone CNF to un-configure punt between VPP and the CNF.
// If cnfMsLabel is empty then microservice label of this CNF is assumed (returned by ServiceLabel plugin).
func (p *Plugin) DelPunt(cnfMsLabel, key, label string) error {
	p.Lock()
	defer p.Unlock()

	cnfMode := p.CnfRegistry.GetCnfMode()
	if cnfMode == cnfreg.CnfMode_STONEWORK_MODULE {
		panic(fmt.Errorf("method DelPunt is not available in the CNF mode %v", cnfMode))
	}

	// check if the punt is created
	if cnfMsLabel == "" {
		cnfMsLabel = p.ServiceLabel.GetAgentLabel()
	}
	id := puntID{
		cnfMsLabel: cnfMsLabel,
		key:        key,
		label:      label,
	}
	punt, exists := p.punts[id]
	if !exists {
		return fmt.Errorf("unknown punt: %v", id)
	}
	puntReq := punt.request
	puntMeta := punt.metadata

	// prepare gRPC clients
	var (
		err             error
		remoteCfgClient client.GenericClient
		remoteCnfClient pb.PuntManagerClient
		remoteTxn       client.ChangeRequest
	)
	if cnfMode == cnfreg.CnfMode_STONEWORK {
		remoteCfgClient, err = p.CnfRegistry.GetCnfCfgClient(cnfMsLabel)
		if err != nil {
			p.Log.Error(err)
			return err
		}
		remoteTxn = remoteCfgClient.ChangeRequest()
		cnfConn, err := p.CnfRegistry.GetCnfGrpcConn(cnfMsLabel)
		if err != nil {
			p.Log.Error(err)
			return err
		}
		remoteCnfClient = pb.NewPuntManagerClient(cnfConn)
	}

	// try to remove punt
	puntHandler, hasHandler := p.puntHandlers[puntReq.GetPuntType()]
	if !hasHandler {
		return fmt.Errorf("punt type %v is not supported", puntReq.GetPuntType())
	}
	localTxn := p.CfgClient.ChangeRequest()
	err = puntHandler.ConfigurePunt(localTxn, id, puntReq, puntMeta.Interconnects, true)
	if err != nil {
		p.Log.Error(err)
		return err
	}

	// try to remove interconnects
	err = p.icManager.DelInterconnects(localTxn, remoteTxn, id)
	if err != nil {
		p.Log.Error(err)
		return err
	}

	// remove metadata from memory
	delete(p.punts, id)
	puntState := pb.PuntState_DELETED

	// send announcement about deleted packet punting into CNF
	if cnfMode == cnfreg.CnfMode_STONEWORK {
		_, err = remoteCnfClient.UpdatePuntState(context.Background(),
			&pb.UpdatePuntStateReq{
				Metadata: puntMeta,
				State:    puntState,
			})
		if err != nil {
			// ignore any errors at this point
			p.Log.Error(err)
			err = nil
		}
	}

	// publish notification already before the punt is removed
	go p.notifDescr.notify(id, true)

	// un-configure punt asynchronously
	go func() {
		// ignore errors - they may get fixed by retrying
		if cnfMode == cnfreg.CnfMode_STONEWORK {
			if err = remoteTxn.Send(context.Background()); err != nil {
				p.Log.Error(err)
			}
		}
		if err = localTxn.Send(context.Background()); err != nil {
			p.Log.Error(err)
		}
	}()
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *Plugin) GetPuntDependencies(cnfMsLabel string, punt *pb.PuntRequest) (deps []kvs.Dependency) {
	if cnfMsLabel != "" && cnfMsLabel != p.ServiceLabel.GetAgentLabel() {
		deps = append(deps, kvs.Dependency{
			Label: punt.GetLabel() + "-cnf-microservice",
			Key:   linux_namespace.MicroserviceKey(cnfMsLabel),
		})
	}
	if puntHandler, hasHandler := p.puntHandlers[punt.GetPuntType()]; hasHandler {
		deps = append(deps, puntHandler.GetPuntDependencies(punt)...)
	}
	return
}

// GetLinuxVrfName returns the name used for Linux VRF device corresponding to the given VPP VRF.
// Method is "static" in the sense that it can be called anytime, regardless of the internal state of the plugin.
func (p *Plugin) GetLinuxVrfName(vrf uint32) string {
	return p.icManager.GetLinuxVrfName(vrf)
}

// UpdatePuntState is called by Punt Manager of StoneWork to notify SW-Module about state change of a punt.
func (p *Plugin) UpdatePuntState(_ context.Context, req *pb.UpdatePuntStateReq) (resp *pb.UpdatePuntStateResp, err error) {
	p.Log.Debugf("Handling UpdatePuntState (%+v)", req)
	resp = &pb.UpdatePuntStateResp{}
	p.Lock()
	defer p.Unlock()

	id := puntIdFromProto(req.Metadata)
	switch req.State {
	case pb.PuntState_UNKNOWN:
		p.Log.Warn("Ignoring unknown punt state")
		return

	case pb.PuntState_INIT:
		_, exists := p.punts[id]
		if exists {
			err = fmt.Errorf("punt %s is already known", id.String())
			return resp, err
		}
		p.punts[id] = &punt{
			state: req.State,
			// request = nil
			metadata: req.Metadata,
		}

	case pb.PuntState_CREATED:
		punt, exists := p.punts[id]
		if !exists {
			err = fmt.Errorf("missing INIT state update for punt %s", id.String())
			return resp, err
		}
		if punt.state != pb.PuntState_INIT {
			p.Log.Warn("Ignoring punt state update (id=%s, state=%v, update=%v)",
				id, punt.state, req.State)
			return
		}
		punt.state = req.State
		p.notifDescr.notify(id, false)

	case pb.PuntState_DELETED:
		punt, exists := p.punts[id]
		if !exists {
			p.Log.Warn("Ignoring punt state update (id=%s, update=%v)",
				id, req.State)
			return
		}
		if punt.state == pb.PuntState_CREATED {
			p.notifDescr.notify(id, true)
		}
		delete(p.punts, id)
	}
	return resp, nil
}

func puntIdFromProto(puntMeta *pb.PuntMetadata) puntID {
	return puntID{
		cnfMsLabel: puntMeta.Id.GetCnfMsLabel(),
		key:        puntMeta.Id.GetKey(),
		label:      puntMeta.Id.GetLabel(),
	}
}

func isSubsetOf(slice1, slice2 []string) bool {
	for _, v1 := range slice1 {
		found := false
		for _, v2 := range slice2 {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type puntChangeRequest struct {
	txn    *client.LazyValTxn
	labels map[string]string
	err    error
}

func newPuntChangeRequest(labels map[string]string) *puntChangeRequest {
	return &puntChangeRequest{
		txn:    client.NewLazyValTxn(local.DefaultRegistry.PropagateChanges),
		labels: labels,
	}
}

func (r *puntChangeRequest) Update(items ...proto.Message) client.ChangeRequest {
	if r.err != nil {
		return r
	}
	for _, item := range items {
		key, err := models.GetKey(item)
		if err != nil {
			r.err = err
			return r
		}
		r.txn.Put(key, client.UpdateItem{Message: item, Labels: r.labels})
	}
	return r
}

func (r *puntChangeRequest) Delete(items ...proto.Message) client.ChangeRequest {
	if r.err != nil {
		return r
	}
	for _, item := range items {
		key, err := models.GetKey(item)
		if err != nil {
			r.err = err
			return r
		}
		r.txn.Delete(key)
	}
	return r
}

func (r *puntChangeRequest) Send(ctx context.Context) error {
	if r.err != nil {
		return r.err
	}
	_, withDataSrc := contextdecorator.DataSrcFromContext(ctx)
	if !withDataSrc {
		ctx = contextdecorator.DataSrcContext(ctx, "localclient")
	}
	return r.txn.Commit(ctx)
}

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
	"fmt"

	"go.ligato.io/cn-infra/v2/logging"
	"google.golang.org/protobuf/proto"

	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"

	pb "go.pantheon.tech/stonework/proto/cnfreg"
	"go.pantheon.tech/stonework/proto/puntmgr"
)

// The function initializes proxy descriptor for every config model exposed by the given CNF.
func (p *Plugin) initCnfProxy(swMod swModule) error {
	// create descriptor for each model
	for _, cnfModel := range swMod.cnfModels {
		spec := models.ToSpec(cnfModel.info.Spec)
		model, err := models.GetModel(spec.ModelName())
		if err != nil {
			p.Log.Errorf("failed to get model %s: %v", spec.ModelName(), err)
			return err
		}
		descrName := swMod.cnfMsLabel + "-" + cnfModel.info.ProtoName
		proxyDescr := &proxyDescriptor{
			log:        p.Log.NewLogger(descrName),
			puntMgr:    p.PuntMgr,
			cnfClient:  swMod.cnfClient,
			cfgClient:  swMod.cfgClient,
			cnfMsLabel: swMod.cnfMsLabel,
			model:      model,
			withPunt:   cnfModel.withPunt,
			withDeps:   cnfModel.withDeps,
			punts:      make(map[string]puntReqsForKey),
		}
		descr := &kvs.KVDescriptor{
			Name:               descrName,
			KeySelector:        model.IsKeyValid,
			ValueTypeName:      model.ProtoName(),
			KeyLabel:           model.StripKeyPrefix,
			NBKeyPrefix:        model.KeyPrefix(),
			Create:             proxyDescr.Create,
			Delete:             proxyDescr.Delete,
			Update:             proxyDescr.Update,
			UpdateWithRecreate: proxyDescr.UpdateWithRecreate,
			Dependencies:       proxyDescr.Dependencies,
		}
		if cnfModel.withRetrieve {
			descr.Retrieve = proxyDescr.Retrieve
		}
		err = p.KVScheduler.RegisterKVDescriptor(descr)
		if err != nil {
			p.Log.Errorf("failed to register proxy descriptor for model %s: %v", spec.ModelName(), err)
			return err
		}
	}
	return nil
}

// punt requests for a single key, identified in the map by labels
type puntReqsForKey map[string]*puntmgr.PuntRequest

// Descriptor that proxies CRUD operations to a remote CNF agent.
type proxyDescriptor struct {
	log        logging.Logger
	puntMgr    PuntManagerAPI
	cnfClient  pb.CnfDiscoveryClient
	cfgClient  client.GenericClient
	cnfMsLabel string
	model      models.KnownModel
	withPunt   bool
	withDeps   bool
	punts      map[string]puntReqsForKey // key -> label -> punt requests
}

// Create operation is proxied over the gRPC client.
func (p *proxyDescriptor) Create(key string, value proto.Message) (metadata kvs.Metadata, err error) {
	if p.withPunt {
		// establish packet punting before creating the configuration item
		puntReqs, err := p.getPuntReqs(value)
		if err != nil {
			p.log.Error(err)
			return nil, err
		}
		p.punts[key] = puntReqs
		for _, puntReq := range puntReqs {
			err = p.puntMgr.AddPunt(p.cnfMsLabel, key, puntReq)
			if err != nil {
				err = fmt.Errorf("AddPunt failed (%s|%s|%s): %w",
					p.cnfMsLabel, key, puntReq.Label, err)
				p.log.Error(err)
				return nil, err
			}
		}
	}
	ctx := context.Background()
	err = p.cfgClient.ChangeRequest().Update(value).Send(ctx)
	return nil, err
}

func (p *proxyDescriptor) getPuntReqs(value proto.Message) (puntReqs puntReqsForKey, err error) {
	puntReqs = make(puntReqsForKey)
	item, err := models.MarshalItem(value)
	if err != nil {
		err = fmt.Errorf("failed to marshal proto message into Item: %w", err)
		p.log.Error(err)
		return nil, err
	}
	ctx := context.Background()
	reqs, err := p.cnfClient.GetPuntRequests(ctx, item)
	if err != nil {
		err = fmt.Errorf("GetPuntRequests failed: %w", err)
		p.log.Error(err)
		return nil, err
	}
	for _, puntReq := range reqs.PuntRequests {
		puntReqs[puntReq.Label] = puntReq
	}
	return puntReqs, nil
}

func (p *proxyDescriptor) getDependencies(value proto.Message) (deps []*pb.ConfigItemDependency, err error) {
	item, err := models.MarshalItem(value)
	if err != nil {
		err = fmt.Errorf("failed to marshal proto message into Item: %w", err)
		p.log.Error(err)
		return nil, err
	}
	ctx := context.Background()
	resp, err := p.cnfClient.GetItemDependencies(ctx, item)
	if err != nil {
		err = fmt.Errorf("GetItemDependencies failed: %w", err)
		p.log.Error(err)
		return nil, err
	}
	deps = resp.Dependencies
	return deps, nil
}

// Delete operation is proxied over the gRPC client.
func (p *proxyDescriptor) Delete(key string, value proto.Message, metadata kvs.Metadata) (err error) {
	err = p.cfgClient.ChangeRequest().Delete(value).Send(context.Background())
	if err != nil {
		return err
	}
	if puntReqs := p.punts[key]; puntReqs != nil {
		for _, puntReq := range puntReqs {
			err = p.puntMgr.DelPunt(p.cnfMsLabel, key, puntReq.Label)
			if err != nil {
				err = fmt.Errorf("DelPunt failed (%s|%s|%s): %w",
					p.cnfMsLabel, key, puntReq.Label, err)
				p.log.Error(err)
				return err
			}
		}
	}
	delete(p.punts, key)
	return nil
}

// Update operation is proxied over the gRPC client.
func (p *proxyDescriptor) Update(key string, oldValue, newValue proto.Message, oldMetadata kvs.Metadata) (
	newMetadata kvs.Metadata, err error) {

	var (
		newPuntReqs  puntReqsForKey
		addPR, delPR []*puntmgr.PuntRequest
	)
	// add new punt configuration
	if p.withPunt {
		newPuntReqs, err = p.getPuntReqs(newValue)
		if err != nil {
			return nil, err
		}
		prevPuntReqs := p.punts[key]
		addPR, delPR, _ = p.diffPuntReqs(prevPuntReqs, newPuntReqs)
		for _, puntReq := range addPR {
			err = p.puntMgr.AddPunt(p.cnfMsLabel, key, puntReq)
			if err != nil {
				err = fmt.Errorf("AddPunt failed (%s|%s|%s): %w",
					p.cnfMsLabel, key, puntReq.Label, err)
				p.log.Error(err)
				return nil, err
			}
		}
	}

	// update configuration over gRPC
	err = p.cfgClient.ChangeRequest().Update(newValue).Send(context.Background())

	// delete obsolete punt configuration
	if p.withPunt {
		for _, puntReq := range delPR {
			err = p.puntMgr.DelPunt(p.cnfMsLabel, key, puntReq.Label)
			if err != nil {
				err = fmt.Errorf("DelPunt failed (%s|%s|%s): %w",
					p.cnfMsLabel, key, puntReq.Label, err)
				p.log.Error(err)
				return nil, err
			}
		}
		p.punts[key] = newPuntReqs
	}
	return nil, err
}

// UpdateWithRecreate returns true if the punt configuration has changed.
func (p *proxyDescriptor) UpdateWithRecreate(key string, oldValue, newValue proto.Message, metadata kvs.Metadata) bool {
	if !p.withPunt {
		return false
	}
	newPuntReqs, err := p.getPuntReqs(newValue)
	if err != nil {
		return true
	}
	prevPuntReqs := p.punts[key]
	addPR, delPR, updatePR := p.diffPuntReqs(prevPuntReqs, newPuntReqs)
	if len(updatePR) != 0 {
		return true
	}
	// we can add prior to updating or delete after updating but cannot do both
	return len(addPR) != 0 && len(delPR) != 0
}

func (p *proxyDescriptor) diffPuntReqs(prev, new puntReqsForKey) (add, del, update []*puntmgr.PuntRequest) {
	if prev == nil && new == nil {
		return
	}
	if prev == nil {
		for _, newPuntReq := range new {
			add = append(add, newPuntReq)
		}
		return
	}
	if new == nil {
		for _, prevPuntReq := range prev {
			del = append(del, prevPuntReq)
		}
	}
	for _, newPuntReq := range new {
		prevPuntReq := prev[newPuntReq.Label]
		if !proto.Equal(newPuntReq, prevPuntReq) {
			if prevPuntReq == nil {
				add = append(add, newPuntReq)
			} else {
				update = append(update, newPuntReq)
			}
		}
	}
	for _, prevPuntReq := range prev {
		if new[prevPuntReq.Label] == nil {
			del = append(del, prevPuntReq)
		}
	}
	return
}

// Retrieve operation is proxied over the gRPC client.
func (p *proxyDescriptor) Retrieve(correlate []kvs.KVWithMetadata) (retrieved []kvs.KVWithMetadata, err error) {
	resp, err := p.cfgClient.DumpState()
	if err != nil {
		return nil, err
	}
	for _, si := range resp {
		model, err := models.GetModelForItem(si.Item)
		if err != nil || model.Name() != p.model.Name() {
			continue
		}
		key, err := models.GetKeyForItem(si.Item)
		if err != nil {
			p.log.Warnf("failed to get key for a dumped item: %w", key, err)
			continue
		}
		value, err := models.UnmarshalItem(si.Item)
		if err != nil {
			p.log.Warnf("failed to unmarshal dumped item %s: %w", key, err)
			continue
		}
		retrieved = append(retrieved, kvs.KVWithMetadata{
			Key:    key,
			Value:  value,
			Origin: kvs.FromNB,
		})
	}
	return retrieved, nil
}

// Dependencies return the list of VPP interfaces that are mentioned in the punt configuration
// + any additional dependencies requested through the model registration.
func (p *proxyDescriptor) Dependencies(key string, value proto.Message) (deps []kvs.Dependency) {
	if p.withPunt {
		puntReqs, err := p.getPuntReqs(value)
		if err != nil {
			return nil
		}
		for _, puntReq := range puntReqs {
			deps = append(deps, p.puntMgr.GetPuntDependencies(p.cnfMsLabel, puntReq)...)
		}
	}
	if p.withDeps {
		extraDeps, err := p.getDependencies(value)
		if err != nil {
			return nil
		}
		for _, dep := range extraDeps {
			switch v := dep.Dep.(type) {
			case *pb.ConfigItemDependency_Key_:
				deps = append(deps, kvs.Dependency{
					Label: dep.GetLabel(),
					Key:   v.Key,
				})
			case *pb.ConfigItemDependency_Anyof:
				deps = append(deps, kvs.Dependency{
					Label: dep.GetLabel(),
					AnyOf: kvs.AnyOfDependency{
						KeyPrefixes: v.Anyof.GetKeyPrefixes(),
					},
				})
			}
		}
	}
	return deps
}

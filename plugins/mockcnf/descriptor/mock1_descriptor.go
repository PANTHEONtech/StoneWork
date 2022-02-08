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

package descriptor

import (
	"context"
	"errors"

	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	"google.golang.org/protobuf/proto"

	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	linux_l3 "go.ligato.io/vpp-agent/v3/proto/ligato/linux/l3"
	vpp_acl "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"

	"go.pantheon.tech/stonework/plugins/cnfreg"
	"go.pantheon.tech/stonework/plugins/mockcnf/descriptor/adapter"
	puntmgr_plugin "go.pantheon.tech/stonework/plugins/puntmgr"
	cnfreg_proto "go.pantheon.tech/stonework/proto/cnfreg"
	"go.pantheon.tech/stonework/proto/mockcnf"
	"go.pantheon.tech/stonework/proto/puntmgr"
)

const (
	Mock1DescrName = "mockcnf1"
	PuntLabel      = "singleton" // only single punt per config item
	mockRouteDst   = "7.7.7.7/32"
)

type MockCnf1Descriptor struct {
	log     logging.Logger
	puntMgr puntmgr_plugin.PuntManagerAPI
	cnfReg  cnfreg.CnfAPI

	config *mockcnf.MockCnf1
}

func NewMockCnf1Descriptor(cnfReg cnfreg.CnfAPI, puntMgr puntmgr_plugin.PuntManagerAPI,
	log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &MockCnf1Descriptor{
		cnfReg:  cnfReg,
		puntMgr: puntMgr,
		log:     log.NewLogger(Mock1DescrName),
	}

	typedDescr := &adapter.MockCnf1Descriptor{
		Name:          Mock1DescrName,
		NBKeyPrefix:   mockcnf.ModelMockCnf.KeyPrefix(),
		ValueTypeName: mockcnf.ModelMockCnf.ProtoName(),
		KeySelector:   mockcnf.ModelMockCnf.IsKeyValid,
		KeyLabel:      mockcnf.ModelMockCnf.StripKeyPrefix,
		Create:        ctx.Create,
		Delete:        ctx.Delete,
		Retrieve:      ctx.Retrieve,
		Dependencies:  ctx.Dependencies,
	}
	return adapter.NewMockCnf1Descriptor(typedDescr)
}

func (d *MockCnf1Descriptor) Create(key string, config *mockcnf.MockCnf1) (metadata interface{}, err error) {
	if d.cnfReg.GetCnfMode() == cnfreg_proto.CnfMode_STANDALONE {
		for _, puntReq := range MockCnf1PuntReqs(config).GetPuntRequests() {
			d.puntMgr.AddPunt("", key, puntReq)
		}
	}

	puntMeta := d.puntMgr.GetPuntMetadata("", key, PuntLabel)
	if puntMeta == nil {
		return nil, errors.New("missing punt metadata")
	}
	if len(puntMeta.Interconnects) == 0 {
		return nil, errors.New("missing punt-interconnect metadata")
	}
	d.log.Debugf("Interconnect metadata: %+v", puntMeta)
	// MockCNF1 configures Linux route via the punted interface
	go func() {
		err = client.LocalClient.ChangeRequest().Update(
			&linux_l3.Route{
				OutgoingInterface: puntMeta.Interconnects[0].CnfInterface.Name,
				Scope:             linux_l3.Route_GLOBAL,
				DstNetwork:        mockRouteDst,
			}).Send(context.Background())
	}()
	d.config = config
	return
}

func (d *MockCnf1Descriptor) Delete(key string, config *mockcnf.MockCnf1, metadata interface{}) (err error) {
	puntMeta := d.puntMgr.GetPuntMetadata("", key, PuntLabel)
	if puntMeta == nil {
		return errors.New("missing punt metadata")
	}
	if len(puntMeta.Interconnects) == 0 {
		return errors.New("missing punt-interconnect metadata")
	}
	go func() {
		err = client.LocalClient.ChangeRequest().Delete(
			&linux_l3.Route{
				OutgoingInterface: puntMeta.Interconnects[0].CnfInterface.Name,
				Scope:             linux_l3.Route_GLOBAL,
				DstNetwork:        mockRouteDst,
			}).Send(context.Background())
	}()
	d.config = nil

	if d.cnfReg.GetCnfMode() == cnfreg_proto.CnfMode_STANDALONE {
		for _, puntReq := range MockCnf1PuntReqs(config).GetPuntRequests() {
			d.puntMgr.DelPunt("", key, puntReq.GetLabel())
		}
	}
	return
}

func (d *MockCnf1Descriptor) Retrieve(correlate []adapter.MockCnf1KVWithMetadata) (retrieved []adapter.MockCnf1KVWithMetadata, err error) {
	if d.config == nil {
		return nil, nil
	}
	return []adapter.MockCnf1KVWithMetadata{
		{
			Key:    models.Key(d.config),
			Value:  d.config,
			Origin: kvs.FromNB,
		},
	}, nil
}

func (d *MockCnf1Descriptor) Dependencies(key string, config *mockcnf.MockCnf1) (deps []kvs.Dependency) {
	if d.cnfReg.GetCnfMode() == cnfreg_proto.CnfMode_STANDALONE {
		for _, puntReq := range MockCnf1PuntReqs(config).GetPuntRequests() {
			deps = append(deps, d.puntMgr.GetPuntDependencies("", puntReq)...)
		}
	}
	return deps
}

func MockCnf1PuntReqs(configItem proto.Message) *puntmgr.PuntRequests {
	value, ok := configItem.(*mockcnf.MockCnf1)
	if !ok {
		return nil
	}
	return &puntmgr.PuntRequests{
		PuntRequests: []*puntmgr.PuntRequest{
			{
				Label:    PuntLabel,
				PuntType: puntmgr.PuntRequest_ABX,
				Config: &puntmgr.PuntRequest_Abx_{
					Abx: &puntmgr.PuntRequest_Abx{
						Vrf:          value.GetVrf(),
						VppInterface: value.GetVppInterface(),
						IngressAclRules: []*vpp_acl.ACL_Rule_IpRule{
							{
								Ip: &vpp_acl.ACL_Rule_IpRule_Ip{
									DestinationNetwork: "local",
									SourceNetwork:      "any",
									Protocol:           value.IpProtocol,
								},
							},
						},
						EgressAclRules: nil,
					},
				},
			},
		},
	}
}

func MockCnf1ItemDeps(configItem proto.Message) []*cnfreg_proto.ConfigItemDependency {
	value, ok := configItem.(*mockcnf.MockCnf1)
	if !ok {
		return nil
	}
	if value.GetWaitForInterface() != "" {
		return []*cnfreg_proto.ConfigItemDependency{
			{
				Label: "wait-for-interface",
				Dep: &cnfreg_proto.ConfigItemDependency_Key_{
					Key: vpp_interfaces.InterfaceKey(value.GetWaitForInterface()),
				},
			},
		}
	}
	return nil
}

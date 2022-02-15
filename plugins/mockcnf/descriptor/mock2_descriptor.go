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

	"go.ligato.io/vpp-agent/v3/client"
	linux_l3 "go.ligato.io/vpp-agent/v3/proto/ligato/linux/l3"
	"google.golang.org/protobuf/proto"

	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_acl "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"

	"go.pantheon.tech/stonework/plugins/cnfreg"
	"go.pantheon.tech/stonework/plugins/mockcnf/descriptor/adapter"
	puntmgr_plugin "go.pantheon.tech/stonework/plugins/puntmgr"
	cnfreg_proto "go.pantheon.tech/stonework/proto/cnfreg"
	"go.pantheon.tech/stonework/proto/mockcnf"
	"go.pantheon.tech/stonework/proto/puntmgr"
)

const (
	Mock2DescrName = "mockcnf2"

	mockIpAddr = "9.9.9.9"
	mockHwAddr = "02:02:02:02:02:02"
)

type MockCnf2Descriptor struct {
	log     logging.Logger
	puntMgr puntmgr_plugin.PuntManagerAPI
	cnfReg  cnfreg.CnfAPI

	config *mockcnf.MockCnf2
}

func NewMockCnf2Descriptor(cnfReg cnfreg.CnfAPI, puntMgr puntmgr_plugin.PuntManagerAPI,
	log logging.PluginLogger) *kvs.KVDescriptor {
	ctx := &MockCnf2Descriptor{
		cnfReg:  cnfReg,
		puntMgr: puntMgr,
		log:     log.NewLogger(Mock2DescrName),
	}

	typedDescr := &adapter.MockCnf2Descriptor{
		Name:          Mock2DescrName,
		NBKeyPrefix:   mockcnf.ModelMockCnf.KeyPrefix(),
		ValueTypeName: mockcnf.ModelMockCnf.ProtoName(),
		KeySelector:   mockcnf.ModelMockCnf.IsKeyValid,
		KeyLabel:      mockcnf.ModelMockCnf.StripKeyPrefix,
		Create:        ctx.Create,
		Delete:        ctx.Delete,
		Retrieve:      ctx.Retrieve,
		Dependencies:  ctx.Dependencies,
	}
	return adapter.NewMockCnf2Descriptor(typedDescr)
}

func (d *MockCnf2Descriptor) Create(key string, config *mockcnf.MockCnf2) (metadata interface{}, err error) {
	if d.cnfReg.GetCnfMode() == cnfreg_proto.CnfMode_STANDALONE {
		for _, puntReq := range MockCnf2PuntReqs(config).GetPuntRequests() {
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
	// MockCNF2 configures static Linux ARP entry for the punted interface
	go func() {
		err = client.LocalClient.ChangeRequest().Update(
			&linux_l3.ARPEntry{
				Interface: puntMeta.Interconnects[0].CnfInterface.Name,
				IpAddress: mockIpAddr,
				HwAddress: mockHwAddr,
			}).Send(context.Background())
	}()
	d.config = config
	return
}

func (d *MockCnf2Descriptor) Delete(key string, config *mockcnf.MockCnf2, metadata interface{}) (err error) {
	puntMeta := d.puntMgr.GetPuntMetadata("", key, PuntLabel)
	if puntMeta == nil {
		return errors.New("missing punt metadata")
	}
	if len(puntMeta.Interconnects) == 0 {
		return errors.New("missing punt-interconnect metadata")
	}
	go func() {
		err = client.LocalClient.ChangeRequest().Delete(
			&linux_l3.ARPEntry{
				Interface: puntMeta.Interconnects[0].CnfInterface.Name,
				IpAddress: mockIpAddr,
				HwAddress: mockHwAddr,
			}).Send(context.Background())
	}()
	d.config = nil

	if d.cnfReg.GetCnfMode() == cnfreg_proto.CnfMode_STANDALONE {
		for _, puntReq := range MockCnf2PuntReqs(config).GetPuntRequests() {
			d.puntMgr.DelPunt("", key, puntReq.GetLabel())
		}
	}
	return
}

func (d *MockCnf2Descriptor) Retrieve(correlate []adapter.MockCnf2KVWithMetadata) (retrieved []adapter.MockCnf2KVWithMetadata, err error) {
	if d.config == nil {
		return nil, nil
	}
	return []adapter.MockCnf2KVWithMetadata{
		{
			Key:    models.Key(d.config),
			Value:  d.config,
			Origin: kvs.FromNB,
		},
	}, nil
}

func (d *MockCnf2Descriptor) Dependencies(key string, config *mockcnf.MockCnf2) (deps []kvs.Dependency) {
	if d.cnfReg.GetCnfMode() == cnfreg_proto.CnfMode_STANDALONE {
		for _, puntReq := range MockCnf2PuntReqs(config).GetPuntRequests() {
			deps = append(deps, d.puntMgr.GetPuntDependencies("", puntReq)...)
		}
	}
	return deps
}

func MockCnf2PuntReqs(configItem proto.Message) *puntmgr.PuntRequests {
	value, ok := configItem.(*mockcnf.MockCnf2)
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
						VppInterface: value.GetVppInterface(),
						Vrf:          value.GetVrf(),
						IngressAclRules: []*vpp_acl.ACL_Rule_IpRule{
							{
								Ip: &vpp_acl.ACL_Rule_IpRule_Ip{
									DestinationNetwork: "8.8.8.0/24",
									SourceNetwork:      "any",
								},
								Tcp: &vpp_acl.ACL_Rule_IpRule_Tcp{
									DestinationPortRange: &vpp_acl.ACL_Rule_IpRule_PortRange{
										LowerPort: 22,
										UpperPort: 22,
									},
									SourcePortRange: &vpp_acl.ACL_Rule_IpRule_PortRange{
										LowerPort: 0,
										UpperPort: 65535,
									},
								},
							},
						},
						EgressAclRules: []*vpp_acl.ACL_Rule_IpRule{
							{
								Ip: &vpp_acl.ACL_Rule_IpRule_Ip{
									DestinationNetwork: "local",
									SourceNetwork:      "",
								},
								Udp: &vpp_acl.ACL_Rule_IpRule_Udp{
									DestinationPortRange: &vpp_acl.ACL_Rule_IpRule_PortRange{
										LowerPort: 161,
										UpperPort: 162,
									},
									SourcePortRange: &vpp_acl.ACL_Rule_IpRule_PortRange{
										LowerPort: 0,
										UpperPort: 65535,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

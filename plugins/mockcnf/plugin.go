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

//go:generate descriptor-adapter --descriptor-name MockCnf1 --value-type *mockcnf.MockCnf1 --import "go.pantheon.tech/stonework/proto/mockcnf" --output-dir "descriptor"
//go:generate descriptor-adapter --descriptor-name MockCnf2 --value-type *mockcnf.MockCnf2 --import "go.pantheon.tech/stonework/proto/mockcnf" --output-dir "descriptor"

package mockcnf

import (
	"go.ligato.io/cn-infra/v2/infra"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	cnfreg_plugin "go.pantheon.tech/stonework/plugins/cnfreg"
	"go.pantheon.tech/stonework/plugins/mockcnf/descriptor"
	puntmgr_plugin "go.pantheon.tech/stonework/plugins/puntmgr"
	"go.pantheon.tech/stonework/proto/mockcnf"
)

type Plugin struct {
	Deps
	PuntManager puntmgr_plugin.PuntManagerAPI
	CnfRegistry cnfreg_plugin.CnfAPI
}

// Deps lists dependencies.
type Deps struct {
	infra.PluginDeps
	KVScheduler kvs.KVScheduler
}

// Init registers descriptor.
func (p *Plugin) Init() (err error) {
	// different mock-CNF behaviours are prepared for testing
	switch mockcnf.MockCnfIndex() {
	case 1:
		mockCnfDescriptor := descriptor.NewMockCnf1Descriptor(p.CnfRegistry, p.PuntManager, p.Log)
		err = p.KVScheduler.RegisterKVDescriptor(
			mockCnfDescriptor,
		)
		if err != nil {
			return err
		}

		// register the model implemented by CNF
		err = p.CnfRegistry.RegisterCnfModel(mockcnf.ModelMockCnf1, mockCnfDescriptor,
			&cnfreg_plugin.CnfModelCallbacks{
				PuntRequests:     descriptor.MockCnf1PuntReqs,
				ItemDependencies: descriptor.MockCnf1ItemDeps,
			})
		if err != nil {
			return err
		}

	case 2:
		mockCnfDescriptor := descriptor.NewMockCnf2Descriptor(p.CnfRegistry, p.PuntManager, p.Log)
		err = p.KVScheduler.RegisterKVDescriptor(
			mockCnfDescriptor,
		)
		if err != nil {
			return err
		}

		// register the model implemented by CNF
		err = p.CnfRegistry.RegisterCnfModel(mockcnf.ModelMockCnf2, mockCnfDescriptor,
			&cnfreg_plugin.CnfModelCallbacks{PuntRequests: descriptor.MockCnf2PuntReqs})
		if err != nil {
			return err
		}
	}
	return err
}

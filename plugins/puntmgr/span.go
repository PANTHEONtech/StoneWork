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
	"go.ligato.io/vpp-agent/v3/client"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"

	pb "go.pantheon.tech/stonework/proto/puntmgr"
)

// spanPunt implements PuntHandler for PuntRequest_SPAN
type spanPunt struct{}

func NewSpanPuntHandler() PuntHandler {
	return &spanPunt{}
}

// SpanInterfaceSelector is used only by spanPunt because SPAN can be combined with other punt types without conflicts.
func SpanInterfaceSelector(ifaceName string) string {
	return "vpp/span/interface/" + ifaceName
}

// GetInterconnectReqs returns definitions of all interconnects which are required between VPP and CNF
// for this punt request.
func (p *spanPunt) GetInterconnectReqs(punt *pb.PuntRequest) []InterconnectReq {
	if span := punt.GetSpan(); span != nil {
		return []InterconnectReq{
			{
				link:        &InterfaceLink{},
				vppSelector: SpanInterfaceSelector(span.VppInterface),
			},
		}
	}
	return nil
}

// GetPuntDependencies returns dependencies that have to be satisfied before the punt can be added.
func (p *spanPunt) GetPuntDependencies(punt *pb.PuntRequest) (deps []kvs.Dependency) {
	// VPP interface
	if span := punt.GetSpan(); span != nil {
		deps = append(deps,
			kvs.Dependency{
				Label: punt.GetLabel() + "-span-" + span.GetVppInterface(),
				Key:   vpp_interfaces.InterfaceKey(span.GetVppInterface()),
			})
	}
	return deps
}

// CanMultiplex enables interconnection multiplexing for this punting. It could be enabled in certain cases:
// 1. two or more punts of this type can coexist even if they have the same vpp selector
// 2. one or more punts of this type can coexist with other type of punts on the same (TAP-only)
// interconnection if they all have the same vpp selector and cnf selector.
// The TAP-backed interconnection is shared for multiple multiplexing punts with the same cnf selector
// (same network namespace) and vpp selector.
// SPAN key = input + output interface
func (p *spanPunt) CanMultiplex() bool {
	return true
}

// ConfigurePunt prepares txn to (un)configures VPP-side of the punt.
func (p *spanPunt) ConfigurePunt(txn client.ChangeRequest, puntId puntID, puntReq *pb.PuntRequest,
	interconnects []*pb.PuntMetadata_Interconnect, remove bool) error {

	if interconnects[0].Shared {
		// !remove: SPAN already configured
		// remove: SPAN will be un-configured by the last punt request
		return nil
	}
	span := &vpp_interfaces.Span{
		InterfaceFrom: puntReq.GetSpan().VppInterface,
		InterfaceTo:   interconnects[0].VppInterface.Name,
		Direction:     vpp_interfaces.Span_BOTH,
	}
	if remove {
		txn.Delete(span)
	} else {
		txn.Update(span)
	}
	return nil
}

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

package vpp2106

import (
	"go.pantheon.tech/stonework/proto/abx"

	"go.pantheon.tech/stonework/plugins/abx/vppcalls"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/abx"
)

// placeholder for unknown names
const unknownName = "<unknown>"

// DumpABXPolicy retrieves VPP ABX configuration.
func (h *ABXVppHandler) DumpABXPolicy() ([]*vppcalls.ABXDetails, error) {
	// retrieve ABX interfaces
	attachedIfs, err := h.dumpABXInterfaces()
	if err != nil {
		return nil, err
	}

	// retrieve ABX policy
	abxPolicy, err := h.dumpABXPolicy()
	if err != nil {
		return nil, err
	}

	// merge attached interfaces data to policy
	for _, policy := range abxPolicy {
		ifData, ok := attachedIfs[policy.Meta.PolicyID]
		if ok {
			policy.ABX.AttachedInterfaces = ifData
		}
	}

	return abxPolicy, nil
}

func (h *ABXVppHandler) dumpABXInterfaces() (map[uint32][]*vpp_abx.ABX_AttachedInterface, error) {
	// ABX index <-> attached interfaces
	abxIfs := make(map[uint32][]*vpp_abx.ABX_AttachedInterface)

	req := &abx.AbxInterfaceAttachDump{}
	reqCtx := h.callsChannel.SendMultiRequest(req)

	for {
		reply := &abx.AbxInterfaceAttachDetails{}
		last, err := reqCtx.ReceiveReply(reply)
		if err != nil {
			return nil, err
		}
		if last {
			break
		}

		// interface name
		ifName, _, exists := h.ifIndexes.LookupBySwIfIndex(reply.Attach.RxSwIfIndex)
		if !exists {
			ifName = unknownName
		}

		// attached interface entry
		attached := &vpp_abx.ABX_AttachedInterface{
			InputInterface: ifName,
			Priority:       reply.Attach.Priority,
			//			IsIpv6:         uintToBool(reply.Attach.IsIPv6),
		}

		_, ok := abxIfs[reply.Attach.PolicyID]
		if !ok {
			abxIfs[reply.Attach.PolicyID] = []*vpp_abx.ABX_AttachedInterface{}
		}
		abxIfs[reply.Attach.PolicyID] = append(abxIfs[reply.Attach.PolicyID], attached)
	}

	return abxIfs, nil
}

func (h *ABXVppHandler) dumpABXPolicy() ([]*vppcalls.ABXDetails, error) {
	var abxs []*vppcalls.ABXDetails
	req := &abx.AbxPolicyDump{}
	reqCtx := h.callsChannel.SendMultiRequest(req)

	for {
		reply := &abx.AbxPolicyDetails{}
		last, err := reqCtx.ReceiveReply(reply)
		if err != nil {
			return nil, err
		}
		if last {
			break
		}

		// ACL name
		aclName, _, exists := h.aclIndexes.LookupByIndex(reply.Policy.ACLIndex)
		if !exists {
			aclName = unknownName
		}

		abxData := &vppcalls.ABXDetails{
			ABX: &vpp_abx.ABX{
				Index:   reply.Policy.PolicyID,
				AclName: aclName,
				// ForwardingPaths: fwdPaths,
			},
			Meta: &vppcalls.ABXMeta{
				PolicyID: reply.Policy.PolicyID,
			},
		}

		abxs = append(abxs, abxData)
	}

	return abxs, nil
}

func uintToBool(value uint8) bool {
	if value == 0 {
		return false
	}
	return true
}

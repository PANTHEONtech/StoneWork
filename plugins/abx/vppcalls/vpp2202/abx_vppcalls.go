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
	"fmt"
	"net"

	"github.com/go-errors/errors"

	"go.pantheon.tech/stonework/plugins/binapi/vpp2202/abx"
)

// GetAbxVersion retrieves version of the VPP ABX plugin
func (h *ABXVppHandler) GetAbxVersion() (ver string, err error) {
	req := &abx.AbxPluginGetVersion{}
	reply := &abx.AbxPluginGetVersionReply{}

	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d", reply.Major, reply.Minor), nil
}

// AddAbxPolicy creates new ABX entry
func (h *ABXVppHandler) AddAbxPolicy(policyID uint32, aclID uint32, tx_if string, dst_mac string) error {
	if err := h.abxAddDelPolicy(policyID, aclID, tx_if, dst_mac, true); err != nil {
		return errors.Errorf("failed to add ABX policy %d (ACL: %v): %v", policyID, aclID, err)
	}
	return nil
}

// DeleteAbxPolicy removes existing ABX entry
func (h *ABXVppHandler) DeleteAbxPolicy(policyID uint32) error {
	if err := h.abxAddDelPolicy(policyID, 0, "", "", false); err != nil {
		return errors.Errorf("failed to delete ABX policy %d: %v", policyID, err)
	}
	return nil
}

// AbxAttachInterface attaches interface to the ABF
func (h *ABXVppHandler) AbxAttachInterface(policyID, ifIdx, priority uint32) error {
	if err := h.abxAttachDetachInterface(policyID, ifIdx, priority, true); err != nil {
		return errors.Errorf("failed to attach interface %d to ABX policy %d: %v", ifIdx, policyID, err)
	}
	return nil
}

// AbxDetachInterface detaches interface from the ABF
func (h *ABXVppHandler) AbxDetachInterface(policyID, ifIdx, priority uint32) error {
	if err := h.abxAttachDetachInterface(policyID, ifIdx, priority, false); err != nil {
		return errors.Errorf("failed to detach interface %d from ABF policy %d: %v", ifIdx, policyID, err)
	}
	return nil
}

func (h *ABXVppHandler) abxAttachDetachInterface(policyID, ifIdx, priority uint32, isAttach bool) error {
	req := &abx.AbxInterfaceAttachDetach{
		IsAttach: boolToUint(isAttach),
		Attach: abx.AbxInterfaceAttach{
			PolicyID:    policyID,
			Priority:    priority,
			RxSwIfIndex: ifIdx,
		},
	}
	reply := &abx.AbxInterfaceAttachDetachReply{}

	return h.callsChannel.SendRequest(req).ReceiveReply(reply)
}

func (h *ABXVppHandler) abxAddDelPolicy(policyID, aclID uint32, txInterface string, dstMac string, isAdd bool) error {
	var txSwIfIndex uint32
	if isAdd {
		meta, found := h.ifIndexes.LookupByName(txInterface)
		if !found {
			return errors.Errorf("interface %s not found", txInterface)
		}
		txSwIfIndex = meta.SwIfIndex
	}

	req := &abx.AbxPolicyAddDel{
		IsAdd: boolToUint(isAdd),
		Policy: abx.AbxPolicy{
			PolicyID:    policyID,
			ACLIndex:    aclID,
			TxSwIfIndex: txSwIfIndex,
		},
	}
	reply := &abx.AbxPolicyAddDelReply{}

	if dstMac != "" {
		var err error
		macAddr, err := net.ParseMAC(dstMac)
		if err != nil {
			return err
		}
		copy(req.Policy.DstMac[:], macAddr)
	}
	return h.callsChannel.SendRequest(req).ReceiveReply(reply)
}

func boolToUint(input bool) uint8 {
	if input {
		return 1
	}
	return 0
}

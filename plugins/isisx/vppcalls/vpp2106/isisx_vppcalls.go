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

	"github.com/go-errors/errors"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/isisx"
)

// GetISISXVersion retrieves version of the VPP ISISX plugin
func (h *ISISXVppHandler) GetISISXVersion() (ver string, err error) {
	req := &isisx.IsisxPluginGetVersion{}
	reply := &isisx.IsisxPluginGetVersionReply{}

	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d", reply.Major, reply.Minor), nil
}

// AddISISXConnection creates new ISISX unidirectional cross-connection between 2 interfaces
func (h *ISISXVppHandler) AddISISXConnection(inputInterface, outputInterface string) error {
	if err := h.addDeleteISISXConnection(true, inputInterface, outputInterface); err != nil {
		return errors.Errorf("failed to add ISISX cross connection(%s -> %s) due to: %v",
			inputInterface, outputInterface, err)
	}
	return nil
}

// DeleteISISXConnection deletes existing ISISX unidirectional cross-connection between 2 interfaces
func (h *ISISXVppHandler) DeleteISISXConnection(inputInterface, outputInterface string) error {
	if err := h.addDeleteISISXConnection(false, inputInterface, outputInterface); err != nil {
		return errors.Errorf("failed to delete ISISX cross connection(%s -> %s) due to: %v",
			inputInterface, outputInterface, err)
	}
	return nil
}

func (h *ISISXVppHandler) addDeleteISISXConnection(isAdd bool, inputInterface, outputInterface string) error {
	// translate interface name to vpp interface indexes
	meta, found := h.ifIndexes.LookupByName(inputInterface)
	if !found {
		return errors.Errorf("input interface %s not found", inputInterface)
	}
	inputSwIfIndex := meta.SwIfIndex
	meta, found = h.ifIndexes.LookupByName(outputInterface)
	if !found {
		return errors.Errorf("output interface %s not found", outputInterface)
	}
	outputSwIfIndex := meta.SwIfIndex

	// construct request
	req := &isisx.IsisxConnectionAddDel{
		IsAdd: boolToUint(isAdd),
		Connection: isisx.IsisxConnection{
			RxSwIfIndex: inputSwIfIndex,
			TxSwIfIndex: outputSwIfIndex,
		},
	}
	reply := &isisx.IsisxConnectionAddDelReply{}

	// send, wait for and handle reply
	if err := h.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return err
	}
	if reply.Retval != 0 {
		return fmt.Errorf("vpp call %q returned: %d", reply.GetMessageName(), reply.Retval)
	}

	return nil
}

func boolToUint(input bool) uint8 {
	if input {
		return 1
	}
	return 0
}

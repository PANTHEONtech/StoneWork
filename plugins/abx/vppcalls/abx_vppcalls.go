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

package vppcalls

import (
	govppapi "git.fd.io/govpp.git/api"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/plugins/vpp"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/aclplugin/aclidx"
	"go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/ifaceidx"
	abx "go.pantheon.tech/stonework/proto/abx"
)

// ABXDetails contains proto-modeled ABX data together with VPP-related metadata
type ABXDetails struct {
	ABX  *abx.ABX `json:"abx"`
	Meta *ABXMeta `json:"abx_meta"`
}

// ABXMeta contains policy ID (ABX index)
type ABXMeta struct {
	PolicyID uint32 `json:"policy_id"`
}

// ABXVppAPI provides read/write methods required to handle VPP ACL-based forwarding
type ABXVppAPI interface {
	ABXVppRead

	// GetAbxVersion retrieves version of the VPP ABX plugin
	GetAbxVersion() (ver string, err error)
	// AddAbxPolicy creates new ABX entry together with a list of forwarding paths
	AddAbxPolicy(policyID uint32, aclID uint32, tx_if string, dst_mac string) error
	// DeleteAbxPolicy removes existing ABX entry
	DeleteAbxPolicy(policyID uint32) error
	// AbxAttachInterface attaches interface to the ABX
	AbxAttachInterface(policyID uint32, ifIdx, priority uint32) error
	// AbxDetachInterface detaches interface from the ABX
	AbxDetachInterface(policyID uint32, ifIdx, priority uint32) error
}

// ABXVppRead provides read methods for ABX plugin
type ABXVppRead interface {
	// DumpABXPolicy retrieves VPP ABX configuration.
	DumpABXPolicy() ([]*ABXDetails, error)
}

var handler = vpp.RegisterHandler(vpp.HandlerDesc{
	Name:       "abx",
	HandlerAPI: (*ABXVppAPI)(nil),
})

func AddAbxHandlerVersion(version vpp.Version, msgs []govppapi.Message,
	h func(ch govppapi.Channel, aclIdx aclidx.ACLMetadataIndex, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) ABXVppAPI,
) {
	handler.AddVersion(vpp.HandlerVersion{
		Version: version,
		Check: func(c vpp.Client) error {
			ch, err := c.NewAPIChannel()
			if err != nil {
				return err
			}
			return ch.CheckCompatiblity(msgs...)
		},
		NewHandler: func(c vpp.Client, a ...interface{}) vpp.HandlerAPI {
			ch, err := c.NewAPIChannel()
			if err != nil {
				return err
			}
			return h(ch, a[0].(aclidx.ACLMetadataIndex), a[1].(ifaceidx.IfaceMetadataIndex), a[2].(logging.Logger))
		},
	})
}

func CompatibleAbxVppHandler(c vpp.Client, aclIdx aclidx.ACLMetadataIndex, ifIdx ifaceidx.IfaceMetadataIndex,
	log logging.Logger) ABXVppAPI {
	if v := handler.FindCompatibleVersion(c); v != nil {
		return v.NewHandler(c, aclIdx, ifIdx, log).(ABXVppAPI)
	}
	return nil
}

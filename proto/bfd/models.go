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

package bfd

import "go.ligato.io/vpp-agent/v3/pkg/models"

const ModuleName = "vpp.bfd"

var (
	ModelBFD = models.Register(
		&BFD{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "bfd",
			Class:   "config",
		},
		models.WithNameTemplate(
			"{{.Interface}}/peer/{{.PeerIp}}",
		))
)

// BFDKey returns key for the given BFD configuration.
func BFDKey(ifName, peerIP string) string {
	return models.Key(&BFD{
		Interface: ifName,
		PeerIp:    peerIP,
	})
}

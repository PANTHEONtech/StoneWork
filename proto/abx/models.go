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

package vpp_abx

import (
	"strconv"
	"strings"

	"go.ligato.io/vpp-agent/v3/pkg/models"
)

// ModuleName is the name of the module used for models.
const ModuleName = "vpp.abx"

var (
	ModelABX = models.Register(&ABX{}, models.Spec{
		Module:  ModuleName,
		Version: "v1",
		Type:    "abx",
	}, models.WithNameTemplate("{{.Index}}"))
)

// Key returns the prefix used in the ETCD to store VPP ACL-based forwarding
// config of a particular ABX in selected vpp instance.
func Key(index uint32) string {
	return models.Key(&ABX{
		Index: index,
	})
}

const (
	// ABX to interface template is a derived value key
	abxToInterfaceTemplate = "vpp/abx/{abx}/interface/{iface}"
)

const (
	// InvalidKeyPart is used in key for parts which are invalid
	InvalidKeyPart = "<invalid>"
)

// ToABXInterfaceKey returns key for ABX-to-interface
func ToInterfaceKey(abx uint32, iface string) string {
	if iface == "" {
		iface = InvalidKeyPart
	}
	key := abxToInterfaceTemplate
	key = strings.Replace(key, "{abx}", strconv.Itoa(int(abx)), 1)
	key = strings.Replace(key, "{iface}", iface, 1)
	return key
}

// ParseABXToInterfaceKey parses ABX-to-interface key
func ParseToInterfaceKey(key string) (abx, iface string, isABXToInterface bool) {
	parts := strings.Split(key, "/")
	if len(parts) >= 5 &&
		parts[0] == "vpp" && parts[1] == "abx" && parts[3] == "interface" {
		abx = parts[2]
		iface = strings.Join(parts[4:], "/")
		if iface != "" && abx != "" {
			return abx, iface, true
		}
	}
	return "", "", false
}

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
	"fmt"
	"github.com/vishvananda/netns"
	"go.ligato.io/cn-infra/v2/servicelabel"

	"go.ligato.io/vpp-agent/v3/plugins/linux/nsplugin"
	linux_namespace "go.ligato.io/vpp-agent/v3/proto/ligato/linux/namespace"
)

// NetNsRegistry keeps track of all network namespaces used by CNFs.
type NetNsRegistry interface {
	// Get ID representing network namespace referenced by a given microservice label.
	// Network namespace used by multiple microservices will have the same ID regardless of which ms label
	// is used to query it.
	GetNetNsID(msLabel string) (int, error)

	// Each learned network namespace is used by one or more microservices. Label of one of these
	// microservices is designated to represent the namespace.
	GetNetNsLabel(id int) (msLabel string, err error)
}

type netNsRegistry struct {
	nsPlugin     nsplugin.API
	serviceLabel servicelabel.ReaderAPI
	nsByLabel    map[string]netNs // key = ms label
	nsById       map[int]string   // key = id, value = designated ms label
}

type netNs struct {
	id       int
	nsHandle netns.NsHandle
}

func NewNetNsRegistry(nsPlugin nsplugin.API, serviceLabel servicelabel.ReaderAPI) NetNsRegistry {
	return &netNsRegistry{
		nsPlugin:     nsPlugin,
		serviceLabel: serviceLabel,
		nsByLabel:    make(map[string]netNs),
		nsById:       make(map[int]string),
	}
}

// Get ID representing network namespace referenced by a given microservice label.
// Network namespace use by multiple microservices will have the same ID regardless of which ms label
// is used to query it.
func (r *netNsRegistry) GetNetNsID(msLabel string) (id int, err error) {
	if msLabel == "" || msLabel == r.serviceLabel.GetAgentLabel() {
		return 0, nil // = namespace of this CNF
	}
	if nsID, known := r.nsByLabel[msLabel]; known {
		return nsID.id, nil
	}
	nsHandle, err := r.nsPlugin.GetNamespaceHandle(nil,
		&linux_namespace.NetNamespace{
			Type:      linux_namespace.NetNamespace_MICROSERVICE,
			Reference: msLabel,
		})
	if err != nil {
		return 0, err
	}
	maxId := 0
	for _, ns2 := range r.nsByLabel {
		if ns2.nsHandle.Equal(nsHandle) {
			id = ns2.id
			break
		}
		if ns2.id > maxId {
			maxId = ns2.id
		}
	}
	if id == 0 {
		// new namespace
		id = maxId + 1
		r.nsById[id] = msLabel
	}
	r.nsByLabel[msLabel] = netNs{
		id:       id,
		nsHandle: nsHandle,
	}
	return
}

// Each learned network namespace is used by one or more microservices. Label of one of these
// microservices is designated to represent the namespace.
func (r *netNsRegistry) GetNetNsLabel(id int) (msLabel string, err error) {
	if id == 0 {
		return r.serviceLabel.GetAgentLabel(), nil
	}
	msLabel, found := r.nsById[id]
	if !found {
		return "", fmt.Errorf("network namespace with ID %d was not found", id)
	}
	return msLabel, nil
}

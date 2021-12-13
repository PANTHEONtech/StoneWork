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
	"go.ligato.io/cn-infra/v2/logging"
	"strings"

	"github.com/golang/protobuf/proto"
	prototypes "github.com/golang/protobuf/ptypes/empty"

	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
)

const (
	NotifDescriptorName = "punt-notification"
)

// NotificationKey returns key of the SB notification which is sent when the given punt
// is fully created (i.e. metadata are generated and configuration is applied).
func NotificationKey(cnfName, itemKey, puntLabel string) string {
	return fmt.Sprintf("punt/cnf-name/%s/item/%s/punt-label/%s",
		cnfName, itemKey, puntLabel)
}

// NotificationKeyPrefix return prefix of NotificationKey where punt-label
// and potentially some suffix from itemKey are trimmed.
func NotificationKeyPrefix(cnfName, itemKeyOrPrefix string) string {
	return fmt.Sprintf("punt/cnf-name/%s/item/%s",
		cnfName, itemKeyOrPrefix)
}

// puntNotifDescriptor describes punt notifications to KV Scheduler.
type puntNotifDescriptor struct {
	log         logging.Logger
	kvScheduler kvs.KVScheduler
}

func newPuntNotifDescriptor(kvScheduler kvs.KVScheduler, log logging.Logger) (
	*puntNotifDescriptor, *kvs.KVDescriptor) {
	descr := &puntNotifDescriptor{
		log:         log,
		kvScheduler: kvScheduler,
	}
	return descr, &kvs.KVDescriptor{
		Name:        NotifDescriptorName,
		KeySelector: descr.isPuntNotifKey,
	}
}

func (d *puntNotifDescriptor) isPuntNotifKey(key string) bool {
	return strings.HasPrefix(key, "punt/cnf-name/") &&
		strings.Contains(key, "/punt-label/")
}

// notify publishes notification to KVScheduler when punt is fully
// configured and retracts it when it is removed.
func (d *puntNotifDescriptor) notify(puntId puntID, removed bool) {
	var value proto.Message
	if !removed {
		// empty == created, nil == not created
		value = &prototypes.Empty{}
	}
	err := d.kvScheduler.PushSBNotification(kvs.KVWithMetadata{
		Key:   NotificationKey(puntId.cnfMsLabel, puntId.key, puntId.label),
		Value: value,
	})
	if err != nil {
		d.log.Warnf("failed to send notification to KVScheduler: %v", err)
	}
}

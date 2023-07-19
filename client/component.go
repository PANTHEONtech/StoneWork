// SPDX-License-Identifier: Apache-2.0

// Copyright 2023 PANTHEON.tech
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

package client

import (
	"context"
	"fmt"
	"strings"

	"go.ligato.io/vpp-agent/v3/cmd/agentctl/api/types"
	"go.ligato.io/vpp-agent/v3/cmd/agentctl/client"

	"go.ligato.io/vpp-agent/v3/proto/ligato/kvscheduler"
	"go.pantheon.tech/stonework/plugins/cnfreg"
	cnfregpb "go.pantheon.tech/stonework/proto/cnfreg"
)

type ComponentMode int32

const (
	ComponentUnknown ComponentMode = iota

	// Auxiliary means the component is not a CNF and is not managed by StoneWork
	ComponentAuxiliary

	// Standalone means the component is a standalone CNF
	ComponentStandalone

	// ComponentStonework means the component is a StoneWork module managed by StoneWork
	ComponentStoneworkModule

	// Stonework means the component is a StoneWork instance
	ComponentStonework
)

// Component is a component of StoneWork. It can be StoneWork instance itself,
// a CNF connected to it or other Ligato service in connected to StoneWork.
type Component interface {
	GetName() string
	GetMode() ComponentMode
	GetInfo() *cnfreg.Info
	GetMetadata() map[string]string
	ConfigStatus() (*ConfigCounts, error)
}

type component struct {
	agentclient *client.Client
	Name        string
	Mode        ComponentMode
	Info        *cnfreg.Info
	Metadata    map[string]string
}

func (c *component) GetName() string {
	return c.Name
}

func (c *component) GetMode() ComponentMode {
	return c.Mode
}

func (c *component) GetInfo() *cnfreg.Info {
	return c.Info
}

func (c *component) Client() *client.Client {
	return c.agentclient
}

func (c *component) GetMetadata() map[string]string {
	return c.Metadata
}

func (c *component) ConfigStatus() (*ConfigCounts, error) {
	if c.Mode == ComponentAuxiliary || c.Mode == ComponentUnknown {
		return nil, fmt.Errorf("cannot get scheduler values of component %s, this component in not managed by StoneWork", c.Name)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	values, err := c.agentclient.SchedulerValues(ctx, types.SchedulerValuesOptions{})
	if err != nil {
		return nil, err
	}

	var allVals []*kvscheduler.ValueStatus
	for _, baseVal := range values {
		allVals = append(allVals, baseVal.Value)
		allVals = append(allVals, baseVal.DerivedValues...)
	}

	var res ConfigCounts
	for _, val := range allVals {
		switch val.State {
		case kvscheduler.ValueState_INVALID, kvscheduler.ValueState_FAILED:
			res.Err++
		case kvscheduler.ValueState_MISSING:
			res.Missing++
		case kvscheduler.ValueState_PENDING:
			res.Pending++
		case kvscheduler.ValueState_RETRYING:
			res.Retrying++
		case kvscheduler.ValueState_UNIMPLEMENTED:
			res.Unimplemented++
		case kvscheduler.ValueState_CONFIGURED, kvscheduler.ValueState_DISCOVERED, kvscheduler.ValueState_OBTAINED, kvscheduler.ValueState_REMOVED, kvscheduler.ValueState_NONEXISTENT:
			res.Ok++
		}
	}

	return &res, nil
}

type ConfigCounts struct {
	Ok            int
	Err           int
	Missing       int
	Pending       int
	Retrying      int
	Unimplemented int
}

func (cc ConfigCounts) String() string {
	var fields []string
	if cc.Ok != 0 {
		fields = append(fields, fmt.Sprintf("%d OK", cc.Ok))
	}
	if cc.Err != 0 {
		errStr := fmt.Sprintf("%d errors", cc.Ok)
		if cc.Err == 1 {
			errStr = errStr[:len(errStr)-1]
		}
		fields = append(fields, errStr)
	}
	if cc.Missing != 0 {
		fields = append(fields, fmt.Sprintf("%d missing", cc.Missing))
	}
	if cc.Pending != 0 {
		fields = append(fields, fmt.Sprintf("%d pending", cc.Pending))
	}
	if cc.Retrying != 0 {
		fields = append(fields, fmt.Sprintf("%d retrying", cc.Retrying))
	}
	if cc.Unimplemented != 0 {
		fields = append(fields, fmt.Sprintf("%d unimplemented", cc.Unimplemented))
	}
	return strings.Join(fields, ", ")
}

func (c ComponentMode) String() string {
	switch c {
	case ComponentAuxiliary:
		return "auxiliary"
	case ComponentStandalone:
		return "standalone CNF"
	case ComponentStonework:
		return "StoneWork"
	case ComponentStoneworkModule:
		return "StoneWork module"
	default:
		return "unknown"
	}
}

func cnfModeToCompoMode(cm cnfregpb.CnfMode) ComponentMode {
	switch cm {
	case cnfregpb.CnfMode_STANDALONE:
		return ComponentStandalone
	case cnfregpb.CnfMode_STONEWORK_MODULE:
		return ComponentStoneworkModule
	case cnfregpb.CnfMode_STONEWORK:
		return ComponentStonework
	default:
		return ComponentUnknown
	}
}

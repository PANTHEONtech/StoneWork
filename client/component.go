package client

import (
	"context"
	"fmt"

	"go.ligato.io/vpp-agent/v3/cmd/agentctl/api/types"
	"go.ligato.io/vpp-agent/v3/cmd/agentctl/client"

	"go.ligato.io/vpp-agent/v3/proto/ligato/kvscheduler"
	"go.pantheon.tech/stonework/plugins/cnfreg"
	cnfregpb "go.pantheon.tech/stonework/proto/cnfreg"
)

type ComponentMode int32

const (
	// Foreign means the component is not managed by StoneWork
	ComponentForeign ComponentMode = iota

	// StoneworkModule means the component is a StoneWork module managed by StoneWork
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
	SchedulerValues() ([]*kvscheduler.BaseValueStatus, error)
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

func (c *component) SchedulerValues() ([]*kvscheduler.BaseValueStatus, error) {
	if c.Mode == ComponentForeign {
		return nil, fmt.Errorf("cannot get scheduler values of component %s, this component in not managed by StoneWork", c.Name)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	values, err := c.agentclient.SchedulerValues(ctx, types.SchedulerValuesOptions{})
	if err != nil {
		return nil, err
	}
	return values, nil
}

func cnfModeToCompoMode(cm cnfregpb.CnfMode) ComponentMode {
	switch cm {
	case cnfregpb.CnfMode_STANDALONE:
		return ComponentForeign
	case cnfregpb.CnfMode_STONEWORK_MODULE:
		return ComponentStoneworkModule
	case cnfregpb.CnfMode_STONEWORK:
		return ComponentStonework
	default:
		return ComponentForeign
	}
}
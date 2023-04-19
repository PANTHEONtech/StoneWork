package client

import "go.pantheon.tech/stonework/proto/cnfreg"

// Component is a component of StoneWork. It can be StoneWork instance itself,
// a CNF connected to it or other Ligato service in connected to StoneWork.
type Component interface {
	Name() string
	Mode() cnfreg.CnfMode
}

type component struct {
	name string
	mode cnfreg.CnfMode
}

func (c *component) Name() string {
	return c.name
}

func (c *component) Mode() cnfreg.CnfMode {
	return c.mode
}

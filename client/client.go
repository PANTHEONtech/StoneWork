package client

import (
	"fmt"
)

// Option is a function that customizes a Client.
type Option func(*Client)

// API defines client API. It is supposed to be used by various client
// applications, such as swctl or other user applications interacting with
// StoneWork.
type API interface {
	GetComponents() ([]Component, error)
}

// Client implements API interface.
type Client struct {
	components []Component
}

// NewClient creates a new client that implements API. The client can be
// customized by options.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// GetComponents returns list of components.
func (c *Client) GetComponents() ([]Component, error) {
	var components []Component

	// TODO: implement retrieval of components

	return components, fmt.Errorf("NOT IMPLEMENTED YET")

}

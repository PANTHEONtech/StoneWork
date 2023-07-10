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

package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	vppagent "go.ligato.io/vpp-agent/v3/cmd/agentctl/client"
	"go.ligato.io/vpp-agent/v3/cmd/agentctl/client/tlsconfig"

	"go.pantheon.tech/stonework/plugins/cnfreg"
)

const (
	DefaultHost               = "127.0.0.1"
	DefaultHTTPClientTimeout  = 60 * time.Second
	DefaultPortHTTP           = 9191
	DockerComposeServiceLabel = "com.docker.compose.service"
)

// Option is a function that customizes a Client.
type Option func(*Client) error

func WithHost(h string) Option {
	return func(c *Client) error {
		c.host = h
		return nil
	}
}

func WithHTTPPort(p uint16) Option {
	return func(c *Client) error {
		c.httpPort = p
		return nil
	}
}

func WithHTTPTLS(cert, key, ca string, skipVerify bool) Option {
	return func(c *Client) (err error) {
		c.httpTLS, err = withTLS(cert, key, ca, skipVerify)
		return err
	}
}

// API defines client API. It is supposed to be used by various client
// applications, such as swctl or other user applications interacting with
// StoneWork.
type API interface {
	GetComponents() ([]Component, error)
}

// Client implements API interface.
type Client struct {
	dockerClient      *docker.Client
	httpClient        *http.Client
	host              string
	scheme            string
	protocol          string
	httpPort          uint16
	httpTLS           *tls.Config
	customHTTPHeaders map[string]string
}

// NewClient creates a new client that implements API. The client can be
// customized by options.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		host:     DefaultHost,
		scheme:   "http",
		protocol: "tcp",
		httpPort: DefaultPortHTTP,
	}
	var err error
	c.dockerClient, err = docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		if err = o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) DockerClient() *docker.Client {
	return c.dockerClient
}

// HTTPClient returns configured HTTP client.
func (c *Client) HTTPClient() *http.Client {
	if c.httpClient == nil {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		tr.TLSClientConfig = c.httpTLS
		c.httpClient = &http.Client{
			Transport: tr,
			Timeout:   DefaultHTTPClientTimeout,
		}
	}
	return c.httpClient
}

func withTLS(cert, key, ca string, skipVerify bool) (*tls.Config, error) {
	var options []tlsconfig.Option
	if cert != "" && key != "" {
		options = append(options, tlsconfig.CertKey(cert, key))
	}
	if ca != "" {
		options = append(options, tlsconfig.CA(ca))
	}
	if skipVerify {
		options = append(options, tlsconfig.SkipServerVerification())
	}
	return tlsconfig.New(options...)
}

func (c *Client) StatusInfo(ctx context.Context) ([]cnfreg.Info, error) {
	resp, err := c.get(ctx, "/status/info", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	var infos []cnfreg.Info
	if err := json.NewDecoder(resp.body).Decode(&infos); err != nil {
		return nil, fmt.Errorf("decoding reply failed: %w", err)
	}
	return infos, nil
}

func (c *Client) GetComponents() ([]Component, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	infos, err := c.StatusInfo(ctx)
	if err != nil {
		return nil, err
	}

	dc := c.DockerClient()
	containerInfo, err := dc.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}

	var containers []*docker.Container
	for _, container := range containerInfo {
		c, err := dc.InspectContainerWithOptions(docker.InspectContainerOptions{ID: container.ID})
		if err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}

	cnfInfos := make(map[string]cnfreg.Info)
	for _, info := range infos {
		cnfInfos[info.MsLabel] = info
	}

	var components []Component
	var foreignContainers []*docker.Container
	for _, container := range containers {
		after, found := containsPrefix(container.Config.Env, "MICROSERVICE_LABEL=")
		if !found {
			foreignContainers = append(foreignContainers, container)
			continue
		}
		info, ok := cnfInfos[after]
		if !ok {
			foreignContainers = append(foreignContainers, container)
			continue
		}
		client, err := vppagent.NewClientWithOpts(vppagent.WithHost(info.IPAddr), vppagent.WithHTTPPort(info.HTTPPort))
		if err != nil {
			return components, err
		}
		compo := &component{
			agentclient: client,
			Name:        info.MsLabel,
			Mode:        cnfModeToCompoMode(info.CnfMode),
			Info:        &info,
		}
		components = append(components, compo)
	}

	for _, fcontainer := range foreignContainers {
		compo := &component{
			Name: fcontainer.Config.Labels[DockerComposeServiceLabel],
			Mode: ComponentForeign,
		}
		components = append(components, compo)
	}
	return components, nil
}

func containsPrefix(strs []string, prefix string) (string, bool) {
	for _, str := range strs {
		found := strings.HasPrefix(str, prefix)
		if found {
			return strings.TrimPrefix(str, prefix), found
		}
	}
	return "", false
}
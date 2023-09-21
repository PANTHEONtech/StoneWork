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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/consts"
	compose "github.com/docker/compose/v2/pkg/api"
	moby "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
	vppagent "go.ligato.io/vpp-agent/v3/cmd/agentctl/client"
	"go.ligato.io/vpp-agent/v3/cmd/agentctl/client/tlsconfig"

	"go.pantheon.tech/stonework/plugins/cnfreg"
)

const (
	DefaultHTTPClientTimeout = 60 * time.Second
	DefaultPortGRPC          = 9111
	DefaultPortHTTP          = 9191
	StoneWorkServiceName     = "stonework"
)

// Option is a function that customizes a Client.
type Option func(*Client) error

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

func WithComposeFiles(files []string) Option {
	return func(c *Client) error {
		if len(files) > 0 {
			c.composeFiles = files
		}
		return nil
	}
}

// API defines client API. It is supposed to be used by various client
// applications, such as swctl or other user applications interacting with
// StoneWork.
type API interface {
	GetComponents() ([]Component, error)
	GetHost() string
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
	deploymentName    string
	composeFiles      []string
}

// NewClient creates a new client that implements API. The client can be
// customized by options.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		scheme:   "http",
		protocol: "tcp",
		httpPort: DefaultPortHTTP,
	}
	var err error
	for _, o := range opts {
		if err = o(c); err != nil {
			return nil, err
		}
	}

	c.dockerClient, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	if c.deploymentName == "" {
		c.deploymentName, err = resolveDeploymentName(c.composeFiles)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve deployment name: %w", err)
		}
		logrus.Debugf("Deployment name resolved to: %s", c.deploymentName)
	}
	if c.host == "" {
		c.host, err = resolveHostAddr(c.dockerClient, c.deploymentName)
		if err != nil {
			logrus.Warnf("Failed to resolve host address: %v", err)
		} else {
			logrus.Debugf("StoneWork service management IP address resolved to: %s", c.host)
		}
	}
	return c, nil
}

func (c *Client) GetHost() string {
	return c.host
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

	deploymentLabel := filters.Arg("label", fmt.Sprintf("%s=%s", compose.ProjectLabel, c.deploymentName))
	configHashLabel := filters.Arg("label", compose.ConfigHashLabel)
	containerInfo, err := c.dockerClient.ContainerList(ctx, moby.ContainerListOptions{
		Filters: filters.NewArgs(deploymentLabel, configHashLabel),
	})
	if err != nil {
		return nil, err
	}

	var containers []moby.ContainerJSON
	for _, container := range containerInfo {
		c, err := c.dockerClient.ContainerInspect(ctx, container.ID)
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
	for _, container := range containers {
		metadata := make(map[string]string)
		metadata["containerID"] = container.ID
		metadata["containerName"] = container.Name
		metadata["containerServiceName"] = container.Config.Labels[compose.ServiceLabel]
		metadata["dockerImage"] = container.Config.Image
		if container.NetworkSettings.IPAddress != "" {
			metadata["containerIPAddress"] = container.NetworkSettings.IPAddress
		} else {
			for _, nw := range container.NetworkSettings.Networks {
				if nw.IPAddress != "" {
					metadata["containerIPAddress"] = nw.IPAddress
					break
				}
			}
		}

		logrus.Tracef("found metadata for container: %s, data: %+v", container.Name, metadata)

		// TODO: Refactor rest of this function (creation and determining of component type).
		// Rethink standalone CNF detection.
		compo := &component{Metadata: metadata}
		msLabel, found := containsPrefix(container.Config.Env, "MICROSERVICE_LABEL=")
		if !found {
			compo.Name = container.Config.Labels[compose.ServiceLabel]
			compo.Mode = ComponentAuxiliary
			components = append(components, compo)
			continue
		}
		info, ok := cnfInfos[msLabel]
		if !ok {
			client, err := vppagent.NewClientWithOpts(vppagent.WithHost(compo.Metadata["containerIPAddress"]), vppagent.WithHTTPPort(DefaultPortHTTP))
			if err != nil {
				return components, err
			}
			_, err = client.Status(context.Background())
			if err != nil {
				compo.Name = container.Config.Labels[compose.ServiceLabel]
				compo.Mode = ComponentAuxiliary
				components = append(components, compo)
				continue
			}
			compo.Name = container.Config.Labels[compose.ServiceLabel]
			compo.Mode = ComponentStandalone
			compo.agentclient = client
			components = append(components, compo)
			continue
		}
		compo.Name = info.MsLabel
		compo.Info = &info
		compo.Mode = cnfModeToCompoMode(info.CnfMode)

		client, err := vppagent.NewClientWithOpts(vppagent.WithHost(info.IPAddr), vppagent.WithHTTPPort(info.HTTPPort))
		if err != nil {
			return components, err
		}
		compo.agentclient = client
		components = append(components, compo)
	}
	return components, nil
}

// TODO: Docker Compose specific, when context for swctl is added refactor this
// Maybe add a cli flag for user to specify Docker Compose network? Currently this
// function uses default Docker Compose network in form `<DEPLOYMENT_NAME>_default`.
func resolveHostAddr(dc *docker.Client, deploymentName string) (string, error) {
	deployment := filters.Arg("label", fmt.Sprintf("%s=%s", compose.ProjectLabel, deploymentName))
	service := filters.Arg("label", fmt.Sprintf("%s=%s", compose.ServiceLabel, StoneWorkServiceName))
	configHash := filters.Arg("label", compose.ConfigHashLabel)
	containers, err := dc.ContainerList(context.Background(), moby.ContainerListOptions{
		Filters: filters.NewArgs(deployment, service, configHash),
		All:     true,
	})
	if err != nil {
		return "", err
	}
	if len(containers) > 1 {
		return "", fmt.Errorf("multiple StoneWork services found in deployment %s", deploymentName)
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("no StoneWork service found in deployment %s", deploymentName)
	}
	networkName := fmt.Sprintf("%s_default", deploymentName)
	network, ok := containers[0].NetworkSettings.Networks[networkName]
	if !ok {
		return "", fmt.Errorf("Docker Compose network %s not found", networkName)
	}
	return network.IPAddress, nil
}

// TODO: Docker Compose specific, when context for swctl is added refactor this.
// Maybe add a cli flag for user to specify the deployment name?
//
// https://docs.docker.com/compose/reference/#use--p-to-specify-a-project-name
func resolveDeploymentName(composeFiles []string) (string, error) {
	// from env var
	if name := os.Getenv(consts.ComposeProjectName); name != "" {
		return name, nil
	}
	if l := len(composeFiles); l > 0 {
		// from top level `name:` field in last specified compose file
		b, err := os.ReadFile(composeFiles[l-1])
		if err != nil {
			return "", err
		}
		n := struct {
			Name string `yaml:"name"`
		}{}
		if err := yaml.Unmarshal(b, &n); err != nil {
			return "", err
		}
		if n.Name != "" {
			return n.Name, nil
		}
		// from directory of first specified compose file
		return filepath.Base(filepath.Dir(composeFiles[0])), nil
	}
	// from current directory
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(currDir), nil
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

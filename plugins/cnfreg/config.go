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

package cnfreg

import "time"

const (
	// Start of the gRPC port range for StoneWork modules
	defaultSwModGrpcBasePort = 19000
	// Start of the HTTP port range for StoneWork modules
	defaultSwModHttpBasePort = 19100
)

// Config file for CnfRegistry plugin.
type Config struct {
	// Start of the gRPC port range for StoneWork modules
	SwModGrpcBasePort int `json:"sw-module-grpc-base-port"`
	// Start of the HTTP port range for StoneWork modules
	SwModHttpBasePort int `json:"sw-module-http-base-port"`
	// Deprecated. Using this option will result in a warning
	CnfDiscoveryTimeout time.Duration `json:"cnf-discovery-timeout"`
}

// loadConfig returns PuntMgr plugin file configuration if exists.
func (p *Plugin) loadConfig() (*Config, error) {
	cfg := &Config{
		SwModGrpcBasePort: defaultSwModGrpcBasePort,
		SwModHttpBasePort: defaultSwModHttpBasePort,
	}
	found, err := p.Cfg.LoadValue(cfg)
	if err != nil {
		return nil, err
	}
	if !found {
		p.Log.Debug("CnfRegistry config not found")
		return cfg, nil
	}

	p.Log.Debugf("CnfRegistry config found: %+v", cfg)
	var durZeroVal time.Duration
	if cfg.CnfDiscoveryTimeout != durZeroVal {
		p.Log.Warn("Option CnfDiscoveryTimeout is deprecated and setting it has no effect. " +
			"CNF discovery now runs continuously in the background by watching for filesystem changes (pid files created by CNFs). " +
			"The discovery starts immediately after StoneWork initialization without any timeout period.")
	}
	return cfg, err
}

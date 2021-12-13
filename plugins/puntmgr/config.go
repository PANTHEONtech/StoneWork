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

const (
	// CIDR used by default for allocations of /30 subnets for interconnects.
	defaultInterconnectAllocCIDR = "192.168.111.0/24"
)

// Config file for PuntMgr plugin.
type Config struct {
	// InterconnectAllocCIDR defines network from which /30 subnets are allocated for use by VPP<->CNF interconnects.
	InterconnectAllocCIDR string `json:"interconnect-alloc-cidr"`
}

// loadConfig returns PuntMgr plugin file configuration if exists.
func (p *Plugin) loadConfig() (*Config, error) {
	cfg := &Config{
		InterconnectAllocCIDR: defaultInterconnectAllocCIDR,
	}
	found, err := p.Cfg.LoadValue(cfg)
	if err != nil {
		return nil, err
	}
	if !found {
		p.Log.Debug("PuntMgr config not found")
		return cfg, nil
	}

	p.Log.Debugf("PuntMgr config found: %+v", cfg)
	return cfg, nil
}

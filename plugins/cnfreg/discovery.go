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

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vishvananda/netlink"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"go.ligato.io/vpp-agent/v3/client/remoteclient"
	pb "go.pantheon.tech/stonework/proto/cnfreg"
)

const (
	pidFileDir = "/run/stonework/discovery"
	pidFileExt = ".pid"
)

// PidFile written by SW-Module CNF.
type PidFile struct {
	Pid       int    `json:"pid"`
	MsLabel   string `json:"ms-label"`
	IpAddress string `json:"ip-address"`
	GrpcPort  int    `json:"grpc-port"`
	HttpPort  int    `json:"http-port"`
}

// cnfDiscovery discovers all CNFs that should be loaded into StoneWork (as StoneWork modules).
// Currently this is a very simple procedure - StoneWork will just wait configured amount of time and each CNF should
// in the meantime write pid file under the filepath /run/stonework/discovery/<cnf-name>.pid, containing pid, IP address,
// gRPC and http port numbers.
// StoneWork will then connect to each CNF over gRPC and discovers all CNF-specific configuration models.
// Downside of this trivial approach for discovery is that the Init phase will take longer and most of that time nothing
// actually happens (agent sleeps).
func (p *Plugin) cnfDiscovery() error {
	// Give CNFs some time to write pid files.
	time.Sleep(p.config.CnfDiscoveryTimeout)

	// Load all pid files written by SW-Module CNFs.
	p.sw.modules = make(map[string]swModule)
	pidFiles, err := ioutil.ReadDir(pidFileDir)
	if err != nil {
		if os.IsNotExist(err) {
			p.Log.Warnf("Directory \"%s\" does not exist, no CNFs will be loaded", pidFileDir)
			return nil
		}
		p.Log.Errorf("failed to read directory with PID files: %v", err)
		return err
	}
	for _, pidFile := range pidFiles {
		if !strings.HasSuffix(pidFile.Name(), pidFileExt) {
			continue
		}
		content, err := ioutil.ReadFile(path.Join(pidFileDir, pidFile.Name()))
		if err != nil {
			p.Log.Errorf("failed to read PID file %s: %v", pidFile.Name(), err)
			continue
		}
		var pf PidFile
		err = json.Unmarshal(content, &pf)
		if err != nil {
			p.Log.Errorf("failed to parse PID file %s: %v", pidFile.Name(), err)
			continue
		}
		swMod, err := p.getCnfModels(pf.IpAddress, pf.GrpcPort, pf.HttpPort)
		if err != nil {
			p.Log.Errorf("failed to obtain CNF models (pid file: %d): %v", pidFile.Name(), err)
			continue
		}
		p.sw.modules[swMod.cnfMsLabel] = swMod
	}
	return nil
}

func (p *Plugin) getCnfModels(ipAddress string, grpcPort, httpPort int) (swMod swModule, err error) {
	swMod.ipAddress = ipAddress
	swMod.grpcPort = grpcPort
	swMod.httpPort = httpPort

	// connect to the SW-Module CNF over gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	swMod.grpcConn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%d", ipAddress, grpcPort),
		grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return swMod, err
	}

	// call DiscoverCnf to learn the names of proto messages exposed by CNF
	swMod.cnfClient = pb.NewCnfDiscoveryClient(swMod.grpcConn)
	var swGrpcPort, swHttpPort int
	swGrpcPort = p.GRPCPlugin.GetPort()
	if p.HTTPPlugin != nil {
		swHttpPort = p.HTTPPlugin.GetPort()
	}
	ctx = context.Background()
	resp, err := swMod.cnfClient.DiscoverCnf(ctx, &pb.DiscoverCnfReq{
		SwIpAddress: p.ipAddress.String(),
		SwGrpcPort:  uint32(swGrpcPort),
		SwHttpPort:  uint32(swHttpPort),
	})
	if err != nil {
		return swMod, err
	}
	swMod.cnfMsLabel = resp.CnfMsLabel

	// call KnownModels to get meta information about models exposed by CNF
	swMod.cfgClient, err = remoteclient.NewClientGRPC(swMod.grpcConn,
		remoteclient.UseRemoteRegistry("config"))
	if err != nil {
		return swMod, err
	}
	models, err := swMod.cfgClient.KnownModels("config")
	if err != nil {
		return swMod, err
	}

	// for each exposed proto message find the corresponding model
	for _, cfgModel := range resp.ConfigModels {
		var found bool
		for _, model := range models {
			if model.ProtoName == cfgModel.ProtoName {
				swMod.cnfModels = append(swMod.cnfModels, cnfModel{
					info:         model,
					withPunt:     cfgModel.WithPunt,
					withDeps:     cfgModel.WithDeps,
					withRetrieve: cfgModel.WithRetrieve,
				})
				found = true
				break
			}
		}
		if !found {
			p.Log.Warnf("failed to find model info for proto message %s", cfgModel.ProtoName)
		}
	}
	return swMod, nil
}

// discoverMyIP tries to discover the StoneWork/CNF (non-local) (management) IP address.
func (p *Plugin) discoverMyIP() (net.IP, error) {
	links, err := netlink.LinkList()
	if err != nil {
		err = fmt.Errorf("discoverMyIP: LinkList failed: %w", err)
		return nil, err
	}
	var subnet *net.IPNet
	if p.MgmtSubnet != "" {
		_, subnet, err = net.ParseCIDR(p.MgmtSubnet)
		if err != nil {
			err = fmt.Errorf("discoverMyIP: Failed to parse management subnet: %w", err)
			return nil, err
		}
	}
	for _, link := range links {
		linkName := link.Attrs().Name
		if p.MgmtInterface != "" {
			if linkName != p.MgmtInterface {
				continue
			}
		}
		if subnet == nil &&
			(linkName == "lo" || linkName == "docker0" || strings.HasPrefix(linkName, "br-")) {
			// without any hint the local-only interfaces are skipped
			continue
		}
		addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
		if err != nil {
			p.Log.Warnf("discoverMyIP: AddrList failed: %v", err)
			continue
		}
		for _, addr := range addrs {
			if subnet != nil {
				if subnet.Contains(addr.IP) {
					return addr.IP, nil
				}
			} else {
				// without any hint take the first global unicast address
				if addr.IP.IsGlobalUnicast() {
					return addr.IP, nil
				}
			}
		}
	}
	err = errors.New("failed to discover management IP address")
	return nil, err
}

// SW-Module CNF writes pid file with pid, microservice label, IP address, gRPC and http port numbers under a known
// directory for StoneWork to discover it.
func (p *Plugin) writePidFile() error {
	content, err := json.MarshalIndent(PidFile{
		Pid:       os.Getpid(),
		MsLabel:   p.ServiceLabel.GetAgentLabel(),
		IpAddress: p.ipAddress.String(),
		GrpcPort:  p.GetGrpcPort(),
		HttpPort:  p.GetHttpPort(),
	}, "", "  ")
	if err != nil {
		return err
	}
	_ = os.Mkdir(pidFileDir, os.ModeDir)
	err = ioutil.WriteFile(
		path.Join(pidFileDir, p.ServiceLabel.GetAgentLabel()+pidFileExt),
		content, 0644)
	return err
}

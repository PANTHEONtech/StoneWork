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

package cnfreg

import (
	"net/http"
	"os"

	"github.com/unrolled/render"
	"go.ligato.io/cn-infra/v2/rpc/rest"
	"google.golang.org/grpc/connectivity"

	pb "go.pantheon.tech/stonework/proto/cnfreg"
)

type Info struct {
	PID           int
	MsLabel       string
	Mode          pb.CnfMode
	IPAddr        string
	GRPCPort      int
	HTTPPort      int
	GRPCConnState connectivity.State
}

func (p *Plugin) registerHandlers(handlers rest.HTTPHandlers) {
	if handlers == nil {
		p.Log.Debug("No http handler provided, skipping registration of REST handlers")
		return
	}
	if p.cnfMode == pb.CnfMode_STONEWORK {
		handlers.RegisterHTTPHandler("/status/info", p.statusHandler, http.MethodGet)
	}
}

func (p *Plugin) statusHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var infos []*Info

		swInfo := &Info{
			PID:      os.Getpid(),
			MsLabel:  p.ServiceLabel.GetAgentLabel(),
			Mode:     pb.CnfMode_STONEWORK,
			IPAddr:   p.ipAddress.String(),
			GRPCPort: p.GRPCPlugin.GetPort(),
			HTTPPort: p.HTTPPlugin.GetPort(),
		}
		infos = append(infos, swInfo)

		for kv := range p.sw.modules.Iter() {
			swMod := kv.Val
			swModInfo := &Info{
				PID:           swMod.pid,
				MsLabel:       swMod.cnfMsLabel,
				Mode:          pb.CnfMode_STONEWORK_MODULE,
				IPAddr:        swMod.ipAddress,
				GRPCPort:      swMod.grpcPort,
				HTTPPort:      swMod.httpPort,
				GRPCConnState: swMod.grpcConn.GetState(),
			}
			infos = append(infos, swModInfo)
		}
		if err := formatter.JSON(w, http.StatusOK, infos); err != nil {
			p.Log.Error(err)
		}
	}
}
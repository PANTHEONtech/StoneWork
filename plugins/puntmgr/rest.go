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

package puntmgr

// import (
// 	"net/http"

// 	"github.com/unrolled/render"
// 	"go.ligato.io/cn-infra/v2/rpc/rest"
// )

// func (p *Plugin) registerHandlers(handlers rest.HTTPHandlers) {
// 	if handlers == nil {
// 		p.Log.Debug("No http handler provided, skipping registration of REST handlers")
// 		return
// 	}
// 	if p.cnfMode == pb.CnfMode_STONEWORK {
// 		handlers.RegisterHTTPHandler("", )
// 	}
// }

// func (p *Plugin) 
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

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"go.ligato.io/cn-infra/v2/agent"
	"go.ligato.io/cn-infra/v2/logging"

	"go.pantheon.tech/stonework/cmd/stonework/app"

	_ "go.pantheon.tech/stonework/proto/mockcnf"
)

const logo = `
   ______               _      __         __
  / __/ /____  ___  ___| | /| / /__  ____/ /__
 _\ \/ __/ _ \/ _ \/ -_) |/ |/ / _ \/ __/  '_/
/___/\__/\___/_//_/\__/|__/|__/\___/_/ /_/\_\  %s
`

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	if _, err := fmt.Fprintf(os.Stdout, logo, agent.BuildVersion); err != nil {
		logging.DefaultLogger.Fatal(err)
	}

	swAgent := app.New()
	a := agent.NewAgent(agent.AllPlugins(swAgent), agent.StartTimeout(30*time.Second))
	if err := a.Run(); err != nil {
		logging.DefaultLogger.Fatal(err)
	}
}

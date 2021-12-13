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
	"os"

	"github.com/namsral/flag"

	"go.ligato.io/cn-infra/v2/agent"
	sv "go.ligato.io/cn-infra/v2/exec/supervisor"

	"go.pantheon.tech/stonework/pkg/printspec"

	// all configuration models exposed by StoneWork
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/govppmux"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/linux"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/netalloc"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/vpp"
	_ "go.pantheon.tech/stonework/proto/abx"
	_ "go.pantheon.tech/stonework/proto/bfd"
	_ "go.pantheon.tech/stonework/proto/isisx"
	_ "go.pantheon.tech/stonework/proto/nat64"
)

var printSpec = flag.CommandLine.Bool("print-spec", false,
	"only print spec of StoneWork models into stdout and exit")

func main() {
	a := agent.NewAgent(agent.AllPlugins(&sv.DefaultPlugin))
	if *printSpec {
		err := printspec.AllKnownModels(os.Stdout)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err := a.Run(); err != nil {
		panic(err)
	}
}

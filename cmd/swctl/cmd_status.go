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

package main

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"

	"go.pantheon.tech/stonework/client"
)

const statusExample = `
  <white># Show status for all components</>
  $ <yellow>swctl status</>

  <white># Show interface status of StoneWork VPP instance</>
  $ <yellow>swctl status --show-interfaces</>
`

type StatusOptions struct {
	Format         string
	ShowInterfaces bool
}

func (opts *StatusOptions) InstallFlags(flagset *pflag.FlagSet) {
	flagset.StringVar(&opts.Format, "format", "", "Format for the output (yaml, json, go template)")
	flagset.BoolVar(&opts.ShowInterfaces, "show-interfaces", false, "Show interface status of StoneWork VPP instance")
}

func NewStatusCmd(cli Cli) *cobra.Command {
	var opts StatusOptions
	cmd := &cobra.Command{
		Use:     "status [flags]",
		Short:   "Show status of StoneWork components",
		Example: color.Sprint(statusExample),
		Args:    cobra.ArbitraryArgs,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusCmd(cli, opts)
		},
	}
	opts.InstallFlags(cmd.PersistentFlags())
	return cmd
}

type statusInfo struct {
	client.Component
	ConfigCounts *client.ConfigCounts
}

func runStatusCmd(cli Cli, opts StatusOptions) error {
	resp, err := cli.Client().GetComponents()
	if err != nil {
		return err
	}

	if opts.ShowInterfaces {
		for _, compo := range resp {
			if sn, ok := compo.GetMetadata()["containerServiceName"]; ok {
				cmd := fmt.Sprintf("vpp-probe --env=%s --query label=%s=%s discover", defaultVppProbeEnv, client.DockerComposeServiceLabel, sn)
				formatArg := fmt.Sprintf("--format=%s", opts.Format)
				stdout, stderr, err := cli.Exec(cmd, []string{formatArg})
				if err != nil {
					if ee, ok := err.(*exec.ExitError); ok {
						logrus.Tracef("vpp-probe discover failed for service %s with error: %v: %s", sn, ee.String(), ee.Stderr)
						continue
					}
				}
				fmt.Fprintln(cli.Out(), stdout)
				fmt.Fprintln(cli.Err(), stderr)
			}
		}
		return nil
	}

	infos, err := getStatusInfo(resp)
	if err != nil {
		return err
	}

	if opts.Format == "" {
		printStatusTable(cli.Out(), infos, true)
	} else {
		if err := formatAsTemplate(cli.Out(), opts.Format, infos); err != nil {
			return err
		}
	}
	return nil
}

func getStatusInfo(components []client.Component) ([]statusInfo, error) {
	type infoWithErr struct {
		statusInfo
		error
	}
	var infos []statusInfo
	var wg sync.WaitGroup
	infoCh := make(chan infoWithErr)

	for _, compo := range components {
		wg.Add(1)
		go func(compo client.Component) {
			defer wg.Done()
			var counts *client.ConfigCounts
			var err error
			if compo.GetMode() != client.ComponentAuxiliary {
				counts, err = compo.ConfigStatus()
				if err != nil {
					infoCh <- infoWithErr{error: err}
				}
			}
			infoCh <- infoWithErr{
				statusInfo: statusInfo{
					Component:    compo,
					ConfigCounts: counts,
				},
			}
		}(compo)
	}

	go func() {
		wg.Wait()
		close(infoCh)
	}()

	for i := range infoCh {
		if i.error != nil {
			return nil, i.error
		}
		infos = append(infos, i.statusInfo)
	}
	slices.SortFunc(infos, cmpStatus)
	return infos, nil
}

func cmpStatus(a, b statusInfo) bool {
	greater := a.GetMode() > b.GetMode()
	if !greater && a.GetMode() == b.GetMode() {
		greater = a.GetName() > b.GetName()
	}
	return greater
}

func printStatusTable(out io.Writer, infos []statusInfo, useColors bool) {
	table := tablewriter.NewWriter(out)
	header := []string{
		"Name", "Mode", "IP Address", "GRPC Port", "HTTP Port", "Status", "Configuration",
	}
	aleft := tablewriter.ALIGN_LEFT
	acenter := tablewriter.ALIGN_CENTER
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{aleft, aleft, aleft, acenter, acenter, acenter, aleft})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	for _, info := range infos {
		row := []string{info.GetName(), info.GetMode().String()}
		var clrs []tablewriter.Colors
		if info.GetMode() == client.ComponentAuxiliary {
			clrs = []tablewriter.Colors{{}, {}}
			for i := range header[2:] {
				if useColors {
					clrs = append(clrs, []int{tablewriter.FgHiBlackColor})
				}
				row = append(row, strings.Repeat("-", len(header[i+2])))
			}
			table.Rich(row, clrs)
			continue
		}
		config := info.ConfigCounts.String()
		configColor := configColor(info.ConfigCounts)
		compoInfo := info.GetInfo()
		grpcState := compoInfo.GRPCConnState.String()
		var statusClr int
		// gRPC state does not make sense for StoneWork itself
		if info.GetMode() == client.ComponentStonework {
			grpcState = strings.Repeat("-", len("Status"))
			statusClr = tablewriter.FgHiBlackColor
		}
		row = append(row,
			compoInfo.IPAddr,
			strconv.Itoa(compoInfo.GRPCPort),
			strconv.Itoa(compoInfo.HTTPPort),
			grpcState,
			config)

		if useColors {
			clrs = []tablewriter.Colors{{}, {}, {}, {}, {}, {statusClr}, {configColor}}
		} else {
			clrs = []tablewriter.Colors{}
		}
		table.Rich(row, clrs)
	}
	table.Render()
}

func configColor(cc *client.ConfigCounts) int {
	if cc.Err > 0 {
		return tablewriter.FgHiRedColor
	}
	if cc.Retrying > 0 || cc.Pending > 0 {
		return tablewriter.FgYellowColor
	}
	if cc.Unimplemented > 0 {
		return tablewriter.FgMagentaColor
	}
	if cc.Missing > 0 {
		return tablewriter.FgHiYellowColor
	}
	if cc.Ok > 0 {
		return tablewriter.FgGreenColor
	}
	return 0
}

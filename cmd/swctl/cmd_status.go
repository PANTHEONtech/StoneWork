package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"

	"go.ligato.io/vpp-agent/v3/proto/ligato/kvscheduler"

	"go.pantheon.tech/stonework/client"
)

// TODO: improve status overview, show status of components (CNFs)
//   - instead of using raw output from vpp-probe, retrieve the important info
//     about the running/deployed components of StoneWork and show those by default
//   - optionally allow user to set more details which shows the more detailed output
//     similar to vpp-probe discover

const statusExample = `
  <white># Show status for all components</>
  $ <yellow>swctl status</>
`

type StatusCmdOptions struct {
	Args   []string
	Format string
}

func NewStatusCmd(cli Cli) *cobra.Command {
	var opts StatusCmdOptions
	cmd := &cobra.Command{
		Use:     "status [flags]",
		Short:   "Show status of StoneWork components",
		Example: color.Sprint(statusExample),
		Args:    cobra.ArbitraryArgs,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return runStatusCmd(cli, opts)
		},
	}
	return cmd
}

type statusInfo struct {
	client.Component
	ConfigCounts configCounts
}

func runStatusCmd(cli Cli, opts StatusCmdOptions) error {
	resp, err := cli.Client().GetComponents()
	if err != nil {
		return err
	}

	type infoWithErr struct {
		statusInfo
		error
	}
	var infos []statusInfo
	var wg sync.WaitGroup
	infoCh := make(chan infoWithErr)

	type fetched struct {
		values []*kvscheduler.BaseValueStatus
		err    error
	}

	for _, compo := range resp {
		wg.Add(1)
		go func(compo client.Component) {
			defer wg.Done()
			var counts configCounts
			if compo.GetMode() != client.ComponentForeign {
				vals, err := compo.SchedulerValues()
				if err != nil {
					infoCh <- infoWithErr{error: err}
				}
				counts = countConfig(vals)
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
			return i.error
		}
		infos = append(infos, i.statusInfo)
	}
	slices.SortFunc(infos, cmpInfos)
	// if err := formatAsTemplate(cli.Out(), "json", resp); err != nil {
	// 	return err
	// }
	printStatusTable(cli.Out(), infos)
	return nil
}

func countConfig(baseVals []*kvscheduler.BaseValueStatus) configCounts {
	var allVals []*kvscheduler.ValueStatus
	for _, baseVal := range baseVals {
		allVals = append(allVals, baseVal.Value)
		allVals = append(allVals, baseVal.DerivedValues...)
	}

	var res configCounts
	for _, val := range allVals {
		switch val.State {
		case kvscheduler.ValueState_INVALID, kvscheduler.ValueState_FAILED:
			res.Err++
		case kvscheduler.ValueState_MISSING:
			res.Missing++
		case kvscheduler.ValueState_PENDING:
			res.Pending++
		case kvscheduler.ValueState_RETRYING:
			res.Retrying++
		case kvscheduler.ValueState_UNIMPLEMENTED:
			res.Unimplemented++
		case kvscheduler.ValueState_CONFIGURED, kvscheduler.ValueState_DISCOVERED, kvscheduler.ValueState_OBTAINED, kvscheduler.ValueState_REMOVED, kvscheduler.ValueState_NONEXISTENT:
			res.Ok++
		}
	}
	return res
}

type comparable interface {
	constraints.Integer | ~string
}

func less[T comparable](a, b T) bool {
	if a > b {
		return true
	}
	return false
}

func cmpInfos(a, b statusInfo) bool {
	res := less(a.GetMode(), b.GetMode())
	bLessA := less(b.GetMode(), a.GetMode())
	if !(res || bLessA) {
		res = less(a.GetName(), b.GetName())
	}
	return res
}

func printStatusTable(out io.Writer, infos []statusInfo) {
	table := tablewriter.NewWriter(out)
	header := []string{
		"Name", "Mode", "IP Address", "GPRC Port", "HTTP Port", "Status", "Configuration",
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
		row := []string{info.GetName(), compoModeString(info.GetMode())}
		var clrs []tablewriter.Colors
		if info.GetMode() == client.ComponentForeign {
			clrs = []tablewriter.Colors{{}, {}}
			for i := range header[2:] {
				clrs = append(clrs, []int{tablewriter.FgHiBlackColor})
				row = append(row, strings.Repeat("-", len(header[i+2])))
			}
			table.Rich(row, clrs)
			continue
		}
		config := info.ConfigCounts.String()
		configColor := info.ConfigCounts.Color()
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
		clrs = []tablewriter.Colors{
			{}, {}, {}, {}, {}, {statusClr}, {configColor},
		}
		table.Rich(row, clrs)
	}
	table.Render()
}

func compoModeString(c client.ComponentMode) string {
	switch c {
	case client.ComponentForeign:
		return "foreign"
	case client.ComponentStonework:
		return "StoneWork"
	case client.ComponentStoneworkModule:
		return "StoneWork module"
	}
	return "unknown"
}

type configCounts struct {
	Ok            int
	Err           int
	Missing       int
	Pending       int
	Retrying      int
	Unimplemented int
}

func (c configCounts) String() string {
	var fields []string
	if c.Ok != 0 {
		fields = append(fields, fmt.Sprintf("%d OK", c.Ok))
	}
	if c.Err != 0 {
		errStr := fmt.Sprintf("%d errors", c.Ok)
		if c.Err == 1 {
			errStr = errStr[:len(errStr)-1]
		}
		fields = append(fields, errStr)
	}
	if c.Missing != 0 {
		fields = append(fields, fmt.Sprintf("%d missing", c.Missing))
	}
	if c.Pending != 0 {
		fields = append(fields, fmt.Sprintf("%d pending", c.Pending))
	}
	if c.Retrying != 0 {
		fields = append(fields, fmt.Sprintf("%d retrying", c.Retrying))
	}
	if c.Unimplemented != 0 {
		fields = append(fields, fmt.Sprintf("%d unimplemented", c.Unimplemented))
	}
	return strings.Join(fields, ", ")
}

func (c configCounts) Color() int {
	if c.Err > 0 {
		return tablewriter.FgHiRedColor
	}
	if c.Retrying > 0 || c.Pending > 0 {
		return tablewriter.FgYellowColor
	}
	if c.Unimplemented > 0 {
		return tablewriter.FgMagentaColor
	}
	if c.Missing > 0 {
		return tablewriter.FgHiYellowColor
	}
	if c.Ok > 0 {
		return tablewriter.FgGreenColor
	}
	return 0
}
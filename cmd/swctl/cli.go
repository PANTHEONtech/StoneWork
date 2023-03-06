package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/cli/cli/streams"
	"github.com/moby/term"
	"github.com/sirupsen/logrus"

	"go.pantheon.tech/stonework/client"
)

type Cli interface {
	Initialize(opts Options) error
	Client() client.API
	Exec(cmd string, args []string) (string, error)

	Out() *streams.Out
	Err() io.Writer
	In() *streams.In
	Apply(...CliOption) error
}

type CLI struct {
	client client.API

	vppProbePath string

	out *streams.Out
	err io.Writer
	in  *streams.In
}

func NewCli(opt ...CliOption) (*CLI, error) {
	cli := new(CLI)
	if err := cli.Apply(opt...); err != nil {
		return nil, err
	}
	if cli.out == nil || cli.in == nil || cli.err == nil {
		stdin, stdout, stderr := term.StdStreams()
		if cli.in == nil {
			cli.in = streams.NewIn(stdin)
		}
		if cli.out == nil {
			cli.out = streams.NewOut(stdout)
		}
		if cli.err == nil {
			cli.err = stderr
		}
	}
	return cli, nil
}

func (cli *CLI) Initialize(opts Options) (err error) {
	cli.client, err = initClient(opts)
	if err != nil {
		return fmt.Errorf("init error: %w", err)
	}

	if os.Getenv("SWCTL_VPP_PROBE_NO_DOWNLOAD") == "" {
		vppProbePath, err := getVppProbe()
		if err != nil {
			logrus.Errorf("getting vpp-probe failed: %v", err)
		} else {
			cli.vppProbePath = vppProbePath
		}
	} else {
		logrus.Debugf("vpp-probe downloading disabled by user")
	}

	return nil
}

func initClient(opts Options) (*client.Client, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (cli *CLI) Client() client.API {
	return cli.client
}

const programVppProbe = "vpp-probe"

func (cli *CLI) Exec(cmd string, args []string) (string, error) {
	if cli.Out().IsTerminal() {
		if strings.HasPrefix(cmd, programVppProbe) {
			cmd = programVppProbe + " --color=always" + strings.TrimPrefix(cmd, programVppProbe)
		}
	}

	// Use downloaded VPP probe
	if cli.vppProbePath != "" && strings.HasPrefix(cmd, programVppProbe) {
		cmd = fmt.Sprintf("%s %s", cli.vppProbePath, strings.TrimPrefix(cmd, programVppProbe))
	}

	return execCmd(cmd, args)
}

func (cli *CLI) Apply(opt ...CliOption) error {
	for _, o := range opt {
		if err := o(cli); err != nil {
			return err
		}
	}
	return nil
}

func (cli *CLI) Out() *streams.Out {
	return cli.out
}

func (cli *CLI) Err() io.Writer {
	return cli.err
}

func (cli *CLI) In() *streams.In {
	return cli.in
}

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

// Cli is a client API for CLI application.
type Cli interface {
	Initialize(opts Options) error
	Apply(...CliOption) error
	Client() client.API
	Entities() []Entity
	Exec(cmd string, args []string) (string, error)

	Out() *streams.Out
	Err() io.Writer
	In() *streams.In
}

// CLI implements Cli interface.
type CLI struct {
	client client.API

	entities     []Entity
	vppProbePath string

	out *streams.Out
	err io.Writer
	in  *streams.In
}

// NewCli returns a new CLI instance. It accepts CliOption for customization.
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

	// load entity files
	cli.entities, err = loadEntityFiles(opts.EntityFiles)
	if err != nil {
		return fmt.Errorf("loading entity files failed: %v", err)
	}

	// get vpp-probe
	vppProbePath, err := initVppProbe()
	if err != nil {
		logrus.Errorf("vpp-probe error: %v", err)
	} else {
		cli.vppProbePath = vppProbePath
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

func initVppProbe() (string, error) {
	if os.Getenv(EnvVarVppProbeNoDownload) != "" {
		logrus.Debugf("vpp-probe download disabled by user")
		return "", fmt.Errorf("downloading disabled by user")
	}

	vppProbePath, err := downloadVppProbe()
	if err != nil {
		return "", fmt.Errorf("downloading vpp-probe failed: %w", err)
	}

	return vppProbePath, nil
}

func (cli *CLI) Client() client.API {
	return cli.client
}

func (cli *CLI) Entities() []Entity {
	return cli.entities
}

const programVppProbe = "vpp-probe"

func (cli *CLI) Exec(cmd string, args []string) (string, error) {

	if strings.HasPrefix(cmd, programVppProbe) {
		if cli.Out().IsTerminal() {
			cmd = programVppProbe + " --color=always" + strings.TrimPrefix(cmd, programVppProbe)
		}

		// Use downloaded VPP probe
		if cli.vppProbePath != "" {
			cmd = fmt.Sprintf("%s %s", cli.vppProbePath, strings.TrimPrefix(cmd, programVppProbe))
		}
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

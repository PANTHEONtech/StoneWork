package app

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/docker/cli/cli/streams"
	"github.com/moby/term"
	"go.pantheon.tech/stonework/client"
)

// Cli is a client API for CLI application.
type Cli interface {
	Initialize(opts *GlobalOptions) error
	Apply(...CliOption) error
	Client() client.API
	Entities() []Entity
	GlobalOptions() *GlobalOptions
	Exec(cmd string, args []string, liveOutput bool) (stdout string, stderr string, err error)
	AppName() string

	Out() *streams.Out
	Err() io.Writer
	In() *streams.In
}

// CLI implements Cli interface.
type CLI struct {
	client client.API

	entities      []Entity
	globalOptions *GlobalOptions

	out *streams.Out
	err io.Writer
	in  *streams.In

	appName string
	// customizations is the generic way how to pass CLI customizations without extending the API. It should
	// be used for small modifications or changes that are not worthy to change the CLI API.
	customizations map[string]interface{}
}

// NewCli returns a new CLI instance. It accepts CliOption for customization.
func NewCli(appName string, opt ...CliOption) (*CLI, error) {
	cli := new(CLI)
	if err := cli.Apply(opt...); err != nil {
		return nil, err
	}
	cli.appName = appName
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

func (cli *CLI) Initialize(opts *GlobalOptions) (err error) {
	InitGlobalOptions(cli, opts)
	cli.globalOptions = opts

	cli.client, err = initClient(client.WithComposeFiles(cli.globalOptions.ComposeFiles))
	if err != nil {
		return fmt.Errorf("init client error: %w", err)
	}

	// load entity files
	cli.entities, err = loadEntityFiles(opts.EntityFiles)
	if err != nil {
		return fmt.Errorf("loading entity files failed: %w", err)
	}
	if cli.entities == nil {
		cli.entities, err = loadEmbeddedEntities(opts.EmbeddedEntityByte)
	}
	if err != nil {
		return fmt.Errorf("loading embedded entity files failed: %w", err)
	}

	return nil
}

func initClient(opts ...client.Option) (*client.Client, error) {
	c, err := client.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (cli *CLI) Client() client.API {
	return cli.client
}

func (cli *CLI) Entities() []Entity {
	return cli.entities
}

func (cli *CLI) GlobalOptions() *GlobalOptions {
	return cli.globalOptions
}

func (cli *CLI) Exec(cmd string, args []string, liveOutput bool) (string, string, error) {
	if cmd == "" {
		return "", "", errors.New("cannot execute empty command")
	}
	cmdParts := strings.Split(cmd, " ")
	if len(cmdParts) > 1 {
		args = append(cmdParts[1:], args...)
	}
	ecmd := newExternalCmd(externalExe(cmdParts[0]), args, cli)

	res, err := ecmd.exec(liveOutput)
	if err != nil {
		return "", "", err
	}
	return res.Stdout, res.Stderr, nil
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

func (cli *CLI) AppName() string {
	return cli.appName
}

package main

import (
	"io"

	"github.com/docker/cli/cli/streams"
	"github.com/moby/term"

	"go.pantheon.tech/stonework/client"
)

type CliOption func(cli *CLI) error

// WithStandardStreams sets a cli in, out and err streams with the standard streams.
func WithStandardStreams() CliOption {
	return func(cli *CLI) error {
		// Set terminal emulation based on platform as required.
		stdin, stdout, stderr := term.StdStreams()
		cli.in = streams.NewIn(stdin)
		cli.out = streams.NewOut(stdout)
		cli.err = stderr
		return nil
	}
}

// WithCombinedStreams uses the same stream for the output and error streams.
func WithCombinedStreams(combined io.Writer) CliOption {
	return func(cli *CLI) error {
		cli.out = streams.NewOut(combined)
		cli.err = combined
		return nil
	}
}

// WithInputStream sets a cli input stream.
func WithInputStream(in io.ReadCloser) CliOption {
	return func(cli *CLI) error {
		cli.in = streams.NewIn(in)
		return nil
	}
}

// WithOutputStream sets a cli output stream.
func WithOutputStream(out io.Writer) CliOption {
	return func(cli *CLI) error {
		cli.out = streams.NewOut(out)
		return nil
	}
}

// WithErrorStream sets a cli error stream.
func WithErrorStream(err io.Writer) CliOption {
	return func(cli *CLI) error {
		cli.err = err
		return nil
	}
}

// WithClient sets an APIClient.
func WithClient(c client.API) CliOption {
	return func(cli *CLI) error {
		cli.client = c
		return nil
	}
}

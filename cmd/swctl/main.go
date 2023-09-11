package main

import (
	"github.com/sirupsen/logrus"

	"go.pantheon.tech/stonework/cmd/swctl/app"
)

func main() {
	Execute()
}

// Execute executes a root command using default behavior
func Execute() {
	cli, err := app.NewCli("swctl")
	if err != nil {
		logrus.Fatalf("ERROR: %v", err)
	}

	root := app.NewRootCmd(cli)

	if err := root.Execute(); err != nil {
		logrus.Fatalf("ERROR: %v", err)
	}
}

package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	Execute()
}

// Execute executes a root command using default behavior
func Execute() {
	cli, err := NewCli()
	if err != nil {
		logrus.Fatalf("ERROR: %v", err)
	}

	root := NewRootCmd(cli)

	if err := root.Execute(); err != nil {
		logrus.Fatalf("ERROR: %v", err)
	}
}

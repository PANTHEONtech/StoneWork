package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

func execCmd(cmd string, args []string) (string, error) {
	cmdParts := strings.Split(cmd, " ")
	command := cmdParts[0]

	var commandArgs []string
	if len(cmdParts) > 1 {
		for _, part := range cmdParts[1:] {
			if strings.TrimSpace(part) == "" {
				continue
			}
			commandArgs = append(commandArgs, part)
		}
	}
	commandArgs = append(commandArgs, args...)

	c := exec.Command(command, commandArgs...)

	var stdout, stderr bytes.Buffer
	//c.Stdin = os.Stdin
	c.Stdout = &stdout
	c.Stderr = &stderr
	//env := os.Environ()
	//c.Env = env

	logrus.Tracef("[%s] %q", color.Gray.Sprint("EXEC"), c.String())

	t := time.Now()

	err := c.Run()
	out := strings.TrimRight(stdout.String(), "\n")
	took := time.Since(t).Seconds()
	l := logrus.WithField("took", fmt.Sprintf("%.3f", took))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			ee.Stderr = stderr.Bytes()
		}

		l.Tracef("[%s] %q (%v)\n%s\n", color.Red.Sprint("ERR"), c.Args, color.Red.Sprint(err), color.LightRed.Sprint(stderr.String()))
		return stdout.String(), err
	} else {
		l.Tracef("[%s] %q\n%s\n", color.Green.Sprint("OK"), c.Args, color.FgGray.Sprint(out))
	}

	return stdout.String(), nil
}

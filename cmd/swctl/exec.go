package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

type externalExe string

// known external programs that swctl can execute
const (
	cmdDocker   externalExe = "docker"
	cmdVppProbe externalExe = "vpp-probe"
	cmdAgentCtl externalExe = "agentctl"
)

type externalCmd struct {
	cli   *CLI
	exe   externalExe
	name  string
	args  []string
	env   []string
	color bool
}

func newExternalCmd(cmd externalExe, args []string, cli *CLI, swctlOpts GlobalOptions) *externalCmd {
	ecmd := &externalCmd{
		exe:   cmd,
		name:  string(cmd),
		cli:   cli,
		color: swctlOpts.Color != "never" && cli.Out().IsTerminal(),
	}
	ecmd.setDebugArg(swctlOpts.Debug)
	ecmd.setLogLevelArg(logrus.GetLevel())
	ecmd.setColorEnv()
	ecmd.setMiscFlags()
	ecmd.args = append(ecmd.args, args...)
	return ecmd
}

func (ec *externalCmd) prependArg(arg string, val string, aliases ...string) {
	if !hasAnyPrefix(ec.args, arg) && !hasAnyPrefix(ec.args, aliases...) {
		if val == "" {
			ec.args = append([]string{arg}, ec.args...)
		} else {
			ec.args = append([]string{fmt.Sprintf("%s=%s", arg, val)}, ec.args...)
		}
	}
}

func (ec *externalCmd) setDebugArg(debug bool) {
	switch ec.exe {
	case cmdAgentCtl, cmdDocker, cmdVppProbe:
		if debug {
			ec.prependArg("--debug", "", "-D")
		}
	}
}

func (ec *externalCmd) setLogLevelArg(loglvl logrus.Level) {
	if loglvl < logrus.FatalLevel {
		loglvl = logrus.FatalLevel
	}
	if loglvl > logrus.DebugLevel {
		loglvl = logrus.DebugLevel
	}
	switch ec.exe {
	case cmdAgentCtl, cmdDocker:
		ec.prependArg("--log-level", loglvl.String(), "-l")
	case cmdVppProbe:
		ec.prependArg("--loglevel", loglvl.String(), "-L")
	}
}

// https://no-color.org/
func (ec *externalCmd) setColorEnv() {
	if !ec.color {
		ec.env = append(ec.env, "NO_COLOR=1")
	}
}

func (ec *externalCmd) setMiscFlags() {
	switch ec.exe {
	case cmdVppProbe:
		if ec.color {
			ec.prependArg("--color", "always")
		}
		if ec.cli.vppProbePath != "" {
			ec.name = ec.cli.vppProbePath
		}
	case cmdAgentCtl:
		ec.prependArg("--host", ec.cli.client.GetHost(), "-H")
	}
}

type ExecResult struct {
	Took   time.Duration
	Status int
	Stdout string
	Stderr string
}

func (ec *externalCmd) exec() (*ExecResult, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(ec.name, ec.args...)
	cmd.Env = ec.env
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	now := time.Now()
	logrus.Tracef("[%s] %q", color.Gray.Sprint("EXEC"), cmd.String())
	err := cmd.Run()
	took := time.Since(now)
	l := logrus.WithField("took", took.Round(1*time.Millisecond).String())

	var execRes ExecResult
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			l.Tracef("[%s] %q (%v)\n%s\n", color.Red.Sprint("ERR"), cmd.Args, color.Red.Sprint(err), color.LightRed.Sprint(stderr.String()))
			execRes.Status = ee.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			return nil, err
		}
	}

	execRes.Took = took
	execRes.Stdout = strings.TrimRight(stdout.String(), "\n")
	execRes.Stderr = strings.TrimRight(stderr.String(), "\n")
	l.Tracef("[%s] %q\n%s\n", color.Green.Sprint("OK"), cmd.Args, color.FgGray.Sprint(execRes.Stdout))
	return &execRes, nil
}

func anyPrefixIndex(elems []string, prefixes ...string) int {
	for i, elem := range elems {
		for _, pref := range prefixes {
			if strings.HasPrefix(elem, pref) {
				return i
			}
		}
	}
	return -1
}

func hasAnyPrefix(elems []string, prefixes ...string) bool {
	return anyPrefixIndex(elems, prefixes...) >= 0
}

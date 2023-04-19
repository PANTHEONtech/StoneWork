package main

import (
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	EnvVarDebug              = "SWCTL_DEBUG"
	EnvVarLogLevel           = "SWCTL_LOGLEVEL"
	EnvVarVppProbeNoDownload = "SWCTL_VPP_PROBE_NO_DOWNLOAD"
)

var (
	glob GlobalOptions
)

type GlobalOptions struct {
	Debug    bool
	LogLevel string
	Color    string
	// TODO: support config file
	// Config string
}

func InitGlobalOptions(cli Cli, opts *GlobalOptions) {
	// color mode
	if opts.Color == "" && os.Getenv("NO_COLOR") != "" {
		// https://no-color.org/
		opts.Color = "never"
	}
	switch strings.ToLower(opts.Color) {
	case "auto", "":
		if !cli.Out().IsTerminal() {
			color.Disable()
		}
	case "on", "enabled", "always", "1", "true":
		color.Enable = true
	case "off", "disabled", "never", "0", "false":
		color.Disable()
	default:
		logrus.Fatalf("invalid color mode: %q", opts.Color)
	}

	// debug mode
	if os.Getenv(EnvVarDebug) != "" {
		opts.Debug = true
	}

	// log level
	if loglvl := os.Getenv(EnvVarLogLevel); loglvl != "" {
		opts.LogLevel = loglvl
	}
	if opts.LogLevel != "" {
		if lvl, err := logrus.ParseLevel(opts.LogLevel); err == nil {
			logrus.SetLevel(lvl)
			if lvl >= logrus.TraceLevel {
				logrus.SetReportCaller(true)
				//infralogrus.DefaultLogger().SetLevel(logging.LogLevel(lvl))
			}
			logrus.Tracef("log level set to %v", lvl)
		} else {
			logrus.Fatalf("log level invalid: %v", err)
		}
	} else if opts.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		//infralogrus.DefaultLogger().SetLevel(logging.ErrorLevel)
	}
}

func (glob *GlobalOptions) InstallFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&glob.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&glob.LogLevel, "log-level", "L", "", "Set logging level")
	flags.StringVar(&glob.Color, "color", "", "Color mode; auto/always/never")
}

type Options struct {
	ComposeFiles []string
	EntityFiles  []string
}

func (opts *Options) InstallFlags(flags *pflag.FlagSet) {
	flags.StringSliceVar(&opts.ComposeFiles, "composefile", nil, "Docker Compose configuration files")
	must(cobra.MarkFlagFilename(flags, "composefile", "yaml", "yml"))
	flags.StringSliceVar(&opts.EntityFiles, "entityfile", nil, "Entity configuration files")
	must(cobra.MarkFlagFilename(flags, "entityfile", "yaml", "yml"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

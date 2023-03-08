# SWCTL User Guide

## Installation

Ensure that Go is installed on your machine.

Run the following command to install swctl:

```bash
go install go.pantheon.tech/stonework/cmd/swctl@latest
```

## Usage

Run `swctl -h` to display the help menu:

```bash
$ swctl -h                                                                 

 ███████╗██╗    ██╗ ██████╗████████╗██╗     
 ██╔════╝██║    ██║██╔════╝╚══██╔══╝██║     
 ███████╗██║ █╗ ██║██║        ██║   ██║       stonework v22.10.0-5-gcb14aa9-dirty
 ╚════██║██║███╗██║██║        ██║   ██║       Wed Mar  8 13:02:38 CET 2023 (just now)
 ███████║╚███╔███╔╝╚██████╗   ██║   ███████╗  ondrej@devmachine (go1.20 linux/amd64)
 ╚══════╝ ╚══╝╚══╝  ╚═════╝   ╚═╝   ╚══════╝

Usage:
  swctl [command]

Available Commands:
  config      Manage config of StoneWork components
  deployment  Manage deployments of StoneWork
  help        Help about any command
  status      Show status of StoneWork components
  support     Export support data
  trace       Trace packets across data path

Flags:
  -f, --composefile strings   Docker Compose configuration files
  -D, --debug                 Enable debug mode
  -L, --loglevel string       Set logging level
      --color string          Color mode; auto/always/never
  -v, --version               version for swctl

Use "swctl [command] --help" for more information about a command.
```

This will display the basic usage of the swctl command and a list of available 
subcommands. To learn more about a specific subcommand, use the `--help` flag 
with the subcommand:

```bash
$ swctl deployment --help
```

## Options

You can pass flags to the swctl command to customize its behavior. Some of the 
most commonly used flags are:

```
    -D or --debug: Enables debug mode, which prints additional debugging information to the console.
    -L or --loglevel: Sets the logging level. Valid values are debug, info, warning, and error.
    --color: Sets the color mode of the output. Valid values are auto, always, and never.
```

## Commands

The available subcommands are:

* config: manage the configuration of StoneWork components
* deployment: manage deployments of StoneWork
* status: show the status of StoneWork components
* support: export support data
* trace: trace packets across data path

For example, to display the status of StoneWork components, run:

```bash
$ swctl status
```

To trace packets across the data path for troubleshooting network connectivity 
issues, run:

```bash
$ swctl trace
```
swctl also accepts one or more Docker Compose configuration files to specify the 
location of the Docker Compose configuration files using the `-f` or `--composefile` 
flag.

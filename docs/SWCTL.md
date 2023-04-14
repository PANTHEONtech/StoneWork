# SWCTL User Guide

## Installation

### Prerequisites

- Go 1.20+

Ensure that Go is installed on your machine.

### Install

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
 ███████║╚███╔███╔╝╚██████╗   ██║   ███████╗  user@machine (go1.20 linux/amd64)
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
```

This will display the basic usage of the swctl command and a list of available 
subcommands. To learn more about a specific subcommand, use `swctl [command] --help`.

### Flags

You can pass flags to the swctl command to customize its behavior. Some of the 
most commonly used flags are:

* `-D` or `--debug`: Enables debug mode, which prints additional debugging information to the console.
* `-L` or `--loglevel`: Sets the logging level. Valid values are debug, info, warning, and error.
* `--color`: Sets the color mode of the output. Valid values are auto, always, and never.

### Commands

The available subcommands for `swctl` are:

* `deployment`: manage deployments of StoneWork
* `manage`: perform complex config changes using customizable entities 
* `config`: manage the configuration of StoneWork components
* `status`: show the status of StoneWork components
* `trace`: trace packets across data path
* `support`: collect and export all relevant support info about StoneWork

#### Deployment

To manage the deployment of StoneWork and its components, run:

```bash
# Create and start StoneWork deployment
swctl deploy up

# Stop and remove StoneWork deployment
swctl deploy down
```

Additionally, to print various info relevant to the StoneWork deployment, run:

```bash
# Print StoneWork deployment information
swctl deploy config
swctl deploy info
swctl deploy images
swctl deploy services
```

> **Note**
> The `config` command is a simple wrapper for `docker compose` and expects `docker-compose.yaml` file.

#### Manage

For managing complex configuration changes, the manage command uses _entities_ loaded from entity file. The _entity_ is a special config template that uses parameters as input. The parameters use templating for their value to automatically render a value or let user override it. This allows for very quick 
config generation of any complexity.

By default, the entities are loaded from entity file - `entities.yaml` file in current working directory when running `swctl manage`. The expected format of the entity file is defined as:

```
---
entities:
  - name: ENTITY_NAME
    description: ENTITY_DESCRIPTION
    options:
      - name: OPTION_NAME
        value: OPTION_VALUE
    config: |
      ENTITY_CONFIG
# - name: ENTITY2_NAME
#   ...     
```

The `OPTION_VALUE` and `ENTITY_CONFIG` interpolate any references formatted as `${OPTION_NAME}` with values of options set eariler. After the interpolation, the Go `text/template` is used to render the value.

There are two options that are pre-defined for all entities, these are `ID`, starts from `1`, and `IDX` thaat starts from `0`. These options can be used to automatically increment or allocate values of other options.

The `ENTITY_CONFIG` uses the same format as StoneWork startup configuration file.

```bash
# List all available entities
swctl manage

# Print details about specific entity
swctl manage ENTITY
```

To generate config for a single entity, run:

```bash
# Generate entity config
swctl manage ENTITY add

# Generate entity config with an offset for IDX & ID
swctl manage ENTITY add --offset=100 
```

To generate combined config for multiple entites, run:

```bash
# Generate config for multiple entities
swctl manage ENTITY add --count=5
```

To merge an existing config file with generated config, run:

```bash
swctl manage ENTITY add --target=config.yaml
```

#### Config

To manage the configuration of StoneWork and its components, run:

```bash
# Get configuration
swctl config get

# Update configuration
swctl config update config.yaml

# Show history of configuration
swctl config history
```

> **Note**
> The `config` command is a simple wrapper for `agentctl config`.

#### Status

To display the status of StoneWork components and their interfaces, run:

```bash
swctl status
```

> **Note**
> The `config` command is a simple wrapper for `vpp-prove discover`.

#### Trace

To trace packets across the data path for troubleshooting network connectivity 
issues, run:

```bash
swctl trace
```

> **Note**
> The `config` command is a simple wrapper for `vpp-prove trace`.

#### Support

To export support file, run:

```bash
swctl support
```

> **Note**
> The `config` command is a simple wrapper for `agentctl report`.

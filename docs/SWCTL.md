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
 ███████╗██║ █╗ ██║██║        ██║   ██║       stonework v23.02
 ╚════██║██║███╗██║██║        ██║   ██║       Wed May  7 06:02:38 CET 2023
 ███████║╚███╔███╔╝╚██████╗   ██║   ███████╗  ondrej@beast (go1.20 linux/amd64)
 ╚══════╝ ╚══╝╚══╝  ╚═════╝   ╚═╝   ╚══════╝

Usage:
  swctl [command]

Available Commands:
  config      Manage config of StoneWork components
  deployment  Manage deployments of StoneWork
  help        Help about any command
  manage      Manage config changes with entities
  status      Show status of StoneWork components
  support     Export support data
  trace       Trace packets across data path

Flags:
  -f, --composefile strings   Docker Compose configuration files
      --entityfile strings    Entity configuration files
  -D, --debug                 Enable debug mode
  -L, --log-level string      Set logging level
      --color string          Set color mode (auto/always/never)
  -v, --version               Print swctl version
```

This will display the basic usage of the swctl command and a list of available 
subcommands. To learn more about a specific subcommand, use `swctl [command] --help`.

### Flags

You can pass flags to the swctl command to customize its behavior. Some of the 
most commonly used flags are:

* `-D` or `--debug`: Enables debug mode, which prints additional debugging information to the console.
* `-L` or `--log-level`: Sets the logging level. Valid values are debug, info, warning, and error.
* `--color`: Sets the color mode of the output. Valid values are auto, always, and never.

### Commands

The available subcommands for `swctl` are:

* `deployment`: manage deployments of StoneWork
* `manage`: perform config changes using customizable entities 
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

The manage command offers user-friendly and flexible way to manage configuration. Manage command helps users with building their initial config file, updating running config or performing various operations with config. It uses _entities_ loaded from _entity file_. The _entity_ is a special config template that uses variables as input. The parameters use templating for their value to automatically render a value or let user override it. This allows for very quick config generation of any complexity.

```go
// Entity is a blueprint for an object defined with a config template of related parts.
type Entity struct {
	Name        string      `json:"name"`
	Plural      string      `json:"plural"`
	Description string      `json:"description"`
	Vars        []EntityVar `json:"vars"`
	Config      string      `json:"config"`
	Single      bool        `json:"single"`
}

// EntityVar is a variable of an entity defined with a template to render its value.
type EntityVar struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"default"`
	Type        string `json:"type"`
	When        string `json:"when"`
}
```

By default, the entities are loaded from entity file - `entities.yaml` file in current working directory when running `swctl manage`. The expected format of the entity file is defined as:

```yaml
---
entities:

  - name: ENTITY_NAME
    description: ENTITY_DESCRIPTION
    vars:
      - name: VAR_NAME
        value: VAR_VALUE
    config: |
      ENTITY_CONFIG

# - name: ENTITY2_NAME
#   ...     
```

The `VAR_VALUE` and `ENTITY_CONFIG` interpolate any references formatted as `${VAR_NAME}` with values of variables set eariler. After the interpolation, the Go `text/template` is used to render the value.

There are two variables that are pre-defined for all entities (except those set as _single_). These variables are `ID` - starts from `1` and `IDX` - starts from `0`. These variables are used as automatic reference for other variables.

The templates used to render var values and config support special functions:

 - `add`: Takes two integers as arguments and returns their sum.
 - `inc`: Takes an integer as an argument and returns the integer incremented by 1.
 - `dec`: Takes an integer as an argument and returns the integer decremented by 1.
 - `previp`: Takes an IP address and a decrement integer as arguments, and returns the previous IP address by decrementing the provided IP address by the given integer. If an error occurs, it returns an error message.
 - `nextip`: Takes an IP address and an increment integer as arguments, and returns the next IP address by incrementing the provided IP address by the given integer. If an error occurs, it returns an error message.
 - subnet: Takes a CIDR (IP address with subnet mask) and an increment integer as arguments, and returns a new subnet based on the original subnet and the increment. If an error occurs, it returns an error message.
 - `trimsuffix`: Takes two strings as arguments, a main string and a suffix, and returns the main string with the specified suffix removed. If the suffix does not exist in the main string, the main string remains unchanged.
 - `trimprefix`: Takes two strings as arguments, a main string and a prefix, and returns the main string with the specified prefix removed. If the prefix does not exist in the main string, the main string remains unchanged.

Here are some examples of variable values:

- static: `10` - renders as `10` (if no override)
- interpolated: `abc-${ID}` - renders as `abc-1` for ID=1, `abc-2` for ID=2, etc.
- template: `{{ add ${ID} 100 }}` - renders `101` for ID=2, `102` for ID=2, etc.

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

# Generate entity config for specific instance ID
swctl manage ENTITY add --id=100 
```

To generate merged config for multiple entites, run:

```bash
# Generate config for multiple entities
swctl manage ENTITY add --count=5
```

To merge an existing config file with generated config, run:

```bash
swctl manage ENTITY add --target=config.yaml
```

To override value(s) of a specific entity variable, run:

```bash
swctl manage ENTITY add --var MY_VAR="my-value"
```

To set the variables using interactive mode, run: 

```bash
swctl manage ENTITY add --interactive
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

By default, `swctl config get` outputs only user defined configuration, the configuration that is created internally by StoneWork (data paths for punted packets from VPP to CNFs) is hidden.
Complete configuration can be shown by using the `--show-internal` flag: `swctl config get --show-internal`.

> **Note**
> The `config` command is a wrapper for `agentctl config`.

#### Status

To display the status of StoneWork components, run:

```bash
swctl status
```

Sample output:

```
NAME       MODE              IP ADDRESS  GPRC PORT  HTTP PORT  STATUS  CONFIGURATION  
stonework  StoneWork         172.20.0.2    9111       9191     ------  52 OK          
mockcnf2   StoneWork module  172.20.0.3    19002      19102    READY   17 OK          
mockcnf1   StoneWork module  172.20.0.3    19001      19101    READY   17 OK          
router-ns  auxiliary         ----------  ---------  ---------  ------  -------------  
```

The _MODE_ column can have values of _StoneWork_, _StoneWork module_, _standalone CNF_ or _auxiliary_.
Mode _auxiliary_ means the component is not a CNF and that it is not managed by a StoneWork service.  
The _STATUS_ column shows the state of GRPC connection between StoneWork service and the given StoneWork module.  
The _CONFIGURATION_ column shows short overview of counts and state of configuration items corresponding to the given component.

> **Note**
> When used with `--show-interfaces` flag the `status` command calls `vpp-probe discover`.

#### Trace

To trace packets across the data path for troubleshooting network connectivity 
issues, run:

```bash
swctl trace
```

> **Note**
> The `trace` command is a simple wrapper for `vpp-probe trace`.

#### Support

To export support file, run:

```bash
swctl support
```

> **Note**
> The `support` command is a simple wrapper for `agentctl report`.

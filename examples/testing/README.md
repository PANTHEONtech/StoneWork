Instructions For Running The Testing Examples
=============================================

The examples included here are used as automated tests, but they can also be explored manually.
The fisrt example is very simple and subsequent examples are progressively more complex,
so you can use these examples as an introduction into configuring StoneWork.

It is recommended to read the [Getting started examlpe][getting-started] before proceeding to these examples.

This document describes common prerequisities and instructions for running the examples.
Specific information about individual examples can be found in their respective README files.

Prerequisities
--------------

1. 
   Make sure that `stonework` production image and testing-only image `tester` are already built.
   In the top-level directory of the repository, trigger the makefile target `images` with:
   ```
   $ make images
   ```
2. 
   **Optional**: StoneWork can be managed through `agentctl`, a CLI provided by the Ligato framework.
   You can install it with
   ```
   $ go get go.ligato.io/vpp-agent/v3/cmd/agentctl
   ```
   This will enable you to use `agentctl` to access StoneWork directly from the host machine.

   If you skip this step, you will still be able to use `docker` to access `agentctl`, that comes pre-installed in the `stonework` container.

Running The Examples: The Automated Way
---------------------------------------

The examples can be easily controlled using their respective Makefile.

Just navigate to the directory of the chosen example and run `make test`
to run the example, perform automated tests and dump logs in case of failure.  
By default the tests are run with [the official StoneWork image from ghcr](https://github.com/PANTHEONtech/StoneWork/pkgs/container/stonework).
To use a custom image, you need to set the `STONEWORK_IMAGE` variable:
```
$ STONEWORK_IMAGE=stonework:<version> make test
```

Alternatively, to have more control over the execution of the example, you can use these commands:

- First, run `make start-example` to deploy the containers and apply the startup configuration. You can optionally specify a custom StoneWork image using `STONEWORK_IMAGE=stonework:<version> make start-example`.
- When the containers are running, run `make test-stonework` to run automated tests.
- When the containers are running, run `make dump-logs` to dump StoneWork logs.
  The logs will be saved into a file `example.log` (if the file exists, it will be overwritten).
- Finally, shut down and clean up everything (except dumped logs) with `make stop-example`.

Running The Examples: The Manual Way
------------------------------------

### Start The Example

Navigate to the directory of the chosen example and start the deployment with
```
$ docker-compose up -d
```
or using a custom StoneWork image:
```
STONEWORK_IMAGE=stonework:<version> docker-compose up -d
```
This command starts the example topology (it is defined in `docker-compose.yaml`).

The initial configuration (`config/day0-config.yaml`) is applied automatically (more info about this in [Getting started example][getting-started]).

### Explore and Manage StoneWork

You can use `agentctl` to manage StoneWork.

If you have installed `agentctl` on the host, you can simply run
```
$ agentctl <command>
```
Otherwise, you can run
```
$ docker-compose exec stonework agentctl <command>
```
For more information about `agentctl` commands, see [Getting started example][getting-started].

You can explore any container in the topology using
```
$ docker-compose exec -it <container-name> bash
$ # Example:
$ docker-compose exec -it stonework bash
```

**Note**: Automated tests for StoneWork can be found in `test-stonework.sh`.

You can use this file as a source of inspiration for running your own commands.

Or you can run the tests using:
```
$ ./test-stonework.sh
```

### StoneWork Logs

To study StoneWork logs, run
```
$ docker-compose logs stonework
```

### Stop The Example

When you are done, bring the deployment down with
```
$ docker-compose down
```

[getting-started]: ../getting-started/README.md

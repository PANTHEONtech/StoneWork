[Example] StoneWork as a Cross-Connect
============================

This example demonstrates how to use StoneWork as a cross-connect.

Network Diagram
---------------

Boxes in the diagram below denote Docker containers.
The interfaces attached to `stonework` are cross-connected in both directions.
```
+---------+                  +-----------+                  +---------+
|         |                  |           |                  |         |
| tester1 +------------------+ stonework +------------------+ tester2 |
|         | 10.10.1.1/24     |           |     10.10.1.2/24 |         |
+---------+                  +-----------+                  +---------+
```

Prerequisities
--------------

1. 
   Make sure that `stonework` production image and testing-only image `tester` are already built and tagged as `latest`.
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

Running The Example: The Automated Way
----------------------------------

The example can be easily controlled using the provided Makefile.

Just run `make test` to run the example, perform automated tests and dump logs in case of failure.

Alternatively, to have more control over the execution of the example, you can use these commands:

- First, run `make start-example` to deploy the containers and apply the startup configuration.
- When the containers are running, run `make test-stonework` to run automated tests.
- When the containers are running, run `make dump-logs` to dump StoneWork logs.
  The logs will be saved into a file `example.log` (if the file exists, it will be overwritten).
- Finally, shut down and clean up everything (except dumped logs) with `make stop-example`.

Running The Example: The Manual Way
-------------------------------

Start the deployment with
```
$ docker-compose up -d
```
This command starts the example topology (it is defined in `docker-compose.yaml`).

The initial configuration (`config/day0-config.yaml`) is applied automatically (more info about this in [Getting started example][getting-started]).

You can use `agentctl` to manage StoneWork.

If you have installed `agentctl` on the host, you can simply run
```
$ agentctl <command>
```
Otherwise, you can run
```
$ docker exec stonework agentctl <command>
```
For more information about `agentctl` commands, see [Getting started example][getting-started].

You can explore any container in the topology using
```
$ docker exec -it <container-name> bash
$ # Example:
$ docker exec -it stonework bash
```
## Automated Tests for StoneWork

Automated tests for StoneWork can be found in `test-stonework.sh`.

You can use this file as a source of inspiration for running your own commands.

Or you can run the tests using:
```
$ ./test-stonework.sh
```

To study StoneWork logs, run
```
$ docker logs stonework
```

When you are done, bring the deployment down with
```
$ docker-compose down
```

[getting-started]: ../../getting-started/README.md

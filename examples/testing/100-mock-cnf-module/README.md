[Example] StoneWork deployment based on Docker-Compose
======================================================

Prerequisites
--------------

Make sure that `stonework` production image and testing-only image `mockcnf` are already built.
In the top-level directory of the repository, trigger makefile target `images`:
```
$ make images
```

Running The Example
---------------

The example can be easily controlled using the provided Makefile. 

1. Run `make start-example` to deploy the containers and apply the startup configuration
2. Trigger `make test-stonework` to verify that the dynamic integration of mock CNF with StoneWork works as expected
3. Shutdown and clean everything up with `make stop-example`.

Manual Verification
-------------------

Start the deployment with:
```
$ docker compose up -d
```

StoneWork can be managed through the CLI, provided by the Ligato framework:
```
$ go get go.ligato.io/vpp-agent/v3/cmd/agentctl
```

An example configuration can be found here:
```
./config/day0-config.yaml
```

Change the initial configuration with:
```
$ agentctl config update --replace ./config/running-config.yaml
```

Use the commands provided by `agentctl` and also study logs collected by Docker for `stonework`, `mockcnf1` and `mockcnf2`
to verify, that CNFs were successfully integrated with StoneWork.

Automated checks can be found in `test-stonework.sh`.

Bring the deployment down with:
```
$ docker compose down
```

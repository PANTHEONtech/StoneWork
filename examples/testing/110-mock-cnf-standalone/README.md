[Example] Standalone CNF (running its own VPP instance)
=============================================

Prerequisites
--------------

Make sure that the testing-only image, `mockcnf`, is already built.
In the top-level directory of the repository, trigger the makefile target `images`:
```
$ make images
```

Run The Example
---------------

The example can be easily controlled using the provided Makefile. 

1. Run `make start-example` to deploy the containers and apply the startup configuration
2. Trigger `make test-stonework` to verify that the dynamic integration of mock CNF with StoneWork works as expected
3. Shutdown and clean everything up with `make stop-example`.

Manual Verification
-------------------

Start the deployment with:
```
$ docker-compose up -d
```

CNF can be managed through CLI provided by the Ligato framework:
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

Use command provided by `agentctl` and also study logs collected by docker for the CNF to verify that the CNFs operates in the standalone mode, without any issues.

Automated checks can be found in `test-cnf.sh`.

Bring the deployment down with:
```
$ docker-compose down
```

Standalone CNF (running its own VPP instance)
=============================================

Prerequisites
--------------

Make sure that testing-only image `mockcnf` is already built.
In the top-level directory of the repository, trigger makefile target `images`:
```
$ make images
```

Run the example
---------------

The example can be easily controlled using the provided Makefile. Run `make start-example` to deploy the container
and apply the startup configuration, then trigger `make test-cnf` to verify that the CNF operates in the standalone
mode without any issues and finally shutdown and clean up everything with `make stop-example`.

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

Use command provided by `agentctl` and also study logs collected by docker for the CNF to verify
that CNFs operates in the standalone mode without any issues.
Automated checks can be found in `test-cnf.sh`.

Bring the deployment down with:
```
$ docker-compose down
```

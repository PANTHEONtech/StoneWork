[Example] StoneWork as an IPv6 Router
=====================================

This example demonstrates how to use StoneWork as an **IPv6 router**.

This example contains two configuration files.

The file `config/day0-config.yaml` is used in automated tests. It contains a minimal configuration, required to make the example work.

**Note:** The first one or two pings may fail because of *not-configured ARPs*. During these pings, the ARPs get configured automatically, so subsequent pings work as expected.

The file `config/config-with-arps.yaml` is not used in tests. It can be used when running the example manually. It has all the contents from `day0-config.yaml` and it also has configuration for ARPs.

When using [this configuration](config/config-with-arps.yaml), all the pings, including the first one, should be successful.

Network Diagram
---------------
Boxes in the diagram below denote Docker containers.

Configured routes are shown below the containers.

Note that IPv6 forwarding needs to be enabled in the `tester2` container (see `docker-compose.yaml`).
```
+---------+ 2001:0:0:1::/64 +-----------+ 2001:0:0:2::/64 +---------+ 2001:0:0:3::/64 +---------+
|         |                 |           |                 |         |                 |         |
| tester1 +-----------------+ stonework +-----------------+ tester2 +-----------------+ tester3 |
|         | ::1         ::2 |           | ::1         ::2 |         | ::1         ::2 |         |
+---------+                 +-----------+                 +---------+                 +---------+

default                     2001:0:0:3::/64               2001:0:0:1::/64             default
via 2001:0:0:1::2           via 2001:0:0:2::2             via 2001:0:0:2::1           via 2001:0:0:3::1
```

Running The Example
-------------------

Prerequisities and instructions for running the example can be found in the
[README.md][readme-testing-examples] in `testing` directory.

**Note:** to use the alternate configuration file, start the example as usual and then run
```
$ docker exec stonework agentctl config update --replace /etc/stonework/config/config-with-arps.yaml
```

[readme-testing-examples]: ../README.md

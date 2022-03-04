[Example] StoneWork as an IPv4 Router
=====================================

This example demonstrates how to use StoneWork as an **IPv4 router**.

Network Diagram
---------------
Boxes in the diagram below denote Docker containers.

The configured routes are shown below the containers.

The topology here is bigger than the topologies in previous examples.

The reason for this is to have a container, that is **not** directly connected to StoneWork, so we can demonstrate adding routes to StoneWork.

Note that the default routes, pre-configured by Linux, need to be removed (see `docker-compose.yaml`).

Otherwise, they would conflict with the default routes that StoneWork tries to configure.
```
+---------+ 10.10.1.0/24 +-----------+ 10.10.2.0/24 +---------+ 10.10.3.0/24 +---------+
|         |              |           |              |         |              |         |
| tester1 +--------------+ stonework +--------------+ tester2 +--------------+ tester3 |
|         | .1        .2 |           | .1        .2 |         | .1        .2 |         |
+---------+              +-----------+              +---------+              +---------+

default                  10.10.3.0/24               10.10.1.0/24             default
via 10.10.1.2            via 10.10.2.2              via 10.10.2.1            via 10.10.3.1
```

Running The Example
-------------------

Prerequisities and instructions for running the example can be found in the
[README.md][readme-testing-examples] in `testing` directory.

[readme-testing-examples]: ../README.md
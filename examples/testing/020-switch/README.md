[Example] StoneWork as a Switch
===============================

This example demonstrates how to use StoneWork as a switch.

Network diagram
---------------

Boxes in the diagram below denote Docker containers.
The interfaces attached to `stonework` are connected using a bridge domain.
```
+---------+                  +-----------+                  +---------+
|         |                  |           |                  |         |
| tester1 +------------------+ stonework +------------------+ tester2 |
|         | 10.10.1.1/24     |           |     10.10.1.2/24 |         |
+---------+                  +-----------+                  +---------+
```

Running The Example
-------------------

Prerequisities and instructions for running the example can be found in the
[README.md][readme-testing-examples] in `testing` directory.

[readme-testing-examples]: ../README.md
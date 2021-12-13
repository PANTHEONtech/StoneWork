# ISISX (Intermediate System to Intermediate System Xconnect)
=============================================================

Intermediate System to Intermediate System Xconnect (ISISX) is a plugin designed for FD.io VPP.
Plugin introduces a way of handling ISIS protocol based packets what is not supported natively.
The filtered packets can be then routed between interfaces, hence the Xconnect.


# 1. Building the plugin
========================

## 1.1 Building directly from FD.io VPP sources
===============================================

Cloning an official FD.io VPP repository is mandatory `git clone https://gerrit.fd.io/r/vpp`,
however make sure appropriate version of VPP is checkouted to prevent any compatibility
issues with ISISX plugin that might possibly occur. Afterwards, just copy directory containing
plugin source files from `vpp21xx/isisx` into `src/plugins/` located inside root of cloned
VPP workspace. Due to recent build system improvements no other files need to be touched,
simply do `make build` what will automaticaly detect and compile integrated ISISX plugin.


## 1.2 Building externally via build script
===========================================

It is also possible to build ISISX plugin externally without compiling whole VPP what takes
significantly longer. This method assumes you have cloned FD.io VPP sources of the same
version and VPP is installed on your system because of required headers. If all requirements
are met, use provided build script `./build.sh /path/to/vpp/workspace`. Successful build will
produce .so and .api.json files. For complete integration with VPP, these artifacts need to be
copied to standard paths, e.g. `/usr/lib/x86_64-linux-gnu/vpp_plugins/`.

To verify whether the ISISX plugin was successfully loaded at startup by FD.io VPP,
run `vppctl show plugins`. If everything went correctly, ISISX plugin should appear
on the outputted plugin list.


# 2. Example usage
==================

Once the ISISX plugin integration is done with FD.io VPP, a quick setup can be done to get
plugin working. Example usage of binary API won't be covered since it wouldn't be very 
efficient when VPP supports multiple APIs. Instead, usage of CLI commands will be shown,
which implements the same behavior.


## 2.1 Interfaces setup
=======================

Before an Xconnect can be established between two interfaces, they need to be created:

`create tap host-ip4-addr 192.168.1.1/24`
`create tap host-ip4-addr 192.168.2.1/24`

A VPP should now recognize two created tap interfaces. Both have their unique name and interface index what can be verified:

`show interface`

Both interface states are by default set to 'down' and need to be 'up' to receive and transmit packets:

`set int state tap0 up`
`set int state tap1 up`


## 2.2 ISISX connection
=======================

Once both interfaces are ready a configuration can be added into the plugin:

`isisx connection add tap0 via tap1`

VPP will now route all ISIS packets received on tap0 to specified tap1 interface.

To dump all active ISISX configurations with their responding interfaces:
`show isisx connect`

When there's no more need to have ISIS enabled on specific interface, configuration can be removed:
`isisx connection del tap0`


## 2.3 Packets trace
====================

For verification purposes it is possible to check whether all ISIS incoming packets are transmitted
to configured interface. Of course, trace needs to be added before any packets are received:

`trace add virtio-input 50`
`show trace`

Most bottom nodes should report ISISX plugin and interface-output.


## 3. How does ISISX plugin work
================================

In overall, the ISISX is a simple VPP plugin that has pretty straightforward usage. The filtering
of ISIS packets is done by registering OSI input protocol on plugin's node. Internally, plugin
keeps reference of hash table that holds data about configured interfaces. RX interface is
the primary key that is used for accessing elements in the hash table. It is worth mentioning
that it is possible to have only one configuration per each RX interface. That way it is possible
to check whether configuration exists for specific RX interface and then forward the packet
in plugin's node. Modifying hash table configurations is done either by binary API or CLI.

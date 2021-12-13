Punt Manager Plugin
===================

To **"punt"** can mean different things to different people. In VPP the data-plane punts when a packet cannot be handled
by any further nodes. Punt differs from drop, in that VPP is giving other elements of the system the opportunity
to handle this packet.

For StoneWork the meaning of punt is to send packets to the user/control-plane of a CNF (typically a 3rd party
open-source software packaged alongside a CNF). This is specific option of the more general case described above, where VPP
is handing the packet to the control-plane for further prosessing.

**Punt Manager** plugin allows for multiple ligato plugins and even distributed agents to request packet punting
between a shared VPP and the same or distinct Linux network namespace(s) using TAPs or between the VPP and CNFs
directly using memifs or AF-UNIX sockets. Unless there is a conflict between punt requests,
the manager will ensure that common configuration items are shared and properly updated (e.g. ABX rules, TAP
connection, etc.). The manager supports different kinds of packet punting approaches for L2 or L3 source VPP
interfaces, with memifs, TAPs or AF-UNIX sockets used to deliver packets to the Linux network stack / user-space
application.

The plugin can be used by:
  - *Standalone CNF* (even for a single punt it is a good practise to use the plugin),
  - *StoneWork* to orchestrate punt between the all-in-one VPP and every *SW-Module*,
  - and by a *SW-Module* to learn the metadata about a created punt configuration.

Supported Punt Types
--------------------

Multiple different types of packet punting methods and topologies are supported to satisfy the wide-range of
requirements from present and future CNFs:

  - **HAIRPIN_XCONNECT**: create an L2 "hairpin x-connect" using TAPs or MEMIFs as follows:
    ```
    vpp_interface1 <-> vpp tap/memif 1 <-> linux tap/memif 1 -- CNF -- linux tap/memif 2 <-> vpp tap/memif 2 <-> vpp_interface2
    ```
     (i.e. hairpinning over linux network stack or via memif-enabled CNF)
  - **HAIRPIN**: like HAIRPIN x-connect except that while one side is attached to an existing L2 VPP interface,
    the other side is created as memif or TAP with given attributes. Basically it is like a feature attached
    to VPP interface (in the form of a new interface linked with an existing one, just like tunnel interfaces),
    which causes all traffic arriving/leaving via that interface to also flow through a CNF/Linux network stack before
    entering/exiting VPP. Unlike HAIRPIN x-connect it is therefore possible to attach further processing
    to this traffic (x-connect just forwards it through VPP unprocessed).
  - **SPAN**: copy traffic arriving and/or leaving via L2/L3 interface and send it to Linux or memif-enabled CNF.
  - **ABX**: effectively replicate L3 VPP interface in Linux using ACL-based xConnect as follows:
    ```
    vpp-interface with IP  <-- ABX --> unnumbered vpp memif/tap interface <-> Linux Tap / CNF memif
    ```
    Only packets matched by ACL associated with the ABX are punted.\
    Note: ABX is a proprietary feature developed by PANTHEON.tech.
  - **PUNT_TO_SOCKET**: Punt traffic matching given conditions (received through any interface) and punt it
     over a AF_UNIX socket.
  - **DHCP_PROXY**: Proxy DHCP requests for a given (L3) VRF into the Linux network stack or into a memif-enabled CNF.
  - **ISISX**: effectively replicate L3 VPP interface in Linux for ISIS protocol packets using xConnect as follows:
    ```
    vpp-interface with IP  <-- ISISX --> unnumbered vpp memif/tap interface <-> Linux Tap / CNF memif
    ```
    Basically it has the same goal as ABX, but ABX can't be used for ISIS protocol packets as packets
    for this protocol get dropped in VPP before reaching ACL VPP node.

The following diagram visually depicts all supported packet punting methods:

![Punt type][punt-types-diagram]


[punt-types-diagram]: img/punt-types.png

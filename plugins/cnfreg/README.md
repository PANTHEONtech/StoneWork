CNF Registry Plugin
===================

CNF Registry plugin allows to load a CNF module into the StoneWork (all-in-one VPP distribution; SW for short)
during the Init phase. CNF can be built as another image and run as a separate container.
This allows to enable/disable CNF without having to rebuild StoneWork docker image or the StoneWork agent binary.

Apart from single common VPP it is also possible to share network namespace between all/some CNFs and therefore
integrate different network functions inside the Linux network stack (e.g. to use OSPF-learned routes to connect
with a BGP peer).

The plugin operates in one of the 3 following modes depending on the value of the `CNF_MODE` environment variable:
 1. **STANDALONE** (default, i.e. assumed if the variable is not defined):
     - CNF is used on its own, potentially chained with other CNFs using for example NSM
       (i.e. each VPP-based CNF runs its own VPP instance)
     - The CNF Registry plugin is also used by a Standalone CNF, but merely to keep track of CNF Index ID.
 2. **STONEWORK_MODULE**:
     - CNF used as a SW-module
     - VPP-based CNFs do not run VPP inside their container, instead they connect with the all-in-one VPP of StoneWork
     - in this mode the Registry acts as a client of the Registry running by the StoneWork agent
     - internally the plugin uses gRPC to exchange all the information needed between the Registries of CNF and SW
       to load the CNF and use with the all-in-one VPP
     - CNF should use only those methods of the plugin which are defined by the `CnfAPI` interface
 3. **STONEWORK**:
     - CNF Registry plugin is used by StoneWork firstly to discover all the enabled CNFs and then to collect
       all the information about them to be able to integrate them with the all-in-one VPP
     - for each CNF, StoneWork needs to learn the NB-facing configuration models, traffic Punting to use (some
       subset of the traffic typically needs to be diverted from VPP into the Linux network stack via TAPs
       or directly into CNF using memifs for further processing)
     - StoneWork should use only those methods of the plugin which are defined by the `StoneWorkAPI` interface

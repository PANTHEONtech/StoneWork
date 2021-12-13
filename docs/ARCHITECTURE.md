StoneWork Architecture
======================

CNF Integration
---------------

Each CNF remains to be packaged and delivered as a single docker image. The same image can be used to deploy the
CNF either as **Standalone** (i.e. inside a chain/mesh of CNFs) or as a **StoneWork Module** or SW-Module for short.
Simply by setting the environment variable `CNF_MODE` to `STANDALONE` or `STONEWORK_MODULE`, the CNF will either start
its own instance of VPP or connect with the single shared VPP instance managed by Stonework, respectively.

StoneWork image itself consists of VPP (from upstream with some additional PANTHEON.tech plugins) and a control agent,
which is effectively [Ligato VPP-Agent][ligato-vpp-agent] extended with only two additional plugins -- 
[PuntManager][punt-manager-plugin] and [CNFRegistry][cnf-registry-plugin], working together to find and dynamically
load all enabled CNFs.
Neither control-plane nor data-plane features of any CNF are built-in to the StoneWork image directly. Not even
the protobuf configuration models. This means that the StoneWork image remains quite small and doesn't have
to be rebuilt even if a new CNF feature is added into the StoneWork. Instead every CNF runs as a separate container
from its own image, running its own control-plane agent that communicates with StoneWork to cooperatively
manage the single shared data-plane. All (enabled) CNFs and the StoneWork itself are typically orchestrated by
docker-compose. To enable a new CNF feature is then as easy as to add a container entry into docker-compose.yaml
for the CNF, restart the deployment and let the StoneWork to discover it.
If instead for a particular deployment a given CNF is never used, it doesn't have to be mentioned in docker-compose.yaml
and the image itself does not have to be shipped to the target device. This means that only CNFs that are actually
needed will use the system resources.

Such dynamic integration is possible because [CNFRegistry][cnf-registry-plugin] allows StoneWork to discover all CNFs,
learn their configuration models over gRPC and prepare KVDescriptors for them that will proxy CRUD operations
between StoneWork and CNFs. [PuntManager][punt-manager-plugin] then allows to establish shared or separate data paths
for punted packets (i.e. packets received by VPP but forwarded to the data-plane of a CNF for further processing).

The following diagram visually depicts StoneWork deployment consisting of one StoneWork container and one container
per CNF. The established data path links are marked with blue color, while control-plane links are orange colored.
The proxying of CRUD operations is highlighted using the red color:

![StoneWork Diagram][stonework-diagram]


Control-Flow
------------

The following diagram shows the sequence of communication between StoneWork and CNFs during the initialization
as well as during run-time when a (CNF) configuration change is requested:

![Control-Flow Diagram][control-flow-diagram]


[ligato-vpp-agent]: https://github.com/ligato/vpp-agent
[stonework-diagram]: img/stonework.png
[control-flow-diagram]: img/control-flow.png
[punt-manager-plugin]: plugins/puntmgr/README.md
[cnf-registry-plugin]: plugins/cnfreg/README.md



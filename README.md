# StoneWork

 [![CI](https://github.com/PANTHEONtech/StoneWork/actions/workflows/ci.yml/badge.svg)](https://github.com/PANTHEONtech/StoneWork/actions/workflows/ci.yml)
 [![stable](https://img.shields.io/github/release/PANTHEONtech/StoneWork.svg?label=latest%20release&logo=github)](https://github.com/PANTHEONtech/StoneWork/releases/latest)
 [![ligato/vpp-agent](https://img.shields.io/badge/image-ghcr.io/pantheontech/stonework-blue.svg?logo=docker&logoColor=white)](https://github.com/PANTHEONtech/StoneWork/pkgs/container/stonework)

A **high-performance** data plane, **modular** control plane solution.

StoneWork is used by PANTHEON.tech to integrate [its cloud-native network functions][cdnf-io] on top of a single shared
[FD.io VPP][VPP] data plane instance, to achieve the *best possible* resource
utilization. 

This network appliance, however, is not a step back from distributed chained/meshed microservices, to monolithic architecture. 

Instead, the integration is:
- Dynamic 
- Based on container orchestration
- CNF discovery
- Sharing of network namespaces and 
- Re-use of data paths for packet punting between CNFs

## Features

* High-performance [VPP][VPP]-based **data plane**
* Management agent build on top of [Ligato VPP-Agent][ligato-vpp-agent]
* Suitable for both **cloud & bare-metal** deployments
* Can be deployed as either *multiple* interconnected instances (service function
  chaining), or *a set* of control/management plane microservices that use a
  single VPP instance for data plane (this is a trade-off between flexibility
  and resource utilization)
* **Northbound APIs** are modeled with protobuf and accessible over `gRPC`, `REST`,
  `K8s CRD`or through a key-value DB (`etcd`, `redis`, ...)
* Wide-range of **networking features**, natively implemented in VPP, e.g.:
    * High-performance device drivers (DPDK, RDMA, virtio)
    * Routing, switching
    * Tunneling (VXLAN, GRE, IP-IP)
    * ACL-based filtering and routing
    * NAT44, NAT64
    * Segment routing
    * VPN (Wireguard, IPSec)
    * Bridge domains, VRFs (multi-tenancy)
* **Management features** provided by the Linux network stack:
    * Routes, ARPs
    * iptables
    * Namespaces, VRFs (multi-tenancy)
* Dynamically (at run-time) **extensible** with additional features provided by
  [CNFs from PANTHEON.tech][cdnf-io]

## Examples

Before using StoneWork, we recommend reading this README and related
documentation in the StoneWork [distribution folder][docs]. 

If you are **new to StoneWork**, it may be easier to first
explore and run the provided examples, rather than trying to create deployment
manifests from scratch.

Examples of deployment manifests and configurations for various use-cases can be found under the [examples sub-directory][examples].

The [Getting Started][getting-started] example will guide you through your first StoneWork
deployment.


## Configuration

Configuration for StoneWork consists of two tasks:

#### 1. VPP Startup Configuration

The VPP Startup Configuration comprises configuration options, which are set
before VPP is started. They *cannot* be changed at the run-time, either by a
management plane API or the VPP CLI). For StoneWork, the default VPP startup
configuration file is packaged in the image, under `/etc/vpp/vpp.conf`. 

Some of the [examples][examples] override the default configuration with a customized
version of `vpp.conf`mounted into the container using volumes. Typically, the
only configuration section that may require customization is the `dpdk` stanza,
where PCI addresses of NICs, used by VPP, should be listed. 

Run the `lshw-class network -businfo` command to view the available network devices
and their respective PCI addresses. For example, if the PCI addresses of
interfaces were `0000:00:08.0` and `0000:00:09.0` (e.g. inside and outside
network), then the `dpdk` configuration would be:
```
dpdk {
    dev 0000:00:08.0 {
        name eth0
    }
    dev 0000:00:09.0 {
        name eth1
    }
}
```
Interface names can be selected arbitrarily, for example `eth0` and `eth1`'
in the above example.

More information about attaching physical interfaces into VPP can be found
[here][vpp-pci].

#### 2. Protobuf-modeled Network Configuration

StoneWork's network configuration (VPP, Linux, CNFs) is modeled using
Google Protocol Buffers. 

A summary of all configuration items and their
attributes, with descriptions, can be found [here][config] (in markdown; also
available as a single [PDF document][config-pdf]). 

A [JSON Schema][config-jsonschema] is provided as well, and can be used to validate input
configuration before it is submitted. 

Some text editors, for example
[VS Code][vscode-jsonschema], can even load the Schema and provide
autocomplete suggestions based on it, thus making the process of preparing
input configuration a lot easier. 

The original protobuf files, from which the
documentation and schema were generated, can be found in the `/api` folder inside
the StoneWork distribution. There is also the `/api/models.spec.yaml` file,
which contains one YAML document with metadata for every configuration model.

These metadata are used to associate a configuration model with the corresponding
protobuf definitions.

Network configuration is submitted into the control-plane agent either via:

- a **[CLI][agentctl]** (YAML formatted), written into a key-value datastore (e.g.
`etcd`; JSON-formatted) 

- or applied programmatically over **gRPC** (serialized by
protobuf) or REST (JSON) APIs. The initial configuration that should be applied
immediately after StoneWork starts up can be mounted into the container under
`/etc/stonework/config/day0-config.yaml` (YAML formatted).

Each of the attached [examples][examples] has a sub-directory named `config`,
where you can find configuration stanzas to learn from. Each example contains
the startup configuration `day0-config.yaml`. 

Additional `*.yaml` files are used to show how run-time configuration can be modified over CLI. Please
refer to each examples `README.md` file for more information.


## Installation

The following steps will guide you through the StoneWork **installation process**.
The distribution package contains the **StoneWork Docker image** (`stonework.image`),
documentation (`*.md`) and some examples to get you started.

#### Requirements

1. StoneWork requires an **Ubuntu VM** or a **bare-metal server** running Ubuntu, preferably version **20.04 (Focal Fossa)**.


2. Next, Docker and Docker Compose plugin must be installed.

   Official manual for installing Docker and Docker Compose can be found [here][install-docker] and [here][install-compose] respectively.

3. **(DPDK Only)** Install/Enable Drivers
   
   Depending on the type of NICs that VPP of StoneWork should bind to, you may
   have to install/enable the corresponding drivers. 
   
   For example, in a VM environment, the [Virtual Function I/O (VFIO)][vfio] is preferred over the
   UIO framework for better performance and more security. In order to load a VFIO
   driver, run:
   ```
   $ modprobe vfio-pci
   $ echo "vfio-pci" > /etc/modules-load.d/vfio.conf
   ```
   Check with:
   ```
   $ lsmod | grep vfio_pci
   vfio_pci               45056  0
   ```
   More information about Linux network I/O drivers that are compatible with
   DPDK (used by VPP), can be found [here][dpdk-linux-drivers].


4. **(DPDK Only)** Check Network Interfaces
   
   Make sure that the network interfaces are not already used by the Linux
   kernel, or else VPP/DPDK will not be able to grab them. Run `ip link set
   dev {device} down` for each device to un-configure it from Linux. Preferably
   disable the interfaces using configuration files to make the changes
   persistent (e.g. inside `/etc/network/interfaces`).


5. **(DPDK Only)** Huge Pages
   
   In order to optimize memory access, VPP/DPDK uses [Huge Pages][hugepages],
   which have to be allocated before deploying StoneWork.
   For example, to allocate 512 Huge Pages (1024MiB memory for default 2M
   hugepage size), run:
   ```
   $ echo "vm.nr_hugepages=512" >> /etc/sysctl.conf
   $ sysctl -p
   ```
   Detailed recommendations on allocations of Huge Pages for VPP can be found
   [here][vpp-hugepages].


6. Finally, the StoneWork image has to be loaded so that
   Docker/Docker Compose/K8s is able to provision a container instance. Run:
   ```
   $ docker load <./stonework.image
   ```

## Deployment

StoneWork is deployed using [Docker Compose][compose] version 3.3 or
newer. StoneWork itself is only a single container (with VPP and StoneWork agent
inside), but every CNF that is deployed alongside it runs in a **separate
container**, hence the use of Compose. 

The following is a template for the
`docker-compose.yaml` file, used to describe deployment in the language of
Docker Compose. The template contains detailed comments, that explain the meaning
of attributes contained in the template and how they work in StoneWork. 

Angle brackets are used to mark placeholders that have to be replaced with appropriate
actual values in the target deployment.

```yaml
version: '3.3'

# Volume shared between StoneWork and every CNF deployed alongside it.
# CNFs and StoneWork use it to discover each other.
volumes:
  runtime_data: {}

services:
  stonework:
    container_name: stonework
    image: "ghcr.io/pantheontech/stonework:22.10"
    # StoneWork runs in the privileged mode to be able to perform administrative network operations.
    privileged: true
    # StoneWork runs in the PID namespace of the host so that it can read PIDs of CNF processes.
    pid: "host"
    environment:
      # Set log level (i.e. only log entries with that severity or anything above it will be printed).
      # Supported values: Trace, Debug, Info, Warning, Error, Fatal and Panic.
      INITIAL_LOGLVL: "debug"
      # MICROSERVICE_LABEL is used to mark container with StoneWork.
      MICROSERVICE_LABEL: "stonework"
      # By default etcd datastore is used as the source of the configuration.
      # Env. variable ETCD_CONFIG with empty value is used to disable etcd
      # and use CLI (agentctl) or gRPC as the primary source of the configuration.
      ETCD_CONFIG: ""
    ports:
      # Expose HTTP and gRPC APIs.
      - "9111:9111"
      - "9191:9191"
    volumes:
      # /run/stonework must be shared between StoneWork and every CNF.
      - runtime_data:/run/stonework
      # /sys/bus/pci and /dev are mounted for StoneWork to be able to access PCI devices over DPDK.
      - /sys/bus/pci:/sys/bus/pci
      - /dev:/dev
      # Docker socket is mounted so that StoneWork can obtain container metadata for every CNF.
      - /run/docker.sock:/run/docker.sock
      # To customize vpp startup configuration, create your own version of vpp.conf (here called vpp-startup.conf),
      # put it next to this docker-compose.yaml and mount it under /etc/vpp/vpp.conf.
      # Otherwise remove this mount.
      - ./vpp-startup.conf:/etc/vpp/vpp.conf
      # To start StoneWork with some initial configuration, create day0-config.yaml under the config
      # sub-directory, placed next to this docker-compose.yaml and mount it under /etc/stonework/config
      # Otherwise remove this mount.
      - ./config:/etc/stonework/config

  # Multiple CNFs may share the same Linux network namespace. This is in some case needed
  # if CNFs are to work together (e.g. BGP peering established over OSPF-learned routes).
  # The common network namespace is represented by a separate container (similar to the
  # sandbox container of a K8s Pod).
  router-ns:
    container_name: router-ns
    image: "busybox:1.29.3"
    command: tail -f /dev/null

  # CNF running alongside StoneWork (i.e. using the VPP of StoneWork as data-plane).
  # Name the container such that it is clear what services CNF provides (e.g. "cnf-dhcp").
  <cnf-name>:
    container_name: <cnf-name>
    image: "<cnf-image-name>"
    depends_on:
      - stonework
    # CNFs typically require privileges to perform administrative network operations.
    privileged: true
    # <cnf-name>-license.env is file that is obtained when the license of CNF is purchased.
    # Put <cnf-name>-license.env into the same directory as docker-compose.yaml.
    # It contains single line:
    # LICENSE=<signed license content>
    env_file:
      - <cnf-name>-license.env
    volumes:
      # /run/stonework must be shared between StoneWork and every CNF.
      - runtime_data:/run/stonework
    environment:
      INITIAL_LOGLVL: "debug"
      # MICROSERVICE_LABEL is effectively used to mark the container with CNF name.
      # StoneWork is then able to identify the CNF container among all containers.
      MICROSERVICE_LABEL: "<cnf-name>"
      ETCD_CONFIG: ""
      # If CNF runs alongside StoneWork (and not standalone), env. variable "CNF_MODE"
      # must be defined with value "STONEWORK_MODULE".
      CNF_MODE: "STONEWORK_MODULE"
    # Multiple CNFs may share the same Linux network namespace.
    # Use network_mode and point a group of CNFs to the same container (acting just like sandbox
    # container of a K8s Pod).
    network_mode: "service:router-ns"

  # here list other CNFs...
```

## Development
[![Go Reference](https://pkg.go.dev/badge/go.pantheon.tech/stonework.svg)](https://pkg.go.dev/go.pantheon.tech/stonework)

- **Build**: Build instruction for StoneWork can be found [here][build].
- **Architecture**: StoneWork architecture is described in detail [here][architecture].
- **CNF Compatibility**: A guide on how to make CNF compatible with StoneWork can be found [here][cnf-how-to].
- **GNS3 & StoneWork**: StoneWork GNS3 VM development documentation is [here][gns3-vm-docs].

[architecture]: docs/ARCHITECTURE.md
[build]: docs/BUILD.md
[cnf-how-to]: docs/CNF_HOW_TO.md
[config]: docs/config/STONEWORK-CONFIG.md
[config-pdf]: docs/config/STONEWORK-CONFIG.pdf
[config-jsonschema]: docs/config/STONEWORK-CONFIG.jsonschema
[vscode-jsonschema]: https://dev.to/brpaz/how-to-create-your-own-auto-completion-for-json-and-yaml-files-on-vs-code-with-the-help-of-json-schema-k1i
[compose]: https://docs.docker.com/compose/
[install-docker]: https://docs.docker.com/engine/install/ubuntu/
[install-compose]: https://docs.docker.com/compose/install/linux/
[vpp]: https://wiki.fd.io/view/VPP
[ligato-vpp-agent]: https://github.com/ligato/vpp-agent
[cdnf-io]: https://cdnf.io/cnf_list/
[examples]: examples/README.md
[getting-started]: examples/getting-started/README.md
[agentctl]: https://docs.ligato.io/en/latest/user-guide/agentctl/
[vpp-pci]: https://wiki.fd.io/view/VPP/How_To_Connect_A_PCI_Interface_To_VPP
[vfio]: https://www.kernel.org/doc/Documentation/vfio.txt
[dpdk-linux-drivers]: https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html
[hugepages]: https://wiki.debian.org/Hugepages
[vpp-hugepages]: https://s3-docs.fd.io/vpp/21.10.1/gettingstarted/users/configuring/hugepages.html
[gns3-vm-docs]: docs/GNS3_APPLIANCES.md

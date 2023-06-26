StoneWork: Getting Started
==========================

If you are new to StoneWork, this is a minimalistic example of the StoneWork
deployment manifest that you can **try right away**, without having to make any
configuration adjustments for your environment. 

In this example, StoneWork is not attached to any of the physical interfaces, instead only a **TAP**-based,
L3 inter-connection is configured between the host and the VPP instance running inside
the StoneWork container. 

This is enough to successfully: 
- Ping the VPP instance
- Try StoneWork & VPP CLIs 
- Discover the provided REST APIs 
- Learn to read StoneWork logs
- and more

Deployment Description
----------------------

This example consists of two YAML-formatted files:
 - `docker-compose.yaml`: describes how to deploy the StoneWork container in the
   [Docker-Compose][docker-compose] language. While in this example there is
   only **one container** - StoneWork itself - it is still advised to take this
   opportunity and learn how to work with docker-compose. Extending StoneWork's
   feature-set with any of the CNFs from the PANTHEON.tech
   [cloud-native network functions portfolio][cdnf-portfolio] requires deploying additional containers
   alongside StoneWork, hence the use of *Compose*. The content of
   `docker-compose.yaml` is described in the [top level README.md][readme].


 - `config/day0-config.yaml`: notice that `docker-compose.yaml` mounts the
   adjacent `./config` directory into the StoneWork container under
   `/etc/stonework/config`. When StoneWork starts, it will look for the
   `/etc/stonework/config/day0-config.yaml` configuration file. If the file is
   found, StoneWork will apply it immediately after it transits from *init* to
   *ready* state. This file is not mandatory. StoneWork can start with an empty
   config state and receive a configuration later.

In this example, there is only a TAP interconnection configured between the host
and VPP. Both sides of the TAP are configured separately, one under the
`vppConfig`, the other under `linuxConfig`. The configured model is
described in [detail here][config].

### Network Diagram

Depicted below is the network topology of this very simple example:

```
 +-------+  192.168.222.0/30  +------------+
 |       |       (TAP)        |            |
 | host  +--------------------+ StoneWork  |
 |       | .2              .1 |            |
 +-------+                    +------------+
```

Interacting with StoneWork
--------------------------

In order to deploy StoneWork, simply run (from within this directory):
```
$ docker-compose up -d
```
The StoneWork container should be present almost immediately:
```
$ docker ps
CONTAINER ID   IMAGE                                  COMMAND                  CREATED         STATUS         PORTS     NAMES
83f034b16cd8   ghcr.io/pantheontech/stonework:23.06   "/bin/sh -c 'rm -f /…"   6 seconds ago   Up 5 seconds             stonework
```

It shouldn't take too long for StoneWork to initialize and apply the "day0"
configuration (check after 5 sec).

First, you will notice that the `vpp` interface appeared in the host - this is our
TAP between the host and VPP:
```
$ ifconfig vpp
vpp: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.222.2  netmask 255.255.255.252  broadcast 192.168.222.3
        inet6 fe80::fe:64ff:fe1c:2660  prefixlen 64  scopeid 0x20<link>
        ether 02:fe:64:1c:26:60  txqueuelen 1000  (Ethernet)
        RX packets 4  bytes 336 (336.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 67  bytes 9830 (9.8 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

You can try to ping the VPP side of the TAP:
```
$ ping -c 3 192.168.222.1
PING 192.168.222.1 (192.168.222.1) 56(84) bytes of data.
64 bytes from 192.168.222.1: icmp_seq=1 ttl=64 time=0.839 ms
64 bytes from 192.168.222.1: icmp_seq=2 ttl=64 time=0.392 ms
64 bytes from 192.168.222.1: icmp_seq=3 ttl=64 time=0.247 ms

--- 192.168.222.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2626ms
rtt min/avg/max/mdev = 0.247/0.492/0.839/0.251 ms
```

Let's enter the VPP CLI:
```
$ docker-compose exec -it stonework vppctl
    _______    _        _   _____  ___ 
 __/ __/ _ \  (_)__    | | / / _ \/ _ \
 _/ _// // / / / _ \   | |/ / ___/ ___/
 /_/ /____(_)_/\___/   |___/_/  /_/    

vpp#
```

Use the VPP CLI for read-only operations (i.e. ping, show ...). Configuration
requests should be entered over the StoneWork APIs, as shown below. Before
configuring, let's use the VPP CLI to display the TAP interface (named `tap0` in
VPP) and view the Rx/Tx counters:
```
vpp# show interface 
              Name               Idx    State  MTU (L3/IP4/IP6/MPLS)     Counter          Count     
local0                            0     down          0/0/0/0       
tap0                              1      up          9000/0/0/0     rx packets                    43
                                                                    rx bytes                    6182
                                                                    tx packets                     4
                                                                    tx bytes                     336
                                                                    drops                         39
                                                                    ip4                           23
                                                                    ip6                           19
```
Use `q` to exit from the VPP shell:
```
vpp# q
```

## StoneWork Config & CLI

To add or change the StoneWork configuration, use either a human-friendly CLI, such
as the `agentctl`binary inside the StoneWork container, or a programmatic API,
such as `gRPC` and `REST`. It is also possible to push configuration over a
key-value datastore, such as etcd, that acts as a persistent storage for the
desired configuration state (just like in K8s).

Note, that using key-value datastore is **not covered** in this simple example.

Let's try the StoneWork CLI. It doesn't come with its own shell. Instead, every
command is a separate execution of the `agentctl` binary, installed inside the
StoneWork container.

Obtain the currently running configuration with:
```
$ docker-compose exec stonework agentctl config get 2>/dev/null
netallocConfig: {}
linuxConfig:
  interfaces:
  - name: linux-tap
    type: TAP_TO_VPP
    hostIfName: vpp
    enabled: true
    ipAddresses:
    - 192.168.222.2/30
    tap:
      vppTapIfName: vpp-tap
vppConfig:
  interfaces:
  - name: vpp-tap
    type: TAP
    enabled: true
    ipAddresses:
    - 192.168.222.1/30
    tap:
      version: 2
```
This is actually the applied *desired state*. To obtain the actual *running state*,
run `agentctl config retrieve`instead. The output of `config get` and
`config retrieve` differs if StoneWork failed to apply some configuration items
(shouldn't be the case here).

To change the configuration, prepare the new, desired configuration state first.
Either edit `config/day0-config.yaml`, or more preferably, make a copy
of the file and make the following edits there:
```
$ cd config
$ cp day0-config.yaml new-config.yaml

### edit new-config.yaml in your preferred editor, e.g. change IPs (here increased by 2):

$ cat new-config.yaml
vppConfig:
  interfaces:
    # VPP-side of the TAP interface
    - name: vpp-tap
      type: TAP
      enabled: true
      ipAddresses:
        - 192.168.222.3/30
      tap:
        version: 2 # virtio-based TAP interface

linuxConfig:
  interfaces:
    # Linux-side of the TAP interface
    - name: linux-tap
      type: TAP_TO_VPP
      hostIfName: vpp # name of the interface in the host
      enabled: true
      ipAddresses:
        - 192.168.222.4/30
      tap:
        vppTapIfName: vpp-tap # name of the VPP-side
```
Remember, that the config dir is mounted into
the container under `/etc/stonework/config`. Apply the new, desired config with:
```
$ docker-compose exec stonework agentctl config update --replace /etc/stonework/config/new-config.yaml
```
Observe the performed config changes:
```
$ ifconfig vpp
vpp: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.222.4  netmask 255.255.255.252  broadcast 192.168.222.7
        inet6 fe80::fe:e9ff:fe72:479  prefixlen 64  scopeid 0x20<link>
        ether 02:fe:e9:72:04:79  txqueuelen 1000  (Ethernet)
        RX packets 4  bytes 336 (336.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 77  bytes 11885 (11.8 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
$ docker-compose exec stonework vppctl show interface address
local0 (dn):
tap0 (up):
  L3 192.168.222.3/30
```

### Testing the REST API

Next, let's test the REST API. StoneWork REST APIs are (by default) exposed on
port 9191. To see the index of (almost) all provided REST APIs, open
`http://localhost:9191/`. Note that the `docker-compose.yaml` file contains the line
`network_mode: "host"`, which means that StoneWork runs inside the network
namespace of the host, hence the REST API is available on `localhost`.

For example, to obtain the currently applied running configuration (equivalent
to `agentctl config get`), send a GET request for the `/configuration` path:
```
$ curl http://localhost:9191/configuration
netallocConfig: {}
linuxConfig:
interfaces:
- name: linux-tap
  type: TAP_TO_VPP
  hostIfName: vpp
  enabled: true
  ipAddresses:
   - 192.168.222.2/30
     tap:
     vppTapIfName: vpp-tap
     vppConfig:
     interfaces:
- name: vpp-tap
  type: TAP
  enabled: true
  ipAddresses:
   - 192.168.222.1/30
     tap:
     version: 2
```

If you change the `config/new-config.yaml` again, apply the change over REST API
with a PUT request (notice the `replace=true` argument):
```
$ curl -v --header "Content-Type: application/yaml" --request PUT --data-binary "@config/new-config.yaml" localhost:9191/configuration?replace=true
```

### StoneWork Logs

To observe what is happening behind the scenes and to debug any potential issues,
obtain StoneWork logs with:
```
$ docker-compose logs stonework
```

### Shutdown & Undeploy

Finally, to shutdown and undeploy StoneWork container, simply run:
```
$ docker-compose down
```


[docker-compose]: https://docs.docker.com/compose/
[config]: ../../docs/config/STONEWORK-CONFIG.md
[readme]: ../../README.md
[cdnf-portfolio]: https://cdnf.io/cnf_list/

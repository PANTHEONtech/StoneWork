StoneWork User-Guide
====================

StoneWork is a control plane for FD.io VPP, a fast data plane capable of running on commodity hardware.

StoneWork introduces a **complete routing platform**, based on VPP, running well on both bare metal and cloud natively.

Dependencies
============

First of all we need to install Docker and *docker-compose*:
```
$ apt-get install docker.io docker-compose
```
First Steps w/ StoneWork
========================

StoneWork images are publicly available and can be pulled, [as described here](https://github.com/orgs/PANTHEONtech/packages/container/package/stonework).

## Run
To try out the StoneWork, we can simply run it as a Docker container:
```
$ docker run -d --rm --name stonework -e ETCD_CONFIG="" ghcr.io/pantheontech/stonework:23.06
```
This will run a Docker container named *stonework* in the background. 

An empty etcd config means that etcd won't be used as a source of configuration, but we will use a configuration file instead.

## Logs
To view StoneWork logs:
```
$ docker logs stonework
```

When working with StoneWork, it is useful to see the logs of control plane transactions. These are actions
done by StoneWork to set its dataplane - FD.io VPP. 

By looking at the logs, it can be verified whether concrete transaction were successful or not.

To access the StoneWork container:
```
$ docker exec -it stonework bash
```

Working w/ agentctl
===================

Now, when we are inside the container, we can manage StoneWork using [agentctl tool][agentctl-link]. This tool comes
from [ligato][ligato-link] and is the main tool for StoneWork management, using CLI.

It can be used for configuration retrieval, setup, getting models which serves as a regulations for the configuration and lots more.

Lets try to get the configuration first:
```
root@c95743f3e2e0:/# agentctl config get
netallocConfig: {}
linuxConfig: {}
vppConfig: {}
```

As we started a fresh StoneWork instance without any startup config, all of its root elements are empty.

To configure the StoneWork, we need to provide it with configuration file.

## Create StoneWork Config
Lets create a configuration file.
```
mkdir /etc/stonework/config/
```
```
  echo "vppConfig:
  interfaces:
  - name: loop-test-1
    type: SOFTWARE_LOOPBACK
    enabled: true
    ip_addresses:
    - 10.10.1.3/24
    mtu: 1500" >> /etc/stonework/config/my.conf
```
This example just declares the **loopback interface**, sets it up and sets its IP address and mtu.

## Update StoneWork Config
To update the StoneWork configuration, use following command, providing it with our newly created configuration:
```
$ agentctl config update --replace /etc/stonework/config/my.conf
```
To verify that the operation was successfull, try:
```
$ agentctl config history
```
The most recent entries are on the bottom. In this output, you usually see just the result, i.e. whether the operation was successful or not.
Sometimes we need to inspect the transactions more deeply, for example if our configuration file is corrupt or there are some unmet dependencies.

To do so, we should inspect Docker logs of StoneWork from outside of the container:
```
$ docker logs stonework
```
Where we can see much more information.

To remove the configuration, simply delete it or part of it from the configuration file and rerun the agentctl config
update command.

Working w/ vppctl
=================

Next we may want to check the state even in our data plane.
To access it use:
```
$ vppctl
```
Lets see the interfaces, to see the result of what we configured using StoneWork, showing our interface settings:

```
vpp# show int addr
local0 (dn):
loop0 (up):
  L3 10.10.1.3/24
```
Indeed, there is a new loopback interface with the specified IPv4 address and it's set to UP.

**IMPORTANT:** Please note that vppctl tool and all of its commands should be used only for troubleshooting.
To see the VPP state and running commands like ping, but not for configuration. This is important to keep
the system in a consistent state - all of the configuration has to be performed only from StoneWork control plane.

More VPP commands can be found [here][vpp-cli-guide].

TAP Interface
=============

To create tap interface, we need to make the StoneWork container privileged, so running it as:
```
$ docker run -d --rm --name stonework --privileged -e ETCD_CONFIG="" ghcr.io/pantheontech/stonework:23.06
```

Then we will use the configuration as:
```
vppConfig:
  interfaces:
  - name: vpp-tap1
    type: TAP
    enabled: true
    ip_addresses:
    - 10.10.10.1/24
    tap:
      version: 2
```
The new interface can be safely added to the existing configuration yaml file and the configuration
can be updated using agentctl as shown earlier:
```
$ agentctl config update --replace /etc/stonework/config/my.conf
```
Now we can verify that the tap is created in data plane:
```
root@44480538e5b0:/# vppctl sh int
              Name               Idx    State  MTU (L3/IP4/IP6/MPLS)     Counter          Count
local0                            0     down          0/0/0/0
tap0                              1      up          9000/0/0/0
```
and also inside the container:
```
root@44480538e5b0:/# ip add
...
2: tap-430978008: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UNKNOWN group default qlen 1000
    link/ether 02:fe:c1:1e:c8:d7 brd ff:ff:ff:ff:ff:ff
```
But this tap interface is present only inside our StoneWork container. What if we wanted to expose the tap to our host system?

Then we need to use host network mode in our docker run one liner:
```
$ docker run -it --rm --name stonework --privileged --network="host" -e ETCD_CONFIG="" ghcr.io/pantheontech/stonework:23.06
```
For more details about tap interfaces and VPP, take a look [at this example][tap-example].

docker-compose Manifest
=======================

As our Docker run one-liner grows, it becomes better to use docker-compose, in terms of readability and maintainability.

Lets rewrite the above mentioned docker run command, i.e.:
```
$ docker run -it --rm --name stonework --privileged --network="host" -e ETCD_CONFIG="" ghcr.io/pantheontech/stonework:23.06
```
into *docker-compose.yaml*:
```
version: '3.3'

services:
  stonework:
    container_name: stonework
    image: "ghcr.io/pantheontech/stonework:23.06"
    privileged: true
    network_mode: "host"
    environment:
      ETCD_CONFIG: ""
```
All of the fields should be obvious, as it is a 1-2-1 translation of the command we used so far.

The StoneWork docker container is then turned on/off by the following commands from the same directory as our docker-compose file is located:
```
$ docker-compose up -d
```
```
$ docker-compose down
```

AF_PACKET Interface
===================

Using this interface type, we can easily start using interfaces of the host system in StoneWork.

The downside of the af_packet interface is its performance, so it should be used mainly for testing purposes.

To use the host interfaces using AF_PACKET, the same docker-compose.yaml can be used as the one for TAPs.

Now, suppose we have an interface called ens33 and its MAC address is 00:0c:29:0a:93:4d.

As shown by:
```
$ ip link
...
2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP mode DEFAULT group default qlen 1000
    link/ether 00:0c:29:0a:93:4d brd ff:ff:ff:ff:ff:ff
```
Then our StoneWork config will look like this:
```
vppConfig:
  interfaces:
  - name: "my-iface"
    type: AF_PACKET
    enabled: true
    phys_address: "00:0c:29:0a:93:4d"
    afpacket:
      host_if_name: "ens33"
```
Please note that we need to set the **same MAC address** for the AF_PACKET as used by the underlying host interface, otherwise packets copied between VPP and the host would get dropped due to MAC mismatch.

DPDK Interface
==============

StoneWorks data plane - FD.io VPP - is designed for high-performance and for the same reason uses [DPDK library][dpdk].

This type of interfaces is the right choice for real world use cases, which requires big throughput. However, it's a bit
more complicated to set them.

Prerequisities
--------------

To enable DPDK interfaces in data plane, few specifics steps must be done, those are well described in the
[StoneWork README][readme] file, in the [Installation section.](/README.md#Installation)

docker-compose
--------------

To enable DPDK interfaces, we will need to extend our docker-compose a bit more.
```
version: '3.3'

services:
  stonework:
    container_name: stonework
    image: "ghcr.io/pantheontech/stonework:23.06"
    privileged: true
    network_mode: "host"
    environment:
      ETCD_CONFIG: ""
    volumes:                            # new
      - /sys/bus/pci:/sys/bus/pci       # new
      - /dev:/dev                       # new
      - ./vpp.conf:/etc/vpp/vpp.conf    # new
      - ./config:/etc/stonework/config  # new
```

First 2 volumes, i.e. `/sys/bus/pci:/sys/bus/pci` and `/dev:/dev` are mounted for DPDK to be able to access PCI devices and
to grab the interfaces from the system. Next volume contains VPP config, that one we are going to create in next section.
And finally, the config volume will contain StoneWork configuration, we will return to this one later and leave it empty for
now.

VPP Configuration
-----------------

Lets start with the following configuration (`vpp.conf`):
```
unix {
    interactive
    cli-no-pager
    cli-listen /run/vpp/cli.sock
    log /tmp/vpp.log
    coredump-size unlimited
    full-coredump

    # (!) Comment out for the best performance (CPU utilization will increase considerably).
    poll-sleep-usec 50
}

#dpdk {
#    dev 0000:02:06.0 {
#        name my-dpdk-iface
#    }
#}

api-trace {
    on
}

socksvr {
    default
}

statseg {
    default
    per-node-counters on
}

punt {
    socket /run/stonework/vpp/punt-to-vpp.sock
}
```
For now, we are only going to enable DPDK interface, understanding all of its parts is out of
scope of this tutorial, but you can visit [VPP docs][vpp-startup] for more information.

So to enable the interface, first we need to know its PCI address, lets view it by:
```
$ sudo lshw -class network -businfo
Bus info          Device      Class          Description
========================================================
pci@0000:02:01.0  ens33       network        82545EM Gigabit Ethernet Controller (Copper)
pci@0000:02:06.0  ens38       network        82545EM Gigabit Ethernet Controller (Copper)
                  virbr0-nic  network        Ethernet interface
```
Now lets say, we want to use just one interface - ens38, which PCI address is 0000:02:06.0.

To do so, just uncomment the dpdk part of the config above.

An arbitrary amount of interfaces can be used.

If we now start the docker container, VPP will contain our new DPDK interface:
```
$ docker exec -it stonework vppctl show int
              Name               Idx    State  MTU (L3/IP4/IP6/MPLS)     Counter          Count
local0                            0     down          0/0/0/0
my-dpdk-iface                     1     down         9000/0/0/0
```
So now, the data plane knows the interface and all we need to do now is just start using it in a control plane - StoneWork.

Configuring DPDK Interfaces in StoneWork
----------------------------------------

Basic StoneWork configuration for DPDK interfaces looks like this:
```
vppConfig:
  # Physical interfaces.
  interfaces:
    - name: my-dpdk-iface
      type: DPDK
      enabled: true
      ip_addresses:
        - 192.168.1.1/24
```
Lets place it in the config/ directory we left empty previously and call it `config/day0-config.yaml`.

Note, that if StoneWork finds the configuration at `/etc/stonework/config/day0-config.yaml` inside its container, it will
use it as its startup config and you don't need to use agentctl to apply it.

As probably obvious, this configuration is telling StoneWork to set the DPDK interface called `my-dpdk-iface` up and to assign it an IP address.

Now, it's time to wake up StoneWork:
```
$ docker-compose up -d
```
And to verify that everything was set properly in the data plane:
```
$ docker exec -it stonework vppctl show int addr
local0 (dn):
my-dpdk-iface (up):
  L3 192.168.1.1/24
```
How To Write the StoneWork Configuration
========================================

So far, we have seen few simple examples of StoneWork configuration. Now we will learn how to write a custom configuration.

Lets take a look at the [configuration model][conf-model]. The configuration is organized as a tree. We can see, that at
the very beginning there is a *root*, which contains 3 top most elements of StoneWork: 
- LinuxConfig
- NetallocConfig
- VppConfig

One additional configuration subtree will be available for every CNF deployed additionally alongside StoneWork - [see below](#adding-cnfs).

Also, if we consider our DPDK interface config, we can see *VppConfig* at the top most level:
```
vppConfig:
  # Physical interfaces.
  interfaces:
    - name: my-dpdk-iface
      type: DPDK
      enabled: true
      ip_addresses:
        - 192.168.1.1/24
```
To see what the options are inside vppConfig, we search the [configuration model][conf-model] for:
- "stonework.Root.VppConfig" (with quotes)

And a little bellow, we see that it contains lots of entries as ACLs, FIBs and of course interfaces.

To see, what options we have under interfaces, search further for:
- "ligato.vpp.interfaces.Interface"

Now, we can see all the options under interface, including name, type, ip_addresses and others. 

Note that there are 2 kinds of entries, those as name, are of concrete type as bool, string, or similar and other refers to the next
subtrees, which can be followed in the same manner.

You can use this method for both navigating through the StoneWork configuration and creation of custom one.

StoneWork Enterprise
====================

With StoneWork Enterprise you'll get access to additional control plane features. Enterprise features are packaged
as container images. To use them you need a valid license. 

First the CNF image must be loaded into docker as:
```
docker load -i <path/to/image/file>
```
Then, the `docker-compose.yaml` is updated.

Lets demonstrate it on cnf-bgp, this is how the new docker-compose will look like:
```
version: '3.3'

volumes:                                    # new
  runtime_data: {}                          # new

services:
  stonework:
    container_name: stonework
    image: "ghcr.io/pantheontech/stonework:23.06"
    privileged: true
    network_mode: "host"
    environment:
      INITIAL_LOGLVL: "debug"               # new
      MICROSERVICE_LABEL: "stonework"       # new
      ETCD_CONFIG: ""
    volumes:
      - runtime_data:/run/stonework         # new
      - /run/docker.sock:/run/docker.sock   # new
      - /sys/bus/pci:/sys/bus/pci
      - /dev:/dev
      - ./vpp.conf:/etc/vpp/vpp.conf
      - ./config:/etc/stonework/config

  bgp:                                      # new whole bgp service
    container_name: bgp
    image: "cnf-bgp-docker.pantheon.tech/cnf-bgp:21.01"
    depends_on:
      - stonework
    privileged: true
    env_file:
      - ./license.env
    environment:
      CNF_MODE: "STONEWORK_MODULE"
      INITIAL_LOGLVL: "debug"
      MICROSERVICE_LABEL: "bgp"
      ETCD_CONFIG: ""
    volumes:
      - runtime_data:/run/stonework
```
## New Entries Explanation

`MICROSERVICE_LABEL`  has to be unique for every CNF, thus the best practise is to use the same value for:
- service
- container_name
- MICROSERVICE_LABEL 

for sake of consistency and readability.

- `/run/stonework` directory, which is shared between stonework and all its CNFs, is for CNF <-> StoneWork discovery and for
sharing of memif sockets between StoneWork and CNFs.

- `/run/docker.sock` volume is to enable StoneWork communication with Docker, to be able to obtain the network namespace
handle for every CNF (stonework doesn't support any other container runtime yet).

- The BGP container name, image and dependency on StoneWork should be obvious.

- The `env_file` is the path to the license file. In StoneWork every CNF has separate license file which is needed to
provide this way to run the CNF.

- `CNF_MODE` environment variable is mandatory and basically sets the CNF to work in cooperation with StoneWork and StoneWork data plane instead of initiating its own data plane.

The rest of environment variables have the same meaning as for StoneWork.

### Start StoneWork

Now start the StoneWork with:
```
docker-compose up -d
```
And finally verify that the CNF was successfully discovered by examining the StoneWork logs as:
```
$ docker logs stonework 2>&1 | grep "Discovered CNFs"
time="2021-04-13 08:52:42.24204" level=debug msg="Discovered CNFs: map[bgp:{cnfMsLabel:bgp ipAddress:172.19.0.2 ...
```
As we can see, the `cnf-bgp` was discovered by StoneWork.

CNF Configuration Example
-------------------------

To start using the CNF functionality, we also need to update the StoneWork configuration file. 

The user should obtain the options inside the [configuration model][conf-model] file. This can be used following the same procedure as described in the **How to Write the StoneWork Configuration** section. Thus lets jump right into an example, the snippet of the `cnf-bgp` configuration.
```
bgpConfig:
  Server_list:
    - routerId: 1.1.1.1
      autonomousSystem: 100
      vrfId: 1
      neighbors:
        - network: 100.10.0.0/24

  SessionEndpoint_list:
    - vppInterface: memif-vrf1
      vrfId: 1
      remoteAsn: 100
      peers:
        - ipAddress: 10.10.0.10/24
```
So the `cnf-bgp` introduces new root level entry called bgpConfig, to describe it a bit:

- **1.1.1.1** is just an identifier, which should have format of IP and should be unique
- **100.10.0.0/24** is the network that our BGP server will advertize
- **10.10.0.10/24** is some other BGP server to connect to and to exchange the routes with
- **memif-vrf1** is interface towards 10.10.0.10/24. Note, that not neccessarily directly to that network, there might be next hops in between

[vpp-hugepages]: https://fd.io/docs/vpp/master/gettingstarted/users/configuring/hugepages.html
[vpp-startup]: https://my-vpp-docs.readthedocs.io/en/latest/gettingstarted/users/configuring/startup.html
[agentctl-link]: https://docs.ligato.io/en/latest/user-guide/agentctl/
[ligato-link]: https://ligato.io/
[conf-model]: config/STONEWORK-CONFIG.md
[vpp-cli-guide]: https://s3-docs.fd.io/vpp/22.06/cli-reference/gettingstarted/index.html
[tap-example]: ../examples/getting-started/README.md
[dpdk]: https://www.dpdk.org/
[readme]: ../README.md

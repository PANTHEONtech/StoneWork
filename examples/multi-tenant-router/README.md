[Example] StoneWork as a MultiTenant router
===========================================

Example of using the Stonework as multitenant router. Each customer/tenant has its VLAN memif 
subinterface and its own VRF table. This enables us to configure routes for each customer separately and 
achieve customer configuration separation. This example also uses NAT for separating customer subnets and 
inner subnets that could provide services to customers.  

Topology
--------

```
                                                                                                                                        
                                                                                          Stonework(router)                                                          
                                                                     +------------------------------------------------------+                                
                                                                     |                         +-----------+                |                                
   Customer1                      VSwitch(PE device)                 |             +-----------|           |                |                                
 +------------+                +----------------------+              |             |           |           +-------------+  |                                
 |            |172.16.10.10    |     +--------+VLAN 2 |              |   VLAN 2   /|   VRF2    |100.100.1.2|             |  |7.7.0.0/24      Server                   
 |            ------------------------   BD   |\      |              |172.16.10.1/ |           |           |             |  |             +----------+          
 +------------+                |     +--------+ \     |    memif     |          /  +-----------|           |             |.2|  tap     .1 |          | 
   Customer2                   |                 -------------------------------               |    NAT    |    VRF0     ------------------          |       
 +------------+                |     +--------+ /     |              |          \  +-----------|           |             |  |             |          |       
 |            ------------------------   BD   |/      |              |           \ |           |           |             |  |             +----------+                   
 |            |172.16.11.10    |     +--------+VLAN 3 |              |   VLAN 3   \|   VRF3    |100.100.1.3|             |  |                                
 +------------+                +----------------------+              |172.16.11.1  |           |           |             |  |                                
                                                                     |             +-----------|           +-------------+  |                                
                                                                     |                         +-----------+                |                                
                                                                     +------------------------------------------------------+                                
```
We simulate 2 customers by simple docker containers. They are both connected to VSwitch. The VSwitch is 
the simulation of the provider edge router(PE router/device). It forwards all traffic to Stonework, the 
multitenant router. The traffic is customer-separated in VSwitch by using dedicated bridge-domain and 
dedicated memif subinterfaces(VLAN) for each customer. The traffic passes from VSwitch to Stonework via 
high performance memif interface. At Stonework, the traffic is redirected per customer to separate 
subinterfaces (VLAN 2 and VLAN 3). Each subinterface is connected to new VRF table dedicated for that customer.
This enables us to achieve multitenancy, make route configuration separate for each customer. Then NAT will handle 
traffic to server by changing them to enter inner subnet (that means changing of the source IP addresses) where server 
is located. The NAT also redirects traffic to tap interface going to server. The NAT has different inner subnet 
IP addresses for both customers. This means that multitenancy is holding until traffic entering the server.

The traffic on the way back takes the same path, but it's destination IP address is not customer at first. It 
is the translated address in NAT (100.100.1.2 or 100.100.1.3). Then after translating it back by NAT it has 
destination IP address of customer. 

Note that customer interface and VLAN in Stonework are in the same IP subnet. It is because the path between them 
is connected through Bridge domain. That means that forwarding of packets is handled purely by L2. From the 
perspective of L3 this path looks like a single link. What goes into VLAN 2 goes out to customer1. What goes 
into VLAN 3 goes out to customer 2. The same can be said about the opposite direction.

Running the Example
-------------------

Setup the topology:
```
docker compose up --exit-code-from stonework --remove-orphans --renew-anon-volumes --force-recreate
```

The server contains simple Netcat server. We can test the traffic flow by using `wget` from customer side. 
Run this in another console: 
```
docker compose exec customer1 wget --show-progress -O /dev/null http://7.7.0.1/
```

Destroy the topology when done with it:
```
docker compose down --remove-orphans --volumes
```
Note: You can also use `Ctrl+c` in console where you run the topology setup, if you intend to setup 
the topology again. However, this won't clean everything and i.e you will see exited docker containers 
in `docker ps -a` command output.

## Packet Tracing

Packet tracing in VPP can be useful to debug problems with routing in VPP.

We will use the `vpp-probe` tool that has simpler usage and nicer output than the traditional VPP tracing 
with enabling and showing trace commands. You can install it like this:
```
go install go.ligato.io/vpp-probe@v0.2.0
```

You can trace simple ping from customer 1 to server and back with one command (it will container traces for both
VPPs, the one in the Vswitch and the one in Stonework):
```
vpp-probe --env=docker trace --print -- docker compose exec customer1 ping -c 1 -W 1 "7.7.0.1"
```
Note: If you run this as first traffic after topology setup, ping drop occur for the first 2 ping packets 
as the topology still needs to learn about their neighbors (ARP packet learning). So either run wget before
pinging or change the ping count to at least 3.

The expected output of the command above should be:
```
$ vpp-probe --env=docker trace --print -- docker compose exec customer1 ping -c 1 -W 1 "7.7.0.1"

PING 7.7.0.1 (7.7.0.1) 56(84) bytes of data.
64 bytes from 7.7.0.1: icmp_seq=1 ttl=63 time=1.36 ms

--- 7.7.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 1.360/1.360/1.360/0.000 ms

----------------------------------------------------------------------------------------------------------------------------------
 container: multi-tenant-router-stonework-1 | image: ghcr.io/pantheontech/stonework:23.06 | id: 0a97827a3299
----------------------------------------------------------------------------------------------------------------------------------
  2 packets traced
  
  # Packet 1 | ⏲  00:00:49.81400 | memif-input  ￫  tap0-tx | took 87µs | nodes 11
    - memif-input
    	memif: hw_if_index 2 next-index 4
    	  slot: ring 0
    - ethernet-input (+55µs)
    	frame: flags 0x1, hw-if-index 2, sw-if-index 2
    	IP4: 02:fe:05:99:11:e7 -> 02:fe:b0:be:c7:79 802.1q vlan 2
    - ip4-input (+61µs)
    	ICMP: 172.16.10.10 -> 7.7.0.1
    	  tos 0x00, ttl 64, length 84, checksum 0x3d48 dscp CS0 ecn NON_ECN
    	  fragment id 0x403f, flags DONT_FRAGMENT
    	ICMP echo_request checksum 0x7af1 id 86
    - ip4-sv-reassembly-feature (+65µs)
    	[not-fragmented]
    - nat-pre-in2out (+66µs)
    	in2out next_index 2 arc_next_index 11
    - nat44-ed-in2out (+68µs)
    	NAT44_IN2OUT_ED_FAST_PATH: sw_if_index 3, next index 3
    	search key local 172.16.10.10:86 remote 7.7.0.1:86 proto ICMP fib 2 thread-index 0 session-index 0
    - nat44-ed-in2out-slowpath (+71µs)
    	NAT44_IN2OUT_ED_SLOW_PATH: sw_if_index 3, next index 11, session 1, translation result 'success' via i2of
    	i2of match: saddr 172.16.10.10 sport 86 daddr 7.7.0.1 dport 86 proto ICMP fib_idx 2 rewrite: saddr 100.100.1.2 daddr 7.7.0.1 icmp-id 63327 txfib 0
    	o2if match: saddr 7.7.0.1 sport 63327 daddr 100.100.1.2 dport 63327 proto ICMP fib_idx 0 rewrite: daddr 172.16.10.10 icmp-id 86 txfib 2
    - ip4-lookup (+81µs)
    	fib 0 dpo-idx 8 flow hash: 0x00000000
    	ICMP: 100.100.1.2 -> 7.7.0.1
    	  tos 0x00, ttl 64, length 84, checksum 0x8dfc dscp CS0 ecn NON_ECN
    	  fragment id 0x403f, flags DONT_FRAGMENT
    	ICMP echo_request checksum 0x83e7 id 63327
    - ip4-rewrite (+84µs)
    	tx_sw_if_index 1 dpo-idx 8 : ipv4 via 7.7.0.1 tap0: mtu:9000 next:3 flags:[] 02fec839b87602fea887f8860800 flow hash: 0x00000000
    	00000000: 02fec839b87602fea887f886080045000054403f40003f018efc646401020707
    	00000020: 0001080083e7f75f000165f7f6640000000058880900000000001011
    - tap0-output (+85µs)
    	tap0 flags 0x10180005
    	IP4: 02:fe:a8:87:f8:86 -> 02:fe:c8:39:b8:76
    	ICMP: 100.100.1.2 -> 7.7.0.1
    	  tos 0x00, ttl 63, length 84, checksum 0x8efc dscp CS0 ecn NON_ECN
    	  fragment id 0x403f, flags DONT_FRAGMENT
    	ICMP echo_request checksum 0x83e7 id 63327
    - tap0-tx (+87µs)
    	buffer 0x97684: current data 4, length 98, buffer-pool 0, ref-count 1, trace handle 0x0
    	                vlan-1-deep l2-hdr-offset 0 l3-hdr-offset 18
    	  hdr-sz 0 l2-hdr-offset 4 l3-hdr-offset 14 l4-hdr-offset 0 l4-hdr-sz 0
    	  IP4: 02:fe:a8:87:f8:86 -> 02:fe:c8:39:b8:76
    	  ICMP: 100.100.1.2 -> 7.7.0.1
    	tos 0x00, ttl 63, length 84, checksum 0x8efc dscp CS0 ecn NON_ECN
    	fragment id 0x403f, flags DONT_FRAGMENT
    	  ICMP echo_request checksum 0x83e7 id 63327
  
  # Packet 2 | ⏲  00:00:49.81500 | virtio-input  ￫  memif1/1-output | took 12µs | nodes 9
    - virtio-input
    	virtio: hw_if_index 1 next-index 4 vring 0 len 98
    	  hdr: flags 0x00 gso_type 0x00 hdr_len 0 gso_size 0 csum_start 0 csum_offset 0 num_buffers 1
    - ethernet-input (+5µs)
    	IP4: 02:fe:c8:39:b8:76 -> 02:fe:a8:87:f8:86
    - ip4-input (+6µs)
    	ICMP: 7.7.0.1 -> 100.100.1.2
    	  tos 0x00, ttl 64, length 84, checksum 0x72ee dscp CS0 ecn NON_ECN
    	  fragment id 0x9b4d
    	ICMP echo_reply checksum 0x8be7 id 63327
    - ip4-sv-reassembly-feature (+7µs)
    	[not-fragmented]
    - nat-pre-out2in (+7µs)
    	out2in next_index 6 arc_next_index 11
    - nat44-ed-out2in (+8µs)
    	NAT44_OUT2IN_ED_FAST_PATH: sw_if_index 1, next index 11, session 1, translation result 'success' via o2if
    	i2of match: saddr 172.16.10.10 sport 86 daddr 7.7.0.1 dport 86 proto ICMP fib_idx 2 rewrite: saddr 100.100.1.2 daddr 7.7.0.1 icmp-id 63327 txfib 0
    	o2if match: saddr 7.7.0.1 sport 63327 daddr 100.100.1.2 dport 63327 proto ICMP fib_idx 0 rewrite: daddr 172.16.10.10 icmp-id 86 txfib 2
    	search key local 7.7.0.1:63327 remote 100.100.1.2:63327 proto ICMP fib 0 thread-index 0 session-index 0
    	no reason for slow path
    - ip4-lookup (+10µs)
    	fib 2 dpo-idx 9 flow hash: 0x00000000
    	ICMP: 7.7.0.1 -> 172.16.10.10
    	  tos 0x00, ttl 64, length 84, checksum 0x223a dscp CS0 ecn NON_ECN
    	  fragment id 0x9b4d
    	ICMP echo_reply checksum 0x82f1 id 86
    - ip4-rewrite (+11µs)
    	tx_sw_if_index 3 dpo-idx 9 : ipv4 via 172.16.10.10 memif1/1.2: mtu:9000 next:4 flags:[] 02fe059911e702feb0bec779810000020800 flow hash: 0x00000000
    	00000000: 02fe059911e702feb0bec779810000020800450000549b4d00003f01233a0707
    	00000020: 0001ac100a0a000082f10056000165f7f66400000000588809000000
    - memif1/1-output (+12µs)
    	memif1/1.2 flags 0x00180005
    	IP4: 02:fe:b0:be:c7:79 -> 02:fe:05:99:11:e7 802.1q vlan 2
    	ICMP: 7.7.0.1 -> 172.16.10.10
    	  tos 0x00, ttl 63, length 84, checksum 0x233a dscp CS0 ecn NON_ECN
    	  fragment id 0x9b4d
    	ICMP echo_reply checksum 0x82f1 id 86
  
----------------------------------------------------------------------------------------------------------------------------------
 container: multi-tenant-router-vswitch-1 | image: ligato/vpp-agent:v3.4.0 | id: 58a8a0871a16
----------------------------------------------------------------------------------------------------------------------------------
  2 packets traced
  
  # Packet 1 | ⏲  00:00:49.83000 | virtio-input  ￫  memif1/1-output | took 15µs | nodes 7
    - virtio-input
    	virtio: hw_if_index 1 next-index 4 vring 0 len 98
    	  hdr: flags 0x00 gso_type 0x00 hdr_len 0 gso_size 0 csum_start 0 csum_offset 0 num_buffers 1
    - ethernet-input (+7µs)
    	IP4: 02:fe:05:99:11:e7 -> 02:fe:b0:be:c7:79
    - l2-input (+9µs)
    	l2-input: sw_if_index 1 dst 02:fe:b0:be:c7:79 src 02:fe:05:99:11:e7 [l2-learn l2-fwd l2-flood ]
    - l2-learn (+11µs)
    	l2-learn: sw_if_index 1 dst 02:fe:b0:be:c7:79 src 02:fe:05:99:11:e7 bd_index 1
    - l2-fwd (+13µs)
    	l2-fwd:   sw_if_index 1 dst 02:fe:b0:be:c7:79 src 02:fe:05:99:11:e7 bd_index 1 result [0x1000000000004, 4] none
    - l2-output (+14µs)
    	l2-output: sw_if_index 4 dst 02:fe:b0:be:c7:79 src 02:fe:05:99:11:e7 data 81 00 00 02 08 00 45 00 00 54 40 3f
    - memif1/1-output (+15µs)
    	memif1/1.2
    	IP4: 02:fe:05:99:11:e7 -> 02:fe:b0:be:c7:79 802.1q vlan 2
    	ICMP: 172.16.10.10 -> 7.7.0.1
    	  tos 0x00, ttl 64, length 84, checksum 0x3d48 dscp CS0 ecn NON_ECN
    	  fragment id 0x403f, flags DONT_FRAGMENT
    	ICMP echo_request checksum 0x7af1 id 86
  
  # Packet 2 | ⏲  00:00:49.83100 | memif-input  ￫  tap0-tx | took 12µs | nodes 9
    - memif-input
    	memif: hw_if_index 3 next-index 4
    	  slot: ring 0
    - ethernet-input (+5µs)
    	IP4: 02:fe:b0:be:c7:79 -> 02:fe:05:99:11:e7 802.1q vlan 2
    - l2-input (+7µs)
    	l2-input: sw_if_index 4 dst 02:fe:05:99:11:e7 src 02:fe:b0:be:c7:79 [l2-input-vtr l2-learn l2-fwd l2-flood ]
    - l2-input-vtr (+8µs)
    	l2-input-vtr: sw_if_index 4 dst 02:fe:05:99:11:e7 src 02:fe:b0:be:c7:79 data 08 00 45 00 00 54 9b 4d 00 00 3f 01
    - l2-learn (+9µs)
    	l2-learn: sw_if_index 4 dst 02:fe:05:99:11:e7 src 02:fe:b0:be:c7:79 bd_index 1
    - l2-fwd (+10µs)
    	l2-fwd:   sw_if_index 4 dst 02:fe:05:99:11:e7 src 02:fe:b0:be:c7:79 bd_index 1 result [0x1000000000001, 1] none
    - l2-output (+10µs)
    	l2-output: sw_if_index 1 dst 02:fe:05:99:11:e7 src 02:fe:b0:be:c7:79 data 08 00 45 00 00 54 9b 4d 00 00 3f 01
    - tap0-output (+11µs)
    	tap0
    	IP4: 02:fe:b0:be:c7:79 -> 02:fe:05:99:11:e7
    	ICMP: 7.7.0.1 -> 172.16.10.10
    	  tos 0x00, ttl 63, length 84, checksum 0x233a dscp CS0 ecn NON_ECN
    	  fragment id 0x9b4d
    	ICMP echo_reply checksum 0x82f1 id 86
    - tap0-tx (+12µs)
    	buffer 0x8f04d: current data 4, length 98, buffer-pool 0, ref-count 1, trace handle 0x1
    	                l2-hdr-offset 4 l3-hdr-offset 18
    	  hdr-sz 0 l2-hdr-offset 4 l3-hdr-offset 14 l4-hdr-offset 0 l4-hdr-sz 0
    	  IP4: 02:fe:b0:be:c7:79 -> 02:fe:05:99:11:e7
    	  ICMP: 7.7.0.1 -> 172.16.10.10
    	tos 0x00, ttl 63, length 84, checksum 0x233a dscp CS0 ecn NON_ECN
    	fragment id 0x9b4d
    	  ICMP echo_reply checksum 0x82f1 id 86 
```

Note: you can alternatively use also the [vpp trace script][vpptrace-script] or just simply use the VPP CLI commands:
```
docker compose exec stonework vppctl trace add memif-input 100
docker compose exec stonework vppctl trace add virtio-input 100
docker compose exec customer1 ping -c 1 -W 1 "7.7.0.1"
docker compose exec stonework vppctl sh trace
```


[vpptrace-script]: ../../docker/vpptrace.sh
[Example] StoneWork as a NAT Gateway
====================================

Ths document describes an example configuration for StoneWork, used as a NAT
gateway. 

Here, we assume that StoneWork runs on a **bare-metal server** or a VM with
two (data-plane) DPDK-supported physical interfaces, `PCI:0000:00:08.0` and
`PCI:0000:00:09.0` connected to the private and public network, respectively.

- The **private** network is configured with the `192.168.1.0/24` IP subnet
- The GW is assigned the `192.168.1.1` IP address
- The **public** network is configured with `80.80.80.0/24` IP subnet 
- The GW is assigned `80.80.80.1` IP address

Traffic initiated from the private side is S-NATed by StoneWork to `80.80.80.1`, before it is sent to the public network.


Custom Configuration
--------------------

You will need to change the provided example configuration to adapt it to your
target environment, as described below.

- First, the PCI addresses of the physical interfaces should be listed in the
  attached `vpp-startup.conf`, section `dpdk { ... }`.

  To obtain the PCI addresses, run `lshw` command:
  ```
  $ sudo lshw -class network -businfo
  Bus info          Device     Class      Description
  ===================================================
  pci@0000:00:08.0  eth0       network    82540EM Gigabit Ethernet Controller
  pci@0000:00:09.0  eth1       network    82540EM Gigabit Ethernet Controller
  ```
  Open the attached `vpp-startup.conf` and update the PCI addresses from the
example to the actual values.


- Next, the IP addresses of the interfaces need to be edited in
  `./config/day0-config.yaml` to match the target environment. Replace
  `80.80.80.1/24` with the actual external IP subnet and `192.168.1.1/24`
  with the actual IP subnet of the private network. 
  
  Next, the hop of the default route has to be updated accordingly.


- Finally, open `config/add-nat-config.yaml` and change the IP pool used by S-NAT
  from `80.80.80.1` (`/32`) to the actual IP address of the StoneWork
  deployment on the public side.

The existing configuration can be further extended for more complex use-cases.
For example, different private networks or different kinds of the traffic could
be separated into different VLANs, each with its own sub-interface and IP subnet.


Running the Example
-------------------

Once the configuration is updated to match the actual environment, deploy
StoneWork using:
```
$ docker compose up -d
```

Initially, only physical interfaces should be configured.

To verify this, run:
```
$ docker compose exec stonework vppctl show interface address
gbe-private-net (up):
  L3 192.168.1.1/24
gbe-public-net (up):
  L3 80.80.80.1/24
local0 (dn):
```

Explore the CLI provided by StoneWork:
```
$ docker compose exec stonework agentctl --help
```

## Add & Enable NAT Config

For example, to **add and enable NAT configuration**, run:
```
$ docker compose exec stonework agentctl config update /etc/stonework/config/add-nat-config.yaml
```
`add-nat-config.yaml` is mounted from `./config/add-nat-config.yaml`.

Feel free to experiment and make some configuration changes on your own. 

Once applied, the NAT gateway should be operational and clients from the private
network should be able to access servers in the public network(s).

## Observing NAT / Packet Tracing

Observe the NAT in process using packet tracing on VPP (only `SYN` packet shown):
```
$ docker compose exec stonework vpptrace.sh -i dpdk

00:04:40:604246: dpdk-input
  gbe-private-net rx queue 0
  buffer 0x9b52e: current data 0, length 74, buffer-pool 0, ref-count 1, totlen-nifb 0, trace handle 0x3
                  ext-hdr-valid 
                  l4-cksum-computed l4-cksum-correct 
  PKT MBUF: port 0, nb_segs 1, pkt_len 74
    buf_len 2176, data_len 74, ol_flags 0x0, data_off 128, phys_addr 0xa22d4c00
    packet_type 0x0 l2_len 0 l3_len 0 outer_l2_len 0 outer_l3_len 0
    rss 0x0 fdir.hi 0x0 fdir.lo 0x0
  IP4: 08:00:27:ae:0a:e2 -> 08:00:27:d4:65:b6
  TCP: 192.168.1.2 -> 80.80.80.2
    tos 0x00, ttl 64, length 60, checksum 0x0947 dscp CS0 ecn NON_ECN
    fragment id 0xcf78, flags DONT_FRAGMENT
  TCP: 34356 -> 8080
    seq. 0x825b6b36 ack 0x00000000
    flags 0x02 SYN, tcp header: 40 bytes
    window 29200, checksum 0x0cca
00:04:40:604285: ethernet-input
  frame: flags 0x3, hw-if-index 1, sw-if-index 1
  IP4: 08:00:27:ae:0a:e2 -> 08:00:27:d4:65:b6
00:04:40:604292: ip4-input-no-checksum
  TCP: 192.168.1.2 -> 80.80.80.2
    tos 0x00, ttl 64, length 60, checksum 0x0947 dscp CS0 ecn NON_ECN
    fragment id 0xcf78, flags DONT_FRAGMENT
  TCP: 34356 -> 8080
    seq. 0x825b6b36 ack 0x00000000
    flags 0x02 SYN, tcp header: 40 bytes
    window 29200, checksum 0x0cca
00:04:40:604295: ip4-sv-reassembly-feature
  [not-fragmented]
00:04:40:604296: nat44-in2out
  NAT44_IN2OUT_FAST_PATH: sw_if_index 1, next index 3, session -1
00:04:40:604299: nat44-in2out-slowpath
  NAT44_IN2OUT_SLOW_PATH: sw_if_index 1, next index 0, session -1
00:04:40:604301: ip4-lookup
  fib 0 dpo-idx 2 flow hash: 0x00000000
  TCP: 192.168.1.2 -> 80.80.80.2
    tos 0x00, ttl 64, length 60, checksum 0x0947 dscp CS0 ecn NON_ECN
    fragment id 0xcf78, flags DONT_FRAGMENT
  TCP: 34356 -> 8080
    seq. 0x825b6b36 ack 0x00000000
    flags 0x02 SYN, tcp header: 40 bytes
    window 29200, checksum 0x0cca
00:04:40:604305: ip4-rewrite
  tx_sw_if_index 2 dpo-idx 2 : ipv4 via 80.80.80.2 gbe-public-net: mtu:9000 next:5 080027f5da05080027488c2f0800 flow hash: 0x00000000
  00000000: 080027f5da05080027488c2f08004500003ccf7840003f060a47c0a801025050
  00000020: 500286341f90825b6b3600000000a00272100cca0000020405b40402
00:04:40:604308: ip4-sv-reassembly-output-feature
  [not-fragmented]
00:04:40:604318: nat44-in2out-output
  NAT44_IN2OUT_FAST_PATH: sw_if_index 1, next index 3, session -1
00:04:40:604318: nat44-in2out-output-slowpath
  NAT44_IN2OUT_SLOW_PATH: sw_if_index 1, next index 0, session 0
00:04:40:604377: gbe-public-net-output
  gbe-public-net 
  IP4: 08:00:27:48:8c:2f -> 08:00:27:f5:da:05
  TCP: 80.80.80.1 -> 80.80.80.2
    tos 0x00, ttl 63, length 60, checksum 0x2ba0 dscp CS0 ecn NON_ECN
    fragment id 0xcf78, flags DONT_FRAGMENT
  TCP: 63327 -> 8080
    seq. 0x825b6b36 ack 0x00000000
    flags 0x02 SYN, tcp header: 40 bytes
    window 29200, checksum 0xbcf7
00:04:40:604379: gbe-public-net-tx
  gbe-public-net tx queue 0
  buffer 0x9b52e: current data 0, length 74, buffer-pool 0, ref-count 1, totlen-nifb 0, trace handle 0x3
                  ext-hdr-valid 
                  l4-cksum-computed l4-cksum-correct natted l2-hdr-offset 0 l3-hdr-offset 14 
  PKT MBUF: port 0, nb_segs 1, pkt_len 74
    buf_len 2176, data_len 74, ol_flags 0x0, data_off 128, phys_addr 0xa22d4c00
    packet_type 0x0 l2_len 0 l3_len 0 outer_l2_len 0 outer_l3_len 0
    rss 0x0 fdir.hi 0x0 fdir.lo 0x0
  IP4: 08:00:27:48:8c:2f -> 08:00:27:f5:da:05
  TCP: 80.80.80.1 -> 80.80.80.2
    tos 0x00, ttl 63, length 60, checksum 0x2ba0 dscp CS0 ecn NON_ECN
    fragment id 0xcf78, flags DONT_FRAGMENT
  TCP: 63327 -> 8080
    seq. 0x825b6b36 ack 0x00000000
    flags 0x02 SYN, tcp header: 40 bytes
    window 29200, checksum 0xbcf7

...
```

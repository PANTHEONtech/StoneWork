This document describes BFD plugin, its API, features and data representation.

BFD UDP
-------

BFD plugin supports VPP single hop UDP based Bidirectional Forwarding Detection based on RFC 5880 and RFC 5881.
BFD configuration requires source interface, local IP address and peer IP address. In order to successfully establish
BFD session, both sides must be configured.

BFD local configuration key ID and BFD key ID (as carried in BFD control frames) is generated
and assigned by the plugin itself.

Other configurable parameters:

* **Desired minimum TX interval** named `min_tx_interval` in the API is an interval in microseconds where BFD
  transmits control packets. May not be a zero value.  
* **Required minimum RX interval** named `min_rx_interval` in the API is an interval in microseconds between
  received control packets. Can be set to zero - in that case the system does not expect control packets
  from the remote router.
* **Detect multiplier** is a value multiplying the negotiated transmit, defining final detection time.

The BFD plugin does not support authentication.

A server key/data representation based on the CNF Protobuf model:
```
Key: /vnf-agent/<microservice_label>/config/vpp/bfd/v1/<interface-name>/peer/<peer-ip>
Data: 
{
    "if_name": "<interface-name>", 
    "local_ip": "<local-interface-ip>"
    "peer_ip": "<peer-ip>"
    "min_tx_interval": <desired-min-tx-interval>
    "min_rx_interval": <required-min-rx-interval>
    "detect_multiplier": <value>
}
``` 

The VPP supports CLI commands to show configured BFD sessions:
```
vpp# sh bfd sessions   
   Index               Property                  Local value         Remote value    
     0     IPv4 address                                10.10.0.5           10.10.0.10
           Session state                                      Up                   Up
           Diagnostic code                         No Diagnostic        No Diagnostic
           Detect multiplier                                   1                    1
           Required Min Rx Interval (usec)               1000000              1000000
           Desired Min Tx Interval (usec)                1000000              1000000
           Transmit interval                             1000000
           Last control frame tx                        .56s ago
           Last control frame rx                        .39s ago
           Min Echo Rx Interval (usec)                         1                    1
           Demand mode                                        no                   no
           Poll state                        BFD_POLL_NOT_NEEDED
     1     IPv4 address                                20.10.0.5           20.10.0.10
           Session state                                      Up                   Up
           Diagnostic code                         No Diagnostic        No Diagnostic
           Detect multiplier                                   1                    1
           Required Min Rx Interval (usec)               1000000              1000000
           Desired Min Tx Interval (usec)                1000000              1000000
           Transmit interval                             1000000
           Last control frame tx                        .34s ago
           Last control frame rx                        .39s ago
           Min Echo Rx Interval (usec)                         1                    1
           Demand mode                                        no                   no
           Poll state                        BFD_POLL_NOT_NEEDED
Number of configured BFD sessions: 2
vpp# 
```

The example above shows two BFD sessions with local and peer addresses,
session state and configured intervals.


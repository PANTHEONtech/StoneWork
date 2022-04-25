Configuration Model of StoneWork
====================================

## Table of Contents

- *stonework-root.proto*
    - [Root](#stonework.Root)
    - [Root.LinuxConfig](#stonework.Root.LinuxConfig)
    - [Root.NetallocConfig](#stonework.Root.NetallocConfig)
    - [Root.VppConfig](#stonework.Root.VppConfig)
  
- *ligato/vpp/wireguard/wireguard.proto*
    - [Peer](#ligato.vpp.wireguard.Peer)
  
- *ligato/vpp/vpp.proto*
    - [ConfigData](#ligato.vpp.ConfigData)
    - [Notification](#ligato.vpp.Notification)
    - [Stats](#ligato.vpp.Stats)
  
- *ligato/vpp/stn/stn.proto*
    - [Rule](#ligato.vpp.stn.Rule)
  
- *ligato/vpp/srv6/srv6.proto*
    - [LocalSID](#ligato.vpp.srv6.LocalSID)
    - [LocalSID.End](#ligato.vpp.srv6.LocalSID.End)
    - [LocalSID.EndAD](#ligato.vpp.srv6.LocalSID.EndAD)
    - [LocalSID.EndDT4](#ligato.vpp.srv6.LocalSID.EndDT4)
    - [LocalSID.EndDT6](#ligato.vpp.srv6.LocalSID.EndDT6)
    - [LocalSID.EndDX2](#ligato.vpp.srv6.LocalSID.EndDX2)
    - [LocalSID.EndDX4](#ligato.vpp.srv6.LocalSID.EndDX4)
    - [LocalSID.EndDX6](#ligato.vpp.srv6.LocalSID.EndDX6)
    - [LocalSID.EndT](#ligato.vpp.srv6.LocalSID.EndT)
    - [LocalSID.EndX](#ligato.vpp.srv6.LocalSID.EndX)
    - [Policy](#ligato.vpp.srv6.Policy)
    - [Policy.SegmentList](#ligato.vpp.srv6.Policy.SegmentList)
    - [SRv6Global](#ligato.vpp.srv6.SRv6Global)
    - [Steering](#ligato.vpp.srv6.Steering)
    - [Steering.L2Traffic](#ligato.vpp.srv6.Steering.L2Traffic)
    - [Steering.L3Traffic](#ligato.vpp.srv6.Steering.L3Traffic)
  
- *ligato/vpp/punt/punt.proto*
    - [Exception](#ligato.vpp.punt.Exception)
    - [IPRedirect](#ligato.vpp.punt.IPRedirect)
    - [Reason](#ligato.vpp.punt.Reason)
    - [ToHost](#ligato.vpp.punt.ToHost)
  
    - [L3Protocol](#ligato.vpp.punt.L3Protocol)
    - [L4Protocol](#ligato.vpp.punt.L4Protocol)
  
- *ligato/vpp/nat/nat.proto*
    - [DNat44](#ligato.vpp.nat.DNat44)
    - [DNat44.IdentityMapping](#ligato.vpp.nat.DNat44.IdentityMapping)
    - [DNat44.StaticMapping](#ligato.vpp.nat.DNat44.StaticMapping)
    - [DNat44.StaticMapping.LocalIP](#ligato.vpp.nat.DNat44.StaticMapping.LocalIP)
    - [Nat44AddressPool](#ligato.vpp.nat.Nat44AddressPool)
    - [Nat44Global](#ligato.vpp.nat.Nat44Global)
    - [Nat44Global.Address](#ligato.vpp.nat.Nat44Global.Address)
    - [Nat44Global.Interface](#ligato.vpp.nat.Nat44Global.Interface)
    - [Nat44Interface](#ligato.vpp.nat.Nat44Interface)
    - [VirtualReassembly](#ligato.vpp.nat.VirtualReassembly)
  
    - [DNat44.Protocol](#ligato.vpp.nat.DNat44.Protocol)
    - [DNat44.StaticMapping.TwiceNatMode](#ligato.vpp.nat.DNat44.StaticMapping.TwiceNatMode)
  
- *ligato/vpp/l3/vrrp.proto*
    - [VRRPEntry](#ligato.vpp.l3.VRRPEntry)
  
- *ligato/vpp/l3/vrf.proto*
    - [VrfTable](#ligato.vpp.l3.VrfTable)
    - [VrfTable.FlowHashSettings](#ligato.vpp.l3.VrfTable.FlowHashSettings)
  
    - [VrfTable.Protocol](#ligato.vpp.l3.VrfTable.Protocol)
  
- *ligato/vpp/l3/teib.proto*
    - [TeibEntry](#ligato.vpp.l3.TeibEntry)
  
- *ligato/vpp/l3/route.proto*
    - [Route](#ligato.vpp.l3.Route)
  
    - [Route.RouteType](#ligato.vpp.l3.Route.RouteType)
  
- *ligato/vpp/l3/l3xc.proto*
    - [L3XConnect](#ligato.vpp.l3.L3XConnect)
    - [L3XConnect.Path](#ligato.vpp.l3.L3XConnect.Path)
  
    - [L3XConnect.Protocol](#ligato.vpp.l3.L3XConnect.Protocol)
  
- *ligato/vpp/l3/l3.proto*
    - [DHCPProxy](#ligato.vpp.l3.DHCPProxy)
    - [DHCPProxy.DHCPServer](#ligato.vpp.l3.DHCPProxy.DHCPServer)
    - [IPScanNeighbor](#ligato.vpp.l3.IPScanNeighbor)
    - [ProxyARP](#ligato.vpp.l3.ProxyARP)
    - [ProxyARP.Interface](#ligato.vpp.l3.ProxyARP.Interface)
    - [ProxyARP.Range](#ligato.vpp.l3.ProxyARP.Range)
  
    - [IPScanNeighbor.Mode](#ligato.vpp.l3.IPScanNeighbor.Mode)
  
- *ligato/vpp/l3/arp.proto*
    - [ARPEntry](#ligato.vpp.l3.ARPEntry)
  
- *ligato/vpp/l2/xconnect.proto*
    - [XConnectPair](#ligato.vpp.l2.XConnectPair)
  
- *ligato/vpp/l2/fib.proto*
    - [FIBEntry](#ligato.vpp.l2.FIBEntry)
  
    - [FIBEntry.Action](#ligato.vpp.l2.FIBEntry.Action)
  
- *ligato/vpp/l2/bridge_domain.proto*
    - [BridgeDomain](#ligato.vpp.l2.BridgeDomain)
    - [BridgeDomain.ArpTerminationEntry](#ligato.vpp.l2.BridgeDomain.ArpTerminationEntry)
    - [BridgeDomain.Interface](#ligato.vpp.l2.BridgeDomain.Interface)
  
- *ligato/vpp/ipsec/ipsec.proto*
    - [SecurityAssociation](#ligato.vpp.ipsec.SecurityAssociation)
    - [SecurityPolicy](#ligato.vpp.ipsec.SecurityPolicy)
    - [SecurityPolicyDatabase](#ligato.vpp.ipsec.SecurityPolicyDatabase)
    - [SecurityPolicyDatabase.Interface](#ligato.vpp.ipsec.SecurityPolicyDatabase.Interface)
    - [SecurityPolicyDatabase.PolicyEntry](#ligato.vpp.ipsec.SecurityPolicyDatabase.PolicyEntry)
    - [TunnelProtection](#ligato.vpp.ipsec.TunnelProtection)
  
    - [CryptoAlg](#ligato.vpp.ipsec.CryptoAlg)
    - [IntegAlg](#ligato.vpp.ipsec.IntegAlg)
    - [SecurityAssociation.IPSecProtocol](#ligato.vpp.ipsec.SecurityAssociation.IPSecProtocol)
    - [SecurityPolicy.Action](#ligato.vpp.ipsec.SecurityPolicy.Action)
    - [SecurityPolicyDatabase.PolicyEntry.Action](#ligato.vpp.ipsec.SecurityPolicyDatabase.PolicyEntry.Action)
  
- *ligato/vpp/ipfix/ipfix.proto*
    - [IPFIX](#ligato.vpp.ipfix.IPFIX)
    - [IPFIX.Collector](#ligato.vpp.ipfix.IPFIX.Collector)
  
- *ligato/vpp/ipfix/flowprobe.proto*
    - [FlowProbeFeature](#ligato.vpp.ipfix.FlowProbeFeature)
    - [FlowProbeParams](#ligato.vpp.ipfix.FlowProbeParams)
  
- *ligato/vpp/interfaces/interface.proto*
    - [AfpacketLink](#ligato.vpp.interfaces.AfpacketLink)
    - [BondLink](#ligato.vpp.interfaces.BondLink)
    - [BondLink.BondedInterface](#ligato.vpp.interfaces.BondLink.BondedInterface)
    - [GreLink](#ligato.vpp.interfaces.GreLink)
    - [GtpuLink](#ligato.vpp.interfaces.GtpuLink)
    - [IPIPLink](#ligato.vpp.interfaces.IPIPLink)
    - [IPSecLink](#ligato.vpp.interfaces.IPSecLink)
    - [Interface](#ligato.vpp.interfaces.Interface)
    - [Interface.IP6ND](#ligato.vpp.interfaces.Interface.IP6ND)
    - [Interface.RxMode](#ligato.vpp.interfaces.Interface.RxMode)
    - [Interface.RxPlacement](#ligato.vpp.interfaces.Interface.RxPlacement)
    - [Interface.Unnumbered](#ligato.vpp.interfaces.Interface.Unnumbered)
    - [MemifLink](#ligato.vpp.interfaces.MemifLink)
    - [RDMALink](#ligato.vpp.interfaces.RDMALink)
    - [SubInterface](#ligato.vpp.interfaces.SubInterface)
    - [TapLink](#ligato.vpp.interfaces.TapLink)
    - [VmxNet3Link](#ligato.vpp.interfaces.VmxNet3Link)
    - [VxlanLink](#ligato.vpp.interfaces.VxlanLink)
    - [VxlanLink.Gpe](#ligato.vpp.interfaces.VxlanLink.Gpe)
    - [WireguardLink](#ligato.vpp.interfaces.WireguardLink)
  
    - [BondLink.LoadBalance](#ligato.vpp.interfaces.BondLink.LoadBalance)
    - [BondLink.Mode](#ligato.vpp.interfaces.BondLink.Mode)
    - [GreLink.Type](#ligato.vpp.interfaces.GreLink.Type)
    - [GtpuLink.NextNode](#ligato.vpp.interfaces.GtpuLink.NextNode)
    - [IPIPLink.Mode](#ligato.vpp.interfaces.IPIPLink.Mode)
    - [IPSecLink.Mode](#ligato.vpp.interfaces.IPSecLink.Mode)
    - [Interface.RxMode.Type](#ligato.vpp.interfaces.Interface.RxMode.Type)
    - [Interface.Type](#ligato.vpp.interfaces.Interface.Type)
    - [MemifLink.MemifMode](#ligato.vpp.interfaces.MemifLink.MemifMode)
    - [RDMALink.Mode](#ligato.vpp.interfaces.RDMALink.Mode)
    - [SubInterface.TagRewriteOptions](#ligato.vpp.interfaces.SubInterface.TagRewriteOptions)
    - [VxlanLink.Gpe.Protocol](#ligato.vpp.interfaces.VxlanLink.Gpe.Protocol)
  
- *ligato/vpp/interfaces/state.proto*
    - [InterfaceNotification](#ligato.vpp.interfaces.InterfaceNotification)
    - [InterfaceState](#ligato.vpp.interfaces.InterfaceState)
    - [InterfaceState.Statistics](#ligato.vpp.interfaces.InterfaceState.Statistics)
    - [InterfaceStats](#ligato.vpp.interfaces.InterfaceStats)
    - [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter)
  
    - [InterfaceNotification.NotifType](#ligato.vpp.interfaces.InterfaceNotification.NotifType)
    - [InterfaceState.Duplex](#ligato.vpp.interfaces.InterfaceState.Duplex)
    - [InterfaceState.Status](#ligato.vpp.interfaces.InterfaceState.Status)
  
- *ligato/vpp/interfaces/span.proto*
    - [Span](#ligato.vpp.interfaces.Span)
  
    - [Span.Direction](#ligato.vpp.interfaces.Span.Direction)
  
- *ligato/vpp/interfaces/dhcp.proto*
    - [DHCPLease](#ligato.vpp.interfaces.DHCPLease)
  
- *ligato/vpp/dns/dns.proto*
    - [DNSCache](#ligato.vpp.dns.DNSCache)
  
- *ligato/vpp/acl/acl.proto*
    - [ACL](#ligato.vpp.acl.ACL)
    - [ACL.Interfaces](#ligato.vpp.acl.ACL.Interfaces)
    - [ACL.Rule](#ligato.vpp.acl.ACL.Rule)
    - [ACL.Rule.IpRule](#ligato.vpp.acl.ACL.Rule.IpRule)
    - [ACL.Rule.IpRule.Icmp](#ligato.vpp.acl.ACL.Rule.IpRule.Icmp)
    - [ACL.Rule.IpRule.Icmp.Range](#ligato.vpp.acl.ACL.Rule.IpRule.Icmp.Range)
    - [ACL.Rule.IpRule.Ip](#ligato.vpp.acl.ACL.Rule.IpRule.Ip)
    - [ACL.Rule.IpRule.PortRange](#ligato.vpp.acl.ACL.Rule.IpRule.PortRange)
    - [ACL.Rule.IpRule.Tcp](#ligato.vpp.acl.ACL.Rule.IpRule.Tcp)
    - [ACL.Rule.IpRule.Udp](#ligato.vpp.acl.ACL.Rule.IpRule.Udp)
    - [ACL.Rule.MacIpRule](#ligato.vpp.acl.ACL.Rule.MacIpRule)
  
    - [ACL.Rule.Action](#ligato.vpp.acl.ACL.Rule.Action)
  
- *ligato/vpp/abf/abf.proto*
    - [ABF](#ligato.vpp.abf.ABF)
    - [ABF.AttachedInterface](#ligato.vpp.abf.ABF.AttachedInterface)
    - [ABF.ForwardingPath](#ligato.vpp.abf.ABF.ForwardingPath)
  
- *ligato/netalloc/netalloc.proto*
    - [ConfigData](#ligato.netalloc.ConfigData)
    - [IPAllocation](#ligato.netalloc.IPAllocation)
  
    - [IPAddressForm](#ligato.netalloc.IPAddressForm)
    - [IPAddressSource](#ligato.netalloc.IPAddressSource)
  
- *ligato/linux/punt/punt.proto*
    - [PortBased](#ligato.linux.punt.PortBased)
    - [Proxy](#ligato.linux.punt.Proxy)
    - [SocketBased](#ligato.linux.punt.SocketBased)
  
    - [PortBased.L3Protocol](#ligato.linux.punt.PortBased.L3Protocol)
    - [PortBased.L4Protocol](#ligato.linux.punt.PortBased.L4Protocol)
  
- *ligato/linux/namespace/namespace.proto*
    - [NetNamespace](#ligato.linux.namespace.NetNamespace)
  
    - [NetNamespace.ReferenceType](#ligato.linux.namespace.NetNamespace.ReferenceType)
  
- *ligato/linux/linux.proto*
    - [ConfigData](#ligato.linux.ConfigData)
    - [Notification](#ligato.linux.Notification)
  
- *ligato/linux/l3/route.proto*
    - [Route](#ligato.linux.l3.Route)
  
    - [Route.Scope](#ligato.linux.l3.Route.Scope)
  
- *ligato/linux/l3/arp.proto*
    - [ARPEntry](#ligato.linux.l3.ARPEntry)
  
- *ligato/linux/iptables/iptables.proto*
    - [RuleChain](#ligato.linux.iptables.RuleChain)
  
    - [RuleChain.ChainType](#ligato.linux.iptables.RuleChain.ChainType)
    - [RuleChain.Policy](#ligato.linux.iptables.RuleChain.Policy)
    - [RuleChain.Protocol](#ligato.linux.iptables.RuleChain.Protocol)
    - [RuleChain.Table](#ligato.linux.iptables.RuleChain.Table)
  
- *ligato/linux/interfaces/state.proto*
    - [InterfaceNotification](#ligato.linux.interfaces.InterfaceNotification)
    - [InterfaceState](#ligato.linux.interfaces.InterfaceState)
    - [InterfaceState.Statistics](#ligato.linux.interfaces.InterfaceState.Statistics)
  
    - [InterfaceNotification.NotifType](#ligato.linux.interfaces.InterfaceNotification.NotifType)
    - [InterfaceState.Status](#ligato.linux.interfaces.InterfaceState.Status)
  
- *ligato/linux/interfaces/interface.proto*
    - [Interface](#ligato.linux.interfaces.Interface)
    - [TapLink](#ligato.linux.interfaces.TapLink)
    - [VethLink](#ligato.linux.interfaces.VethLink)
    - [VrfDevLink](#ligato.linux.interfaces.VrfDevLink)
  
    - [Interface.Type](#ligato.linux.interfaces.Interface.Type)
    - [VethLink.ChecksumOffloading](#ligato.linux.interfaces.VethLink.ChecksumOffloading)
  
- *ligato/kvscheduler/value_status.proto*
    - [BaseValueStatus](#ligato.kvscheduler.BaseValueStatus)
    - [ValueStatus](#ligato.kvscheduler.ValueStatus)
  
    - [TxnOperation](#ligato.kvscheduler.TxnOperation)
    - [ValueState](#ligato.kvscheduler.ValueState)
  
- *ligato/govppmux/metrics.proto*
    - [Metrics](#ligato.govppmux.Metrics)
  
- *ligato/generic/options.proto*
    - [File-level Extensions](#ligato/generic/options.proto-extensions)
  
- *ligato/generic/model.proto*
    - [ModelDetail](#ligato.generic.ModelDetail)
    - [ModelDetail.Option](#ligato.generic.ModelDetail.Option)
    - [ModelSpec](#ligato.generic.ModelSpec)
  
- *ligato/generic/meta.proto*
    - [KnownModelsRequest](#ligato.generic.KnownModelsRequest)
    - [KnownModelsResponse](#ligato.generic.KnownModelsResponse)
    - [ProtoFileDescriptorRequest](#ligato.generic.ProtoFileDescriptorRequest)
    - [ProtoFileDescriptorResponse](#ligato.generic.ProtoFileDescriptorResponse)
  
    - [MetaService](#ligato.generic.MetaService)
  
- *ligato/generic/manager.proto*
    - [ConfigItem](#ligato.generic.ConfigItem)
    - [ConfigItem.LabelsEntry](#ligato.generic.ConfigItem.LabelsEntry)
    - [Data](#ligato.generic.Data)
    - [DumpStateRequest](#ligato.generic.DumpStateRequest)
    - [DumpStateResponse](#ligato.generic.DumpStateResponse)
    - [GetConfigRequest](#ligato.generic.GetConfigRequest)
    - [GetConfigResponse](#ligato.generic.GetConfigResponse)
    - [Item](#ligato.generic.Item)
    - [Item.ID](#ligato.generic.Item.ID)
    - [ItemStatus](#ligato.generic.ItemStatus)
    - [Notification](#ligato.generic.Notification)
    - [SetConfigRequest](#ligato.generic.SetConfigRequest)
    - [SetConfigResponse](#ligato.generic.SetConfigResponse)
    - [StateItem](#ligato.generic.StateItem)
    - [StateItem.MetadataEntry](#ligato.generic.StateItem.MetadataEntry)
    - [SubscribeRequest](#ligato.generic.SubscribeRequest)
    - [SubscribeResponse](#ligato.generic.SubscribeResponse)
    - [Subscription](#ligato.generic.Subscription)
    - [UpdateItem](#ligato.generic.UpdateItem)
    - [UpdateItem.LabelsEntry](#ligato.generic.UpdateItem.LabelsEntry)
    - [UpdateResult](#ligato.generic.UpdateResult)
  
    - [UpdateResult.Operation](#ligato.generic.UpdateResult.Operation)
  
    - [ManagerService](#ligato.generic.ManagerService)
  
- *ligato/configurator/statspoller.proto*
    - [PollStatsRequest](#ligato.configurator.PollStatsRequest)
    - [PollStatsResponse](#ligato.configurator.PollStatsResponse)
    - [Stats](#ligato.configurator.Stats)
  
    - [StatsPollerService](#ligato.configurator.StatsPollerService)
  
- *ligato/configurator/configurator.proto*
    - [Config](#ligato.configurator.Config)
    - [DeleteRequest](#ligato.configurator.DeleteRequest)
    - [DeleteResponse](#ligato.configurator.DeleteResponse)
    - [DumpRequest](#ligato.configurator.DumpRequest)
    - [DumpResponse](#ligato.configurator.DumpResponse)
    - [GetRequest](#ligato.configurator.GetRequest)
    - [GetResponse](#ligato.configurator.GetResponse)
    - [Notification](#ligato.configurator.Notification)
    - [NotifyRequest](#ligato.configurator.NotifyRequest)
    - [NotifyResponse](#ligato.configurator.NotifyResponse)
    - [UpdateRequest](#ligato.configurator.UpdateRequest)
    - [UpdateResponse](#ligato.configurator.UpdateResponse)
  
    - [ConfiguratorService](#ligato.configurator.ConfiguratorService)
  
- *ligato/annotations.proto*
    - [LigatoOptions](#ligato.LigatoOptions)
    - [LigatoOptions.IntRange](#ligato.LigatoOptions.IntRange)
  
    - [LigatoOptions.Type](#ligato.LigatoOptions.Type)
  
    - [File-level Extensions](#ligato/annotations.proto-extensions)
  
- *nat64/nat64.proto*
    - [Nat64AddressPool](#nat64.Nat64AddressPool)
    - [Nat64IPv6Prefix](#nat64.Nat64IPv6Prefix)
    - [Nat64Interface](#nat64.Nat64Interface)
    - [Nat64StaticBIB](#nat64.Nat64StaticBIB)
  
    - [Nat64Interface.Type](#nat64.Nat64Interface.Type)
    - [Nat64StaticBIB.Protocol](#nat64.Nat64StaticBIB.Protocol)
  
- *isisx/isisx.proto*
    - [ISISXConnection](#vpp.isisx.ISISXConnection)
  
- *bfd/bfd.proto*
    - [BFD](#bfd.BFD)
    - [BFDEvent](#bfd.BFDEvent)
    - [WatchBFDEventsRequest](#bfd.WatchBFDEventsRequest)
  
    - [BFDEvent.SessionState](#bfd.BFDEvent.SessionState)
  
    - [BFDWatcher](#bfd.BFDWatcher)
  
- *abx/abx.proto*
    - [ABX](#vpp.abx.ABX)
    - [ABX.AttachedInterface](#vpp.abx.ABX.AttachedInterface)
  



<a name="stonework-root.proto"></a>

## stonework-root.proto
Proto file with the configuration model of StoneWork.


<a name="stonework.Root"></a>

### Root
Configuration root wrapping all models supported by StoneWork.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| linuxConfig | [Root.LinuxConfig](#stonework.Root.LinuxConfig) |  |  |
| netallocConfig | [Root.NetallocConfig](#stonework.Root.NetallocConfig) |  |  |
| vppConfig | [Root.VppConfig](#stonework.Root.VppConfig) |  |  |






<a name="stonework.Root.LinuxConfig"></a>

### Root.LinuxConfig



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| RuleChain_list | [ligato.linux.iptables.RuleChain](#ligato.linux.iptables.RuleChain) | repeated |  |
| arp_entries | [ligato.linux.l3.ARPEntry](#ligato.linux.l3.ARPEntry) | repeated |  |
| interfaces | [ligato.linux.interfaces.Interface](#ligato.linux.interfaces.Interface) | repeated |  |
| routes | [ligato.linux.l3.Route](#ligato.linux.l3.Route) | repeated |  |






<a name="stonework.Root.NetallocConfig"></a>

### Root.NetallocConfig



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ip_addresses | [ligato.netalloc.IPAllocation](#ligato.netalloc.IPAllocation) | repeated |  |






<a name="stonework.Root.VppConfig"></a>

### Root.VppConfig



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ABX_list | [vpp.abx.ABX](#vpp.abx.ABX) | repeated |  |
| BFD_list | [bfd.BFD](#bfd.BFD) | repeated |  |
| DNSCache | [ligato.vpp.dns.DNSCache](#ligato.vpp.dns.DNSCache) |  |  |
| ISISXConnection_list | [vpp.isisx.ISISXConnection](#vpp.isisx.ISISXConnection) | repeated |  |
| Nat64AddressPool_list | [nat64.Nat64AddressPool](#nat64.Nat64AddressPool) | repeated |  |
| Nat64IPv6Prefix_list | [nat64.Nat64IPv6Prefix](#nat64.Nat64IPv6Prefix) | repeated |  |
| Nat64Interface_list | [nat64.Nat64Interface](#nat64.Nat64Interface) | repeated |  |
| Nat64StaticBIB_list | [nat64.Nat64StaticBIB](#nat64.Nat64StaticBIB) | repeated |  |
| Rule_list | [ligato.vpp.stn.Rule](#ligato.vpp.stn.Rule) | repeated |  |
| VRRPEntry_list | [ligato.vpp.l3.VRRPEntry](#ligato.vpp.l3.VRRPEntry) | repeated |  |
| abfs | [ligato.vpp.abf.ABF](#ligato.vpp.abf.ABF) | repeated |  |
| acls | [ligato.vpp.acl.ACL](#ligato.vpp.acl.ACL) | repeated |  |
| arps | [ligato.vpp.l3.ARPEntry](#ligato.vpp.l3.ARPEntry) | repeated |  |
| bridge_domains | [ligato.vpp.l2.BridgeDomain](#ligato.vpp.l2.BridgeDomain) | repeated |  |
| dhcp_proxies | [ligato.vpp.l3.DHCPProxy](#ligato.vpp.l3.DHCPProxy) | repeated |  |
| dnat44s | [ligato.vpp.nat.DNat44](#ligato.vpp.nat.DNat44) | repeated |  |
| fibs | [ligato.vpp.l2.FIBEntry](#ligato.vpp.l2.FIBEntry) | repeated |  |
| interfaces | [ligato.vpp.interfaces.Interface](#ligato.vpp.interfaces.Interface) | repeated |  |
| ipfix_flowprobe_params | [ligato.vpp.ipfix.FlowProbeParams](#ligato.vpp.ipfix.FlowProbeParams) |  |  |
| ipfix_flowprobes | [ligato.vpp.ipfix.FlowProbeFeature](#ligato.vpp.ipfix.FlowProbeFeature) | repeated |  |
| ipfix_global | [ligato.vpp.ipfix.IPFIX](#ligato.vpp.ipfix.IPFIX) |  |  |
| ipscan_neighbor | [ligato.vpp.l3.IPScanNeighbor](#ligato.vpp.l3.IPScanNeighbor) |  |  |
| ipsec_sas | [ligato.vpp.ipsec.SecurityAssociation](#ligato.vpp.ipsec.SecurityAssociation) | repeated |  |
| ipsec_spds | [ligato.vpp.ipsec.SecurityPolicyDatabase](#ligato.vpp.ipsec.SecurityPolicyDatabase) | repeated |  |
| ipsec_sps | [ligato.vpp.ipsec.SecurityPolicy](#ligato.vpp.ipsec.SecurityPolicy) | repeated |  |
| ipsec_tunnel_protections | [ligato.vpp.ipsec.TunnelProtection](#ligato.vpp.ipsec.TunnelProtection) | repeated |  |
| l3xconnects | [ligato.vpp.l3.L3XConnect](#ligato.vpp.l3.L3XConnect) | repeated |  |
| nat44_global | [ligato.vpp.nat.Nat44Global](#ligato.vpp.nat.Nat44Global) |  |  |
| nat44_interfaces | [ligato.vpp.nat.Nat44Interface](#ligato.vpp.nat.Nat44Interface) | repeated |  |
| nat44_pools | [ligato.vpp.nat.Nat44AddressPool](#ligato.vpp.nat.Nat44AddressPool) | repeated |  |
| proxy_arp | [ligato.vpp.l3.ProxyARP](#ligato.vpp.l3.ProxyARP) |  |  |
| punt_exceptions | [ligato.vpp.punt.Exception](#ligato.vpp.punt.Exception) | repeated |  |
| punt_ipredirects | [ligato.vpp.punt.IPRedirect](#ligato.vpp.punt.IPRedirect) | repeated |  |
| punt_tohosts | [ligato.vpp.punt.ToHost](#ligato.vpp.punt.ToHost) | repeated |  |
| routes | [ligato.vpp.l3.Route](#ligato.vpp.l3.Route) | repeated |  |
| spans | [ligato.vpp.interfaces.Span](#ligato.vpp.interfaces.Span) | repeated |  |
| srv6_global | [ligato.vpp.srv6.SRv6Global](#ligato.vpp.srv6.SRv6Global) |  |  |
| srv6_localsids | [ligato.vpp.srv6.LocalSID](#ligato.vpp.srv6.LocalSID) | repeated |  |
| srv6_policies | [ligato.vpp.srv6.Policy](#ligato.vpp.srv6.Policy) | repeated |  |
| srv6_steerings | [ligato.vpp.srv6.Steering](#ligato.vpp.srv6.Steering) | repeated |  |
| teib_entries | [ligato.vpp.l3.TeibEntry](#ligato.vpp.l3.TeibEntry) | repeated |  |
| vrfs | [ligato.vpp.l3.VrfTable](#ligato.vpp.l3.VrfTable) | repeated |  |
| wg_peers | [ligato.vpp.wireguard.Peer](#ligato.vpp.wireguard.Peer) | repeated |  |
| xconnect_pairs | [ligato.vpp.l2.XConnectPair](#ligato.vpp.l2.XConnectPair) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/wireguard/wireguard.proto"></a>

## ligato/vpp/wireguard/wireguard.proto



<a name="ligato.vpp.wireguard.Peer"></a>

### Peer



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| public_key | [string](#string) |  | Public-key base64 |
| port | [uint32](#uint32) |  | Peer UDP port |
| persistent_keepalive | [uint32](#uint32) |  | Keepalive interval (sec) |
| endpoint | [string](#string) |  | Endpoint IP |
| wg_if_name | [string](#string) |  | The name of the wireguard interface to which this peer belongs |
| flags | [uint32](#uint32) |  | Flags WIREGUARD_PEER_STATUS_DEAD = 0x1 |
| allowed_ips | [string](#string) | repeated | Allowed IPs |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/vpp.proto"></a>

## ligato/vpp/vpp.proto



<a name="ligato.vpp.ConfigData"></a>

### ConfigData
ConfigData holds the entire VPP configuration.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interfaces | [interfaces.Interface](#ligato.vpp.interfaces.Interface) | repeated |  |
| spans | [interfaces.Span](#ligato.vpp.interfaces.Span) | repeated |  |
| acls | [acl.ACL](#ligato.vpp.acl.ACL) | repeated |  |
| abfs | [abf.ABF](#ligato.vpp.abf.ABF) | repeated |  |
| bridge_domains | [l2.BridgeDomain](#ligato.vpp.l2.BridgeDomain) | repeated |  |
| fibs | [l2.FIBEntry](#ligato.vpp.l2.FIBEntry) | repeated |  |
| xconnect_pairs | [l2.XConnectPair](#ligato.vpp.l2.XConnectPair) | repeated |  |
| routes | [l3.Route](#ligato.vpp.l3.Route) | repeated |  |
| arps | [l3.ARPEntry](#ligato.vpp.l3.ARPEntry) | repeated |  |
| proxy_arp | [l3.ProxyARP](#ligato.vpp.l3.ProxyARP) |  |  |
| ipscan_neighbor | [l3.IPScanNeighbor](#ligato.vpp.l3.IPScanNeighbor) |  |  |
| vrfs | [l3.VrfTable](#ligato.vpp.l3.VrfTable) | repeated |  |
| l3xconnects | [l3.L3XConnect](#ligato.vpp.l3.L3XConnect) | repeated |  |
| dhcp_proxies | [l3.DHCPProxy](#ligato.vpp.l3.DHCPProxy) | repeated |  |
| teib_entries | [l3.TeibEntry](#ligato.vpp.l3.TeibEntry) | repeated |  |
| nat44_global | [nat.Nat44Global](#ligato.vpp.nat.Nat44Global) |  |  |
| dnat44s | [nat.DNat44](#ligato.vpp.nat.DNat44) | repeated |  |
| nat44_interfaces | [nat.Nat44Interface](#ligato.vpp.nat.Nat44Interface) | repeated |  |
| nat44_pools | [nat.Nat44AddressPool](#ligato.vpp.nat.Nat44AddressPool) | repeated |  |
| ipsec_spds | [ipsec.SecurityPolicyDatabase](#ligato.vpp.ipsec.SecurityPolicyDatabase) | repeated |  |
| ipsec_sas | [ipsec.SecurityAssociation](#ligato.vpp.ipsec.SecurityAssociation) | repeated |  |
| ipsec_tunnel_protections | [ipsec.TunnelProtection](#ligato.vpp.ipsec.TunnelProtection) | repeated |  |
| ipsec_sps | [ipsec.SecurityPolicy](#ligato.vpp.ipsec.SecurityPolicy) | repeated |  |
| punt_ipredirects | [punt.IPRedirect](#ligato.vpp.punt.IPRedirect) | repeated |  |
| punt_tohosts | [punt.ToHost](#ligato.vpp.punt.ToHost) | repeated |  |
| punt_exceptions | [punt.Exception](#ligato.vpp.punt.Exception) | repeated |  |
| srv6_global | [srv6.SRv6Global](#ligato.vpp.srv6.SRv6Global) |  |  |
| srv6_localsids | [srv6.LocalSID](#ligato.vpp.srv6.LocalSID) | repeated |  |
| srv6_policies | [srv6.Policy](#ligato.vpp.srv6.Policy) | repeated |  |
| srv6_steerings | [srv6.Steering](#ligato.vpp.srv6.Steering) | repeated |  |
| ipfix_global | [ipfix.IPFIX](#ligato.vpp.ipfix.IPFIX) |  |  |
| ipfix_flowprobe_params | [ipfix.FlowProbeParams](#ligato.vpp.ipfix.FlowProbeParams) |  |  |
| ipfix_flowprobes | [ipfix.FlowProbeFeature](#ligato.vpp.ipfix.FlowProbeFeature) | repeated |  |
| wg_peers | [wireguard.Peer](#ligato.vpp.wireguard.Peer) | repeated |  |
| dns_cache | [dns.DNSCache](#ligato.vpp.dns.DNSCache) |  |  |






<a name="ligato.vpp.Notification"></a>

### Notification



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [interfaces.InterfaceNotification](#ligato.vpp.interfaces.InterfaceNotification) |  |  |






<a name="ligato.vpp.Stats"></a>

### Stats



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [interfaces.InterfaceStats](#ligato.vpp.interfaces.InterfaceStats) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/stn/stn.proto"></a>

## ligato/vpp/stn/stn.proto



<a name="ligato.vpp.stn.Rule"></a>

### Rule



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ip_address | [string](#string) |  |  |
| interface | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/srv6/srv6.proto"></a>

## ligato/vpp/srv6/srv6.proto



<a name="ligato.vpp.srv6.LocalSID"></a>

### LocalSID



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| sid | [string](#string) |  | segment id (IPv6 Address) |
| installation_vrf_id | [uint32](#uint32) |  | ID of IPv6 VRF table where to install LocalSID routing components (LocalSids with End.AD function ignore this setting due to missing setting in the API. The End.AD functionality is separated from the SRv6 functionality and have no binary API. It has only the CLI API and that doesn't have the installation vrf id (in VPP API called FIB table) setting configurable.) Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |
| base_end_function | [LocalSID.End](#ligato.vpp.srv6.LocalSID.End) |  |  |
| end_function_x | [LocalSID.EndX](#ligato.vpp.srv6.LocalSID.EndX) |  |  |
| end_function_t | [LocalSID.EndT](#ligato.vpp.srv6.LocalSID.EndT) |  |  |
| end_function_dx2 | [LocalSID.EndDX2](#ligato.vpp.srv6.LocalSID.EndDX2) |  |  |
| end_function_dx4 | [LocalSID.EndDX4](#ligato.vpp.srv6.LocalSID.EndDX4) |  |  |
| end_function_dx6 | [LocalSID.EndDX6](#ligato.vpp.srv6.LocalSID.EndDX6) |  |  |
| end_function_dt4 | [LocalSID.EndDT4](#ligato.vpp.srv6.LocalSID.EndDT4) |  |  |
| end_function_dt6 | [LocalSID.EndDT6](#ligato.vpp.srv6.LocalSID.EndDT6) |  |  |
| end_function_ad | [LocalSID.EndAD](#ligato.vpp.srv6.LocalSID.EndAD) |  |  |






<a name="ligato.vpp.srv6.LocalSID.End"></a>

### LocalSID.End
End function behavior of simple endpoint


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| psp | [bool](#bool) |  | use PSP (penultimate segment POP of the SRH) or by default use USP (Ultimate Segment Pop of the SRH) |






<a name="ligato.vpp.srv6.LocalSID.EndAD"></a>

### LocalSID.EndAD
End function behavior of dynamic segment routing proxy endpoint


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| outgoing_interface | [string](#string) |  | name of interface on segment routing proxy side sending data to segment routing unaware service |
| incoming_interface | [string](#string) |  | name of interface on segment routing proxy side receiving data from segment routing unaware service |
| l3_service_address | [string](#string) |  | IPv6/IPv4 address of L3 SR-unaware service (address type depends whether service is IPv4 or IPv6 service), in case of L2 service it must be empty |






<a name="ligato.vpp.srv6.LocalSID.EndDT4"></a>

### LocalSID.EndDT4
End function behavior of endpoint with decapsulation and specific IPv4 table lookup


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | vrf index of IPv4 table that should be used for lookup. vrf_index and fib_table_id should refer to the same routing table. VRF index refer to it from client side and FIB table id from VPP-internal side (index of memory allocated structure from pool)(source: https://wiki.fd.io/view/VPP/Per-feature_Notes). Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |






<a name="ligato.vpp.srv6.LocalSID.EndDT6"></a>

### LocalSID.EndDT6
End function behavior of endpoint with decapsulation and specific IPv6 table lookup


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | vrf index of IPv6 table that should be used for lookup. vrf_index and fib_table_id should refer to the same routing table. VRF index refer to it from client side and FIB table id from VPP-internal side (index of memory allocated structure from pool)(source: https://wiki.fd.io/view/VPP/Per-feature_Notes). Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |






<a name="ligato.vpp.srv6.LocalSID.EndDX2"></a>

### LocalSID.EndDX2
End function behavior of endpoint with decapsulation and Layer-2 cross-connect (or DX2 with egress VLAN rewrite when VLAN notzero - not supported this variant yet)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vlan_tag | [uint32](#uint32) |  | Outgoing VLAN tag |
| outgoing_interface | [string](#string) |  | name of cross-connected outgoing interface |






<a name="ligato.vpp.srv6.LocalSID.EndDX4"></a>

### LocalSID.EndDX4
End function behavior of endpoint with decapsulation and IPv4 cross-connect


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| outgoing_interface | [string](#string) |  | name of cross-connected outgoing interface |
| next_hop | [string](#string) |  | next hop address for cross-connected link |






<a name="ligato.vpp.srv6.LocalSID.EndDX6"></a>

### LocalSID.EndDX6
End function behavior of endpoint with decapsulation and IPv6 cross-connect


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| outgoing_interface | [string](#string) |  | name of cross-connected outgoing interface |
| next_hop | [string](#string) |  | next hop address for cross-connected link |






<a name="ligato.vpp.srv6.LocalSID.EndT"></a>

### LocalSID.EndT
End function behavior of endpoint with specific IPv6 table lookup


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| psp | [bool](#bool) |  | use PSP (penultimate segment POP of the SRH) or by default use USP (Ultimate Segment Pop of the SRH) |
| vrf_id | [uint32](#uint32) |  | vrf index of IPv6 table that should be used for lookup. vrf_index and fib_table_id should refer to the same routing table. VRF index refer to it from client side and FIB table id from VPP-internal side (index of memory allocated structure from pool)(source: https://wiki.fd.io/view/VPP/Per-feature_Notes). Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |






<a name="ligato.vpp.srv6.LocalSID.EndX"></a>

### LocalSID.EndX
End function behavior of endpoint with Layer-3 cross-connect (IPv6)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| psp | [bool](#bool) |  | use PSP (penultimate segment POP of the SRH) or by default use USP (Ultimate Segment Pop of the SRH) |
| outgoing_interface | [string](#string) |  | name of cross-connected outgoing interface |
| next_hop | [string](#string) |  | IPv6 next hop address for cross-connected link |






<a name="ligato.vpp.srv6.Policy"></a>

### Policy
Model for SRv6 policy (policy without at least one policy segment is only cached in ligato and not written to VPP)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| bsid | [string](#string) |  | binding SID (IPv6 Address) |
| installation_vrf_id | [uint32](#uint32) |  | ID of IPv6 VRF table where to install Policy routing components (for loadbalancing/spray are used VPP features that are using VRF table) Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |
| srh_encapsulation | [bool](#bool) |  | are SR headers handled by encapsulation? (no means insertion of SR headers) |
| spray_behaviour | [bool](#bool) |  | spray(multicast) to all policy segments? (no means to use PolicySegment.weight to loadbalance traffic) |
| segment_lists | [Policy.SegmentList](#ligato.vpp.srv6.Policy.SegmentList) | repeated |  |






<a name="ligato.vpp.srv6.Policy.SegmentList"></a>

### Policy.SegmentList
Model for SRv6 Segment List


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| weight | [uint32](#uint32) |  | used for loadbalancing in case of multiple policy segments in routing process (ignored in case of spray policies) |
| segments | [string](#string) | repeated | list of sids creating one segmented road |






<a name="ligato.vpp.srv6.SRv6Global"></a>

### SRv6Global
Global SRv6 config


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| encap_source_address | [string](#string) |  | IPv6 source address for sr encapsulated packets |






<a name="ligato.vpp.srv6.Steering"></a>

### Steering
Model for steering traffic to SRv6 policy


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | globally unique steering identification (used in keys when is steering stored in key-value stores(i.e. ETCD)) |
| policy_bsid | [string](#string) |  | BSID identifier for policy to which we want to steer routing into (policyBSID and policyIndex are mutual exclusive) |
| policy_index | [uint32](#uint32) |  | (vpp-internal)Index identifier for policy to which we want to steer routing into (policyBSID and policyIndex are mutual exclusive) |
| l2_traffic | [Steering.L2Traffic](#ligato.vpp.srv6.Steering.L2Traffic) |  |  |
| l3_traffic | [Steering.L3Traffic](#ligato.vpp.srv6.Steering.L3Traffic) |  |  |






<a name="ligato.vpp.srv6.Steering.L2Traffic"></a>

### Steering.L2Traffic
L2 traffic that should be steered into SR policy


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface_name | [string](#string) |  | name of interface with incoming traffic that should be steered to SR policy |






<a name="ligato.vpp.srv6.Steering.L3Traffic"></a>

### Steering.L3Traffic
L3 traffic that should be steered into SR policy


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| installation_vrf_id | [uint32](#uint32) |  | ID of IPv4/IPv6 VRF table where to install L3 Steering routing components (VRF table type (IPv4/IPv6) is decided by prefix_address value) Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |
| prefix_address | [string](#string) |  | IPv4/IPv6 prefix address(CIRD format) of traffic destination. All traffic with given destination will be steered to given SR policy |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/punt/punt.proto"></a>

## ligato/vpp/punt/punt.proto



<a name="ligato.vpp.punt.Exception"></a>

### Exception
Exception allows specifying punt exceptions used for punting packets.
The type of exception is defined by reason name.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| reason | [string](#string) |  | Name should contain reason name, e.g. `ipsec4-spi-0`. |
| socket_path | [string](#string) |  | SocketPath defines path to unix domain socket used for punt packets to the host. In dumps, it will actually contain the socket defined in VPP config under punt section. |






<a name="ligato.vpp.punt.IPRedirect"></a>

### IPRedirect
IPRedirect allows otherwise dropped packet which destination IP address
matching some of the VPP addresses to redirect to the defined next hop address
via the TX interface.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| l3_protocol | [L3Protocol](#ligato.vpp.punt.L3Protocol) |  | L3 protocol to be redirected |
| rx_interface | [string](#string) |  | Receive interface name. Optional, only redirect traffic incoming from this interface |
| tx_interface | [string](#string) |  | Transmit interface name |
| next_hop | [string](#string) |  | Next hop IP where the traffic is redirected |






<a name="ligato.vpp.punt.Reason"></a>

### Reason
Reason represents punt reason used in exceptions.
List of known exceptions can be retrieved in VPP CLI
with following command:

vpp# show punt reasons
   [0] ipsec4-spi-0 from:[ipsec ]
   [1] ipsec6-spi-0 from:[ipsec ]
   [2] ipsec4-spi-o-udp-0 from:[ipsec ]
   [3] ipsec4-no-such-tunnel from:[ipsec ]
   [4] ipsec6-no-such-tunnel from:[ipsec ]
   [5] VXLAN-GBP-no-such-v4-tunnel from:[vxlan-gbp ]
   [6] VXLAN-GBP-no-such-v6-tunnel from:[vxlan-gbp ]


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Name contains reason name. |






<a name="ligato.vpp.punt.ToHost"></a>

### ToHost
ToHost allows otherwise dropped packet which destination IP address matching
some of the VPP interface IP addresses to be punted to the host.
L3 and L4 protocols can be used for filtering */


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| l3_protocol | [L3Protocol](#ligato.vpp.punt.L3Protocol) |  | L3 destination protocol a packet has to match in order to be punted. |
| l4_protocol | [L4Protocol](#ligato.vpp.punt.L4Protocol) |  | L4 destination protocol a packet has to match. Currently VPP only supports UDP. |
| port | [uint32](#uint32) |  | Destination port |
| socket_path | [string](#string) |  | SocketPath defines path to unix domain socket used for punt packets to the host. In dumps, it will actually contain the socket defined in VPP config under punt section. |





 <!-- end messages -->


<a name="ligato.vpp.punt.L3Protocol"></a>

### L3Protocol
L3Protocol defines Layer 3 protocols.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_L3 | 0 |  |
| IPV4 | 4 |  |
| IPV6 | 6 |  |
| ALL | 10 |  |



<a name="ligato.vpp.punt.L4Protocol"></a>

### L4Protocol
L4Protocol defines Layer 4 protocols.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_L4 | 0 |  |
| TCP | 6 |  |
| UDP | 17 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/nat/nat.proto"></a>

## ligato/vpp/nat/nat.proto



<a name="ligato.vpp.nat.DNat44"></a>

### DNat44
DNat44 defines destination NAT44 configuration.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| label | [string](#string) |  | Unique identifier for the DNAT configuration. |
| st_mappings | [DNat44.StaticMapping](#ligato.vpp.nat.DNat44.StaticMapping) | repeated | A list of static mappings in DNAT. |
| id_mappings | [DNat44.IdentityMapping](#ligato.vpp.nat.DNat44.IdentityMapping) | repeated | A list of identity mappings in DNAT. |






<a name="ligato.vpp.nat.DNat44.IdentityMapping"></a>

### DNat44.IdentityMapping
IdentityMapping defines an identity mapping in DNAT.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | VRF (table) ID. Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto). |
| interface | [string](#string) |  | Name of the interface to use address from; preferred over ip_address. |
| ip_address | [string](#string) |  | IP address. |
| port | [uint32](#uint32) |  | Port (do not set for address mapping). |
| protocol | [DNat44.Protocol](#ligato.vpp.nat.DNat44.Protocol) |  | Protocol used for identity mapping. |






<a name="ligato.vpp.nat.DNat44.StaticMapping"></a>

### DNat44.StaticMapping
StaticMapping defines a list of static mappings in DNAT.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| external_interface | [string](#string) |  | Interface to use external IP from; preferred over external_ip. |
| external_ip | [string](#string) |  | External address. |
| external_port | [uint32](#uint32) |  | Port (do not set for address mapping). |
| local_ips | [DNat44.StaticMapping.LocalIP](#ligato.vpp.nat.DNat44.StaticMapping.LocalIP) | repeated | List of local IP addresses. If there is more than one entry, load-balancing is enabled. |
| protocol | [DNat44.Protocol](#ligato.vpp.nat.DNat44.Protocol) |  | Protocol used for static mapping. |
| twice_nat | [DNat44.StaticMapping.TwiceNatMode](#ligato.vpp.nat.DNat44.StaticMapping.TwiceNatMode) |  | Enable/disable (self-)twice NAT. |
| twice_nat_pool_ip | [string](#string) |  | IP address from Twice-NAT address pool that should be used as source IP in twice-NAT processing. This is override for default behaviour of choosing the first IP address from twice-NAT pool that has available at least one free port (NAT is tracking translation sessions and exhausts free ports for given IP address). This is needed for example in use cases when multiple twice-NAT translations need to use different IP Addresses as source IP addresses. This functionality works with VPP 20.09 and newer. It also needs to have twice_nat set to ENABLED. It doesn't work for load-balanced static mappings (=local_ips has multiple values). |
| session_affinity | [uint32](#uint32) |  | Session affinity. 0 means disabled, otherwise client IP affinity sticky time in seconds. |






<a name="ligato.vpp.nat.DNat44.StaticMapping.LocalIP"></a>

### DNat44.StaticMapping.LocalIP
LocalIP defines a local IP addresses.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | VRF (table) ID. Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto). |
| local_ip | [string](#string) |  | Local IP address). |
| local_port | [uint32](#uint32) |  | Port (do not set for address mapping). |
| probability | [uint32](#uint32) |  | Probability level for load-balancing mode. |






<a name="ligato.vpp.nat.Nat44AddressPool"></a>

### Nat44AddressPool
Nat44AddressPool defines an address pool used for NAT44.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Unique name for address pool |
| vrf_id | [uint32](#uint32) |  | VRF id of tenant, 0xFFFFFFFF means independent of VRF. Non-zero (and not all-ones) VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto). |
| first_ip | [string](#string) |  | First IP address of the pool. |
| last_ip | [string](#string) |  | Last IP address of the pool. Should be higher than first_ip or empty. |
| twice_nat | [bool](#bool) |  | Enable/disable twice NAT. |






<a name="ligato.vpp.nat.Nat44Global"></a>

### Nat44Global
Nat44Global defines global NAT44 configuration.
In VPP version 21.01 and newer the NAT44 plugin has to be explicitly enabled (by default it is
disabled so that it doesn't consume any computational resources). With ligato control-plane
the NAT44 plugin is enabled by submitting the NAT44Global configuration (even default values
will make the plugin enabled). Without Nat44Global, all other NAT44 configuration items
(DNat44, Nat44Interface and Nat44AddressPool) will be in the PENDING state.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| forwarding | [bool](#bool) |  | Enable/disable forwarding. By default it is disabled. |
| endpoint_independent | [bool](#bool) |  | Enable/disable endpoint-independent mode. In endpoint-independent (also known as "simple") mode the VPP NAT plugin holds less information for each session, but only works with outbound NAT and static mappings. In endpoint-dependent mode, which ligato selects as the default, the VPP NAT plugin uses more information to track each session, which in turn enables additional features such as out-to-in-only and twice-nat. In older versions of VPP (<= 20.09) this field is ignored because mode at which the NAT44 plugin operates is given by the VPP startup configuration file (i.e. config created before VPP even starts, therefore not managed by ligato). The endpoint-independent mode is the default and the dependent mode is turned on with this config stanza (included in vpp.conf used by ligato for older VPPs): nat { endpoint-dependent } |
| nat_interfaces | [Nat44Global.Interface](#ligato.vpp.nat.Nat44Global.Interface) | repeated | **Deprecated.** List of NAT-enabled interfaces. Deprecated - use separate Nat44Interface entries instead. |
| address_pool | [Nat44Global.Address](#ligato.vpp.nat.Nat44Global.Address) | repeated | **Deprecated.** Address pool used for source IP NAT. Deprecated - use separate Nat44AddressPool entries instead. |
| virtual_reassembly | [VirtualReassembly](#ligato.vpp.nat.VirtualReassembly) |  | Virtual reassembly for IPv4. |






<a name="ligato.vpp.nat.Nat44Global.Address"></a>

### Nat44Global.Address
Address defines an address to be used for source IP NAT.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| address | [string](#string) |  | IPv4 address. |
| vrf_id | [uint32](#uint32) |  | VRF id of tenant, 0xFFFFFFFF means independent of VRF. Non-zero (and not all-ones) VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto). |
| twice_nat | [bool](#bool) |  | Enable/disable twice NAT. |






<a name="ligato.vpp.nat.Nat44Global.Interface"></a>

### Nat44Global.Interface
Interface defines a network interface enabled for NAT.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Interface name (logical). |
| is_inside | [bool](#bool) |  | Distinguish between inside/outside interface. |
| output_feature | [bool](#bool) |  | Enable/disable output feature. |






<a name="ligato.vpp.nat.Nat44Interface"></a>

### Nat44Interface
Nat44Interface defines a local network interfaces enabled for NAT44.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Interface name (logical). |
| nat_inside | [bool](#bool) |  | Enable/disable NAT on inside. |
| nat_outside | [bool](#bool) |  | Enable/disable NAT on outside. |
| output_feature | [bool](#bool) |  | Enable/disable output feature. |






<a name="ligato.vpp.nat.VirtualReassembly"></a>

### VirtualReassembly
VirtualReassembly defines NAT virtual reassembly settings.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| timeout | [uint32](#uint32) |  | Reassembly timeout. |
| max_reassemblies | [uint32](#uint32) |  | Maximum number of concurrent reassemblies. |
| max_fragments | [uint32](#uint32) |  | Maximum number of fragments per reassembly. |
| drop_fragments | [bool](#bool) |  | If set to true fragments are dropped, translated otherwise. |





 <!-- end messages -->


<a name="ligato.vpp.nat.DNat44.Protocol"></a>

### DNat44.Protocol
Available protocols.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| TCP | 0 |  |
| UDP | 1 |  |
| ICMP | 2 | ICMP is not permitted for load balanced entries. |



<a name="ligato.vpp.nat.DNat44.StaticMapping.TwiceNatMode"></a>

### DNat44.StaticMapping.TwiceNatMode
Available twice-NAT modes.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| DISABLED | 0 |  |
| ENABLED | 1 |  |
| SELF | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/vrrp.proto"></a>

## ligato/vpp/l3/vrrp.proto



<a name="ligato.vpp.l3.VRRPEntry"></a>

### VRRPEntry
VRRPEntry represents Virtual Router desired state.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  | This field refers to logical interface name |
| vr_id | [uint32](#uint32) |  | Should be > 0 and <= 255 |
| priority | [uint32](#uint32) |  | Priority defines which router becomes master. Should be > 0 and <= 255. |
| interval | [uint32](#uint32) |  | VR advertisement interval in milliseconds, should be => 10 and <= 65535. (Later, in implemetation it is converted into centiseconds, so precision may be lost). |
| preempt | [bool](#bool) |  | Controls whether a (starting or restarting) higher-priority Backup router preempts a lower-priority Master router. |
| accept | [bool](#bool) |  | Controls whether a virtual router in Master state will accept packets addressed to the address owner's IPvX address as its own if it is not the IPvX address owner. |
| unicast | [bool](#bool) |  | Unicast mode may be used to take advantage of newer token ring adapter implementations that support non-promiscuous reception for multiple unicast MAC addresses and to avoid both the multicast traffic and usage conflicts associated with the use of token ring functional addresses. |
| ip_addresses | [string](#string) | repeated | Ip address quantity should be > 0 and <= 255. |
| enabled | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/vrf.proto"></a>

## ligato/vpp/l3/vrf.proto



<a name="ligato.vpp.l3.VrfTable"></a>

### VrfTable



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| id | [uint32](#uint32) |  | ID is mandatory identification for VRF table. NOTE: do not confuse with fib index (shown by some VPP CLIs), which is VPP's internal offset in the vector of allocated tables. |
| protocol | [VrfTable.Protocol](#ligato.vpp.l3.VrfTable.Protocol) |  |  |
| label | [string](#string) |  | Label is an optional description for the table. - maximum allowed length is 63 characters - included in the output from the VPP CLI command "show ip fib" - if undefined, then VPP will generate label using the template "<protocol>-VRF:<id>" |
| flow_hash_settings | [VrfTable.FlowHashSettings](#ligato.vpp.l3.VrfTable.FlowHashSettings) |  |  |






<a name="ligato.vpp.l3.VrfTable.FlowHashSettings"></a>

### VrfTable.FlowHashSettings
FlowHashSettings allows tuning of hash calculation of IP flows in the VRF table.
This affects hash table size as well as the stickiness of flows by load-balancing.
If not defined, default settings that are implicitly enabled are:
 - use_src_ip, use_dst_ip, use_src_port, use_dst_port, use_protocol


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| use_src_ip | [bool](#bool) |  |  |
| use_dst_ip | [bool](#bool) |  |  |
| use_src_port | [bool](#bool) |  |  |
| use_dst_port | [bool](#bool) |  |  |
| use_protocol | [bool](#bool) |  |  |
| reverse | [bool](#bool) |  |  |
| symmetric | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="ligato.vpp.l3.VrfTable.Protocol"></a>

### VrfTable.Protocol
Protocol define IP protocol of VRF table.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| IPV4 | 0 |  |
| IPV6 | 1 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/teib.proto"></a>

## ligato/vpp/l3/teib.proto



<a name="ligato.vpp.l3.TeibEntry"></a>

### TeibEntry
TeibEntry represents an tunnel endpoint information base entry.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  | Interface references a tunnel interface this TEIB entry is linked to. |
| peer_addr | [string](#string) |  | IP address of the peer. |
| next_hop_addr | [string](#string) |  | Next hop IP address. |
| vrf_id | [uint32](#uint32) |  | VRF ID used to reach the next hop. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/route.proto"></a>

## ligato/vpp/l3/route.proto



<a name="ligato.vpp.l3.Route"></a>

### Route



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| type | [Route.RouteType](#ligato.vpp.l3.Route.RouteType) |  |  |
| vrf_id | [uint32](#uint32) |  | VRF identifier, field required for remote client. This value should be consistent with VRF ID in static route key. If it is not, value from key will be preffered and this field will be overriden. Non-zero VRF has to be explicitly created (see api/models/vpp/l3/vrf.proto) |
| dst_network | [string](#string) |  | Destination network defined by IP address and prefix (format: <address>/<prefix>). |
| next_hop_addr | [string](#string) |  | Next hop address. |
| outgoing_interface | [string](#string) |  | Interface name of the outgoing interface. |
| weight | [uint32](#uint32) |  | Weight is used for unequal cost load balancing. |
| preference | [uint32](#uint32) |  | Preference defines path preference. Lower preference is preferred. Only paths with the best preference contribute to forwarding (a poor man's primary and backup). |
| via_vrf_id | [uint32](#uint32) |  | Specifies VRF ID for the next hop lookup / recursive lookup |





 <!-- end messages -->


<a name="ligato.vpp.l3.Route.RouteType"></a>

### Route.RouteType


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| INTRA_VRF | 0 | Forwarding is being done in the specified vrf_id only, or according to the specified outgoing interface. |
| INTER_VRF | 1 | Forwarding is being done by lookup into a different VRF, specified as via_vrf_id field. In case of these routes, the outgoing interface should not be specified. The next hop IP address does not have to be specified either, in that case VPP does full recursive lookup in the via_vrf_id VRF. |
| DROP | 2 | Drops the network communication designated for specific IP address. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/l3xc.proto"></a>

## ligato/vpp/l3/l3xc.proto



<a name="ligato.vpp.l3.L3XConnect"></a>

### L3XConnect



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  |  |
| protocol | [L3XConnect.Protocol](#ligato.vpp.l3.L3XConnect.Protocol) |  |  |
| paths | [L3XConnect.Path](#ligato.vpp.l3.L3XConnect.Path) | repeated |  |






<a name="ligato.vpp.l3.L3XConnect.Path"></a>

### L3XConnect.Path



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| outgoing_interface | [string](#string) |  |  |
| next_hop_addr | [string](#string) |  |  |
| weight | [uint32](#uint32) |  |  |
| preference | [uint32](#uint32) |  |  |





 <!-- end messages -->


<a name="ligato.vpp.l3.L3XConnect.Protocol"></a>

### L3XConnect.Protocol


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| IPV4 | 0 |  |
| IPV6 | 1 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/l3.proto"></a>

## ligato/vpp/l3/l3.proto



<a name="ligato.vpp.l3.DHCPProxy"></a>

### DHCPProxy
DHCP Proxy


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| source_ip_address | [string](#string) |  |  |
| rx_vrf_id | [uint32](#uint32) |  |  |
| servers | [DHCPProxy.DHCPServer](#ligato.vpp.l3.DHCPProxy.DHCPServer) | repeated |  |






<a name="ligato.vpp.l3.DHCPProxy.DHCPServer"></a>

### DHCPProxy.DHCPServer



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  |  |
| ip_address | [string](#string) |  |  |






<a name="ligato.vpp.l3.IPScanNeighbor"></a>

### IPScanNeighbor
IP Neighbour Config


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| mode | [IPScanNeighbor.Mode](#ligato.vpp.l3.IPScanNeighbor.Mode) |  |  |
| scan_interval | [uint32](#uint32) |  |  |
| max_proc_time | [uint32](#uint32) |  |  |
| max_update | [uint32](#uint32) |  |  |
| scan_int_delay | [uint32](#uint32) |  |  |
| stale_threshold | [uint32](#uint32) |  |  |






<a name="ligato.vpp.l3.ProxyARP"></a>

### ProxyARP
ARP Proxy


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interfaces | [ProxyARP.Interface](#ligato.vpp.l3.ProxyARP.Interface) | repeated | List of interfaces proxy ARP is enabled for. |
| ranges | [ProxyARP.Range](#ligato.vpp.l3.ProxyARP.Range) | repeated |  |






<a name="ligato.vpp.l3.ProxyARP.Interface"></a>

### ProxyARP.Interface



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  |  |






<a name="ligato.vpp.l3.ProxyARP.Range"></a>

### ProxyARP.Range



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| first_ip_addr | [string](#string) |  |  |
| last_ip_addr | [string](#string) |  |  |
| vrf_id | [uint32](#uint32) |  |  |





 <!-- end messages -->


<a name="ligato.vpp.l3.IPScanNeighbor.Mode"></a>

### IPScanNeighbor.Mode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| DISABLED | 0 |  |
| IPV4 | 1 |  |
| IPV6 | 2 |  |
| BOTH | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l3/arp.proto"></a>

## ligato/vpp/l3/arp.proto



<a name="ligato.vpp.l3.ARPEntry"></a>

### ARPEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  |  |
| ip_address | [string](#string) |  |  |
| phys_address | [string](#string) |  |  |
| static | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l2/xconnect.proto"></a>

## ligato/vpp/l2/xconnect.proto



<a name="ligato.vpp.l2.XConnectPair"></a>

### XConnectPair



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| receive_interface | [string](#string) |  |  |
| transmit_interface | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l2/fib.proto"></a>

## ligato/vpp/l2/fib.proto



<a name="ligato.vpp.l2.FIBEntry"></a>

### FIBEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| phys_address | [string](#string) |  | unique destination MAC address |
| bridge_domain | [string](#string) |  | name of bridge domain this FIB table entry belongs to |
| action | [FIBEntry.Action](#ligato.vpp.l2.FIBEntry.Action) |  | action to tke on matching frames |
| outgoing_interface | [string](#string) |  | outgoing interface for matching frames |
| static_config | [bool](#bool) |  | true if this is a statically configured FIB entry |
| bridged_virtual_interface | [bool](#bool) |  | the MAC address is a bridge virtual interface MAC |





 <!-- end messages -->


<a name="ligato.vpp.l2.FIBEntry.Action"></a>

### FIBEntry.Action


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| FORWARD | 0 | forward the matching frame |
| DROP | 1 | drop the matching frame |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/l2/bridge_domain.proto"></a>

## ligato/vpp/l2/bridge_domain.proto



<a name="ligato.vpp.l2.BridgeDomain"></a>

### BridgeDomain



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | bridge domain name (can be any string) |
| flood | [bool](#bool) |  | enable/disable broadcast/multicast flooding in the BD |
| unknown_unicast_flood | [bool](#bool) |  | enable/disable unknown unicast flood in the BD |
| forward | [bool](#bool) |  | enable/disable forwarding on all interfaces in the BD |
| learn | [bool](#bool) |  | enable/disable learning on all interfaces in the BD |
| arp_termination | [bool](#bool) |  | enable/disable ARP termination in the BD |
| mac_age | [uint32](#uint32) |  | MAC aging time in min, 0 for disabled aging |
| interfaces | [BridgeDomain.Interface](#ligato.vpp.l2.BridgeDomain.Interface) | repeated | list of interfaces |
| arp_termination_table | [BridgeDomain.ArpTerminationEntry](#ligato.vpp.l2.BridgeDomain.ArpTerminationEntry) | repeated | list of ARP termination entries |






<a name="ligato.vpp.l2.BridgeDomain.ArpTerminationEntry"></a>

### BridgeDomain.ArpTerminationEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ip_address | [string](#string) |  | IP address |
| phys_address | [string](#string) |  | MAC address matching to the IP |






<a name="ligato.vpp.l2.BridgeDomain.Interface"></a>

### BridgeDomain.Interface



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | interface name belonging to this bridge domain |
| bridged_virtual_interface | [bool](#bool) |  | true if this is a BVI interface |
| split_horizon_group | [uint32](#uint32) |  | VXLANs in the same BD need the same non-zero SHG |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/ipsec/ipsec.proto"></a>

## ligato/vpp/ipsec/ipsec.proto



<a name="ligato.vpp.ipsec.SecurityAssociation"></a>

### SecurityAssociation
Security Association (SA)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| index | [uint32](#uint32) |  | Numerical security association index, serves as a unique identifier |
| spi | [uint32](#uint32) |  | Security parameter index |
| protocol | [SecurityAssociation.IPSecProtocol](#ligato.vpp.ipsec.SecurityAssociation.IPSecProtocol) |  |  |
| crypto_alg | [CryptoAlg](#ligato.vpp.ipsec.CryptoAlg) |  | Cryptographic algorithm for encryption |
| crypto_key | [string](#string) |  |  |
| crypto_salt | [uint32](#uint32) |  |  |
| integ_alg | [IntegAlg](#ligato.vpp.ipsec.IntegAlg) |  | Cryptographic algorithm for authentication |
| integ_key | [string](#string) |  |  |
| use_esn | [bool](#bool) |  | Use extended sequence number |
| use_anti_replay | [bool](#bool) |  | Use anti replay |
| tunnel_src_addr | [string](#string) |  |  |
| tunnel_dst_addr | [string](#string) |  |  |
| enable_udp_encap | [bool](#bool) |  | Enable UDP encapsulation for NAT traversal |
| tunnel_src_port | [uint32](#uint32) |  |  |
| tunnel_dst_port | [uint32](#uint32) |  |  |






<a name="ligato.vpp.ipsec.SecurityPolicy"></a>

### SecurityPolicy



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| spd_index | [uint32](#uint32) |  | Security policy database index |
| sa_index | [uint32](#uint32) |  | Security association index |
| priority | [int32](#int32) |  |  |
| is_outbound | [bool](#bool) |  |  |
| remote_addr_start | [string](#string) |  |  |
| remote_addr_stop | [string](#string) |  |  |
| local_addr_start | [string](#string) |  |  |
| local_addr_stop | [string](#string) |  |  |
| protocol | [uint32](#uint32) |  |  |
| remote_port_start | [uint32](#uint32) |  |  |
| remote_port_stop | [uint32](#uint32) |  |  |
| local_port_start | [uint32](#uint32) |  |  |
| local_port_stop | [uint32](#uint32) |  |  |
| action | [SecurityPolicy.Action](#ligato.vpp.ipsec.SecurityPolicy.Action) |  |  |






<a name="ligato.vpp.ipsec.SecurityPolicyDatabase"></a>

### SecurityPolicyDatabase
Security Policy Database (SPD)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| index | [uint32](#uint32) |  | Numerical security policy database index, serves as a unique identifier |
| interfaces | [SecurityPolicyDatabase.Interface](#ligato.vpp.ipsec.SecurityPolicyDatabase.Interface) | repeated | List of interfaces belonging to this SPD |
| policy_entries | [SecurityPolicyDatabase.PolicyEntry](#ligato.vpp.ipsec.SecurityPolicyDatabase.PolicyEntry) | repeated | **Deprecated.** List of policy entries belonging to this SPD. Deprecated and actually trying to use this will return an error. Use separate model for Security Policy (SP) defined below. |






<a name="ligato.vpp.ipsec.SecurityPolicyDatabase.Interface"></a>

### SecurityPolicyDatabase.Interface



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Name of the related interface |






<a name="ligato.vpp.ipsec.SecurityPolicyDatabase.PolicyEntry"></a>

### SecurityPolicyDatabase.PolicyEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| sa_index | [uint32](#uint32) |  | Security association index |
| priority | [int32](#int32) |  |  |
| is_outbound | [bool](#bool) |  |  |
| remote_addr_start | [string](#string) |  |  |
| remote_addr_stop | [string](#string) |  |  |
| local_addr_start | [string](#string) |  |  |
| local_addr_stop | [string](#string) |  |  |
| protocol | [uint32](#uint32) |  |  |
| remote_port_start | [uint32](#uint32) |  |  |
| remote_port_stop | [uint32](#uint32) |  |  |
| local_port_start | [uint32](#uint32) |  |  |
| local_port_stop | [uint32](#uint32) |  |  |
| action | [SecurityPolicyDatabase.PolicyEntry.Action](#ligato.vpp.ipsec.SecurityPolicyDatabase.PolicyEntry.Action) |  |  |






<a name="ligato.vpp.ipsec.TunnelProtection"></a>

### TunnelProtection
TunnelProtection allows enabling IPSec tunnel protection on an existing interface
(only IPIP tunnel interfaces are currently supported)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  | Name of the interface to be protected with IPSec. |
| sa_out | [uint32](#uint32) | repeated | Outbound security associations identified by SA index. |
| sa_in | [uint32](#uint32) | repeated | Inbound security associations identified by SA index. |
| next_hop_addr | [string](#string) |  | (Optional) Next hop IP address, used for multipoint tunnels. |





 <!-- end messages -->


<a name="ligato.vpp.ipsec.CryptoAlg"></a>

### CryptoAlg
Cryptographic algorithm for encryption

vpp/src/vnet/ipsec/ipsec_sa.h:22

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| NONE_CRYPTO | 0 |  |
| AES_CBC_128 | 1 |  |
| AES_CBC_192 | 2 |  |
| AES_CBC_256 | 3 |  |
| AES_CTR_128 | 4 |  |
| AES_CTR_192 | 5 |  |
| AES_CTR_256 | 6 |  |
| AES_GCM_128 | 7 |  |
| AES_GCM_192 | 8 |  |
| AES_GCM_256 | 9 |  |
| DES_CBC | 10 |  |
| DES3_CBC | 11 | 3DES_CBC |



<a name="ligato.vpp.ipsec.IntegAlg"></a>

### IntegAlg
Cryptographic algorithm for authentication

vpp/src/vnet/ipsec/ipsec_sa.h:44

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| NONE_INTEG | 0 |  |
| MD5_96 | 1 | RFC2403 |
| SHA1_96 | 2 | RFC2404 |
| SHA_256_96 | 3 | draft-ietf-ipsec-ciph-sha-256-00 |
| SHA_256_128 | 4 | RFC4868 |
| SHA_384_192 | 5 | RFC4868 |
| SHA_512_256 | 6 | RFC4868 |



<a name="ligato.vpp.ipsec.SecurityAssociation.IPSecProtocol"></a>

### SecurityAssociation.IPSecProtocol


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| AH | 0 | Authentication Header, provides a mechanism for authentication only |
| ESP | 1 | Encapsulating Security Payload is for data confidentiality and authentication |



<a name="ligato.vpp.ipsec.SecurityPolicy.Action"></a>

### SecurityPolicy.Action


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| BYPASS | 0 |  |
| DISCARD | 1 |  |
| RESOLVE | 2 | Note: this particular action is unused in VPP |
| PROTECT | 3 |  |



<a name="ligato.vpp.ipsec.SecurityPolicyDatabase.PolicyEntry.Action"></a>

### SecurityPolicyDatabase.PolicyEntry.Action


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| BYPASS | 0 |  |
| DISCARD | 1 |  |
| RESOLVE | 2 | Note: this particular action is unused in VPP |
| PROTECT | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/ipfix/ipfix.proto"></a>

## ligato/vpp/ipfix/ipfix.proto



<a name="ligato.vpp.ipfix.IPFIX"></a>

### IPFIX
IPFIX defines the IP Flow Information eXport (IPFIX) configuration.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| collector | [IPFIX.Collector](#ligato.vpp.ipfix.IPFIX.Collector) |  |  |
| source_address | [string](#string) |  |  |
| vrf_id | [uint32](#uint32) |  |  |
| path_mtu | [uint32](#uint32) |  |  |
| template_interval | [uint32](#uint32) |  |  |






<a name="ligato.vpp.ipfix.IPFIX.Collector"></a>

### IPFIX.Collector



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| address | [string](#string) |  |  |
| port | [uint32](#uint32) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/ipfix/flowprobe.proto"></a>

## ligato/vpp/ipfix/flowprobe.proto



<a name="ligato.vpp.ipfix.FlowProbeFeature"></a>

### FlowProbeFeature



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  |  |
| l2 | [bool](#bool) |  |  |
| ip4 | [bool](#bool) |  |  |
| ip6 | [bool](#bool) |  |  |






<a name="ligato.vpp.ipfix.FlowProbeParams"></a>

### FlowProbeParams



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| active_timer | [uint32](#uint32) |  |  |
| passive_timer | [uint32](#uint32) |  |  |
| record_l2 | [bool](#bool) |  |  |
| record_l3 | [bool](#bool) |  |  |
| record_l4 | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/interfaces/interface.proto"></a>

## ligato/vpp/interfaces/interface.proto



<a name="ligato.vpp.interfaces.AfpacketLink"></a>

### AfpacketLink
AfpacketLink defines configuration for interface type: AF_PACKET


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| host_if_name | [string](#string) |  | Name of the host (Linux) interface to bind to. This type of reference is suitable for scenarios when the target interface is not managed (and should not be touched) by the agent. In such cases the interface does not have logical name in the agent's namespace and can only be referenced by the host interface name (i.e. the name used in the Linux network stack). Please note that agent learns about externally created interfaces through netlink notifications. If, however, the target interface is managed by the agent, then it is recommended to use the alternative reference <linux_interface> (see below), pointing to the interface by its logical name. One advantage of such approach is, that if AF-PACKET and the target Linux interface are requested to be created at the same time, then it can be done inside the same transaction because the agent does not rely on any notification from the Linux. It is mandatory to define either <host_if_name> or <linux_interface>. |
| linux_interface | [string](#string) |  | Logical name of the Linux interface to bind to. This is an alternative interface reference to <host_if_name> and preferred if the target interface is managed by the agent and not created externally (see comments for <host_if_name> for explanation). It is mandatory to define either <host_if_name> or <linux_interface>. |






<a name="ligato.vpp.interfaces.BondLink"></a>

### BondLink
BondLink defines configuration for interface type: BOND_INTERFACE


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| id | [uint32](#uint32) |  |  |
| mode | [BondLink.Mode](#ligato.vpp.interfaces.BondLink.Mode) |  |  |
| lb | [BondLink.LoadBalance](#ligato.vpp.interfaces.BondLink.LoadBalance) |  | Load balance is optional and valid only for XOR and LACP modes |
| bonded_interfaces | [BondLink.BondedInterface](#ligato.vpp.interfaces.BondLink.BondedInterface) | repeated |  |






<a name="ligato.vpp.interfaces.BondLink.BondedInterface"></a>

### BondLink.BondedInterface



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  |  |
| is_passive | [bool](#bool) |  |  |
| is_long_timeout | [bool](#bool) |  |  |






<a name="ligato.vpp.interfaces.GreLink"></a>

### GreLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| tunnel_type | [GreLink.Type](#ligato.vpp.interfaces.GreLink.Type) |  |  |
| src_addr | [string](#string) |  |  |
| dst_addr | [string](#string) |  |  |
| outer_fib_id | [uint32](#uint32) |  |  |
| session_id | [uint32](#uint32) |  |  |






<a name="ligato.vpp.interfaces.GtpuLink"></a>

### GtpuLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| src_addr | [string](#string) |  | Source VTEP address |
| dst_addr | [string](#string) |  | Destination VTEP address |
| multicast | [string](#string) |  | Name of multicast interface |
| teid | [uint32](#uint32) |  | Tunnel endpoint identifier - local |
| remote_teid | [uint32](#uint32) |  | Tunnel endpoint identifier - remote |
| encap_vrf_id | [uint32](#uint32) |  | VRF id for the encapsulated packets |
| decap_next | [GtpuLink.NextNode](#ligato.vpp.interfaces.GtpuLink.NextNode) |  | **Deprecated.** DEPRECATED - use decap_next_node |
| decap_next_node | [uint32](#uint32) |  | Next VPP node after decapsulation |






<a name="ligato.vpp.interfaces.IPIPLink"></a>

### IPIPLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| tunnel_mode | [IPIPLink.Mode](#ligato.vpp.interfaces.IPIPLink.Mode) |  | Mode of the IPIP tunnel |
| src_addr | [string](#string) |  | Source VTEP IP address |
| dst_addr | [string](#string) |  | Destination VTEP IP address |






<a name="ligato.vpp.interfaces.IPSecLink"></a>

### IPSecLink
IPSecLink defines configuration for interface type: IPSEC_TUNNEL
In VPP 21.06 and newer, IPSecLink serves just for creation of the link and thus only tunnel_mode is taken into
account and all of the remaining (deprecated) fields are ignored.
Please use separate SecurityPolicy, SecurityAssociation and TunnelProtection messages from ligato.vpp.ipsec
package to associate SA, SP and tunnel protection with the link.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| tunnel_mode | [IPSecLink.Mode](#ligato.vpp.interfaces.IPSecLink.Mode) |  | Mode of the IPIP tunnel |
| esn | [bool](#bool) |  | **Deprecated.** Extended sequence number |
| anti_replay | [bool](#bool) |  | **Deprecated.** Anti replay option |
| local_ip | [string](#string) |  | **Deprecated.** Local IP address |
| remote_ip | [string](#string) |  | **Deprecated.** Remote IP address |
| local_spi | [uint32](#uint32) |  | **Deprecated.** Local security parameter index |
| remote_spi | [uint32](#uint32) |  | **Deprecated.** Remote security parameter index |
| crypto_alg | [ligato.vpp.ipsec.CryptoAlg](#ligato.vpp.ipsec.CryptoAlg) |  | **Deprecated.** Cryptographic algorithm for encryption |
| local_crypto_key | [string](#string) |  | **Deprecated.**  |
| remote_crypto_key | [string](#string) |  | **Deprecated.**  |
| integ_alg | [ligato.vpp.ipsec.IntegAlg](#ligato.vpp.ipsec.IntegAlg) |  | **Deprecated.** Cryptographic algorithm for authentication |
| local_integ_key | [string](#string) |  | **Deprecated.**  |
| remote_integ_key | [string](#string) |  | **Deprecated.**  |
| enable_udp_encap | [bool](#bool) |  | **Deprecated.**  |






<a name="ligato.vpp.interfaces.Interface"></a>

### Interface
Interface defines a VPP interface.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Name is mandatory field representing logical name for the interface. It must be unique across all configured VPP interfaces. |
| type | [Interface.Type](#ligato.vpp.interfaces.Interface.Type) |  | Type represents the type of VPP interface and it must match the actual Link. |
| enabled | [bool](#bool) |  | Enabled controls if the interface should be UP. |
| phys_address | [string](#string) |  | PhysAddress represents physical address (MAC) of the interface. Random address will be assigned if left empty. |
| ip_addresses | [string](#string) | repeated | IPAddresses define list of IP addresses for the interface and must be defined in the following format: <ipAddress>/<ipPrefix>. Interface IP address can be also allocated via netalloc plugin and referenced here, see: api/models/netalloc/netalloc.proto |
| vrf | [uint32](#uint32) |  | Vrf defines the ID of VRF table that the interface is assigned to. The VRF table must be explicitely configured (see api/models/vpp/l3/vrf.proto). When using unnumbered interface the actual vrf is inherited from the interface referenced by the numbered interface and this field is ignored. |
| set_dhcp_client | [bool](#bool) |  | SetDhcpClient enables DHCP client on interface. |
| ip6_nd | [Interface.IP6ND](#ligato.vpp.interfaces.Interface.IP6ND) |  |  |
| mtu | [uint32](#uint32) |  | Mtu sets MTU (Maximum Transmission Unit) for this interface. If set to zero, default MTU (usually 9216) will be used. |
| unnumbered | [Interface.Unnumbered](#ligato.vpp.interfaces.Interface.Unnumbered) |  |  |
| rx_modes | [Interface.RxMode](#ligato.vpp.interfaces.Interface.RxMode) | repeated |  |
| rx_placements | [Interface.RxPlacement](#ligato.vpp.interfaces.Interface.RxPlacement) | repeated |  |
| sub | [SubInterface](#ligato.vpp.interfaces.SubInterface) |  |  |
| memif | [MemifLink](#ligato.vpp.interfaces.MemifLink) |  |  |
| afpacket | [AfpacketLink](#ligato.vpp.interfaces.AfpacketLink) |  |  |
| tap | [TapLink](#ligato.vpp.interfaces.TapLink) |  |  |
| vxlan | [VxlanLink](#ligato.vpp.interfaces.VxlanLink) |  |  |
| ipsec | [IPSecLink](#ligato.vpp.interfaces.IPSecLink) |  | **Deprecated.** Deprecated in VPP 20.01+. Use IPIP_TUNNEL + ipsec.TunnelProtection instead. |
| vmx_net3 | [VmxNet3Link](#ligato.vpp.interfaces.VmxNet3Link) |  |  |
| bond | [BondLink](#ligato.vpp.interfaces.BondLink) |  |  |
| gre | [GreLink](#ligato.vpp.interfaces.GreLink) |  |  |
| gtpu | [GtpuLink](#ligato.vpp.interfaces.GtpuLink) |  |  |
| ipip | [IPIPLink](#ligato.vpp.interfaces.IPIPLink) |  |  |
| wireguard | [WireguardLink](#ligato.vpp.interfaces.WireguardLink) |  |  |
| rdma | [RDMALink](#ligato.vpp.interfaces.RDMALink) |  |  |






<a name="ligato.vpp.interfaces.Interface.IP6ND"></a>

### Interface.IP6ND
Ip6Nd is used to enable/disable IPv6 ND address autoconfiguration
and setting up default routes


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| address_autoconfig | [bool](#bool) |  | Enable IPv6 ND address autoconfiguration. |
| install_default_routes | [bool](#bool) |  | Enable installing default routes. |






<a name="ligato.vpp.interfaces.Interface.RxMode"></a>

### Interface.RxMode



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| queue | [uint32](#uint32) |  |  |
| mode | [Interface.RxMode.Type](#ligato.vpp.interfaces.Interface.RxMode.Type) |  |  |
| default_mode | [bool](#bool) |  | DefaultMode, if set to true, the <queue> field will be ignored and the <mode> will be used as a default for all the queues. |






<a name="ligato.vpp.interfaces.Interface.RxPlacement"></a>

### Interface.RxPlacement



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| queue | [uint32](#uint32) |  | Select from interval <0, number-of-queues) |
| worker | [uint32](#uint32) |  | Select from interval <0, number-of-workers) |
| main_thread | [bool](#bool) |  | Let the main thread to process the given queue - if enabled, value of <worker> is ignored |






<a name="ligato.vpp.interfaces.Interface.Unnumbered"></a>

### Interface.Unnumbered
Unnumbered is used for inheriting IP address from another interface.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface_with_ip | [string](#string) |  | InterfaceWithIp is the name of interface to inherit IP address from. |






<a name="ligato.vpp.interfaces.MemifLink"></a>

### MemifLink
MemifLink defines configuration for interface type: MEMIF


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| mode | [MemifLink.MemifMode](#ligato.vpp.interfaces.MemifLink.MemifMode) |  |  |
| master | [bool](#bool) |  |  |
| id | [uint32](#uint32) |  | Id is a 32bit integer used to authenticate and match opposite sides of the connection |
| socket_filename | [string](#string) |  | Filename of the socket used for connection establishment |
| secret | [string](#string) |  |  |
| ring_size | [uint32](#uint32) |  | The number of entries of RX/TX rings |
| buffer_size | [uint32](#uint32) |  | Size of the buffer allocated for each ring entry |
| rx_queues | [uint32](#uint32) |  | Number of rx queues (only valid for slave) |
| tx_queues | [uint32](#uint32) |  | Number of tx queues (only valid for slave) |






<a name="ligato.vpp.interfaces.RDMALink"></a>

### RDMALink
https://github.com/FDio/vpp/blob/master/src/plugins/rdma/rdma_doc.rst


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| host_if_name | [string](#string) |  | Linux interface name representing the RDMA-enabled network device to attach into. |
| mode | [RDMALink.Mode](#ligato.vpp.interfaces.RDMALink.Mode) |  | Mode at which the RDMA driver operates. |
| rxq_num | [uint32](#uint32) |  | Number of receive queues. By default only one RX queue is used. |
| rxq_size | [uint32](#uint32) |  | The size of each RX queue. Default is 1024 bytes. |
| txq_size | [uint32](#uint32) |  | The size of each TX queue. Default is 1024 bytes. |






<a name="ligato.vpp.interfaces.SubInterface"></a>

### SubInterface
SubInterface defines configuration for interface type: SUB_INTERFACE


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| parent_name | [string](#string) |  | Name of the parent (super) interface |
| sub_id | [uint32](#uint32) |  | SubInterface ID, used as VLAN |
| tag_rw_option | [SubInterface.TagRewriteOptions](#ligato.vpp.interfaces.SubInterface.TagRewriteOptions) |  | VLAN tag rewrite rule applied for given tag for sub-interface |
| push_dot1q | [bool](#bool) |  | Set ether-type of the first tag to dot1q if true, dot1ad otherwise |
| tag1 | [uint32](#uint32) |  | First tag (required for PUSH1 and any TRANSLATE) |
| tag2 | [uint32](#uint32) |  | Second tag (required for PUSH2 and any TRANSLATE) |






<a name="ligato.vpp.interfaces.TapLink"></a>

### TapLink
TapLink defines configuration for interface type: TAP


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| version | [uint32](#uint32) |  | 1 / unset = use the original TAP interface; 2 = use a fast virtio-based TAP |
| host_if_name | [string](#string) |  | Name of the TAP interface in the host OS; if empty, it will be auto-generated (suitable for combination with TAP_TO_VPP interface from Linux ifplugin, because then this name is only temporary anyway) |
| to_microservice | [string](#string) |  | If TAP connects VPP with microservice, fill this parameter with the target microservice name - should match with the namespace reference of the associated TAP_TO_VPP interface (it is still moved to the namespace by Linux-ifplugin but VPP-ifplugin needs to be aware of this dependency) |
| rx_ring_size | [uint32](#uint32) |  | Rx ring buffer size; must be power of 2; default is 256; only for TAP v.2 |
| tx_ring_size | [uint32](#uint32) |  | Tx ring buffer size; must be power of 2; default is 256; only for TAP v.2 |
| enable_gso | [bool](#bool) |  | EnableGso enables GSO mode for TAP interface. |
| enable_tunnel | [bool](#bool) |  | EnableTunnel enables tunnel mode for TAP interface. |






<a name="ligato.vpp.interfaces.VmxNet3Link"></a>

### VmxNet3Link
VmxNet3Link defines configuration for interface type: VMXNET3_INTERFACE
PCI address (unsigned 32bit int) is derived from vmxnet3 interface name. It is expected that the interface
name is in format `vmxnet3-<d>/<b>/<s>/<f>`, where `d` stands for domain (max ffff), `b` is bus (max ff),
`s` is slot (max 1f) and `f` is function (max 7). All values are base 16


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| enable_elog | [bool](#bool) |  | Turn on elog |
| rxq_size | [uint32](#uint32) |  | Receive queue size (default is 1024) |
| txq_size | [uint32](#uint32) |  | Transmit queue size (default is 1024) |






<a name="ligato.vpp.interfaces.VxlanLink"></a>

### VxlanLink
VxlanLink defines configuration for interface type: VXLAN_TUNNEL


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| src_address | [string](#string) |  | SrcAddress is source VTEP address |
| dst_address | [string](#string) |  | DstAddress is destination VTEP address |
| vni | [uint32](#uint32) |  | Vni stands for VXLAN Network Identifier |
| multicast | [string](#string) |  | Multicast defines name of multicast interface |
| gpe | [VxlanLink.Gpe](#ligato.vpp.interfaces.VxlanLink.Gpe) |  |  |






<a name="ligato.vpp.interfaces.VxlanLink.Gpe"></a>

### VxlanLink.Gpe
Gpe (Generic Protocol Extension) allows encapsulating not only Ethernet frame payload.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| decap_vrf_id | [uint32](#uint32) |  |  |
| protocol | [VxlanLink.Gpe.Protocol](#ligato.vpp.interfaces.VxlanLink.Gpe.Protocol) |  | Protocol defines encapsulated protocol |






<a name="ligato.vpp.interfaces.WireguardLink"></a>

### WireguardLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| private_key | [string](#string) |  | Private-key base64 |
| port | [uint32](#uint32) |  | Listen UDP port |
| src_addr | [string](#string) |  | Source IP address |





 <!-- end messages -->


<a name="ligato.vpp.interfaces.BondLink.LoadBalance"></a>

### BondLink.LoadBalance


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| L2 | 0 |  |
| L34 | 1 |  |
| L23 | 2 |  |
| RR | 3 | Round robin |
| BC | 4 | Broadcast |
| AB | 5 | Active backup |



<a name="ligato.vpp.interfaces.BondLink.Mode"></a>

### BondLink.Mode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| ROUND_ROBIN | 1 |  |
| ACTIVE_BACKUP | 2 |  |
| XOR | 3 |  |
| BROADCAST | 4 |  |
| LACP | 5 |  |



<a name="ligato.vpp.interfaces.GreLink.Type"></a>

### GreLink.Type


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| L3 | 1 | L3 GRE (i.e. this tunnel is in L3 mode) |
| TEB | 2 | TEB - Transparent Ethernet Bridging - the tunnel is in L2 mode |
| ERSPAN | 3 | ERSPAN - the tunnel is for port mirror SPAN output |



<a name="ligato.vpp.interfaces.GtpuLink.NextNode"></a>

### GtpuLink.NextNode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| DEFAULT | 0 | The default next node is l2-input |
| L2 | 1 | l2-input |
| IP4 | 2 | ip4-input |
| IP6 | 3 | ip6-input |



<a name="ligato.vpp.interfaces.IPIPLink.Mode"></a>

### IPIPLink.Mode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| POINT_TO_POINT | 0 | point-to-point tunnel |
| POINT_TO_MULTIPOINT | 1 | point-to multipoint tunnel (supported starting from VPP 20.05) |



<a name="ligato.vpp.interfaces.IPSecLink.Mode"></a>

### IPSecLink.Mode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| POINT_TO_POINT | 0 | point-to-point tunnel |
| POINT_TO_MULTIPOINT | 1 | point-to multipoint tunnel (supported starting from VPP 20.05) |



<a name="ligato.vpp.interfaces.Interface.RxMode.Type"></a>

### Interface.RxMode.Type
Type definition is from: vpp/include/vnet/interface.h

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| POLLING | 1 |  |
| INTERRUPT | 2 |  |
| ADAPTIVE | 3 |  |
| DEFAULT | 4 |  |



<a name="ligato.vpp.interfaces.Interface.Type"></a>

### Interface.Type
Type defines VPP interface types.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_TYPE | 0 |  |
| SUB_INTERFACE | 1 |  |
| SOFTWARE_LOOPBACK | 2 |  |
| DPDK | 3 |  |
| MEMIF | 4 |  |
| TAP | 5 |  |
| AF_PACKET | 6 |  |
| VXLAN_TUNNEL | 7 |  |
| IPSEC_TUNNEL | 8 | Deprecated in VPP 20.01+. Use IPIP_TUNNEL + ipsec.TunnelProtection instead. |
| VMXNET3_INTERFACE | 9 |  |
| BOND_INTERFACE | 10 |  |
| GRE_TUNNEL | 11 |  |
| GTPU_TUNNEL | 12 |  |
| IPIP_TUNNEL | 13 |  |
| WIREGUARD_TUNNEL | 14 |  |
| RDMA | 15 |  |



<a name="ligato.vpp.interfaces.MemifLink.MemifMode"></a>

### MemifLink.MemifMode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| ETHERNET | 0 |  |
| IP | 1 |  |
| PUNT_INJECT | 2 |  |



<a name="ligato.vpp.interfaces.RDMALink.Mode"></a>

### RDMALink.Mode


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| AUTO | 0 |  |
| IBV | 1 | InfiniBand Verb (using libibverb). |
| DV | 2 | Direct Verb allows the driver to access the NIC HW RX/TX rings directly instead of having to go through libibverb and suffering associated overhead. It will be automatically selected if the adapter supports it. |



<a name="ligato.vpp.interfaces.SubInterface.TagRewriteOptions"></a>

### SubInterface.TagRewriteOptions


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| DISABLED | 0 |  |
| PUSH1 | 1 |  |
| PUSH2 | 2 |  |
| POP1 | 3 |  |
| POP2 | 4 |  |
| TRANSLATE11 | 5 |  |
| TRANSLATE12 | 6 |  |
| TRANSLATE21 | 7 |  |
| TRANSLATE22 | 8 |  |



<a name="ligato.vpp.interfaces.VxlanLink.Gpe.Protocol"></a>

### VxlanLink.Gpe.Protocol


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| IP4 | 1 |  |
| IP6 | 2 |  |
| ETHERNET | 3 |  |
| NSH | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/interfaces/state.proto"></a>

## ligato/vpp/interfaces/state.proto



<a name="ligato.vpp.interfaces.InterfaceNotification"></a>

### InterfaceNotification



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| type | [InterfaceNotification.NotifType](#ligato.vpp.interfaces.InterfaceNotification.NotifType) |  |  |
| state | [InterfaceState](#ligato.vpp.interfaces.InterfaceState) |  |  |






<a name="ligato.vpp.interfaces.InterfaceState"></a>

### InterfaceState



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  |  |
| internal_name | [string](#string) |  |  |
| type | [Interface.Type](#ligato.vpp.interfaces.Interface.Type) |  |  |
| if_index | [uint32](#uint32) |  |  |
| admin_status | [InterfaceState.Status](#ligato.vpp.interfaces.InterfaceState.Status) |  |  |
| oper_status | [InterfaceState.Status](#ligato.vpp.interfaces.InterfaceState.Status) |  |  |
| last_change | [int64](#int64) |  |  |
| phys_address | [string](#string) |  |  |
| speed | [uint64](#uint64) |  |  |
| mtu | [uint32](#uint32) |  |  |
| duplex | [InterfaceState.Duplex](#ligato.vpp.interfaces.InterfaceState.Duplex) |  |  |
| statistics | [InterfaceState.Statistics](#ligato.vpp.interfaces.InterfaceState.Statistics) |  |  |






<a name="ligato.vpp.interfaces.InterfaceState.Statistics"></a>

### InterfaceState.Statistics



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| in_packets | [uint64](#uint64) |  |  |
| in_bytes | [uint64](#uint64) |  |  |
| out_packets | [uint64](#uint64) |  |  |
| out_bytes | [uint64](#uint64) |  |  |
| drop_packets | [uint64](#uint64) |  |  |
| punt_packets | [uint64](#uint64) |  |  |
| ipv4_packets | [uint64](#uint64) |  |  |
| ipv6_packets | [uint64](#uint64) |  |  |
| in_nobuf_packets | [uint64](#uint64) |  |  |
| in_miss_packets | [uint64](#uint64) |  |  |
| in_error_packets | [uint64](#uint64) |  |  |
| out_error_packets | [uint64](#uint64) |  |  |






<a name="ligato.vpp.interfaces.InterfaceStats"></a>

### InterfaceStats



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  |  |
| rx | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| tx | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| rx_unicast | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| rx_multicast | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| rx_broadcast | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| tx_unicast | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| tx_multicast | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| tx_broadcast | [InterfaceStats.CombinedCounter](#ligato.vpp.interfaces.InterfaceStats.CombinedCounter) |  |  |
| rx_error | [uint64](#uint64) |  |  |
| tx_error | [uint64](#uint64) |  |  |
| rx_no_buf | [uint64](#uint64) |  |  |
| rx_miss | [uint64](#uint64) |  |  |
| drops | [uint64](#uint64) |  |  |
| punts | [uint64](#uint64) |  |  |
| ip4 | [uint64](#uint64) |  |  |
| ip6 | [uint64](#uint64) |  |  |
| mpls | [uint64](#uint64) |  |  |






<a name="ligato.vpp.interfaces.InterfaceStats.CombinedCounter"></a>

### InterfaceStats.CombinedCounter



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| packets | [uint64](#uint64) |  |  |
| bytes | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="ligato.vpp.interfaces.InterfaceNotification.NotifType"></a>

### InterfaceNotification.NotifType


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| UPDOWN | 1 |  |
| COUNTERS | 2 |  |



<a name="ligato.vpp.interfaces.InterfaceState.Duplex"></a>

### InterfaceState.Duplex


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN_DUPLEX | 0 |  |
| HALF | 1 |  |
| FULL | 2 |  |



<a name="ligato.vpp.interfaces.InterfaceState.Status"></a>

### InterfaceState.Status


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN_STATUS | 0 |  |
| UP | 1 |  |
| DOWN | 2 |  |
| DELETED | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/interfaces/span.proto"></a>

## ligato/vpp/interfaces/span.proto



<a name="ligato.vpp.interfaces.Span"></a>

### Span



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface_from | [string](#string) |  |  |
| interface_to | [string](#string) |  |  |
| direction | [Span.Direction](#ligato.vpp.interfaces.Span.Direction) |  |  |
| is_l2 | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="ligato.vpp.interfaces.Span.Direction"></a>

### Span.Direction


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| RX | 1 |  |
| TX | 2 |  |
| BOTH | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/interfaces/dhcp.proto"></a>

## ligato/vpp/interfaces/dhcp.proto



<a name="ligato.vpp.interfaces.DHCPLease"></a>

### DHCPLease
DHCPLease is a notification, i.e. flows from SB upwards


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface_name | [string](#string) |  |  |
| host_name | [string](#string) |  |  |
| is_ipv6 | [bool](#bool) |  |  |
| host_phys_address | [string](#string) |  |  |
| host_ip_address | [string](#string) |  | IP addresses in the format <ipAddress>/<ipPrefix> |
| router_ip_address | [string](#string) |  | IP addresses in the format <ipAddress>/<ipPrefix> |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/dns/dns.proto"></a>

## ligato/vpp/dns/dns.proto



<a name="ligato.vpp.dns.DNSCache"></a>

### DNSCache
DNSCache configuration models VPP's DNS cache server functionality. The main goal of this functionality is
to cache DNS records and minimize external DNS traffic.
The presence of this configuration enables the VPP DNS functionality and VPP start to acts as DNS cache Server.
It responds on standard DNS port(53) to DNS requests. Removing of this configuration disables the VPP DNS
functionality.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| upstream_dns_servers | [string](#string) | repeated | List of upstream DNS servers that are contacted by VPP when unknown domain name needs to be resolved. The results are cached and there should be no further upstream DNS server request for the same domain name until cached DNS record expiration. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/acl/acl.proto"></a>

## ligato/vpp/acl/acl.proto



<a name="ligato.vpp.acl.ACL"></a>

### ACL
ACL defines Access Control List.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | The name of an access list. A device MAY restrict the length and value of this name, possibly spaces and special characters are not allowed. |
| rules | [ACL.Rule](#ligato.vpp.acl.ACL.Rule) | repeated |  |
| interfaces | [ACL.Interfaces](#ligato.vpp.acl.ACL.Interfaces) |  |  |






<a name="ligato.vpp.acl.ACL.Interfaces"></a>

### ACL.Interfaces
The set of interfaces that has assigned this ACL on ingres or egress.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| egress | [string](#string) | repeated |  |
| ingress | [string](#string) | repeated |  |






<a name="ligato.vpp.acl.ACL.Rule"></a>

### ACL.Rule
List of access list entries (Rules). Each Access Control Rule has
a list of match criteria and a list of actions.
Access List entry that can define:
- IPv4/IPv6 src ip prefix
- src MAC address mask
- src MAC address value
- can be used only for static ACLs.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| action | [ACL.Rule.Action](#ligato.vpp.acl.ACL.Rule.Action) |  |  |
| ip_rule | [ACL.Rule.IpRule](#ligato.vpp.acl.ACL.Rule.IpRule) |  |  |
| macip_rule | [ACL.Rule.MacIpRule](#ligato.vpp.acl.ACL.Rule.MacIpRule) |  |  |






<a name="ligato.vpp.acl.ACL.Rule.IpRule"></a>

### ACL.Rule.IpRule



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ip | [ACL.Rule.IpRule.Ip](#ligato.vpp.acl.ACL.Rule.IpRule.Ip) |  |  |
| icmp | [ACL.Rule.IpRule.Icmp](#ligato.vpp.acl.ACL.Rule.IpRule.Icmp) |  |  |
| tcp | [ACL.Rule.IpRule.Tcp](#ligato.vpp.acl.ACL.Rule.IpRule.Tcp) |  |  |
| udp | [ACL.Rule.IpRule.Udp](#ligato.vpp.acl.ACL.Rule.IpRule.Udp) |  |  |






<a name="ligato.vpp.acl.ACL.Rule.IpRule.Icmp"></a>

### ACL.Rule.IpRule.Icmp



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| icmpv6 | [bool](#bool) |  | ICMPv6 flag, if false ICMPv4 will be used |
| icmp_code_range | [ACL.Rule.IpRule.Icmp.Range](#ligato.vpp.acl.ACL.Rule.IpRule.Icmp.Range) |  | Inclusive range representing icmp codes to be used. |
| icmp_type_range | [ACL.Rule.IpRule.Icmp.Range](#ligato.vpp.acl.ACL.Rule.IpRule.Icmp.Range) |  |  |






<a name="ligato.vpp.acl.ACL.Rule.IpRule.Icmp.Range"></a>

### ACL.Rule.IpRule.Icmp.Range



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| first | [uint32](#uint32) |  |  |
| last | [uint32](#uint32) |  |  |






<a name="ligato.vpp.acl.ACL.Rule.IpRule.Ip"></a>

### ACL.Rule.IpRule.Ip
IP  used in this Access List Entry.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| destination_network | [string](#string) |  | Destination IPv4/IPv6 network address (<ip>/<network>) |
| source_network | [string](#string) |  | Destination IPv4/IPv6 network address (<ip>/<network>) |
| protocol | [uint32](#uint32) |  | IP protocol number (http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml) Zero value (i.e. undefined protocol) means that the protocol to match will be automatically selected from one of the ICMP/ICMP6/TCP/UDP based on the rule definition. For example, if "icmp" is defined and src/dst addresses are IPv6 then packets of the ICMP6 protocol will be matched, etc. |






<a name="ligato.vpp.acl.ACL.Rule.IpRule.PortRange"></a>

### ACL.Rule.IpRule.PortRange
Inclusive range representing destination ports to be used. When
only lower-port is present, it represents a single port.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| lower_port | [uint32](#uint32) |  |  |
| upper_port | [uint32](#uint32) |  | If upper port is set, it must be greater or equal to lower port |






<a name="ligato.vpp.acl.ACL.Rule.IpRule.Tcp"></a>

### ACL.Rule.IpRule.Tcp



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| destination_port_range | [ACL.Rule.IpRule.PortRange](#ligato.vpp.acl.ACL.Rule.IpRule.PortRange) |  |  |
| source_port_range | [ACL.Rule.IpRule.PortRange](#ligato.vpp.acl.ACL.Rule.IpRule.PortRange) |  |  |
| tcp_flags_mask | [uint32](#uint32) |  | Binary mask for tcp flags to match. MSB order (FIN at position 0). Applied as logical AND to tcp flags field of the packet being matched, before it is compared with tcp-flags-value. |
| tcp_flags_value | [uint32](#uint32) |  | Binary value for tcp flags to match. MSB order (FIN at position 0). Before tcp-flags-value is compared with tcp flags field of the packet being matched, tcp-flags-mask is applied to packet field value. |






<a name="ligato.vpp.acl.ACL.Rule.IpRule.Udp"></a>

### ACL.Rule.IpRule.Udp



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| destination_port_range | [ACL.Rule.IpRule.PortRange](#ligato.vpp.acl.ACL.Rule.IpRule.PortRange) |  |  |
| source_port_range | [ACL.Rule.IpRule.PortRange](#ligato.vpp.acl.ACL.Rule.IpRule.PortRange) |  |  |






<a name="ligato.vpp.acl.ACL.Rule.MacIpRule"></a>

### ACL.Rule.MacIpRule



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| source_address | [string](#string) |  |  |
| source_address_prefix | [uint32](#uint32) |  |  |
| source_mac_address | [string](#string) |  | Before source-mac-address is compared with source mac address field of the packet being matched, source-mac-address-mask is applied to packet field value. |
| source_mac_address_mask | [string](#string) |  | Source MAC address mask. Applied as logical AND with source mac address field of the packet being matched, before it is compared with source-mac-address. |





 <!-- end messages -->


<a name="ligato.vpp.acl.ACL.Rule.Action"></a>

### ACL.Rule.Action


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| DENY | 0 |  |
| PERMIT | 1 |  |
| REFLECT | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/vpp/abf/abf.proto"></a>

## ligato/vpp/abf/abf.proto



<a name="ligato.vpp.abf.ABF"></a>

### ABF
ABF defines ACL based forwarding.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| index | [uint32](#uint32) |  | ABF index (unique identifier) |
| acl_name | [string](#string) |  | Name of the associated access list |
| attached_interfaces | [ABF.AttachedInterface](#ligato.vpp.abf.ABF.AttachedInterface) | repeated |  |
| forwarding_paths | [ABF.ForwardingPath](#ligato.vpp.abf.ABF.ForwardingPath) | repeated |  |






<a name="ligato.vpp.abf.ABF.AttachedInterface"></a>

### ABF.AttachedInterface
List of interfaces attached to the ABF


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| input_interface | [string](#string) |  |  |
| priority | [uint32](#uint32) |  |  |
| is_ipv6 | [bool](#bool) |  |  |






<a name="ligato.vpp.abf.ABF.ForwardingPath"></a>

### ABF.ForwardingPath
List of forwarding paths added to the ABF policy (via)


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| next_hop_ip | [string](#string) |  |  |
| interface_name | [string](#string) |  |  |
| weight | [uint32](#uint32) |  |  |
| preference | [uint32](#uint32) |  |  |
| dvr | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/netalloc/netalloc.proto"></a>

## ligato/netalloc/netalloc.proto



<a name="ligato.netalloc.ConfigData"></a>

### ConfigData
ConfigData wraps all configuration items exported by netalloc.
TBD: MACs, VXLAN VNIs, memif IDs, etc.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ip_addresses | [IPAllocation](#ligato.netalloc.IPAllocation) | repeated |  |






<a name="ligato.netalloc.IPAllocation"></a>

### IPAllocation
IPAllocation represents a single allocated IP address.

To reference allocated address, instead of entering specific IP address
for interface/route/ARP/..., use one of the following string templates
prefixed with netalloc keyword "alloc" followed by colon:
 a) reference IP address allocated for an interface:
       "alloc:<network_name>/<interface_name>"
 b) when interface is given (e.g. when asked for IP from interface model),
    interface_name can be omitted:
       "alloc:<network_name>"
 c) reference default gateway IP address assigned to an interface:
       "alloc:<network_name>/<interface_name>/GW"
 d) when asking for GW IP for interface which is given, interface_name
    can be omitted:
       "alloc:<network_name>/GW"


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| network_name | [string](#string) |  | NetworkName is some label assigned to the network where the IP address was assigned to the given interface. In theory, interface can have multiple IP adresses or there can be multiple address allocators and the network name allows to separate them. The network name is not allowed to contain forward slashes. |
| interface_name | [string](#string) |  | InterfaceName is the logical VPP or Linux interface name for which the address is allocated. |
| address | [string](#string) |  | Address is an IP addres allocated to the interface inside the given network. If the address is specified without a mask, the all-ones mask (/32 for IPv4, /128 for IPv6) will be assumed. |
| gw | [string](#string) |  | Gw is the address of the default gateway assigned to the interface in the given network. If the address is specified without a mask, then either: a) the mask of the <address> is used provided that GW IP falls into the same network IP range, or b) the all-ones mask is used otherwise |





 <!-- end messages -->


<a name="ligato.netalloc.IPAddressForm"></a>

### IPAddressForm
IPAddressForm can be used in descriptors whose models reference allocated IP
addresses, to ask for a specific form in which the address should applied.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_FORM | 0 |  |
| ADDR_ONLY | 1 | ADDR_ONLY = apply address without mask, e.g. 192.168.2.5 |
| ADDR_WITH_MASK | 2 | ADDR_WITH_MASK = apply address including the mask of the network, e.g. 192.168.2.5/24 |
| ADDR_NET | 3 | ADDR_NET = apply network implied by the address, e.g. for 192.168.2.10/24 apply 192.168.2.0/24 |
| SINGLE_ADDR_NET | 4 | SINGLE_ADDR_NET = apply address with an all-ones mask (i.e. /32 for IPv4, /128 for IPv6) |



<a name="ligato.netalloc.IPAddressSource"></a>

### IPAddressSource
IPAddressSource can be used to remember the source of an IP address.
(e.g. to distinguish allocated IP addresses from statically defined ones)

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_SOURCE | 0 |  |
| STATIC | 1 | STATIC is IP address statically assigned in the NB configuration. |
| FROM_DHCP | 2 | FROM_DHCP is set when IP address is obtained from DHCP. |
| ALLOC_REF | 3 | ALLOC_REF is a reference inside NB configuration to an allocated IP address. |
| EXISTING | 4 | EXISTING is set when IP address is assigned to (EXISTING) interface externally (i.e. by a different agent or manually by an administrator). |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/punt/punt.proto"></a>

## ligato/linux/punt/punt.proto



<a name="ligato.linux.punt.PortBased"></a>

### PortBased
Define network socket type


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| l4_protocol | [PortBased.L4Protocol](#ligato.linux.punt.PortBased.L4Protocol) |  |  |
| l3_protocol | [PortBased.L3Protocol](#ligato.linux.punt.PortBased.L3Protocol) |  |  |
| port | [uint32](#uint32) |  |  |






<a name="ligato.linux.punt.Proxy"></a>

### Proxy
Proxy allows to listen on network socket or unix domain socket, and resend to another network/unix domain socket


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| rx_port | [PortBased](#ligato.linux.punt.PortBased) |  |  |
| rx_socket | [SocketBased](#ligato.linux.punt.SocketBased) |  |  |
| tx_port | [PortBased](#ligato.linux.punt.PortBased) |  |  |
| tx_socket | [SocketBased](#ligato.linux.punt.SocketBased) |  |  |






<a name="ligato.linux.punt.SocketBased"></a>

### SocketBased
Define unix domain socket type for IPC


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| path | [string](#string) |  |  |





 <!-- end messages -->


<a name="ligato.linux.punt.PortBased.L3Protocol"></a>

### PortBased.L3Protocol
L3 protocol

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_L3 | 0 |  |
| IPV4 | 1 |  |
| IPV6 | 2 |  |
| ALL | 3 |  |



<a name="ligato.linux.punt.PortBased.L4Protocol"></a>

### PortBased.L4Protocol
L4 protocol

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED_L4 | 0 |  |
| TCP | 6 |  |
| UDP | 17 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/namespace/namespace.proto"></a>

## ligato/linux/namespace/namespace.proto



<a name="ligato.linux.namespace.NetNamespace"></a>

### NetNamespace



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| type | [NetNamespace.ReferenceType](#ligato.linux.namespace.NetNamespace.ReferenceType) |  |  |
| reference | [string](#string) |  | Reference defines reference specific to the namespace type: * namespace ID (NSID) * PID number (PID) * file path (FD) * microservice label (MICROSERVICE) |





 <!-- end messages -->


<a name="ligato.linux.namespace.NetNamespace.ReferenceType"></a>

### NetNamespace.ReferenceType


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED | 0 |  |
| NSID | 1 | named namespace |
| PID | 2 | namespace of a given process |
| FD | 3 | namespace referenced by a file handle |
| MICROSERVICE | 4 | namespace of a docker container running given microservice |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/linux.proto"></a>

## ligato/linux/linux.proto



<a name="ligato.linux.ConfigData"></a>

### ConfigData



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interfaces | [interfaces.Interface](#ligato.linux.interfaces.Interface) | repeated |  |
| arp_entries | [l3.ARPEntry](#ligato.linux.l3.ARPEntry) | repeated |  |
| routes | [l3.Route](#ligato.linux.l3.Route) | repeated |  |






<a name="ligato.linux.Notification"></a>

### Notification



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [interfaces.InterfaceNotification](#ligato.linux.interfaces.InterfaceNotification) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/l3/route.proto"></a>

## ligato/linux/l3/route.proto



<a name="ligato.linux.l3.Route"></a>

### Route



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| outgoing_interface | [string](#string) |  | Outgoing interface logical name (mandatory). |
| scope | [Route.Scope](#ligato.linux.l3.Route.Scope) |  | The scope of the area where the link is valid. |
| dst_network | [string](#string) |  | Destination network address in the format <address>/<prefix> (mandatory) Address can be also allocated via netalloc plugin and referenced here, see: api/models/netalloc/netalloc.proto |
| gw_addr | [string](#string) |  | Gateway IP address (without mask, optional). Address can be also allocated via netalloc plugin and referenced here, see: api/models/netalloc/netalloc.proto |
| metric | [uint32](#uint32) |  | routing metric (weight) |





 <!-- end messages -->


<a name="ligato.linux.l3.Route.Scope"></a>

### Route.Scope


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED | 0 |  |
| GLOBAL | 1 |  |
| SITE | 2 |  |
| LINK | 3 |  |
| HOST | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/l3/arp.proto"></a>

## ligato/linux/l3/arp.proto



<a name="ligato.linux.l3.ARPEntry"></a>

### ARPEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  |  |
| ip_address | [string](#string) |  |  |
| hw_address | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/iptables/iptables.proto"></a>

## ligato/linux/iptables/iptables.proto



<a name="ligato.linux.iptables.RuleChain"></a>

### RuleChain



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | logical name of the rule chain across all configured rule chains (mandatory) |
| namespace | [ligato.linux.namespace.NetNamespace](#ligato.linux.namespace.NetNamespace) |  | network namespace in which this rule chain is applied |
| interfaces | [string](#string) | repeated | list of interfaces referred by the rules (optional) |
| protocol | [RuleChain.Protocol](#ligato.linux.iptables.RuleChain.Protocol) |  | protocol (address family) of the rule chain |
| table | [RuleChain.Table](#ligato.linux.iptables.RuleChain.Table) |  | table the rule chain belongs to |
| chain_type | [RuleChain.ChainType](#ligato.linux.iptables.RuleChain.ChainType) |  | type of the chain |
| chain_name | [string](#string) |  | name of the chain, used only for chains with CUSTOM chain_type |
| default_policy | [RuleChain.Policy](#ligato.linux.iptables.RuleChain.Policy) |  | default policy of the chain. Used for FILTER tables only. |
| rules | [string](#string) | repeated | ordered list of strings containing the match and action part of the rules, e.g. "-i eth0 -s 192.168.0.1 -j ACCEPT" |





 <!-- end messages -->


<a name="ligato.linux.iptables.RuleChain.ChainType"></a>

### RuleChain.ChainType


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| CUSTOM | 0 |  |
| INPUT | 1 |  |
| OUTPUT | 2 |  |
| FORWARD | 3 |  |
| PREROUTING | 4 |  |
| POSTROUTING | 5 |  |



<a name="ligato.linux.iptables.RuleChain.Policy"></a>

### RuleChain.Policy


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| NONE | 0 |  |
| ACCEPT | 1 |  |
| DROP | 2 |  |
| QUEUE | 3 |  |
| RETURN | 4 |  |



<a name="ligato.linux.iptables.RuleChain.Protocol"></a>

### RuleChain.Protocol


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| IPV4 | 0 |  |
| IPV6 | 1 |  |



<a name="ligato.linux.iptables.RuleChain.Table"></a>

### RuleChain.Table


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| FILTER | 0 |  |
| NAT | 1 |  |
| MANGLE | 2 |  |
| RAW | 3 |  |
| SECURITY | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/interfaces/state.proto"></a>

## ligato/linux/interfaces/state.proto



<a name="ligato.linux.interfaces.InterfaceNotification"></a>

### InterfaceNotification



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| type | [InterfaceNotification.NotifType](#ligato.linux.interfaces.InterfaceNotification.NotifType) |  |  |
| state | [InterfaceState](#ligato.linux.interfaces.InterfaceState) |  |  |






<a name="ligato.linux.interfaces.InterfaceState"></a>

### InterfaceState



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  |  |
| internal_name | [string](#string) |  |  |
| type | [Interface.Type](#ligato.linux.interfaces.Interface.Type) |  |  |
| if_index | [int32](#int32) |  |  |
| admin_status | [InterfaceState.Status](#ligato.linux.interfaces.InterfaceState.Status) |  |  |
| oper_status | [InterfaceState.Status](#ligato.linux.interfaces.InterfaceState.Status) |  |  |
| last_change | [int64](#int64) |  |  |
| phys_address | [string](#string) |  |  |
| speed | [uint64](#uint64) |  |  |
| mtu | [uint32](#uint32) |  |  |
| statistics | [InterfaceState.Statistics](#ligato.linux.interfaces.InterfaceState.Statistics) |  |  |






<a name="ligato.linux.interfaces.InterfaceState.Statistics"></a>

### InterfaceState.Statistics



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| in_packets | [uint64](#uint64) |  |  |
| in_bytes | [uint64](#uint64) |  |  |
| out_packets | [uint64](#uint64) |  |  |
| out_bytes | [uint64](#uint64) |  |  |
| drop_packets | [uint64](#uint64) |  |  |
| in_error_packets | [uint64](#uint64) |  |  |
| out_error_packets | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="ligato.linux.interfaces.InterfaceNotification.NotifType"></a>

### InterfaceNotification.NotifType


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN | 0 |  |
| UPDOWN | 1 |  |



<a name="ligato.linux.interfaces.InterfaceState.Status"></a>

### InterfaceState.Status


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNKNOWN_STATUS | 0 |  |
| UP | 1 |  |
| DOWN | 2 |  |
| DELETED | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/linux/interfaces/interface.proto"></a>

## ligato/linux/interfaces/interface.proto



<a name="ligato.linux.interfaces.Interface"></a>

### Interface



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Name is mandatory field representing logical name for the interface. It must be unique across all configured interfaces. |
| type | [Interface.Type](#ligato.linux.interfaces.Interface.Type) |  | Type represents the type of interface and It must match with actual Link. |
| namespace | [ligato.linux.namespace.NetNamespace](#ligato.linux.namespace.NetNamespace) |  | Namespace is a reference to a Linux network namespace where the interface should be put into. |
| host_if_name | [string](#string) |  | Name of the interface in the host OS. If not set, the host name will be the same as the interface logical name. |
| enabled | [bool](#bool) |  | Enabled controls if the interface should be UP. |
| ip_addresses | [string](#string) | repeated | IPAddresses define list of IP addresses for the interface and must be defined in the following format: <ipAddress>/<ipPrefix>. Interface IP address can be also allocated via netalloc plugin and referenced here, see: api/models/netalloc/netalloc.proto |
| phys_address | [string](#string) |  | PhysAddress represents physical address (MAC) of the interface. Random address will be assigned if left empty. Not used (and not supported) by VRF devices. |
| mtu | [uint32](#uint32) |  | MTU is the maximum transmission unit value. |
| veth | [VethLink](#ligato.linux.interfaces.VethLink) |  | VETH-specific configuration |
| tap | [TapLink](#ligato.linux.interfaces.TapLink) |  | TAP_TO_VPP-specific configuration |
| vrf_dev | [VrfDevLink](#ligato.linux.interfaces.VrfDevLink) |  | VRF_DEVICE-specific configuration |
| link_only | [bool](#bool) |  | Configure/Resync link only. IP/MAC addresses are expected to be configured externally - i.e. by a different agent or manually via CLI. |
| vrf_master_interface | [string](#string) |  | Reference to the logical name of a VRF_DEVICE interface. If defined, this interface will be enslaved to the VRF device and will thus become part of the VRF (L3-level separation) that the device represents. Interfaces enslaved to the same VRF_DEVICE master interface therefore comprise single VRF with a separate routing table. |






<a name="ligato.linux.interfaces.TapLink"></a>

### TapLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vpp_tap_if_name | [string](#string) |  | Logical name of the VPP TAP interface (mandatory for TAP_TO_VPP) |






<a name="ligato.linux.interfaces.VethLink"></a>

### VethLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| peer_if_name | [string](#string) |  | Name of the VETH peer, i.e. other end of the linux veth (mandatory for VETH) |
| rx_checksum_offloading | [VethLink.ChecksumOffloading](#ligato.linux.interfaces.VethLink.ChecksumOffloading) |  | Checksum offloading - Rx side (enabled by default) |
| tx_checksum_offloading | [VethLink.ChecksumOffloading](#ligato.linux.interfaces.VethLink.ChecksumOffloading) |  | Checksum offloading - Tx side (enabled by default) |






<a name="ligato.linux.interfaces.VrfDevLink"></a>

### VrfDevLink



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| routing_table | [uint32](#uint32) |  | Routing table associated with the VRF. Table ID is an 8-bit unsigned integer value. Please note that 253, 254 and 255 are reserved values for special routing tables (main, default, local). Multiple VRFs inside the same network namespace should each use a different routing table. For more information, visit: http://linux-ip.net/html/routing-tables.html |





 <!-- end messages -->


<a name="ligato.linux.interfaces.Interface.Type"></a>

### Interface.Type


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED | 0 |  |
| VETH | 1 |  |
| TAP_TO_VPP | 2 | TAP created by VPP to have the Linux-side further configured |
| LOOPBACK | 3 | LOOPBACK is used to attach configuration to an existing "lo" interface, but unlike EXISTING type it is not limited to the default network namespace (i.e. loopbacks in other containers can be referenced also). To create an additional interface which effectively acts as a loopback, use DUMMY interface (see below). |
| EXISTING | 4 | Wait for and potentially attach additional network configuration to an interface created externally (i.e. not by this agent) in the default network namespace (i.e. same as used by the agent). Behaviour of the EXISTING interface depends on the values of ip_addresses and link_only attributes as follows: 1. link_only=false and ip_addresses are empty: agent waits for interface to be created externally and then configures it in the L2-only mode (resync will remove any IP addresses configured from outside of the agent) 2. link_only=false and ip_addresses are non-empty: agent waits for interface to be created externally and then attaches the selected IP addresses to it (resync removes any other IPs added externally) 3. link_only=true and ip_addresses are empty: agent only waits for the interface to exists (it doesn't wait for or change any IP addresses attached to it) 4. link_only=true and ip_addresses are non empty: agent waits for the interface to exists and the selected IP addresses to be assigned (i.e. there will be derived value for each expected IP address in the PENDING state until the address is assigned to the interface externally) |
| VRF_DEVICE | 5 | In Linux, VRF is implemented as yet another type of netdevice (i.e. listed with `ip link show`). Network interfaces are then assigned to VRF simply by enslaving them to the VRF device. For more information, visit: https://www.kernel.org/doc/Documentation/networking/vrf.txt |
| DUMMY | 6 | Create a dummy Linux interface which effectively behaves just like the loopback. |



<a name="ligato.linux.interfaces.VethLink.ChecksumOffloading"></a>

### VethLink.ChecksumOffloading


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| CHKSM_OFFLOAD_DEFAULT | 0 |  |
| CHKSM_OFFLOAD_ENABLED | 1 |  |
| CHKSM_OFFLOAD_DISABLED | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/kvscheduler/value_status.proto"></a>

## ligato/kvscheduler/value_status.proto



<a name="ligato.kvscheduler.BaseValueStatus"></a>

### BaseValueStatus



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| value | [ValueStatus](#ligato.kvscheduler.ValueStatus) |  |  |
| derived_values | [ValueStatus](#ligato.kvscheduler.ValueStatus) | repeated |  |






<a name="ligato.kvscheduler.ValueStatus"></a>

### ValueStatus



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| key | [string](#string) |  |  |
| state | [ValueState](#ligato.kvscheduler.ValueState) |  |  |
| error | [string](#string) |  | error returned by the last operation (none if empty string) |
| last_operation | [TxnOperation](#ligato.kvscheduler.TxnOperation) |  |  |
| details | [string](#string) | repeated | - for invalid value, details is a list of invalid fields - for pending value, details is a list of missing dependencies (labels) |





 <!-- end messages -->


<a name="ligato.kvscheduler.TxnOperation"></a>

### TxnOperation


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNDEFINED | 0 |  |
| VALIDATE | 1 |  |
| CREATE | 2 |  |
| UPDATE | 3 |  |
| DELETE | 4 |  |



<a name="ligato.kvscheduler.ValueState"></a>

### ValueState


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| NONEXISTENT | 0 | ValueState_NONEXISTENT is assigned to value that was deleted or has never existed. |
| MISSING | 1 | ValueState_MISSING is assigned to NB value that was configured but refresh found it to be missing. |
| UNIMPLEMENTED | 2 | ValueState_UNIMPLEMENTED marks value received from NB that cannot be configured because there is no registered descriptor associated with it. |
| REMOVED | 3 | ValueState_REMOVED is assigned to NB value after it was removed or when it is being re-created. The state is only temporary: for re-create, the value transits to whatever state the following Create operation produces, and delete values are removed from the graph (go to the NONEXISTENT state) immediately after the notification about the state change is sent. |
| CONFIGURED | 4 | ValueState_CONFIGURED marks value defined by NB and successfully configured. |
| OBTAINED | 5 | ValueState_OBTAINED marks value not managed by NB, instead created automatically or externally in SB. The KVScheduler learns about the value either using Retrieve() or through a SB notification. |
| DISCOVERED | 6 | ValueState_DISCOVERED marks NB value that was found (=retrieved) by refresh but not actually configured by the agent in this run. |
| PENDING | 7 | ValueState_PENDING represents (NB) value that cannot be configured yet due to missing dependencies. |
| INVALID | 8 | ValueState_INVALID represents (NB) value that will not be configured because it has a logically invalid content as declared by the Validate method of the associated descriptor. The corresponding error and the list of affected fields are stored in the <InvalidValueDetails> structure available via <details> for invalid value. |
| FAILED | 9 | ValueState_FAILED marks (NB) value for which the last executed operation returned an error. The error and the type of the operation which caused the error are stored in the <FailedValueDetails> structure available via <details> for failed value. |
| RETRYING | 10 | ValueState_RETRYING marks unsucessfully applied (NB) value, for which, however, one or more attempts to fix the error by repeating the last operation are planned, and only if all the retries fail, the value will then transit to the FAILED state. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/govppmux/metrics.proto"></a>

## ligato/govppmux/metrics.proto



<a name="ligato.govppmux.Metrics"></a>

### Metrics



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| channels_created | [uint64](#uint64) |  |  |
| channels_open | [uint64](#uint64) |  |  |
| requests_sent | [uint64](#uint64) |  |  |
| requests_done | [uint64](#uint64) |  |  |
| requests_fail | [uint64](#uint64) |  |  |
| request_replies | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/generic/options.proto"></a>

## ligato/generic/options.proto


 <!-- end messages -->

 <!-- end enums -->


<a name="ligato/generic/options.proto-extensions"></a>

### File-level Extensions
| Extension     | Type               | Base        | Number      | Description                                     |
| ------------- | ------------------ | ----------- | ----------- | ----------------------------------------------- |
| model | ModelSpec | .google.protobuf.MessageOptions | 50222 |  |

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/generic/model.proto"></a>

## ligato/generic/model.proto



<a name="ligato.generic.ModelDetail"></a>

### ModelDetail
ModelDetail represents info about model details.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| spec | [ModelSpec](#ligato.generic.ModelSpec) |  | Spec is a specificaiton the model was registered with. |
| proto_name | [string](#string) |  | ProtoName is a name of protobuf message representing the model. |
| options | [ModelDetail.Option](#ligato.generic.ModelDetail.Option) | repeated |  |






<a name="ligato.generic.ModelDetail.Option"></a>

### ModelDetail.Option



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| key | [string](#string) |  |  |
| values | [string](#string) | repeated |  |






<a name="ligato.generic.ModelSpec"></a>

### ModelSpec
ModelSpec defines a model specification to identify a model.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| module | [string](#string) |  | Module describes grouping for the model. |
| version | [string](#string) |  | Version describes version of the model schema. |
| type | [string](#string) |  | Type describes name of type described by this model. |
| class | [string](#string) |  | Class describes purpose for the model. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ligato/generic/meta.proto"></a>

## ligato/generic/meta.proto



<a name="ligato.generic.KnownModelsRequest"></a>

### KnownModelsRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| class | [string](#string) |  |  |






<a name="ligato.generic.KnownModelsResponse"></a>

### KnownModelsResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| known_models | [ModelDetail](#ligato.generic.ModelDetail) | repeated |  |
| active_modules | [string](#string) | repeated |  |






<a name="ligato.generic.ProtoFileDescriptorRequest"></a>

### ProtoFileDescriptorRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| full_proto_file_name | [string](#string) |  | full_proto_file_name is full name of proto file that is needed to identify it. It has the form "<proto package name ('.' replaced with '/')>/<simple file name>" (i.e. for this proto model it is "ligato/generic/meta.proto"). If you are using rpc ProtoFileDescriptor for additional information retrieve for known models from rpc KnownModels call, you can use usually present ModelDetail's generic.ModelDetail_Option for key "protoFile" that is containing full proto file name in correct format. |






<a name="ligato.generic.ProtoFileDescriptorResponse"></a>

### ProtoFileDescriptorResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| file_descriptor | [google.protobuf.FileDescriptorProto](#google.protobuf.FileDescriptorProto) |  | file_descriptor is proto message representing proto file descriptor |
| file_import_descriptors | [google.protobuf.FileDescriptorSet](#google.protobuf.FileDescriptorSet) |  | file_import_descriptors is set of file descriptors that the file_descriptor is using as import. This is needed when converting file descriptor proto to protoreflect.FileDescriptor (using "google.golang.org/protobuf/reflect/protodesc".NewFile(...) ) |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ligato.generic.MetaService"></a>

### MetaService
MetaService defines the RPC methods for managing generic models.

| Method Name        | Request Type        | Response Type         | Description                                  |
| ------------------ | ------------------- | --------------------- | ---------------------------------------------|
| KnownModels | [KnownModelsRequest](#ligato.generic.KnownModelsRequest) | [KnownModelsResponse](#ligato.generic.KnownModelsResponse) | KnownModels returns information about service capabilities including list of models supported by the server. |
| ProtoFileDescriptor | [ProtoFileDescriptorRequest](#ligato.generic.ProtoFileDescriptorRequest) | [ProtoFileDescriptorResponse](#ligato.generic.ProtoFileDescriptorResponse) | ProtoFileDescriptor returns proto file descriptor for proto file identified by full name. The proto file descriptor is in form of proto messages (file descriptor proto and proto of its imports) so there are needed additional steps to join them into protoreflect.FileDescriptor ("google.golang.org/protobuf/reflect/protodesc".NewFile(...)).

This rpc can be used together with knownModels rpc to retrieve additional model information. Message descriptor can be retrieved from file descriptor corresponding to knownModel message and used with proto reflecting to get all kinds of information about the known model.

Due to nature of data retrieval, it is expected that at least one message from that proto file is registered as known model. |

 <!-- end services -->



<a name="ligato/generic/manager.proto"></a>

## ligato/generic/manager.proto



<a name="ligato.generic.ConfigItem"></a>

### ConfigItem



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| item | [Item](#ligato.generic.Item) |  |  |
| status | [ItemStatus](#ligato.generic.ItemStatus) |  |  |
| labels | [ConfigItem.LabelsEntry](#ligato.generic.ConfigItem.LabelsEntry) | repeated |  |






<a name="ligato.generic.ConfigItem.LabelsEntry"></a>

### ConfigItem.LabelsEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="ligato.generic.Data"></a>

### Data
Data represents encoded data for an item.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| any | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ligato.generic.DumpStateRequest"></a>

### DumpStateRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ids | [Item.ID](#ligato.generic.Item.ID) | repeated |  |






<a name="ligato.generic.DumpStateResponse"></a>

### DumpStateResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| items | [StateItem](#ligato.generic.StateItem) | repeated |  |






<a name="ligato.generic.GetConfigRequest"></a>

### GetConfigRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| ids | [Item.ID](#ligato.generic.Item.ID) | repeated |  |






<a name="ligato.generic.GetConfigResponse"></a>

### GetConfigResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| items | [ConfigItem](#ligato.generic.ConfigItem) | repeated |  |






<a name="ligato.generic.Item"></a>

### Item
Item represents single instance described by the Model.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| id | [Item.ID](#ligato.generic.Item.ID) |  |  |
| data | [Data](#ligato.generic.Data) |  |  |






<a name="ligato.generic.Item.ID"></a>

### Item.ID
ID represents identifier for distinguishing items.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| model | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="ligato.generic.ItemStatus"></a>

### ItemStatus
Item status describes status of an item.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| status | [string](#string) |  |  |
| message | [string](#string) |  |  |






<a name="ligato.generic.Notification"></a>

### Notification



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| item | [Item](#ligato.generic.Item) |  |  |
| status | [ItemStatus](#ligato.generic.ItemStatus) |  |  |






<a name="ligato.generic.SetConfigRequest"></a>

### SetConfigRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| updates | [UpdateItem](#ligato.generic.UpdateItem) | repeated |  |
| overwrite_all | [bool](#bool) |  | The overwrite_all can be set to true to overwrite all other configuration (this is also known as Full Resync) |






<a name="ligato.generic.SetConfigResponse"></a>

### SetConfigResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| results | [UpdateResult](#ligato.generic.UpdateResult) | repeated |  |






<a name="ligato.generic.StateItem"></a>

### StateItem



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| item | [Item](#ligato.generic.Item) |  |  |
| metadata | [StateItem.MetadataEntry](#ligato.generic.StateItem.MetadataEntry) | repeated |  |






<a name="ligato.generic.StateItem.MetadataEntry"></a>

### StateItem.MetadataEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="ligato.generic.SubscribeRequest"></a>

### SubscribeRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| subscriptions | [Subscription](#ligato.generic.Subscription) | repeated |  |






<a name="ligato.generic.SubscribeResponse"></a>

### SubscribeResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| notifications | [Notification](#ligato.generic.Notification) | repeated |  |






<a name="ligato.generic.Subscription"></a>

### Subscription



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| id | [Item.ID](#ligato.generic.Item.ID) |  |  |






<a name="ligato.generic.UpdateItem"></a>

### UpdateItem



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| item | [Item](#ligato.generic.Item) |  | The item describes item to be updated. For a delete operation set fields item.Data to nil. |
| labels | [UpdateItem.LabelsEntry](#ligato.generic.UpdateItem.LabelsEntry) | repeated | The labels can be used to define user-defined labels for item. |






<a name="ligato.generic.UpdateItem.LabelsEntry"></a>

### UpdateItem.LabelsEntry



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="ligato.generic.UpdateResult"></a>

### UpdateResult



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| id | [Item.ID](#ligato.generic.Item.ID) |  |  |
| key | [string](#string) |  |  |
| op | [UpdateResult.Operation](#ligato.generic.UpdateResult.Operation) |  |  |
| status | [ItemStatus](#ligato.generic.ItemStatus) |  |  |





 <!-- end messages -->


<a name="ligato.generic.UpdateResult.Operation"></a>

### UpdateResult.Operation


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNSPECIFIED | 0 |  |
| CREATE | 1 |  |
| UPDATE | 2 |  |
| DELETE | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ligato.generic.ManagerService"></a>

### ManagerService
ManagerService defines the RPC methods for managing config
using generic model, allowing extending with custom models.

| Method Name        | Request Type        | Response Type         | Description                                  |
| ------------------ | ------------------- | --------------------- | ---------------------------------------------|
| SetConfig | [SetConfigRequest](#ligato.generic.SetConfigRequest) | [SetConfigResponse](#ligato.generic.SetConfigResponse) | SetConfig is used to update desired configuration. |
| GetConfig | [GetConfigRequest](#ligato.generic.GetConfigRequest) | [GetConfigResponse](#ligato.generic.GetConfigResponse) | GetConfig is used to read the desired configuration. |
| DumpState | [DumpStateRequest](#ligato.generic.DumpStateRequest) | [DumpStateResponse](#ligato.generic.DumpStateResponse) | DumpState is used to retrieve the actual running state. |
| Subscribe | [SubscribeRequest](#ligato.generic.SubscribeRequest) | [SubscribeResponse](#ligato.generic.SubscribeResponse) stream | Subscribe is used for subscribing to events. Notifications are returned by streaming updates. |

 <!-- end services -->



<a name="ligato/configurator/statspoller.proto"></a>

## ligato/configurator/statspoller.proto



<a name="ligato.configurator.PollStatsRequest"></a>

### PollStatsRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| period_sec | [uint32](#uint32) |  | PeriodSec defines polling period (in seconds). Set to zero to return just single polling. |
| num_polls | [uint32](#uint32) |  | NumPolls defines number of pollings. Set to non-zero number to stop the polling after specified number of pollings is reached. |






<a name="ligato.configurator.PollStatsResponse"></a>

### PollStatsResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| poll_seq | [uint32](#uint32) |  | PollSeq defines the sequence number of this polling response. |
| stats | [Stats](#ligato.configurator.Stats) |  | Stats contains polled stats data. |






<a name="ligato.configurator.Stats"></a>

### Stats
Stats defines stats data returned by StatsPollerService.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vpp_stats | [ligato.vpp.Stats](#ligato.vpp.Stats) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ligato.configurator.StatsPollerService"></a>

### StatsPollerService
StatsPollerService provides operations for collecting statistics.

| Method Name        | Request Type        | Response Type         | Description                                  |
| ------------------ | ------------------- | --------------------- | ---------------------------------------------|
| PollStats | [PollStatsRequest](#ligato.configurator.PollStatsRequest) | [PollStatsResponse](#ligato.configurator.PollStatsResponse) stream | PollStats is used for polling stats with specific period and number of pollings. |

 <!-- end services -->



<a name="ligato/configurator/configurator.proto"></a>

## ligato/configurator/configurator.proto



<a name="ligato.configurator.Config"></a>

### Config
Config describes all supported configs into a single config message.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vpp_config | [ligato.vpp.ConfigData](#ligato.vpp.ConfigData) |  |  |
| linux_config | [ligato.linux.ConfigData](#ligato.linux.ConfigData) |  |  |
| netalloc_config | [ligato.netalloc.ConfigData](#ligato.netalloc.ConfigData) |  |  |






<a name="ligato.configurator.DeleteRequest"></a>

### DeleteRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| delete | [Config](#ligato.configurator.Config) |  | Delete is a config data to be deleted. |
| wait_done | [bool](#bool) |  | WaitDone option can be used to block until either config delete is done (non-pending) or request times out.

NOTE: WaitDone is intended to be used for config updates that depend on some event from dataplane to fully configure. Using this with incomplete config updates will require another update request to unblock. |






<a name="ligato.configurator.DeleteResponse"></a>

### DeleteResponse







<a name="ligato.configurator.DumpRequest"></a>

### DumpRequest







<a name="ligato.configurator.DumpResponse"></a>

### DumpResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| dump | [Config](#ligato.configurator.Config) |  | Dump is a running config. |






<a name="ligato.configurator.GetRequest"></a>

### GetRequest







<a name="ligato.configurator.GetResponse"></a>

### GetResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| config | [Config](#ligato.configurator.Config) |  | Config describes desired config retrieved from agent. |






<a name="ligato.configurator.Notification"></a>

### Notification
Notification describes all known notifications into a single message.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vpp_notification | [ligato.vpp.Notification](#ligato.vpp.Notification) |  |  |
| linux_notification | [ligato.linux.Notification](#ligato.linux.Notification) |  |  |






<a name="ligato.configurator.NotifyRequest"></a>

### NotifyRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| idx | [uint32](#uint32) |  |  |
| filters | [Notification](#ligato.configurator.Notification) | repeated |  |






<a name="ligato.configurator.NotifyResponse"></a>

### NotifyResponse



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| next_idx | [uint32](#uint32) |  | Index of next notification |
| notification | [Notification](#ligato.configurator.Notification) |  | Notification contains notification data. |






<a name="ligato.configurator.UpdateRequest"></a>

### UpdateRequest



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| update | [Config](#ligato.configurator.Config) |  | Update is a config data to be updated. |
| full_resync | [bool](#bool) |  | FullResync option can be used to overwrite all existing config with config update.

NOTE: Using FullResync with empty config update will remove all existing config. |
| wait_done | [bool](#bool) |  | WaitDone option can be used to block until either config update is done (non-pending) or request times out.

NOTE: WaitDone is intended to be used for config updates that depend on some event from dataplane to fully configure. Using this with incomplete config updates will require another update request to unblock. |






<a name="ligato.configurator.UpdateResponse"></a>

### UpdateResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ligato.configurator.ConfiguratorService"></a>

### ConfiguratorService
ConfiguratorService provides basic operations for managing configuration
and monitoring actual state.

| Method Name        | Request Type        | Response Type         | Description                                  |
| ------------------ | ------------------- | --------------------- | ---------------------------------------------|
| Get | [GetRequest](#ligato.configurator.GetRequest) | [GetResponse](#ligato.configurator.GetResponse) | Get is used for listing desired config. |
| Update | [UpdateRequest](#ligato.configurator.UpdateRequest) | [UpdateResponse](#ligato.configurator.UpdateResponse) | Update is used for updating desired config. |
| Delete | [DeleteRequest](#ligato.configurator.DeleteRequest) | [DeleteResponse](#ligato.configurator.DeleteResponse) | Delete is used for deleting desired config. |
| Dump | [DumpRequest](#ligato.configurator.DumpRequest) | [DumpResponse](#ligato.configurator.DumpResponse) | Dump is used for dumping running config. |
| Notify | [NotifyRequest](#ligato.configurator.NotifyRequest) | [NotifyResponse](#ligato.configurator.NotifyResponse) stream | Notify is used for subscribing to notifications. |

 <!-- end services -->



<a name="ligato/annotations.proto"></a>

## ligato/annotations.proto



<a name="ligato.LigatoOptions"></a>

### LigatoOptions



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| type | [LigatoOptions.Type](#ligato.LigatoOptions.Type) |  |  |
| int_range | [LigatoOptions.IntRange](#ligato.LigatoOptions.IntRange) |  |  |






<a name="ligato.LigatoOptions.IntRange"></a>

### LigatoOptions.IntRange



| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| minimum | [int64](#int64) |  |  |
| maximum | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="ligato.LigatoOptions.Type"></a>

### LigatoOptions.Type


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| UNSPECIFIED | 0 |  |
| IP | 1 |  |
| IPV4 | 2 |  |
| IPV6 | 3 |  |
| IP_WITH_MASK | 4 |  |
| IPV4_WITH_MASK | 5 |  |
| IPV6_WITH_MASK | 6 |  |
| IP_OPTIONAL_MASK | 7 |  |
| IPV4_OPTIONAL_MASK | 8 |  |
| IPV6_OPTIONAL_MASK | 9 |  |


 <!-- end enums -->


<a name="ligato/annotations.proto-extensions"></a>

### File-level Extensions
| Extension     | Type               | Base        | Number      | Description                                     |
| ------------- | ------------------ | ----------- | ----------- | ----------------------------------------------- |
| ligato_options | LigatoOptions | .google.protobuf.FieldOptions | 2000 | NOTE: used option field index(2000) is in extension index range of descriptor.proto, but is not registered in protobuf global extension registry (https://github.com/protocolbuffers/protobuf/blob/master/docs/options.md) |

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="nat64/nat64.proto"></a>

## nat64/nat64.proto



<a name="nat64.Nat64AddressPool"></a>

### Nat64AddressPool
Nat44AddressPool defines an address pool used for NAT64.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | VRF id of tenant, 0xFFFFFFFF means independent of VRF. Non-zero (and not all-ones) VRF has to be explicitly created (see proto/ligato/vpp/l3/vrf.proto). |
| first_ip | [string](#string) |  | First IP address of the pool. |
| last_ip | [string](#string) |  | Last IP address of the pool. Should be higher than first_ip or empty. |






<a name="nat64.Nat64IPv6Prefix"></a>

### Nat64IPv6Prefix
IPv4-Embedded IPv6 Address Prefix used for NAT64.
If no prefix is configured (at all or for a given VRF), then the well-known prefix (64:ff9b::/96) is used.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | VRF id of tenant. At most one IPv6 prefix can be configured for a given VRF (that's why VRF is part of the key but prefix is not). Non-zero (and not all-ones) VRF has to be explicitly created (see proto/ligato/vpp/l3/vrf.proto). |
| prefix | [string](#string) |  | NAT64 prefix in the <IPv6-Address>/<IPv6-Prefix> format. |






<a name="nat64.Nat64Interface"></a>

### Nat64Interface
Nat64Interface defines a local network interfaces enabled for NAT64.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| name | [string](#string) |  | Interface name (logical). |
| type | [Nat64Interface.Type](#nat64.Nat64Interface.Type) |  |  |






<a name="nat64.Nat64StaticBIB"></a>

### Nat64StaticBIB
Static NAT64 binding allowing IPv4 host from the outside to access IPv6 host from the inside.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| vrf_id | [uint32](#uint32) |  | VRF (table) ID. Non-zero VRF has to be explicitly created (see proto/ligato/vpp/l3/vrf.proto). |
| inside_ipv6_address | [string](#string) |  | IPv6 host from the inside/local network. |
| inside_port | [uint32](#uint32) |  | Inside port number (of the IPv6 host). |
| outside_ipv4_address | [string](#string) |  | IPv4 host from the outside/external network. |
| outside_port | [uint32](#uint32) |  | Outside port number (of the IPv4 host). |
| protocol | [Nat64StaticBIB.Protocol](#nat64.Nat64StaticBIB.Protocol) |  |  |





 <!-- end messages -->


<a name="nat64.Nat64Interface.Type"></a>

### Nat64Interface.Type


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| IPV6_INSIDE | 0 | Interface connecting inside/local network with IPv6 endpoints. |
| IPV4_OUTSIDE | 1 | Interface connecting outside/external network with IPv4 endpoints. |



<a name="nat64.Nat64StaticBIB.Protocol"></a>

### Nat64StaticBIB.Protocol
Protocol to which the binding applies.

| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| TCP | 0 |  |
| UDP | 1 |  |
| ICMP | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="isisx/isisx.proto"></a>

## isisx/isisx.proto



<a name="vpp.isisx.ISISXConnection"></a>

### ISISXConnection
Unidirectional cross-connection between 2 interfaces that will cross-connect only ISIS protocol data traffic


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| input_interface | [string](#string) |  | Name of input interface |
| output_interface | [string](#string) |  | Name of outgoing interface |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="bfd/bfd.proto"></a>

## bfd/bfd.proto



<a name="bfd.BFD"></a>

### BFD
Single-hop UDP-based bidirectional forwarding detection session


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  | Name of the interface the BFD session is attached to. |
| local_ip | [string](#string) |  | Local IP address. The interface must have the same address configured. |
| peer_ip | [string](#string) |  | IP address of the peer, must be the same IP version as the local address. |
| min_tx_interval | [uint32](#uint32) |  | Desired minimum TX interval in milliseconds. |
| min_rx_interval | [uint32](#uint32) |  | Required minimum RX interval in milliseconds. |
| detect_multiplier | [uint32](#uint32) |  | Detect multiplier, must be non-zero value. |






<a name="bfd.BFDEvent"></a>

### BFDEvent
BFDEvent is generated whenever a BFD state changes.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| interface | [string](#string) |  |  |
| local_ip | [string](#string) |  |  |
| peer_ip | [string](#string) |  |  |
| session_state | [BFDEvent.SessionState](#bfd.BFDEvent.SessionState) |  |  |






<a name="bfd.WatchBFDEventsRequest"></a>

### WatchBFDEventsRequest
Request message for the WatchBFDEvents method.


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| subscription_label | [string](#string) |  |  |





 <!-- end messages -->


<a name="bfd.BFDEvent.SessionState"></a>

### BFDEvent.SessionState


| Name                       | Number  | Description                                                    |
| -------------------------- | ------- | -------------------------------------------------------------- |
| Unknown | 0 |  |
| Down | 1 |  |
| Init | 2 |  |
| Up | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="bfd.BFDWatcher"></a>

### BFDWatcher
BFDWatcher provides API to watch for BFD events.

| Method Name        | Request Type        | Response Type         | Description                                  |
| ------------------ | ------------------- | --------------------- | ---------------------------------------------|
| WatchBFDEvents | [WatchBFDEventsRequest](#bfd.WatchBFDEventsRequest) | [BFDEvent](#bfd.BFDEvent) stream | WatchBFDEvents allows to subscribe for BFD events. |

 <!-- end services -->



<a name="abx/abx.proto"></a>

## abx/abx.proto



<a name="vpp.abx.ABX"></a>

### ABX
ACL based xconnect


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| index | [uint32](#uint32) |  | ABX index (unique identifier) |
| acl_name | [string](#string) |  | Name of the associated access list |
| output_interface | [string](#string) |  | Name of outgoing interface |
| dst_mac | [string](#string) |  | Rewrite destination mac address |
| attached_interfaces | [ABX.AttachedInterface](#vpp.abx.ABX.AttachedInterface) | repeated |  |






<a name="vpp.abx.ABX.AttachedInterface"></a>

### ABX.AttachedInterface
List of interfaces attached to the ABX


| Field                | Type                | Label     | Description                                      |
| -------------------- | ------------------- | --------- | ------------------------------------------------ |
| input_interface | [string](#string) |  |  |
| priority | [uint32](#uint32) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes                               | C++          | Java          | Python         | Go           |
| ----------- | ----------------------------------- | ------------ | ------------- | -------------- | ------------ |
| <a name="double" /> double |  | double | double | float | float64 |
| <a name="float" /> float |  | float | float | float | float32 |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte |

// SPDX-License-Identifier: Apache-2.0

// Copyright 2021 PANTHEON.tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// backwardCompatibleNames is mapping from dynamic Config fields (derived from currently known models) to
// hardcoded names (proto field name/json name) in hardcoded configurator.Config. This mapping should allow
// dynamically-created Config to read/write configuration from/to json/yaml files in the same way as it is
// for hardcoded configurator.Config.
var backwardCompatibleNames = map[string]string{
	"netallocConfig.IPAllocation":      "ip_addresses",
	"linuxConfig.Interface":            "interfaces",
	"linuxConfig.ARPEntry":             "arp_entries",
	"linuxConfig.Route":                "routes",
	"vppConfig.ABF":                    "abfs",
	"vppConfig.ACL":                    "acls",
	"vppConfig.SecurityPolicyDatabase": "ipsec_spds",
	"vppConfig.SecurityPolicy":         "ipsec_sps",
	"vppConfig.SecurityAssociation":    "ipsec_sas",
	"vppConfig.TunnelProtection":       "ipsec_tunnel_protections",
	"vppConfig.Interface":              "interfaces",
	"vppConfig.Span":                   "spans",
	"vppConfig.IPFIX":                  "ipfix_global",
	"vppConfig.FlowProbeParams":        "ipfix_flowprobe_params",
	"vppConfig.FlowProbeFeature":       "ipfix_flowprobes",
	"vppConfig.BridgeDomain":           "bridge_domains",
	"vppConfig.FIBEntry":               "fibs",
	"vppConfig.XConnectPair":           "xconnect_pairs",
	"vppConfig.ARPEntry":               "arps",
	"vppConfig.Route":                  "routes",
	"vppConfig.ProxyARP":               "proxy_arp",
	"vppConfig.IPScanNeighbor":         "ipscan_neighbor",
	"vppConfig.VrfTable":               "vrfs",
	"vppConfig.DHCPProxy":              "dhcp_proxies",
	"vppConfig.L3XConnect":             "l3xconnects",
	"vppConfig.TeibEntry":              "teib_entries",
	"vppConfig.Nat44Global":            "nat44_global",
	"vppConfig.DNat44":                 "dnat44s",
	"vppConfig.Nat44Interface":         "nat44_interfaces",
	"vppConfig.Nat44AddressPool":       "nat44_pools",
	"vppConfig.IPRedirect":             "punt_ipredirects",
	"vppConfig.ToHost":                 "punt_tohosts",
	"vppConfig.Exception":              "punt_exceptions",
	"vppConfig.LocalSID":               "srv6_localsids",
	"vppConfig.Policy":                 "srv6_policies",
	"vppConfig.Steering":               "srv6_steerings",
	"vppConfig.SRv6Global":             "srv6_global",
	"vppConfig.Peer":                   "wg_peers",
}

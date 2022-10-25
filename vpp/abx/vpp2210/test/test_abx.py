#!/usr/bin/env python3

# SPDX-License-Identifier: Apache-2.0

# Copyright 2022 PANTHEON.tech
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import unittest

from framework import VppTestCase, VppTestRunner
from vpp_sub_interface import VppDot1QSubint
from vpp_acl import AclRule, VppAcl, VppAclInterface

from ipaddress import IPv4Network, IPv6Network

from scapy.layers.l2 import Ether, Dot1Q
from scapy.layers.inet import IP, TCP


class TestABX(VppTestCase):
    """
    Test ABX plugin

    Testing following scenarios:
        1.  Adding, updating and removing of ABX policies
        2.  Adding and removing of ABX attach configs
        3.  Basic ABX test
        4.  Basic ABX test - change destination MAC
        5.  ABX with destination sub-interface
        6.  ABX with destination sub-interface - change destination MAC
        7.  ABX with source sub-interface
        8.  ABX with source sub-interface - change MAC
        9.  ABX with source and destination sub-interface
        10. ABX with source and destination sub-interface
            - change destination MAC

    TODO:
        1. Double VLAN tagged packets forwarding scenarios
        2. IPv6 scenarios

    """

    # Traffic type
    IP = 0

    # IP version
    IPRANDOM = -1
    IPV4 = 0
    IPV6 = 1

    # rule types
    DENY = 0
    PERMIT = 1

    # Transport layer
    TCP = 0
    UDP = 1

    # Supported protocols
    proto = [[6, 17], [1, 58]]

    # Ether types
    ETHER_TYPE_IPV4 = 0x0800
    ETHER_TYPE_DOT1Q = 0x8100

    # VLAN IDs
    PG0_VLAN_100 = 100
    PG1_VLAN_200 = 200

    @classmethod
    def setUpClass(cls):
        super(TestABX, cls).setUpClass()

    @classmethod
    def tearDownClass(cls):
        super(TestABX, cls).tearDownClass()

    def setUp(self):
        super(TestABX, self).setUp()

        # create 2 pg interfaces
        self.create_pg_interfaces(range(4))

        # setup all interfaces
        for i in self.pg_interfaces:
            i.admin_up()
            i.config_ip4()
            i.config_ip6()
            i.resolve_arp()
            i.generate_remote_hosts(3)

        # Create sub-interfaces with VLANs
        self.pg0_vlan_100 = VppDot1QSubint(self, self.pg0, self.PG0_VLAN_100)
        self.pg0_vlan_100.admin_up()
        self.pg0_vlan_100.config_ip4()
        self.pg0_vlan_100.resolve_arp()
        self.pg0_vlan_100.config_ip6()

        self.pg1_vlan_200 = VppDot1QSubint(self, self.pg1, self.PG1_VLAN_200)
        self.pg1_vlan_200.admin_up()
        self.pg1_vlan_200.config_ip4()
        self.pg1_vlan_200.resolve_arp()
        self.pg1_vlan_200.config_ip6()

    def tearDown(self):
        for i in self.pg_interfaces:
            i.unconfig_ip4()
            i.unconfig_ip6()
            i.admin_down()

        self.pg0_vlan_100.unconfig_ip4()
        self.pg0_vlan_100.unconfig_ip6()
        self.pg0_vlan_100.admin_down()
        self.pg0_vlan_100.remove_vpp_config()

        self.pg1_vlan_200.unconfig_ip4()
        self.pg1_vlan_200.unconfig_ip6()
        self.pg1_vlan_200.admin_down()
        self.pg1_vlan_200.remove_vpp_config()

        super(TestABX, self).tearDown()

    def create_rule(
        self,
        ip=0,
        permit_deny=0,
        ports=-1,
        proto=-1,
        s_prefix=0,
        s_ip=0,
        d_prefix=0,
        d_ip=0,
    ):
        if ip:
            src_prefix = IPv6Network((s_ip, s_prefix))
            dst_prefix = IPv6Network((d_ip, d_prefix))
        else:
            src_prefix = IPv4Network((s_ip, s_prefix))
            dst_prefix = IPv4Network((d_ip, d_prefix))
        return AclRule(
            is_permit=permit_deny,
            ports=ports,
            proto=proto,
            src_prefix=src_prefix,
            dst_prefix=dst_prefix,
        )

    def apply_rules(self, rules, tag=None):
        acl = VppAcl(self, rules, tag=tag)
        acl.add_vpp_config()
        self.logger.debug("Dumped ACL: " + str(acl.dump()))
        # Apply a ACL on the interface as inbound
        for i in self.pg_interfaces:
            acl_if = VppAclInterface(
                self, sw_if_index=i.sw_if_index, n_input=1, acls=[acl]
            )
            acl_if.add_vpp_config()
        return acl.acl_index

    def create_stream(
        self,
        in_if,
        out_if,
        vlan_id=0,
        ttl=64,
    ):
        """
        Create input packet stream for defined interface.
        If VLAN tag if other than 0.
        """
        pkts = []
        if vlan_id == 0:
            p = (
                Ether(dst=in_if.local_mac, src=in_if.remote_mac) /
                IP(dst=out_if.remote_ip4, src=in_if.remote_ip4, ttl=ttl) /
                TCP(
                    sport=3000,
                    dport=3000,
                )
            )
        else:
            p = (
                Ether(dst=in_if.local_mac, src=in_if.remote_mac) /
                Dot1Q(vlan=vlan_id) /
                IP(dst=out_if.remote_ip4, src=in_if.remote_ip4, ttl=ttl) /
                TCP(
                    sport=3000,
                    dport=3000,
                )
            )

        pkts.append(p)
        return pkts

    def create_apply_acl(self):
        """
        Create and apply ACL rule
        Return ACL rule index
        """

        # Create ACL rule
        rules = []
        rules.append(
            self.create_rule(
                self.IPV4, self.PERMIT, 3000, self.proto[self.IP][self.TCP]
            )
        )

        # Apply rules
        return self.apply_rules(rules, "permit all")

    def create_abx_policy(self, policy_id, acl_idx, dest_if, change_mac=True):
        """
        Create and add ABX policy
        Returns ABX config as map.
        """

        policy = {
            "policy_id": policy_id,
            "acl_index": acl_idx,
            "tx_sw_if_index": dest_if.sw_if_index,
            "dst_mac": dest_if.local_mac,
        }

        if change_mac is False:
            del policy["dst_mac"]

        self.vapi.abx_policy_add_del(1, policy)
        return policy

    def attach_abx_policy(self, policy_id, rx_sw_if_index):
        """
        Attach ABX config to interface.
        Returns ABX attach config.
        """

        interface_attach = {
            "policy_id": policy_id,
            "priority": 0,
            "rx_sw_if_index": rx_sw_if_index,
        }

        self.vapi.abx_interface_attach_detach(
            is_attach=1, attach=interface_attach)
        return interface_attach

    def cleanup_abx(self, abx_policy_config, abx_attach_config):
        """
        Cleanup ABX attach config and policy.
        """

        self.vapi.abx_interface_attach_detach(
            is_attach=0, attach=abx_attach_config)
        self.vapi.abx_policy_add_del(0, abx_policy_config)

    def send_and_expect(
        self,
        intf,
        pkts,
        output,
        n_rx=None,
        worker=None,
        trace=True,
    ):
        return super().send_and_expect(intf, pkts, output, n_rx, worker, trace)

    def test_abx_basic_acl_mac_dst_sub_if(self):
        """Destination sub-interface - change destination MAC"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1_vlan_200
        )
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0.sw_if_index)

        pkts = self.create_stream(self.pg0, self.pg1)
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_DOT1Q)
        self.assertEqual(pkt_out[Ether].dst, self.pg1.local_mac)
        self.assertEqual(pkt_out[Dot1Q].vlan, self.PG1_VLAN_200)
        self.assertEqual(pkt_out[Dot1Q].type, self.ETHER_TYPE_IPV4)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl_dst_sub_if(self):
        """Destination sub-interface"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1_vlan_200, False
        )
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0.sw_if_index)

        pkts = self.create_stream(self.pg0, self.pg1)
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_DOT1Q)
        self.assertEqual(pkt_out[Ether].dst, pkt_in[Ether].dst)
        self.assertEqual(pkt_out[Dot1Q].vlan, self.PG1_VLAN_200)
        self.assertEqual(pkt_out[Dot1Q].type, self.ETHER_TYPE_IPV4)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl_mac(self):
        """Change destination MAC"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1)
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0.sw_if_index)

        pkts = self.create_stream(self.pg0, self.pg1)
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_IPV4)
        self.assertEqual(pkt_out[Ether].dst, self.pg1.local_mac)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl(self):
        """Basic test"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1, False)
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0.sw_if_index)

        pkts = self.create_stream(self.pg0, self.pg1)
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_IPV4)
        self.assertEqual(pkt_out[Ether].dst, pkt_in[Ether].dst)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl_mac_src_dst_sub_if(self):
        """Source and destination sub-interface - change destination MAC"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1_vlan_200
        )
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0_vlan_100.sw_if_index
        )

        pkts = self.create_stream(
            self.pg0_vlan_100, self.pg1_vlan_200, self.PG0_VLAN_100
        )
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_DOT1Q)
        self.assertEqual(pkt_out[Ether].dst, self.pg1.local_mac)
        self.assertEqual(pkt_out[Dot1Q].vlan, self.PG1_VLAN_200)
        self.assertEqual(pkt_out[Dot1Q].type, self.ETHER_TYPE_IPV4)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl_src_dst_sub_if(self):
        """Source and destination sub-interface"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1_vlan_200, False
        )
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0_vlan_100.sw_if_index
        )

        pkts = self.create_stream(
            self.pg0_vlan_100, self.pg1_vlan_200, self.PG0_VLAN_100
        )
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_DOT1Q)
        self.assertEqual(pkt_out[Ether].dst, pkt_in[Ether].dst)
        self.assertEqual(pkt_out[Dot1Q].vlan, self.PG1_VLAN_200)
        self.assertEqual(pkt_out[Dot1Q].type, self.ETHER_TYPE_IPV4)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl_src_sub_if(self):
        """Source sub-interface"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1, False)
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0_vlan_100.sw_if_index
        )

        pkts = self.create_stream(
            self.pg0_vlan_100, self.pg1, self.PG0_VLAN_100)
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_IPV4)
        self.assertEqual(pkt_out[Ether].dst, pkt_in[Ether].dst)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_abx_basic_acl_mac_src_sub_if(self):
        """Source sub-interface - change MAC"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1)
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0_vlan_100.sw_if_index
        )

        pkts = self.create_stream(
            self.pg0_vlan_100, self.pg1, self.PG0_VLAN_100)
        pkt_in = pkts[0]
        self.logger.debug(
            "\nIncoming Packet\n {0} \nIncoming Packet\n".format(
                pkt_in.show(dump=True))
        )

        pkts = self.send_and_expect(self.pg0, pkts, self.pg1)

        pkt_out = pkts[0]
        self.logger.debug(
            "\nOutcoming Packet\n {0} \nOutcoming Packet\n".format(
                pkt_out.show(dump=True)
            )
        )

        self.assertEqual(pkt_out[Ether].type, self.ETHER_TYPE_IPV4)
        self.assertEqual(pkt_out[Ether].dst, self.pg1.local_mac)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

    def test_add_remove_policy(self):
        """Adding, updating and removing of policies"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg0)

        policies = self.vapi.abx_policy_dump()
        self.assertEqual(len(policies), 1)
        policy = policies[0].policy

        self.assertEqual(policy.policy_id, policy_id)
        self.assertEqual(policy.acl_index, acl_idx)
        self.assertEqual(policy.tx_sw_if_index, self.pg0.sw_if_index)
        self.assertEqual(policy.dst_mac, self.pg0.local_mac)

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1, False)

        policies = self.vapi.abx_policy_dump()
        self.assertEqual(len(policies), 1)
        policy = policies[0].policy

        self.assertEqual(policy.policy_id, policy_id)
        self.assertEqual(policy.acl_index, acl_idx)
        self.assertEqual(policy.tx_sw_if_index, self.pg1.sw_if_index)
        self.assertEqual(policy.dst_mac, "00:00:00:00:00:00")

        self.vapi.abx_policy_add_del(0, abx_policy_config)

        policies = self.vapi.abx_policy_dump()
        self.assertEqual(len(policies), 0)

    def test_add_remove_attach(self):
        """Adding and removing of attach configs"""

        acl_idx = self.create_apply_acl()
        policy_id = 1

        abx_policy_config = self.create_abx_policy(
            policy_id, acl_idx, self.pg1)
        abx_attach_config = self.attach_abx_policy(
            policy_id, self.pg0.sw_if_index)

        attach_configs = self.vapi.abx_interface_attach_dump()
        self.assertEqual(len(attach_configs), 1)
        attach = attach_configs[0].attach

        self.assertEqual(attach.policy_id, policy_id)
        self.assertEqual(attach.priority, 0)
        self.assertEqual(attach.rx_sw_if_index, self.pg0.sw_if_index)

        self.cleanup_abx(abx_policy_config, abx_attach_config)

        attach_configs = self.vapi.abx_interface_attach_dump()
        self.assertEqual(len(attach_configs), 0)


if __name__ == "__main__":
    unittest.main(testRunner=VppTestRunner)

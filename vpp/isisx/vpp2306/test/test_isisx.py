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

from framework import tag_fixme_vpp_workers
from framework import VppTestCase, VppTestRunner

import sys

from vpp_isis import VppIsis
from scapy.all import *
from scapy.contrib.isis import *


@tag_fixme_vpp_workers
class TestIsisx(VppTestCase):
    """ ISIS Test Case """

    @classmethod
    def setUpClass(cls):
        super(TestIsisx, cls).setUpClass()
    @classmethod
    def tearDownClass(cls):
        super(TestIsisx, cls).tearDownClass()

    def setUp(self):
        super(TestIsisx, self).setUp()

        self.create_pg_interfaces(range(3))

        for pg in self.pg_interfaces:
            pg.admin_up()
            pg.resolve_arp()

    def tearDown(self):
        for pg in self.pg_interfaces:
            pg.admin_down()            

        super(TestIsisx, self).tearDown()

    def send(self, ti, pkts):
        ti.add_stream(pkts)
        self.pg_enable_capture(self.pg_interfaces)
        self.pg_start()

    def test_isisx_connection_create(self):
        vpp = VppIsis(self,
                      self.pg0.sw_if_index,
                      self.pg1.sw_if_index)

        vpp.add_vpp_config()

        dump = vpp.query_vpp_config()
        self.assertTrue(dump)

    def test_isisx_connection_delete(self):
        vpp = VppIsis(self,
                      self.pg0.sw_if_index,
                      self.pg1.sw_if_index)

        vpp.add_vpp_config()
        vpp.remove_vpp_config()
        dump = vpp.query_vpp_config()
        self.assertFalse(dump)

    def test_isisx_connection_deleted_if(self):
        pass

    def test_isisx_connection_send_packets(self):
        vpp = VppIsis(self,
                      self.pg0.sw_if_index,
                      self.pg1.sw_if_index)

        vpp.add_vpp_config()

        self.vapi.cli("trace add virtio-input 1")

        p_isis1 = (Dot3(dst=self.pg0.local_mac, src=self.pg0.remote_mac, len=1500) /
                  LLC(dsap=0xfe, ssap=0xfe, ctrl=3) /
                  ISIS_CommonHdr() /
                  ISIS_L1_LAN_Hello())

        self.send(self.pg0, p_isis1)
        capture = self.pg1.get_capture(1, timeout=10)

        self.assertEqual(capture[0][LLC].dsap, 0xfe)
        self.assertEqual(capture[0][LLC].ssap, 0xfe)
        self.assertEqual(capture[0][LLC].ctrl, 3)

        trace = self.vapi.cli("show trace")
        self.assertTrue("OSI isis" in trace)
        self.assertTrue("isisx: rx_sw_if_index: {}".format(self.pg0.sw_if_index) in trace)
        self.assertTrue("{}-tx".format(self.pg1._name) in trace)

        self.vapi.cli("clear trace")
        self.vapi.cli("trace add virtio-input 1")

        p_isis2 = (Dot3(dst=self.pg0.local_mac, src=self.pg0.remote_mac, len=1500) /
                  LLC(dsap=0xfe, ssap=0xfe, ctrl=3) /
                  ISIS_CommonHdr() /
                  ISIS_L2_LAN_Hello())

        self.send(self.pg0, p_isis2)
        capture = self.pg1.get_capture(1, timeout=10)

        self.assertEqual(capture[0][LLC].dsap, 0xfe)
        self.assertEqual(capture[0][LLC].ssap, 0xfe)
        self.assertEqual(capture[0][LLC].ctrl, 3)

        trace = self.vapi.cli("show trace")
        self.assertTrue("OSI isis" in trace)
        self.assertTrue("isisx: rx_sw_if_index: {}".format(self.pg0.sw_if_index) in trace)
        self.assertTrue("{}-tx".format(self.pg1._name) in trace)


if __name__ == '__main__':
    unittest.main(testRunner=VppTestRunner)

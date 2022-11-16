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

from vpp_object import VppObject

def get_if_dump(dump, rx_sw_if_index):
    for d in dump:
        if (d.connection.rx_sw_if_index == rx_sw_if_index):
            return True
    return False

class VppIsis(VppObject):
    def __init__(self, test, rx_sw_if_index, tx_sw_if_index):
        self._test = test
        self.rx_sw_if_index = rx_sw_if_index
        self.tx_sw_if_index = tx_sw_if_index

    def add_vpp_config(self):
        self._test.vapi.isisx_connection_add_del(is_add=1, connection={'rx_sw_if_index': self.rx_sw_if_index, 'tx_sw_if_index': self.tx_sw_if_index})

    def remove_vpp_config(self):
        self._test.vapi.isisx_connection_add_del(is_add=0, connection={'rx_sw_if_index': self.rx_sw_if_index})

    def object_id(self):
        return "%s:%d" % (self.rx_sw_if_index, self.tx_sw_if_index)

    def query_vpp_config(self):
        dump = self._test.vapi.isisx_connection_dump()
        return get_if_dump(dump, self.rx_sw_if_index)

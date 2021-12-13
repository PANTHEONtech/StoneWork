#!/bin/bash

# SPDX-License-Identifier: Apache-2.0

# Copyright 2021 PANTHEON.tech
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

# generate the stonework config for GNS3 docker appliance

mkdir -p /etc/stonework/config/
CONFIG=/etc/stonework/config/day0-config.yaml

rm -f $CONFIG
touch $CONFIG

printf "vppConfig:\n  interfaces:\n" >> $CONFIG

awkcommand='/eth/{
  split($2, iface, ":");
  print "  - name: \"my-" iface[1] "\""
  print "    type: AF_PACKET"
  print "    enabled: false"
  print "    physAddress: \"" $(NF-2) "\""
  print "    afpacket:"
  print "      hostIfName: \"" iface[1] "\""
}
'

ip -o link | awk -F"[@ ]" "$awkcommand" >> $CONFIG


# run standard stonework container initialization

rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api
mkdir -p /run/vpp /run/stonework/vpp
service nginx start
exec stonework-init

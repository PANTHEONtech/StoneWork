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

systemctl enable serial-getty@ttyS0.service
systemctl start serial-getty@ttyS0.service
systemctl enable docker
systemctl start docker

# reserve hugepages for DPDK
echo "vm.nr_hugepages=512" >> /etc/sysctl.conf
sysctl -p

# load docker images
docker load -i /root/*-docker-img.tar
rm /root/stonework-docker-img.tar

# set all interfaces with pci addresses down so VPP can take them
lshw -class network -businfo | awk -F" " \
'/pci@/{ system("ip link set " $2 " down") }'

# generate vpp.conf
printf "\ndpdk {\n" >> /root/vpp.conf

lshw -class network -businfo | awk -F "[@ ]" \
'/pci@/{ print "    dev " $2 " {\n        name " $4 "\n    }" >> "/root/vpp.conf" }'

printf "}\n" >> /root/vpp.conf

# generate StoneWork config
lshw -class network -businfo | awk -F "[@ ]" \
'/pci@/{ print "    - name: " $4"\n      type: DPDK\n      enabled: false" >> "/root/config/day0-config.yaml" }'

# run stonework
cd /root && docker-compose up -d

#!/usr/bin/env bash

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

trap "exit 1" TERM
export TOP_PID=$$

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

function check_rv { # parameters: actual rv, expected rv, error message
    if [ $1 -ne $2 ]; then
        echo "Fail"
        echo "------------------------------------------------"
        echo -e "${RED}[FAIL] ${3}${NC}"
        echo "------------------------------------------------"
        kill -s TERM $TOP_PID
    else
        echo "OK"
    fi
}

function check_in_sync {
    echo -n "Checking if StoneWork is in-sync ... "
    docker-compose exec -T stonework agentctl config resync --verbose 2>&1 | grep -qi -E "Executed|error"
    check_rv $? 1 "StoneWork is not in-sync"
}

check_in_sync

echo -n "Checking cross-connect configuration in StoneWork ... "
[ $(docker-compose exec -T stonework agentctl values 2>/dev/null | grep -c "vpp.l2.xconnect.*CONFIGURED") -eq 2 ]
check_rv $? 0 "Cross-connect configuration missing in StoneWork" 

echo -n "Checking cross-connect in VPP ... "
[ $(docker-compose exec -T stonework vppctl sh int addr | grep -c "L2 xconnect") -eq 2 ]
check_rv $? 0 "Cross-connect not configured"

docker-compose exec -T stonework vppctl trace add virtio-input 10

echo -n "Pinging ... "
docker-compose exec -T tester1 ping -c 1 -w 1 10.10.1.2 >/dev/null
check_rv $? 0 "Ping failed"

echo -n "Checking if ping request went through StoneWork ... "
docker-compose exec -T stonework vppctl show trace | grep -q "ICMP: 10.10.1.1 -> 10.10.1.2"
check_rv $? 0 "Ping request did not go through StoneWork"

echo -n "Checking if ping response went through StoneWork ... "
docker-compose exec -T stonework vppctl show trace | grep -q "ICMP: 10.10.1.2 -> 10.10.1.1"
check_rv $? 0 "Ping response did not go through StoneWork"

echo "------------------------------------------------"
echo -e "${GREEN}[OK] Test successful${NC}"
echo "------------------------------------------------"

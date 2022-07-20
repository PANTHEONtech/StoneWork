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
    docker-compose exec -T stonework curl -X POST localhost:9191/scheduler/downstream-resync?verbose=1 2>&1 \
        | grep -qi -E "Executed|error"
    check_rv $? 1 "StoneWork is not in-sync"
}

check_in_sync

# test JSON schema
#schema=$(docker-compose exec stonework curl localhost:9191/info/configuration/jsonschema 2>/dev/null)
schema=$(curl localhost:9191/info/configuration/jsonschema 2>/dev/null)

echo -n "Checking mock CNF 1 model in JSON schema ... "
echo $schema | grep -q '"mock1Config": {'
check_rv $? 0 "Mock CNF 1 model is missing in JSON schema"

echo -n "Checking mock CNF 2 model in JSON schema ... "
echo $schema | grep -q '"mock2Config": {'
check_rv $? 0 "Mock CNF 2 model is missing in JSON schema"

echo -n "Checking route in mock CNF 1 ... "
docker-compose exec -T mockcnf1 ip route show table 1 | grep -q "7.7.7.7 dev tap"
check_rv $? 0 "Mock CNF 1 has not configured route"

echo -n "Checking ARP entry in mock CNF 2 ... "
docker-compose exec -T mockcnf2 arp -a | grep -qe "9\.9\.9\.9.*02:02:02:02:02:02"
check_rv $? 0 "Mock CNF 2 has not configured the ARP entry"

echo -n "Updating config ... "
docker-compose exec -T stonework agentctl config update --replace /etc/stonework/config/running-config.yaml >/dev/null 2>&1
check_rv $? 0 "Config update failed"

../utils.sh waitForAgentConfig stonework 74 10 # mock CNFs make changes asynchronously

check_in_sync

echo -n "Checking re-configured route in mock CNF 1 ... "
docker-compose exec -T mockcnf1 ip route show table 2 | grep -q "7.7.7.7 dev tap"
check_rv $? 0 "Mock CNF 1 has not re-configured route"

echo -n "Checking if mock CNF 2 removed ARP entry ... "
docker-compose exec -T mockcnf2 arp -a | grep -qe "9\.9\.9\.9.*02:02:02:02:02:02"
check_rv $? 1 "Mock CNF 2 has not removed the ARP entry"

echo "------------------------------------------------"
echo -e "${GREEN}[OK] Test successful${NC}"
echo "------------------------------------------------"

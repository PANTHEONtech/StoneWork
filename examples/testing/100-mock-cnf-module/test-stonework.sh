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
YELLOW='\033[0;33m'
MAGENTA='\033[0;35m'
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
    docker exec stonework curl -X POST localhost:9191/scheduler/downstream-resync?verbose=1 2>&1 \
        | grep -qi -E "Executed|error"
    check_rv $? 1 "StoneWork is not in-sync"
}

function check_container_log {
  STEP=0
  local -n SW_LOG=$1
  while true; do
    SW_LOG=$(docker logs $2 2>/dev/null)
    SEARCH_KEY="CREATE [OBTAINED]"
    COUNT=$(echo "$SW_LOG" | grep -cFe "$SEARCH_KEY")
    echo -n "$COUNT > "
#    echo -n "."
    if [ "$COUNT" -ge $3 ]; then # need >= $2 lines containing $SEARCH_KEY
      break
    fi
#    echo -n "$SW_LOG" >$3
    STEP=$((STEP + 1))
    if [ $STEP -ge 15 ]; then
      break
    fi
    sleep 2
  done
}

function check_stonework_status {
  echo -e -n "${YELLOW}Checking for StoneWork status ... ${NC}"
  STATUS_POST="curl -X GET localhost:9191/scheduler/status"
  STEP=0
  RV=0
  local -n STATUS_JSON=$1
  while true; do
    STATUS_JSON=$(docker exec stonework $STATUS_POST 2>/dev/null)
    RV=$?
    SEARCH_KEY="PENDING"
    COUNT=$(echo "$STATUS_JSON" | grep -cFe "$SEARCH_KEY")
    if [ "$COUNT" -ge 1 ]; then
      echo -n "." # "$STEP $SEARCH_KEY $COUNT"
#      echo -n "$STATUS_JSON" >$2
    else
      break
    fi
    STEP=$((STEP + 1))
    if [ $STEP -ge 15 ]; then
      break
    fi
    sleep 1
  done
  check_rv $RV 0 "$STATUS_POST"
}
check_in_sync

# test JSON schema
schema=$(docker exec stonework curl localhost:9191/info/configuration/jsonschema 2>/dev/null)

echo -n "Checking mock CNF 1 model in JSON schema ... "
echo $schema | grep -q '"mock1Config": {'
check_rv $? 0 "Mock CNF 1 model is missing in JSON schema"

echo -n "Checking mock CNF 2 model in JSON schema ... "
echo $schema | grep -q '"mock2Config": {'
check_rv $? 0 "Mock CNF 2 model is missing in JSON schema"

echo -n "Checking route in mock CNF 1 ... "
docker exec mockcnf1 ip route show table 1 | grep -q "7.7.7.7 dev tap"
check_rv $? 0 "Mock CNF 1 has not configured route"

echo -n "Checking ARP entry in mock CNF 2 ... "
docker exec mockcnf2 arp -a | grep -qe "9\.9\.9\.9.*02:02:02:02:02:02"
check_rv $? 0 "Mock CNF 2 has not configured the ARP entry"

echo -e -n "${YELLOW}Checking StoneWork logs ... ${NC}"
SW_LOG_BEFORE_CONFIG_UPDATE=""
check_container_log SW_LOG_BEFORE_CONFIG_UPDATE "stonework" 6
echo "OK"

echo -e -n "${YELLOW}Checking mockcnf1 logs ... ${NC}"
MOCK1_LOG_BEFORE_CONFIG_UPDATE=""
check_container_log MOCK1_LOG_BEFORE_CONFIG_UPDATE "mockcnf1" 6
echo "OK"

echo -e -n "${YELLOW}Checking mockcnf2 logs ... ${NC}"
MOCK2_LOG_BEFORE_CONFIG_UPDATE=""
check_container_log MOCK2_LOG_BEFORE_CONFIG_UPDATE "mockcnf2" 6
echo "OK"

STATUS_JSON_BEFORE=""
check_stonework_status STATUS_JSON_BEFORE

echo -e -n "Updating config ... "
docker exec stonework agentctl config update --replace /etc/stonework/config/running-config.yaml >/dev/null 2>&1
check_rv $? 0 "Config update failed"

echo -e -n "${MAGENTA}Checking StoneWork logs... ${NC}"
SW_LOG_AFTER_CONFIG_UPDATE=""
check_container_log SW_LOG_AFTER_UCONFIG_PDATE "stonework" 10
echo "OK"

echo -e -n "${MAGENTA}Checking mockcnf1 logs ... ${NC}"
MOCK1_LOG_AFTER_CONFIG_UPDATE=""
check_container_log MOCK1_LOG_AFTER_CONFIG_UPDATE "mockcnf1" 13
echo "OK"

echo -e -n "${MAGENTA}Checking mockcnf2 logs ... ${NC}"
MOCK2_LOG_AFTER_CONFIG_UPDATE=""
check_container_log MOCK2_LOG_AFTER_CONFIG_UPDATE "mockcnf2" 12
echo "OK"

../utils.sh waitForAgentConfig stonework 73 10 # mock CNFs make changes asynchronously

STATUS_JSON_AFTER=""
check_stonework_status STATUS_JSON_AFTER

echo -e -n "${YELLOW}Checking if StoneWork is in-sync ... ${NC}"
#check_in_sync
RESYNK_POST="curl -X POST localhost:9191/scheduler/downstream-resync?verbose=1"
RESYNK_JSON=$(docker exec stonework $RESYNK_POST 2>/dev/null)
RV=$?
grep -qi -E "Executed|error" <<< "$RESYNK_JSON"
if [ $? -ne 1 ]; then
  SW_LOG_AFTER_RESYNC=$(docker logs stonework 2>/dev/null)
  echo -n "$SW_LOG_AFTER_RESYNC" >after-downstream-resync.log
  echo -n "$RESYNK_JSON" >after-downstream-resync.json
fi
check_rv $RV 0 "$RESYNK_POST"

echo -n "Checking re-configured route in mock CNF 1 ... "
docker exec mockcnf1 ip route show table 2 | grep -q "7.7.7.7 dev tap"
check_rv $? 0 "Mock CNF 1 has not re-configured route"

echo -n "Checking if mock CNF 2 removed ARP entry ... "
docker exec mockcnf2 arp -a | grep -qe "9\.9\.9\.9.*02:02:02:02:02:02"
check_rv $? 1 "Mock CNF 2 has not removed the ARP entry"

echo "------------------------------------------------"
echo -e "${GREEN}[OK] Test successful${NC}"
echo "------------------------------------------------"

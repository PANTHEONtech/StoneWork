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

test:
	@make exec-test || (make dump-logs-ignore-errors && make stop-example-ignore-errors && exit 1)

dump-logs-ignore-errors:
	@make dump-logs || true

exec-test: stop-example-ignore-errors start-example test-cnf stop-example

start-example:
	@docker compose up -d
	@../utils.sh waitForAgentConfig mockcnf1 49 20

stop-example:
	@docker compose down --remove-orphans -v

stop-example-ignore-errors:
	@docker compose down --remove-orphans -v || true

test-cnf:
	@./test-cnf.sh

LOGFILE=./example.log

dump-logs:
	@sh -c "(\
		echo ===============================================================================;\
		echo Mock CNF logs;\
		echo ===============================================================================;\
		docker compose logs mockcnf1 || true;\
		echo ===============================================================================;\
		echo Mock CNF running configuration;\
		echo ===============================================================================;\
		docker compose exec -T mockcnf1 agentctl config retrieve || true;\
	) >${LOGFILE} 2>&1"

/*
 * Copyright (c) 2021 PANTHEON.tech s.r.o.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef __isisx_connect_H__
#define __isisx_connect_H__

#include <vnet/vnet.h>

/**
 * An ISIS protocol based Xconnect.
 * This comprises the RX interface index to match against and forward to the interface.
 *
 * ISISx connections consist of two interfaces 'attached' to each other. An input feature
 * will select the connection from the database of connections and a match will divert the packet,
 * if all miss then the packet dropped
 */

/**
 * Get an ISISx tx_sw_if_index from its rx_sw_if_index key
 *
 * @param rx_sw_if_index Client defined RX interface ID
 * @return TX interface ID
 */
const u32 isisx_get_tx_by_rx (const u32 rx_sw_if_index);

/**
 * Create or update an ISISx connection
 *
 * @param rx_sw_if_index Client defined RX interface ID
 * @param tx_sw_if_index The connection TX interface to match on
 */
void isisx_connection_update (const u32 rx_sw_if_index,
			const u32 tx_sw_if_index);

/**
 * Remove a connection from an ISISx
 *
 * @param connect_id Client defined RX interface ID
 */
int isisx_connection_delete (const u32 rx_sw_if_index);

/**
 * Callback function invoked during a walk of all ISISx connections
 */
typedef walk_rc_t (*isisx_connection_cb_t) (u32 rx_sw_if_index, u32 tx_sw_if_index, void *ctx);
	
/**
 * Walk/visit each of the ISISx connections
 */
extern void isisx_connection_walk (isisx_connection_cb_t cb, void *ctx);

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

#endif

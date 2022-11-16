/*
 * Copyright (c) 2022 PANTHEON.tech s.r.o.
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

#include <isisx/isisx.h>
#include <isisx/isisx_connect.h>

#include <vlib/vlib.h>
#include <vnet/plugin/plugin.h>
#include <vnet/fib/fib_path_list.h>
#include <vnet/fib/fib_walk.h>

const u32 isisx_get_tx_by_rx (const u32 rx_sw_if_index)
{  
  uword *ptr = NULL;
  ptr = hash_get ((&isisx_main)->isisx_connection_db, rx_sw_if_index);
  if (ptr)
  { 
    return ptr[0];
  }
  return INDEX_INVALID;
}

inline void
isisx_connection_update (const u32 rx_sw_if_index, const u32 tx_sw_if_index)
{  
  hash_set ((&isisx_main)->isisx_connection_db, rx_sw_if_index, tx_sw_if_index);
}

int
isisx_connection_delete (const u32 rx_sw_if_index)
{
  const u32 tx_sw_if_index = isisx_get_tx_by_rx (rx_sw_if_index);
  if (tx_sw_if_index == INDEX_INVALID)
  {
      /*
       * no such connection
       */
      return 0;
  }
  hash_unset ((&isisx_main)->isisx_connection_db, rx_sw_if_index);
  return 1;
}

void
isisx_connection_walk (isisx_connection_cb_t cb, void *ctx)	
{
	u32 rx_sw_if_index, tx_sw_if_index;

  /* *INDENT-OFF* */
  hash_foreach (rx_sw_if_index, tx_sw_if_index, (&isisx_main)->isisx_connection_db,
  {
    if (!cb (rx_sw_if_index, tx_sw_if_index, ctx)) break;
  });
  /* *INDENT-ON* */
}

static clib_error_t *
isisx_sw_interface_add_del (vnet_main_t * vnm, u32 sw_if_index, u32 is_add)
{
  if (is_add)
    return 0;

  u32 rx_sw_if_index, tx_sw_if_index;

  /* *INDENT-OFF* */
  hash_foreach (rx_sw_if_index, tx_sw_if_index, (&isisx_main)->isisx_connection_db,
  {
    if (rx_sw_if_index == sw_if_index || tx_sw_if_index == sw_if_index)
    {
      isisx_connection_delete(rx_sw_if_index);
    }
  });
  /* *INDENT-ON* */

  return 0;
}

VNET_SW_INTERFACE_ADD_DEL_FUNCTION (isisx_sw_interface_add_del);
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

#include <isisx/isisx.h>
#include <isisx/isisx_connect.h>

#include <vlibapi/api.h>
#include <vlibmemory/api.h>

#include <vnet/ip/ip_types_api.h>

/* define message IDs */
#include <isisx/isisx.api_enum.h>
#include <isisx/isisx.api_types.h>
#include <vnet/format_fns.h>

#include <vlibapi/api_helper_macros.h>

#undef REPLY_MSG_ID_BASE
#define REPLY_MSG_ID_BASE isisx_main.msg_id_base

/* API message handler */

static void
vl_api_isisx_plugin_get_version_t_handler (vl_api_isisx_plugin_get_version_t * mp)
{
  vl_api_isisx_plugin_get_version_reply_t *rmp;
  vl_api_registration_t *rp;

  rp = vl_api_client_index_to_registration (mp->client_index);
  if (!rp)
    return;

  rmp = vl_msg_api_alloc (sizeof (*rmp));
  rmp->_vl_msg_id =
    ntohs (VL_API_ISISX_PLUGIN_GET_VERSION_REPLY + REPLY_MSG_ID_BASE);
  rmp->context = mp->context;
  rmp->major = htonl (ISISX_PLUGIN_VERSION_MAJOR);
  rmp->minor = htonl (ISISX_PLUGIN_VERSION_MINOR);
  vl_api_send_msg (rp, (u8 *) rmp);
}

static void
vl_api_isisx_connection_add_del_t_handler (vl_api_isisx_connection_add_del_t * mp)
{
  vl_api_isisx_connection_add_del_reply_t *rmp;
  int rv = 0;

  u32 rx_sw_if_index = htonl (mp->connection.rx_sw_if_index);
  if (mp->is_add)
    {
      u32 tx_sw_if_index = htonl (mp->connection.tx_sw_if_index);  

      isisx_connection_update (rx_sw_if_index, tx_sw_if_index);
    }
  else
    {
      isisx_connection_delete (rx_sw_if_index);
    }
  REPLY_MACRO (VL_API_ISISX_CONNECTION_ADD_DEL_REPLY);
}

typedef struct isisx_connection_walk_ctx_t_
{
  vl_api_registration_t *reg;
  u32 context;
} isisx_connection_walk_ctx_t;

static walk_rc_t
isisx_connection_send_details (
  u32 rx_sw_interface_index, u32 tx_sw_interface_index, void *args)
{
  vl_api_isisx_connection_details_t *mp = vl_msg_api_alloc (sizeof (*mp));
  isisx_connection_walk_ctx_t *ctx = args;
  clib_memset (mp, 0, sizeof (*mp));

  mp->_vl_msg_id = htons (VL_API_ISISX_CONNECTION_DETAILS + REPLY_MSG_ID_BASE);
  mp->context = ctx->context;
  mp->connection.rx_sw_if_index = htonl (rx_sw_interface_index);
  mp->connection.tx_sw_if_index = htonl (tx_sw_interface_index);
  vl_api_send_msg (ctx->reg, (u8 *) mp);
  
  return WALK_CONTINUE;
}

static void
vl_api_isisx_connection_dump_t_handler (vl_api_isisx_connection_dump_t * mp)
{
  vl_api_registration_t *reg;

  reg = vl_api_client_index_to_registration (mp->client_index);
  if (!reg)
    return;

  isisx_connection_walk_ctx_t ctx = {	
    .reg = reg,
    .context = mp->context,
  };

 isisx_connection_walk (isisx_connection_send_details, &ctx);
}

/* Set up the API message handling tables */

#include <isisx/isisx.api.c>
static clib_error_t *
isisx_plugin_api_hookup (vlib_main_t * vm)
{
  isisx_main_t *pm = &isisx_main;

  /* Ask for a correctly-sized block of API message decode slots */
  pm->msg_id_base = setup_message_id_table ();

  return 0;
}

VLIB_API_INIT_FUNCTION (isisx_plugin_api_hookup);

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

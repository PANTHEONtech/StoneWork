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

#include <stddef.h>

#include <vnet/vnet.h>
#include <vnet/plugin/plugin.h>
#include <abx/abx.h>
#include <abx/abx_policy.h>
#include <abx/abx_if_attach.h>
#include <vnet/ip/ip_types_api.h>
#include <vnet/ethernet/ethernet_types_api.h>
#include <plugins/acl/acl.h>


#include <vlibapi/api.h>
#include <vlibmemory/api.h>

/* define message IDs */
#include <vnet/format_fns.h>
#include <abx/abx.api_enum.h>
#include <abx/abx.api_types.h>

/**
 * Base message ID for the plugin
 */
static u32 abx_base_msg_id;

#define REPLY_MSG_ID_BASE (abx_base_msg_id)
#include <vlibapi/api_helper_macros.h>

/* API message handler */

static void
vl_api_abx_plugin_get_version_t_handler (vl_api_abx_plugin_get_version_t * mp)
{
  vl_api_abx_plugin_get_version_reply_t *rmp;
  vl_api_registration_t *rp;

  rp = vl_api_client_index_to_registration (mp->client_index);
  if (rp == 0)
    return;

  rmp = vl_msg_api_alloc (sizeof (*rmp));
  rmp->_vl_msg_id =
    ntohs (VL_API_ABX_PLUGIN_GET_VERSION_REPLY + abx_base_msg_id);
  rmp->context = mp->context;
  rmp->major = htonl (ABX_PLUGIN_VERSION_MAJOR);
  rmp->minor = htonl (ABX_PLUGIN_VERSION_MINOR);

  vl_api_send_msg (rp, (u8 *) rmp);
}

static void
vl_api_abx_policy_add_del_t_handler (vl_api_abx_policy_add_del_t * mp)
{
  vl_api_abx_policy_add_del_reply_t *rmp;
  int rv = 0;
  u32 policy_id = htonl (mp->policy.policy_id);
  mac_address_t mac;
  
  if (mp->is_add)
    {
      u32 acl_index = htonl (mp->policy.acl_index);
      u32 tx_sw_if_index = htonl (mp->policy.tx_sw_if_index);
      mac_address_decode (mp->policy.dst_mac, &mac);

      abx_policy_update (policy_id, acl_index, tx_sw_if_index, &mac);
    }
  else
    {
      abx_policy_delete (policy_id);
    }
  REPLY_MACRO (VL_API_ABX_POLICY_ADD_DEL_REPLY);
}

static void
vl_api_abx_interface_attach_detach_t_handler (
  vl_api_abx_interface_attach_detach_t * mp)
{
  vl_api_abx_interface_attach_detach_reply_t *rmp;
  int rv = 0;
  u32 rx_sw_if_index = htonl (mp->attach.rx_sw_if_index);
  u32 priority = htonl (mp->attach.priority);
  u32 policy_id = htonl (mp->attach.policy_id);

  if (mp->is_attach)
    {
      rv = abx_if_attach (policy_id, priority, rx_sw_if_index);
    }
  else
    {
      rv = abx_if_detach (policy_id, rx_sw_if_index);
    }
  REPLY_MACRO (VL_API_ABX_INTERFACE_ATTACH_DETACH_REPLY);
}

typedef struct abx_policy_walk_ctx_t_
{
  vl_api_registration_t *reg;
  u32 context;
} abx_policy_walk_ctx_t;

static walk_rc_t
abx_policy_send_details (
  u32 api, void *args)
{
  vl_api_abx_policy_details_t *mp;
  abx_policy_walk_ctx_t *ctx;
  abx_policy_t *ap;

  ctx = args;
  ap = abx_policy_get (api);

  mp = vl_msg_api_alloc (sizeof (*mp));
  clib_memset (mp, 0, sizeof (*mp));
  mp->_vl_msg_id = htons (VL_API_ABX_POLICY_DETAILS + abx_base_msg_id);

  mp->context = ctx->context;
  mp->policy.policy_id = htonl (ap->ap_id);
  mp->policy.acl_index = htonl (ap->ap_acl);
  mp->policy.tx_sw_if_index = htonl (ap->ap_tx_sw_if_index);
  //  mac_address_encode (ap->ap_mac, &mp->policy.dst_mac);

  vl_api_send_msg (ctx->reg, (u8 *) mp);

  return (WALK_CONTINUE);
}

static void
vl_api_abx_policy_dump_t_handler (vl_api_abx_policy_dump_t * mp)
{
  vl_api_registration_t *reg;

  reg = vl_api_client_index_to_registration (mp->client_index);
  if (!reg)
    return;

  abx_policy_walk_ctx_t ctx = {
    .reg = reg,
    .context = mp->context,
  };

 abx_policy_walk (abx_policy_send_details, &ctx);
}

typedef struct abx_if_attach_walk_ctx_t_
{
  vl_api_registration_t *reg;
  u32 context;
} abx_if_attach_walk_ctx_t;

static walk_rc_t
abx_interface_attach_details (
  u32 aiai, void *args)
{
  vl_api_abx_interface_attach_details_t *mp;
  abx_if_attach_walk_ctx_t *ctx;
  abx_if_attach_t *aia;
  abx_policy_t *ap;

  ctx = args;
  aia = abx_if_attach_get (aiai);
  ap = abx_policy_get (aia->aia_abx);

  mp = vl_msg_api_alloc (sizeof (*mp));
  mp->_vl_msg_id = ntohs (VL_API_ABX_INTERFACE_ATTACH_DETAILS + abx_base_msg_id);

  mp->context = ctx->context;
  mp->attach.policy_id = htonl (ap->ap_id);
  mp->attach.priority = htonl (aia->aia_priority);
  mp->attach.rx_sw_if_index = htonl (aia->aia_sw_if_index);

  vl_api_send_msg (ctx->reg, (u8 *) mp);

  return (WALK_CONTINUE);
}

static void
vl_api_abx_interface_attach_dump_t_handler (vl_api_abx_interface_attach_dump_t * mp)
{
  vl_api_registration_t *reg;

  reg = vl_api_client_index_to_registration (mp->client_index);
  if (!reg)
    return;

  abx_if_attach_walk_ctx_t ctx = {
    .reg = reg,
    .context = mp->context,
  };

 abx_if_attach_walk (abx_interface_attach_details, &ctx);
}

#include <abx/abx.api.c>

static clib_error_t *
abx_init_api (vlib_main_t * vm)
{
  /* Ask for a correctly-sized block of API message decode slots */
  abx_base_msg_id = setup_message_id_table ();

  return 0;
}

VLIB_INIT_FUNCTION (abx_init_api);

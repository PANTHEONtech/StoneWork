/*
 * Copyright (c) 2019 PANTHEON.tech s.r.o.
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

#include <abx/abx_policy.h>

#include <vlib/vlib.h>
#include <vnet/plugin/plugin.h>
#include <vnet/fib/fib_path_list.h>
#include <vnet/fib/fib_walk.h>

/**
 * Pool of ABX objects
 */
static abx_policy_t *abx_policy_pool;

/**
 * DB of ABX policy objects
 *  - policy ID to index conversion.
 */
static uword *abx_policy_db;

inline abx_policy_t *
abx_policy_get (u32 index)
{
  return (pool_elt_at_index (abx_policy_pool, index));
}

static abx_policy_t *
abx_policy_find_id (u32 policy_id)
{
  u32 api;

  api = abx_policy_find (policy_id);

  if (INDEX_INVALID != api)
    return (abx_policy_get (api));

  return (NULL);
}

u32
abx_policy_find (u32 policy_id)
{
  uword *p;

  p = hash_get (abx_policy_db, policy_id);

  if (NULL != p)
    return (p[0]);

  return (INDEX_INVALID);
}

void
abx_policy_update (const u32 policy_id, const u32 acl_index,
		   const u32 tx_sw_if_index, const mac_address_t * mac)
{
  abx_policy_t *ap;
  u32 api;

  api = abx_policy_find (policy_id);
  if (INDEX_INVALID == api)
    {
      /*
       * create a new policy
       */
      pool_get (abx_policy_pool, ap);

      api = ap - abx_policy_pool;
      ap->ap_acl = acl_index;
      ap->ap_id = policy_id;
      ap->ap_tx_sw_if_index = tx_sw_if_index;
      if (NULL != mac && !mac_address_is_zero (mac))
	{
	  memcpy (ap->ap_dst_mac.bytes, mac, 6);
	}

      /*
       * add this new policy to the DB
       */
      hash_set (abx_policy_db, policy_id, api);
    }
  else
    {
      ap = abx_policy_get (api);
      ap->ap_acl = acl_index;
      ap->ap_id = policy_id;
      ap->ap_tx_sw_if_index = tx_sw_if_index;
    }
}

int
abx_policy_delete (const u32 policy_id)
{
  abx_policy_t *ap;
  u32 api;

  api = abx_policy_find (policy_id);

  if (INDEX_INVALID == api)
    {
      /*
       * no such policy
       */
      return (-1);
    }
  else
    {
      ap = abx_policy_get (api);
      hash_unset (abx_policy_db, ap->ap_id);
      pool_put (abx_policy_pool, ap);
    }

  return (0);
}

static clib_error_t *
abx_policy_cmd (vlib_main_t * vm,
		unformat_input_t * main_input, vlib_cli_command_t * cmd)
{
  unformat_input_t _line_input, *line_input = &_line_input;
  mac_address_t mac = ZERO_MAC_ADDRESS;
  u32 acl_index, policy_id, tx_sw_if_index;
  u32 is_del;
  vnet_main_t *vnm = vnet_get_main ();

  is_del = 0;
  acl_index = INDEX_INVALID;
  policy_id = INDEX_INVALID;
  tx_sw_if_index = INDEX_INVALID;

  /* Get a line of input. */
  if (!unformat_user (main_input, unformat_line_input, line_input))
    return 0;

  while (unformat_check_input (line_input) != UNFORMAT_END_OF_INPUT)
    {
      if (unformat (line_input, "acl %d", &acl_index))
	;
      else if (unformat (line_input, "id %d", &policy_id))
	;
      else if (unformat (line_input, "del"))
	is_del = 1;
      else if (unformat (line_input, "add"))
	is_del = 0;
      else if (unformat (line_input, "via %U",
			 unformat_vnet_sw_interface, vnm, &tx_sw_if_index))
	;
      else if (unformat (line_input, "dst-mac-rewrite %U",
			 unformat_mac_address, &mac))
	;
      else
	return (clib_error_return (0, "unknown input '%U'",
				   format_unformat_error, line_input));
    }

  if (INDEX_INVALID == policy_id)
    {
      vlib_cli_output (vm, "Specify a Policy ID");
      return 0;
    }

  if (!is_del)
    {
      if (INDEX_INVALID == acl_index)
	{
	  vlib_cli_output (vm, "ACL index must be set");
	  return 0;
	}
      if (INDEX_INVALID == tx_sw_if_index)
	{
	  vlib_cli_output (vm, "TX sw_if_index index must be set");
	  return 0;
	}
      abx_policy_update (policy_id, acl_index, tx_sw_if_index, &mac);
    }
  else
    {
      abx_policy_delete (policy_id);
    }

  unformat_free (line_input);
  return (NULL);
}

/* *INDENT-OFF* */
/**
 * Create an ABX policy.
 */
VLIB_CLI_COMMAND (abx_policy_cmd_node, static) = {
  .path = "abx policy",
  .function = abx_policy_cmd,
  .short_help = "abx policy [add|del] id <index> acl <index> via <sw_if_index> [dst-mac-rewrite <MAC>]",
  .is_mp_safe = 1,
};
/* *INDENT-ON* */

static u8 *
abx_format_mac_address (u8 * s, va_list * args)
{
  u8 *a = va_arg (*args, u8 *);
  return format (s, "%02x:%02x:%02x:%02x:%02x:%02x",
		 a[0], a[1], a[2], a[3], a[4], a[5]);
}

static u8 *
format_abx_policy (u8 * s, va_list * args)
{
  abx_policy_t *ap = va_arg (*args, abx_policy_t *);

  s = format (s, "  abx[%d]: policy:%d acl:%d via-sw-if-index: %d",
	      ap - abx_policy_pool, ap->ap_id, ap->ap_acl,
	      ap->ap_tx_sw_if_index);
  if (!mac_address_is_zero (&ap->ap_dst_mac))
    {
      s = format (s, " dst-mac-rewrite: %U",
		  abx_format_mac_address, &ap->ap_dst_mac);
    }
  return (s);
}

static clib_error_t *
abx_show_policy_cmd (vlib_main_t * vm,
		     unformat_input_t * input, vlib_cli_command_t * cmd)
{
  u32 policy_id;
  abx_policy_t *ap;

  policy_id = INDEX_INVALID;

  while (unformat_check_input (input) != UNFORMAT_END_OF_INPUT)
    {
      if (unformat (input, "%d", &policy_id))
	;
      else
	return (clib_error_return (0, "unknown input '%U'",
				   format_unformat_error, input));
    }

  if (INDEX_INVALID == policy_id)
    {
      /* *INDENT-OFF* */
      pool_foreach(ap, abx_policy_pool)
      {
        vlib_cli_output(vm, "%U", format_abx_policy, ap);
      }
      /* *INDENT-ON* */
    }
  else
    {
      ap = abx_policy_find_id (policy_id);

      if (NULL != ap)
	vlib_cli_output (vm, "%U", format_abx_policy, ap);
      else
	vlib_cli_output (vm, "Invalid policy ID:%d", policy_id);
    }

  return (NULL);
}

/* *INDENT-OFF* */
VLIB_CLI_COMMAND (abx_policy_show_policy_cmd_node, static) = {
  .path = "show abx policy",
  .function = abx_show_policy_cmd,
  .short_help = "show abx policy <policy-id>",
  .is_mp_safe = 1,
};
/* *INDENT-ON* */

void
abx_policy_walk (abx_policy_cb_t cb, void *ctx)
{
  u32 api;

  pool_foreach_index (api, abx_policy_pool)
  {
     if (!cb (api, ctx)) break;
  }
  /* *INDENT-ON* */
}

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

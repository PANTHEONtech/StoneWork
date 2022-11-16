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

#include <abx/abx_if_attach.h>
#include <plugins/acl/exports.h>

/**
 * Pool of ABX interface attachment objects
 */
abx_if_attach_t *abx_if_attach_pool;

/**
 * A per interface vector of attached policies. Used in the data-plane
 */
u32 **abx_per_if;

/**
 * Per interface values of ACL lookup context IDs. Used in the data-plane
 */
static u32 *abx_alctx_per_if;

/**
 * ABX ACL module user id returned during the initialization
 */
static u32 abx_acl_user_id;

/*
 * ACL plugin method vtable
 */

static acl_plugin_methods_t acl_plugin;

/**
 * A DB of attachments; key={abx_index,sw_if_index}
 */
static uword *abx_if_attach_db;

static u64
abx_if_attach_mk_key (u32 abx_index, u32 sw_if_index)
{
  u64 key;

  key = abx_index;
  key = key << 32;
  key |= sw_if_index;

  return (key);
}

static abx_if_attach_t *
abx_if_attach_db_find (u32 abx_index, u32 sw_if_index)
{
  uword *p;
  u64 key;

  key = abx_if_attach_mk_key (abx_index, sw_if_index);

  p = hash_get (abx_if_attach_db, key);

  if (NULL != p)
    return (pool_elt_at_index (abx_if_attach_pool, p[0]));

  return (NULL);
}

static void
abx_if_attach_db_add (u32 abx_index, u32 sw_if_index, abx_if_attach_t * aia)
{
  u64 key;

  key = abx_if_attach_mk_key (abx_index, sw_if_index);

  hash_set (abx_if_attach_db, key, aia - abx_if_attach_pool);
}

static void
abx_if_attach_db_del (u32 abx_index, u32 sw_if_index)
{
  u64 key;

  key = abx_if_attach_mk_key (abx_index, sw_if_index);

  hash_unset (abx_if_attach_db, key);
}

static int
abx_cmp_attach_for_sort (void *v1, void *v2)
{
  const abx_if_attach_t *aia1;
  const abx_if_attach_t *aia2;

  aia1 = abx_if_attach_get (*(u32 *) v1);
  aia2 = abx_if_attach_get (*(u32 *) v2);

  return (aia1->aia_priority - aia2->aia_priority);
}

void
abx_setup_acl_lc (u32 sw_if_index)
{
  u32 *acl_vec = 0;
  u32 *aiai;
  abx_if_attach_t *aia;

  if (~0 == abx_alctx_per_if[sw_if_index])
    return;

  vec_foreach (aiai, abx_per_if[sw_if_index])
  {
    aia = abx_if_attach_get (*aiai);
    vec_add1 (acl_vec, aia->aia_acl);
  }
  acl_plugin.set_acl_vec_for_context (abx_alctx_per_if[sw_if_index], acl_vec);
  vec_free (acl_vec);
}

static int
abx_if_enable_disable (u32 sw_if_index, u8 enable_disable)
{
  vnet_feature_enable_disable ("ip4-unicast", "abx-ip4",
			       sw_if_index, enable_disable, NULL, 0);
  vnet_feature_enable_disable ("ip4-multicast", "abx-ip4",
			       sw_if_index, enable_disable, NULL, 0);
  vnet_feature_enable_disable ("ip6-unicast", "abx-ip6",
			       sw_if_index, enable_disable, NULL, 0);
  vnet_feature_enable_disable ("ip6-multicast", "abx-ip6",
			       sw_if_index, enable_disable, NULL, 0);
  return 0;
}

int
abx_if_attach (u32 policy_id, u32 priority, u32 sw_if_index)
{
  abx_if_attach_t *aia;
  abx_policy_t *ap;
  u32 api;

  api = abx_policy_find (policy_id);
  if (INDEX_INVALID == api)
    return (VNET_API_ERROR_NO_SUCH_ENTRY);

  ap = abx_policy_get (api);

  /*
   * check this is not a duplicate
   */
  aia = abx_if_attach_db_find (policy_id, sw_if_index);
  if (NULL != aia)
    return (VNET_API_ERROR_ENTRY_ALREADY_EXISTS);

  /*
   * construct a new attachment object
   */
  pool_get (abx_if_attach_pool, aia);

  aia->aia_priority = priority;
  aia->aia_acl = ap->ap_acl;
  aia->aia_abx = api;
  aia->aia_sw_if_index = sw_if_index;
  abx_if_attach_db_add (policy_id, sw_if_index, aia);

  /*
   * Insert the policy on the interfaces list.
   */
  vec_validate_init_empty (abx_per_if, sw_if_index, NULL);
  vec_add1 (abx_per_if[sw_if_index], aia - abx_if_attach_pool);
  if (1 == vec_len (abx_per_if[sw_if_index]))
    {
      abx_if_enable_disable (sw_if_index, 1);

      /* if this is the first ABX policy, we need to acquire an ACL lookup context */
      vec_validate_init_empty (abx_alctx_per_if, sw_if_index, ~0);

      abx_alctx_per_if[sw_if_index] =
	acl_plugin.get_lookup_context_index (abx_acl_user_id, sw_if_index, 0);
    }
  else
    {
      vec_sort_with_function (abx_per_if[sw_if_index],
			      abx_cmp_attach_for_sort);
    }

  /* Prepare and set the list of ACLs for lookup within the context */
  abx_setup_acl_lc (sw_if_index);

  return (0);
}

int
abx_if_detach (u32 policy_id, u32 sw_if_index)
{
  abx_if_attach_t *aia;
  u32 index;
  vnet_main_t *vnm = vnet_get_main ();
  
  /*
   * check this is a valid interface
   */
  if (pool_is_free_index (vnm->interface_main.sw_interfaces,
                          sw_if_index))
    return (VNET_API_ERROR_INVALID_SW_IF_INDEX);

  /*
   * check this is a valid attachment
   */
  aia = abx_if_attach_db_find (policy_id, sw_if_index);

  if (NULL == aia)
    return (VNET_API_ERROR_NO_SUCH_ENTRY);

  /*
   * first remove from the interface's vector
   */
  ASSERT (abx_per_if[sw_if_index]);

  index = vec_search (abx_per_if[sw_if_index], aia - abx_if_attach_pool);

  ASSERT (index != ~0);
  vec_del1 (abx_per_if[sw_if_index], index);

  if (0 == vec_len (abx_per_if[sw_if_index]))
    {
      abx_if_enable_disable (sw_if_index, 0);

      /* Return the lookup context, invalidate its id in our records */
      acl_plugin.put_lookup_context_index (abx_alctx_per_if[sw_if_index]);
      abx_alctx_per_if[sw_if_index] = ~0;
    }
  else
    {
      vec_sort_with_function (abx_per_if[sw_if_index],
			      abx_cmp_attach_for_sort);
    }

  /* Prepare and set the list of ACLs for lookup within the context */
  abx_setup_acl_lc (sw_if_index);

  /*
   * remove the attachment from the DB
   */
  abx_if_attach_db_del (policy_id, sw_if_index);

  /*
   * return the object
   */
  pool_put (abx_if_attach_pool, aia);

  return (0);
}

static clib_error_t *
abx_if_attach_cmd (vlib_main_t * vm,
		    unformat_input_t * input, vlib_cli_command_t * cmd)
{
  u32 policy_id, sw_if_index;
  u32 is_del, priority;
  vnet_main_t *vnm;

  is_del = 0;
  sw_if_index = policy_id = ~0;
  vnm = vnet_get_main ();
  priority = 0;

  while (unformat_check_input (input) != UNFORMAT_END_OF_INPUT)
    {
      if (unformat (input, "del"))
	is_del = 1;
      else if (unformat (input, "add"))
	is_del = 0;
      else if (unformat (input, "policy %d", &policy_id))
	;
      else if (unformat (input, "priority %d", &priority))
	;
      else if (unformat (input, "%U",
			 unformat_vnet_sw_interface, vnm, &sw_if_index))
	;
      else
	return (clib_error_return (0, "unknown input '%U'",
				   format_unformat_error, input));
    }

  if (~0 == sw_if_index)
    return (clib_error_return (0, "invalid interface name"));
  if (~0 == policy_id)
    return (clib_error_return (0, "invalid policy ID:%d", policy_id));
  if (~0 == abx_policy_find (policy_id))
    return (clib_error_return (0, "invalid policy ID:%d", policy_id));

  if (is_del)
    abx_if_detach (policy_id, sw_if_index);
  else
    abx_if_attach (policy_id, priority, sw_if_index);

  return (NULL);
}

/* *INDENT-OFF* */
/**
 * Attach an ABX policy to an interface.
 */
VLIB_CLI_COMMAND (abx_if_attach_cmd_node, static) = {
  .path = "abx attach",
  .function = abx_if_attach_cmd,
  .short_help = "abx attach [del] policy <index> priority <prio> <interface>",
  // this is not MP safe
};
/* *INDENT-ON* */

static u8 *
format_abx_if_attach (u8 * s, va_list * args)
{
  abx_if_attach_t *aia = va_arg (*args, abx_if_attach_t *);
  abx_policy_t *ap;

  ap = abx_policy_get (aia->aia_abx);
  s = format (s, "  abx-interface-attach: policy:%d priority:%d",
	      ap->ap_id, aia->aia_priority);

  return (s);
}

static clib_error_t *
abx_show_attach_cmd (vlib_main_t * vm,
		     unformat_input_t * input, vlib_cli_command_t * cmd)
{
  const abx_if_attach_t *aia;
  u32 sw_if_index, *aiai;
  vnet_main_t *vnm;

  sw_if_index = ~0;
  vnm = vnet_get_main ();

  while (unformat_check_input (input) != UNFORMAT_END_OF_INPUT)
    {
      if (unformat (input, "%U", unformat_vnet_sw_interface, vnm, &sw_if_index))
	;
      else
	return (clib_error_return (0, "unknown input '%U'", 
            format_unformat_error, input));
    }

  if (~0 == sw_if_index)
    return (clib_error_return (0, "specify an interface"));

  vec_validate_init_empty (abx_per_if, sw_if_index, NULL);
  if (0 < vec_len (abx_per_if[sw_if_index]))
  {
    /* *INDENT-OFF* */
    vec_foreach(aiai, abx_per_if[sw_if_index])
      {
        aia = pool_elt_at_index(abx_if_attach_pool, *aiai);
        vlib_cli_output(vm, "%U", format_abx_if_attach, aia);
      }
    /* *INDENT-ON* */
  }

  return (NULL);
}

/* *INDENT-OFF* */
VLIB_CLI_COMMAND (abx_show_attach_cmd_node, static) = {
  .path = "show abx attach",
  .function = abx_show_attach_cmd,
  .short_help = "show abx attach <interface>",
  .is_mp_safe = 1,
};
/* *INDENT-ON* */

static clib_error_t *
abx_if_bond_init (vlib_main_t * vm)
{
  clib_error_t *acl_init_res = acl_plugin_exports_init (&acl_plugin);
  if (acl_init_res)
    return (acl_init_res);

  abx_acl_user_id =
    acl_plugin.register_user_module ("ABX plugin", "sw_if_index", NULL);

  return (NULL);
}

VLIB_INIT_FUNCTION (abx_if_bond_init) =
{
	.runs_after = VLIB_INITS("acl_init"),
};

inline u32 get_abx_alctx_per_if(u32 sw_if_index)
{
  /*
   * check this is a valid interface
   */
  if (pool_is_free_index (vnet_get_main ()->interface_main.sw_interfaces,
                          sw_if_index))
    return INDEX_INVALID;

  if (vec_len (abx_alctx_per_if) < sw_if_index)
    return INDEX_INVALID;

  return abx_alctx_per_if[sw_if_index];
}

inline void *get_abx_acl_main()
{
  return acl_plugin.p_acl_main;
}

void
abx_if_attach_walk (abx_if_attach_cb_t cb, void *ctx)
{
  u32 aia;

  /* *INDENT-OFF* */
  pool_foreach_index(aia, abx_if_attach_pool)
  {
    if (!cb(aia, ctx))
      break;
  }
  /* *INDENT-ON* */
}

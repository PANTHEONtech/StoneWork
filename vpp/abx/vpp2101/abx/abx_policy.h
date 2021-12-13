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

#ifndef __ABX_POLICY_H__
#define __ABX_POLICY_H__

#include <vnet/ethernet/mac_address.h>
#include <vnet/vnet.h>

/**
 * An ACL based Xconnect 'policy'.
 * This comprises the ACL index to match against and forward to interface.
 *
 * ABX policies are then 'attached' to interfaces. An input feature
 * will run through the list of policies a match will divert the packet,
 * if all miss then we continues down the interface's feature arc
 */
typedef struct abx_policy_t_
{
  /**
   * The policy ID - as configured by the client
   */
  u32 ap_id;
  /**
   * ACL index to match
   */
  u32 ap_acl;

  /**
   * Tx interface
   */
  u32 ap_tx_sw_if_index;

  /**
   * Tx destination mac address
   */
  mac_address_t ap_dst_mac;
} abx_policy_t;

/**
 * Get an ABX object from its VPP index
 */
abx_policy_t *abx_policy_get (u32 index);

/**
 * Find a ABX object from the client's policy IDwalk_rc_t
 *
 * @param policy_id Client's defined policy ID
 * @return VPP's object index
 */
u32 abx_policy_find (u32 policy_id);

/**
 * Create or update an ABX Policy
 *
 * @param policy_id User defined Policy ID
 * @param acl_index The ACL the policy with match on
 */
void abx_policy_update (const u32 policy_id, const u32 acl_index,
			const u32 tx_sw_if_index, const mac_address_t * mac);

/**
 * Delete paths from an ABX Policy. If no more paths exist, the policy
 * is deleted.
 *
 * @param policy_id User defined Policy ID
 */
int abx_policy_delete (const u32 policy_id);

/**
 * Callback function invoked during a walk of ABX all policies
 */
typedef walk_rc_t (*abx_policy_cb_t) (u32 api, void *ctx);

/**
 * Walk/visit each of the ABX policies
 */
extern void abx_policy_walk (abx_policy_cb_t cb, void *ctx);

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

#endif

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

#ifndef __ABX_ITF_ATTACH_H__
#define __ABX_ITF_ATTACH_H__

#include <abx/abx_policy.h>

/**
 * Attachment data for an ABX policy to an interface
 */
typedef struct abx_if_attach_t_
{
  CLIB_CACHE_LINE_ALIGN_MARK (marker);
  /**
   * ACL index to match
   */
  u32 aia_acl;

  /**
   * The VPP index of the ABX policy
   */
  u32 aia_abx;

  /**
   * The interface for the attachment
   */
  u32 aia_sw_if_index;

  /**
   * The priority of this policy for attachment.
   * The lower the value the higher the priority.
   * The higher priority policies are matched first.
   */
  u32 aia_priority;
} abx_if_attach_t;

extern abx_if_attach_t *abx_if_attach_pool;
extern u32 **abx_per_if;

static inline abx_if_attach_t *
abx_if_attach_get (u32 index)
{
  return (pool_elt_at_index (abx_if_attach_pool, index));
}

static inline u32 *
abx_attachments_per_if_get (u32 index)
{
  return (abx_per_if[index]);
}

extern int abx_if_attach (u32 policy_id, u32 priority, u32 sw_if_index);

extern int abx_if_detach (u32 policy_id, u32 sw_if_index);

extern u32 get_abx_alctx_per_if (u32 sw_if_index);

extern void *get_abx_acl_main ();

typedef walk_rc_t (*abx_if_attach_cb_t) (u32 aiai, void *ctx);

extern void abx_if_attach_walk (abx_if_attach_cb_t cb, void *ctx);

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

#endif

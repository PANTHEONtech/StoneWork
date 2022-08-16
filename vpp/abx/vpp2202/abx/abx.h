
/*
 * abx.h - skeleton vpp engine plug-in header file
 *
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
#ifndef __included_abx_h__
#define __included_abx_h__

#include <abx/abx_policy.h>
#include <vnet/vnet.h>
#include <vnet/ip/ip.h>
#include <vnet/ethernet/ethernet.h>
#include <vnet/format_fns.h>

#include <vppinfra/hash.h>
#include <vppinfra/error.h>

#define ABX_PLUGIN_VERSION_MAJOR 1
#define ABX_PLUGIN_VERSION_MINOR 2

typedef struct
{
  /* API message ID base */
  u16 msg_id_base;

  /* convenience */
  vlib_main_t *vlib_main;
  vnet_main_t *vnet_main;
  ethernet_main_t *ethernet_main;
} abx_main_t;

extern abx_main_t abx_main;

extern vlib_node_registration_t abx_node;

#endif /* __included_abx_h__ */

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

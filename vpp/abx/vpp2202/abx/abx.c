/*
 * abx.c - skeleton vpp engine plug-in
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

#include <vnet/vnet.h>
#include <vnet/plugin/plugin.h>
#include <abx/abx.h>

#include <vlibapi/api.h>
#include <vlibmemory/api.h>

/* define message IDs */
#include <vnet/format_fns.h>
#include <abx/abx.api_enum.h>
#include <abx/abx.api_types.h>

abx_main_t abx_main;

static clib_error_t *
abx_init (vlib_main_t * vm)
{
  abx_main_t *abx = &abx_main;
  clib_error_t *error = 0;

  abx->vlib_main = vm;
  abx->vnet_main = vnet_get_main ();

  return error;
}

VLIB_INIT_FUNCTION (abx_init);

/* *INDENT-OFF* */
VNET_FEATURE_INIT (abx_ip4_unicast, static) = {
  .arc_name = "ip4-unicast",
  .node_name = "abx-ip4",
  .runs_before = VNET_FEATURES("nat-pre-out2in", "nat-pre-in2out", "nat44-ed-out2in", "nat44-ed-in2out"),
};

VNET_FEATURE_INIT (abx_ip4_multicast, static) = {
  .arc_name = "ip4-multicast",
  .node_name = "abx-ip4",
};

VNET_FEATURE_INIT (abx_ip6_unicast, static) = {
  .arc_name = "ip6-unicast",
  .node_name = "abx-ip6",
};

VNET_FEATURE_INIT (abx_ip6_multicast, static) = {
  .arc_name = "ip6-multicast",
  .node_name = "abx-ip6",
};
/* *INDENT-ON */

/* *INDENT-OFF* */
VLIB_PLUGIN_REGISTER () =
{
  .description = "ACL Based Xconnect",
};
/* *INDENT-ON* */

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

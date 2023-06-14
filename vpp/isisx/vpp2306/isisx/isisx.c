/*
 * isisx.c - skeleton vpp engine plug-in
 *
 * Copyright (c) <current-year> <your-organization>
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
#include <isisx/isisx.h>

#include <vlibapi/api.h>
#include <vlibmemory/api.h>
#include <stdbool.h>

#include <isisx/isisx.api_enum.h>
#include <isisx/isisx.api_types.h>

#define REPLY_MSG_ID_BASE pmp->msg_id_base
#include <vlibapi/api_helper_macros.h>

isisx_main_t isisx_main;

static clib_error_t * isisx_init (vlib_main_t * vm)
{
  isisx_main_t * pmp = &isisx_main;
  clib_error_t * error = 0;

  pmp->vlib_main = vm;
  pmp->vnet_main = vnet_get_main();

  /* Register ISIS protocol to be handled by this plugin node */
  vlib_node_t *node = vlib_get_node_by_name (vm, (u8 *) "isisx");
  osi_register_input_protocol (OSI_PROTOCOL_isis, node->index);

  pmp->log_class = vlib_log_register_class("isisx_plugin", 0);
  return error;
}

VLIB_INIT_FUNCTION (isisx_init);

/* *INDENT-OFF* */
VNET_FEATURE_INIT (isisx, static) =
{
  .arc_name = "device-input",
  .node_name = "isisx",
  .runs_before = VNET_FEATURES ("ethernet-input"),

};
/* *INDENT-ON */

/* *INDENT-OFF* */
VLIB_PLUGIN_REGISTER () =
{
  .description = "A simple ISIS protocol based Xconnect",
};
/* *INDENT-ON* */

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

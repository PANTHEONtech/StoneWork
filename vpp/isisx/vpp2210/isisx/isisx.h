
/*
 * isisx.h - skeleton vpp engine plug-in header file
 *
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
#ifndef __included_isisx_h__
#define __included_isisx_h__

#include <vnet/vnet.h>
#include <vnet/ip/ip.h>
#include <vnet/ethernet/ethernet.h>
#include <vnet/osi/osi.h>

#include <vppinfra/hash.h>
#include <vppinfra/error.h>

#define ISISX_PLUGIN_VERSION_MAJOR 1
#define ISISX_PLUGIN_VERSION_MINOR 0

typedef struct {
    /* API message ID base */
    u16 msg_id_base;

    /* on/off switch for the periodic function */
    u8 periodic_timer_enabled;
    /* Node index, non-zero if the periodic process has been created */
    u32 periodic_node_index;

    /* convenience */
    vlib_main_t * vlib_main;
    vnet_main_t * vnet_main;
    ethernet_main_t * ethernet_main;

    /* Logger class */
    vlib_log_class_t log_class;

    /* Hash table of ISISx connections */
    uword * isisx_connection_db;
} isisx_main_t;

extern isisx_main_t isisx_main;

extern vlib_node_registration_t isisx_node;

#define isisx_log_warn(f, ...) do {                           \
    vlib_log(VLIB_LOG_LEVEL_WARNING, isisx_main.log_class, f, \
             ##__VA_ARGS__);                                  \
} while (0)

#endif /* __included_isisx_h__ */

/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

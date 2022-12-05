/*
 * node.c - skeleton vpp engine plug-in dual-loop node skeleton
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
#include <vlib/vlib.h>
#include <vnet/llc/llc.h>
#include <vnet/vnet.h>
#include <vnet/pg/pg.h>
#include <vppinfra/error.h>

#include <isisx/isisx.h>
#include <isisx/isisx_connect.h>

typedef struct
{
  u32 next_index;
  u32 sw_if_index;
  u8 new_src_mac[6];
  u8 new_dst_mac[6];
} isisx_trace_t;

#ifndef CLIB_MARCH_VARIANT

/* packet trace format function */
static u8 * format_isisx_trace (u8 * s, va_list * args)
{
  CLIB_UNUSED (vlib_main_t * vm) = va_arg (*args, vlib_main_t *);
  CLIB_UNUSED (vlib_node_t * node) = va_arg (*args, vlib_node_t *);
  isisx_trace_t * t = va_arg (*args, isisx_trace_t *);

  s = format (s, "isisx: rx_sw_if_index: %d (%U)",
              t->sw_if_index,
              format_vnet_sw_if_index_name, isisx_main.vnet_main, t->sw_if_index);

  return s;
}

vlib_node_registration_t isisx_node;

#endif /* CLIB_MARCH_VARIANT */

#define foreach_isisx_error \
_(FORWARDED, "isisx forwarded packets")

typedef enum {
#define _(sym,str) ISISX_ERROR_##sym,
  foreach_isisx_error
#undef _
  isisx_N_ERROR,
} isisx_error_t;

#ifndef CLIB_MARCH_VARIANT
static char * isisx_error_strings[] =
{
#define _(sym,string) string,
  foreach_isisx_error
#undef _
};
#endif /* CLIB_MARCH_VARIANT */

typedef enum
{
  ISISX_NEXT_DROP,
  ISISX_NEXT_INTERFACE_OUTPUT,
  ISISX_N_NEXT,
} isisx_next_t;

VLIB_NODE_FN (isisx_node) (vlib_main_t * vm,
		  vlib_node_runtime_t * node,
		  vlib_frame_t * frame)
{
  u32 n_left_from, * from, * to_next;
  isisx_next_t next_index;
  u32 pkts_forwarded = 0;

  from = vlib_frame_vector_args (frame);
  n_left_from = frame->n_vectors;
  next_index = node->cached_next_index;

  while (n_left_from > 0)
    {
      u32 n_left_to_next;

      vlib_get_next_frame (vm, node, next_index,
			   to_next, n_left_to_next);

      while (n_left_from > 0 && n_left_to_next > 0)
	  {
      
        u32 bi0;
	      vlib_buffer_t * b0;
        u32 next0 = ISISX_NEXT_DROP;
        u32 sw_if_index0;

        /* speculatively enqueue b0 to the current next frame */
        bi0 = from[0];
        to_next[0] = bi0;
        from += 1;
        to_next += 1;
        n_left_from -= 1;
        n_left_to_next -= 1;

        b0 = vlib_get_buffer (vm, bi0);
        sw_if_index0 = vnet_buffer(b0)->sw_if_index[VLIB_RX];

        /* Tracing */
        isisx_trace_t *t = NULL;
        
        if (PREDICT_FALSE((node->flags & VLIB_NODE_FLAG_TRACE)
            && (b0->flags & VLIB_BUFFER_IS_TRACED)))
        {
          t = vlib_add_trace (vm, node, b0, sizeof (*t));
          t->sw_if_index = sw_if_index0;
          t->next_index = next0;
        }

        u32 tx_sw_if_index = isisx_get_tx_by_rx(sw_if_index0);
        if(tx_sw_if_index != INDEX_INVALID)
        {
          /* Send pkt back out the RX interface */
          next0 = ISISX_NEXT_INTERFACE_OUTPUT;
          /* reset buffer offset */
          vlib_buffer_advance (b0, (sizeof (ethernet_header_t) + sizeof (llc_header_t)) * (-1));
          //t->next_index = next0;
          vnet_buffer(b0)->sw_if_index[VLIB_TX] = tx_sw_if_index;
          pkts_forwarded++;
        }

        /* verify speculative enqueue, maybe switch current next frame */
	      vlib_validate_buffer_enqueue_x1 (vm, node, next_index,
				  to_next, n_left_to_next,
					bi0, next0);
  	}

      vlib_put_next_frame (vm, node, next_index, n_left_to_next);
    }

  vlib_node_increment_counter (vm, isisx_node.index,
                               ISISX_ERROR_FORWARDED, pkts_forwarded);
  return frame->n_vectors;
}

/* *INDENT-OFF* */
#ifndef CLIB_MARCH_VARIANT
VLIB_REGISTER_NODE (isisx_node) =
{
  .name = "isisx",
  .vector_size = sizeof (u32),
  .format_trace = format_isisx_trace,
  .type = VLIB_NODE_TYPE_INTERNAL,

  .n_errors = ARRAY_LEN(isisx_error_strings),
  .error_strings = isisx_error_strings,

  .n_next_nodes = ISISX_N_NEXT,

  .next_nodes = {
        [ISISX_NEXT_DROP] = "error-drop",
        [ISISX_NEXT_INTERFACE_OUTPUT] = "interface-output",
  },
};
#endif /* CLIB_MARCH_VARIANT */

/* *INDENT-ON* */
/*
 * fd.io coding-style-patch-verification: ON
 *
 * Local Variables:
 * eval: (c-set-style "gnu")
 * End:
 */

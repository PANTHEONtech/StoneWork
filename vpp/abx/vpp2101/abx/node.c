/*
 * node.c - skeleton vpp engine plug-in dual-loop node skeleton
 *
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
#include <abx/abx_if_attach.h>
#include <plugins/acl/exports.h>

#include <vlib/vlib.h>
#include <vnet/vnet.h>
#include <vnet/pg/pg.h>
#include <vppinfra/error.h>
#include <abx/abx.h>

typedef struct
{
  u32 next_index;
  u32 sw_if_index;

  bool matched;
  u32 acl_pos;
  u32 aia_sw_if_index;
  u32 ap_id;
} abx_trace_t;

#ifndef CLIB_MARCH_VARIANT
/* packet trace format function */
static u8 *
format_abx_trace (u8 * s, va_list * args)
{
  CLIB_UNUSED (vlib_main_t * vm) = va_arg (*args, vlib_main_t *);
  CLIB_UNUSED (vlib_node_t * node) = va_arg (*args, vlib_node_t *);
  abx_trace_t *t = va_arg (*args, abx_trace_t *);

  s = format (s, "ABX: %s, sw_if_index %d, next index %d",
	      (t->matched) ? "match" : "no-match",
	      t->sw_if_index, t->next_index);
  if (t->matched)
    {
      s = format (s, ", acl_pos %d, attached_if %d, policy_id %d",
		  t->acl_pos, t->aia_sw_if_index, t->ap_id);
    }
  return s;
}

vlib_node_registration_t abx_node;

#endif /* CLIB_MARCH_VARIANT */

#define foreach_abx_error \
_(FWD, "ABX forwarded packets")

typedef enum
{
#define _(sym,str) ABX_ERROR_##sym,
  foreach_abx_error
#undef _
    ABX_N_ERROR,
} abx_error_t;

#ifndef CLIB_MARCH_VARIANT
static char *abx_error_strings[] = {
#define _(sym,string) string,
  foreach_abx_error
#undef _
};
#endif /* CLIB_MARCH_VARIANT */

typedef enum
{
  ABX_NEXT_INTERFACE_OUTPUT,
  ABX_N_NEXT,
} abx_next_t;

static inline uword
abx_node_fn_inline (vlib_main_t * vm,
			vlib_node_runtime_t * node,
			vlib_frame_t * frame,
			u8 is_ip6)
{
  u32 n_left_from, *from, *to_next;
  abx_next_t next_index;
  u32 pkts_abx_fwd = 0;

  from = vlib_frame_vector_args (frame);
  n_left_from = frame->n_vectors;
  next_index = node->cached_next_index;

  while (n_left_from > 0)
    {
      u32 n_left_to_next;

      vlib_get_next_frame (vm, node, next_index, to_next, n_left_to_next);

      while (n_left_from > 0 && n_left_to_next > 0)
	{
	  u32 match_acl_index = ~0;
	  u32 match_acl_pos = ~0;
	  u32 match_rule_index = ~0;
	  u32 trace_bitmap = 0;
	  u8 action;
	  u32 lc_index = ~0;
	  fa_5tuple_opaque_t fa_5tuple0;
	  bool matched = false;
	  const u32 *attachments0;
	  const abx_if_attach_t *aia0 = NULL;
	  const abx_policy_t *ap0 = NULL;

	  u32 bi0;
	  vlib_buffer_t *b0;
	  u32 next0, arc_next0;
	  u32 sw_if_index0;

	  /* speculatively enqueue b0 to the current next frame */
	  bi0 = from[0];
	  to_next[0] = bi0;
	  from += 1;
	  to_next += 1;
	  n_left_from -= 1;
	  n_left_to_next -= 1;

	  b0 = vlib_get_buffer (vm, bi0);
	  vnet_feature_next (&arc_next0, b0);
	  next0 = arc_next0;
	  sw_if_index0 = vnet_buffer (b0)->sw_if_index[VLIB_RX];

	  attachments0 = abx_attachments_per_if_get (sw_if_index0);
	  /*
	   * check if any of the policies attached to this interface matches.
	   */
	  lc_index = get_abx_alctx_per_if (sw_if_index0);
	  if (INDEX_INVALID != lc_index)
	    {
	      acl_plugin_fill_5tuple_inline (get_abx_acl_main (), lc_index,
					     b0, is_ip6, 1, 0, &fa_5tuple0);
	      ip4_header_t *ip4 = vlib_buffer_get_current (b0);
	      ip6_header_t *ip6 = vlib_buffer_get_current (b0);

	      if (acl_plugin_match_5tuple_inline
		  (get_abx_acl_main (), lc_index, &fa_5tuple0, is_ip6,
		   &action, &match_acl_pos, &match_acl_index,
		   &match_rule_index, &trace_bitmap))
		{
		  ethernet_header_t *new_eth, *orig_eth;
		  u8 l2off_valid = b0->flags & VNET_BUFFER_F_L2_HDR_OFFSET_VALID;
		  u16 l2off = l2off_valid ? vnet_buffer (b0)->l2_hdr_offset : 0;
		  orig_eth = (ethernet_header_t *)(b0->data + l2off);
		  vlib_buffer_advance (b0, -sizeof (ethernet_header_t));
		  new_eth = vlib_buffer_get_current (b0);

		  if (ethernet_frame_is_tagged (clib_net_to_host_u16 (orig_eth->type)))
		    {
			  /* Remove VLAN tag. */
			  memmove(new_eth->src_address, orig_eth->src_address, 6);
			  memmove(new_eth->dst_address, orig_eth->dst_address, 6);
			  vnet_buffer (b0)->l2_hdr_offset = b0->current_data;
		    }

		  next0 = ABX_NEXT_INTERFACE_OUTPUT;
		  matched = true;

		  aia0 = abx_if_attach_get (attachments0[match_acl_pos]);
		  if (aia0)
		    {
		      ap0 = abx_policy_get (aia0->aia_abx);
		      if (ap0)
			{
			  /* Overwrite dest mac address if not muticast */
			  if (!mac_address_is_zero (&ap0->ap_dst_mac))
			    {
			      if ((!is_ip6
				   && (ip4->dst_address.as_u8[0] >> 4 != 0xe))
				  || (is_ip6
				      && (ip6->
					   dst_address.as_u16[0] != 0xff)))
				{
				  clib_memcpy (new_eth->dst_address, ap0->ap_dst_mac.bytes, 6);
				}
			    }

			  vnet_buffer (b0)->sw_if_index[VLIB_TX] =
			    ap0->ap_tx_sw_if_index;
			  pkts_abx_fwd++;
			}
		    }
		}
	    }
	  /* Send pkt back out the RX interface */
	  if (PREDICT_FALSE ((node->flags & VLIB_NODE_FLAG_TRACE)
			     && (b0->flags & VLIB_BUFFER_IS_TRACED)))
	    {
	      abx_trace_t *t = vlib_add_trace (vm, node, b0, sizeof (*t));
	      t->sw_if_index = sw_if_index0;
	      t->next_index = next0;
	      t->matched = matched;
	      if (matched && aia0 && ap0)
		{
		  t->acl_pos = match_acl_pos;
		  t->aia_sw_if_index = aia0->aia_sw_if_index;
		  t->ap_id = ap0->ap_id;
		}
	    }

	  /* verify speculative enqueue, maybe switch current next frame */
	  vlib_validate_buffer_enqueue_x1 (vm, node, next_index,
					   to_next, n_left_to_next, bi0,
					   next0);
	}

      vlib_put_next_frame (vm, node, next_index, n_left_to_next);
    }

  vlib_node_increment_counter (vm, abx_node.index, ABX_ERROR_FWD,
			       pkts_abx_fwd);

  return frame->n_vectors;
}

VLIB_NODE_FN (abx_ip4_node) (vlib_main_t * vm,
				     vlib_node_runtime_t * node,
				     vlib_frame_t * frame)
{
  return abx_node_fn_inline (vm, node, frame, 0 /* is_ip6 */);
}

VLIB_NODE_FN (abx_ip6_node) (vlib_main_t * vm,
				     vlib_node_runtime_t * node,
				     vlib_frame_t * frame)
{
  return abx_node_fn_inline (vm, node, frame, 1 /* is_ip6 */);
}

/* *INDENT-OFF* */
#ifndef CLIB_MARCH_VARIANT
VLIB_REGISTER_NODE (abx_ip4_node) =
{
  .name = "abx-ip4",
  .vector_size = sizeof (u32),
  .format_trace = format_abx_trace,
  .type = VLIB_NODE_TYPE_INTERNAL,

  .n_errors = ARRAY_LEN(abx_error_strings),
  .error_strings = abx_error_strings,

  .n_next_nodes = ABX_N_NEXT,
  .next_nodes = {
        [ABX_NEXT_INTERFACE_OUTPUT] = "interface-output",
  },
};

VLIB_REGISTER_NODE (abx_ip6_node) =
{
  .name = "abx-ip6",
  .vector_size = sizeof (u32),
  .format_trace = format_abx_trace,
  .type = VLIB_NODE_TYPE_INTERNAL,

  .n_errors = ARRAY_LEN(abx_error_strings),
  .error_strings = abx_error_strings,

  .n_next_nodes = ABX_N_NEXT,
  .next_nodes = {
        [ABX_NEXT_INTERFACE_OUTPUT] = "interface-output",
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

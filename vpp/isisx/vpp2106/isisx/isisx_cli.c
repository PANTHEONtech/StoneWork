/*
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

#include <isisx/isisx.h>
#include <isisx/isisx_connect.h>

static clib_error_t *
isisx_add_del_connection (vlib_main_t * vm,
		unformat_input_t * main_input, vlib_cli_command_t * cmd)
{
  unformat_input_t _line_input, *line_input = &_line_input;
  u32 rx_sw_if_index = INDEX_INVALID, tx_sw_if_index = INDEX_INVALID, is_del = 0;
  vnet_main_t *vnm = vnet_get_main ();

  /* Get a line of input. */
  if (!unformat_user (main_input, unformat_line_input, line_input))
    return 0;

  while (unformat_check_input (line_input) != UNFORMAT_END_OF_INPUT)
    {
      if (unformat (line_input, "del"))
	is_del = 1;
      else if (unformat (line_input, "add"))
	is_del = 0;
      else if (unformat (line_input, " %U",
			 unformat_vnet_sw_interface, vnm, &rx_sw_if_index))
	;
      else if (unformat (line_input, "via %U",
			 unformat_vnet_sw_interface, vnm, &tx_sw_if_index))
	;
      else
	return (clib_error_return (0, "unknown input '%U'",
				   format_unformat_error, line_input));
    }

  if (rx_sw_if_index == INDEX_INVALID)
    {
      vlib_cli_output (vm, "Specify a RX interface");
      return 0;
    }

  if (!is_del)
    {
      if (rx_sw_if_index == INDEX_INVALID)
	{
	  vlib_cli_output (vm, "rx_sw_if_index must be set");
	  return 0;
	}
      if (tx_sw_if_index == INDEX_INVALID)
	{
	  vlib_cli_output (vm, "tx_sw_if_index must be set");
	  return 0;
	}
      isisx_connection_update (rx_sw_if_index, tx_sw_if_index);
    }
  else
    {
      isisx_connection_delete (rx_sw_if_index);
    }

  unformat_free (line_input);
  return NULL;
}

/* *INDENT-OFF* */
/**
 * Create an ISISx connection.
 */
VLIB_CLI_COMMAND (cli_isisx_add_del_connection, static) = {
  .path = "isisx connection",
  .function = isisx_add_del_connection,
  .short_help = "isisx connection [add|del] <rx_sw_if> via <tx_sw_if>",
  .is_mp_safe = 1,
};
/* *INDENT-ON* */

static u8 *
format_isisx_connection (u8 * s, va_list * args)
{
  u32 * rx_sw_if_index = va_arg (*args, u32 *);
  u32 * tx_sw_if_index = va_arg (*args, u32 *);

  s = format (s, "  rx interface: %U, sw_if_index:%d -> tx interface %U, sw_if_index:%d",
	      format_vnet_sw_if_index_name, isisx_main.vnet_main, rx_sw_if_index, rx_sw_if_index,
        format_vnet_sw_if_index_name, isisx_main.vnet_main, tx_sw_if_index, tx_sw_if_index);

  return s;
}

static clib_error_t *
show_isisx_connection (vlib_main_t * vm,
		     unformat_input_t * input, vlib_cli_command_t * cmd)
{
  u32
    rx_sw_if_index = INDEX_INVALID,
    tx_sw_if_index = INDEX_INVALID;

  while (unformat_check_input (input) != UNFORMAT_END_OF_INPUT)
    {
      if (unformat (input, "%d", &rx_sw_if_index))
	;
      else
	return (clib_error_return (0, "unknown input '%U'",
				   format_unformat_error, input));
    }

  vlib_cli_output (vm, "ISISx connections:");
  if (rx_sw_if_index == INDEX_INVALID)
    {
      /* *INDENT-OFF* */
      hash_foreach (rx_sw_if_index, tx_sw_if_index, (&isisx_main)->isisx_connection_db,
      {
        vlib_cli_output (vm, "%U", format_isisx_connection, rx_sw_if_index, tx_sw_if_index);
      });
      /* *INDENT-ON* */
    }
  else
    {
      tx_sw_if_index = isisx_get_tx_by_rx (rx_sw_if_index);
      if (tx_sw_if_index != INDEX_INVALID)
	      vlib_cli_output (vm, "%U", format_isisx_connection, rx_sw_if_index, tx_sw_if_index);
      else
	      vlib_cli_output (vm, "Invalid rx interface index: %d", rx_sw_if_index);
    }
  return NULL;
}

/* *INDENT-OFF* */
VLIB_CLI_COMMAND (cli_show_isisx_connection, static) = {
  .path = "show isisx connection",
  .function = show_isisx_connection,
  .short_help = "show isisx connection [rx-interface-id]",
  .is_mp_safe = 1,
};
/* *INDENT-ON* */

/*
 * SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2020 PANTHEON.tech. All Rights Reserved.
 */

package dhcp4

import (
	"go.ligato.io/vpp-agent/v3/pkg/models"
)

// ModuleName is the module name used for models.
const (
	ModuleName = "dhcp"

	/* Kea DHCP Server (derived kv) */
	KeaKey = "dhcp/dhcp4/kea"
)

var (
	ModelDHCP4 = models.Register(&Dhcp4{}, models.Spec{
		Module:  ModuleName,
		Version: "v1",
		Type:    "dhcp4",
	})
)

// Dhcp4Key returns the key used in ETCD to store the configuration for the Dhcp4-server.
func Dhcp4Key() string {
	return models.Key(
		&Dhcp4{},
	)
}

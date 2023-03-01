/*
 * SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2020 PANTHEON.tech. All Rights Reserved.
 */

package ospf

import "go.ligato.io/vpp-agent/v3/pkg/models"

const (
	ModuleName = "ospf"

	// The prefix for the derived key with configuration of an OSPF-enabled interface
	// that is applied into the FRR's ospf daemon.
	FrrOspfdInterfaceKeyPrefix = "frr/ospfd/interface/"
)

var (
	ModelOspfRouter = models.Register(
		&Router{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "router",
			Class:   "config",
		})

	ModelOspfInterface = models.Register(
		&Interface{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "interface",
			Class:   "config",
		},
		models.WithNameTemplate(
			"{{.VppInterface}}",
		))
)

// RouterKey returns global OSPF router key
func RouterKey() string {
	return models.Key(&Router{})
}

// InterfaceKey returns key representing given OSPF-enabled interface.
func InterfaceKey(vppIfName string) string {
	return models.Key(&Interface{VppInterface: vppIfName})
}

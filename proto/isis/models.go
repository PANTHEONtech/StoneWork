/*
 * SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2020 PANTHEON.tech. All Rights Reserved.
 */

package isis

import "go.ligato.io/vpp-agent/v3/pkg/models"

const (
	ModuleName = "isis"
	// FrrIsisdInterfaceKeyPrefix is the prefix for the derived key with configuration of an ISIS-enabled interface
	// that is applied into the FRR's isis daemon.
	FrrIsisdInterfaceKeyPrefix = "frr/isisd/interface/"
)

var (
	ModelIsisRouter = models.Register(
		&Router{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "router",
			Class:   "config",
		})
	ModelIsisInterface = models.Register(
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

// RouterKey returns global ISIS router key
func RouterKey() string {
	return models.Key(&Router{})
}

// InterfaceKey returns key representing given ISIS-enabled interface.
func InterfaceKey(vppIfName string) string {
	return models.Key(&Interface{VppInterface: vppIfName})
}

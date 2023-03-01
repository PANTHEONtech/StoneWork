/*
 * SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2020 PANTHEON.tech. All Rights Reserved.
 */

package bgp

import (
	"strconv"
	"strings"

	"go.ligato.io/vpp-agent/v3/pkg/models"
)

const ModuleName = "bgp"

var (
	ModelBgpRouter = models.Register(
		&Router{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "router",
			Class:   "config",
		},
		models.WithNameTemplate(
			`{{if .MultiVrf}}multi-vrf{{else}}vrf/{{.VrfId}}{{end}}`,
		))

	ModelBgpVrf = models.Register(
		&Vrf{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "vrf",
			Class:   "config",
		},
		models.WithNameTemplate(
			"{{.Id}}",
		))

	ModelBgpSession = models.Register(
		&SessionEndpoint{},
		models.Spec{
			Module:  ModuleName,
			Version: "v1",
			Type:    "session-endpoint",
			Class:   "config",
		},
		models.WithNameTemplate(
			"vrf/{{.VrfId}}/interface/{{.VppInterface}}",
		))
)

// RouterKey returns global OSPF router key
func RouterKey() string {
	return models.Key(&Router{})
}

func MultiVrfServerKey() string {
	return models.Key(&Router{MultiVrf: true})
}

func PerVrfServerKey(vrf uint32) string {
	return models.Key(&Router{VrfId: vrf})
}

func SessionEpKey(vppIfName string, vrf uint32) string {
	return models.Key(&SessionEndpoint{VppInterface: vppIfName, VrfId: vrf})
}

func SessionEpVrfKeyPrefix(vrf uint32) string {
	key := models.Key(&SessionEndpoint{VrfId: vrf})
	lastIdx := strings.Index(key, "/interface")
	return key[:lastIdx]
}

// // Derived Keys
const (
	NeighborNetworkKeyPrefix   = "cnf/bgp/server/neighbor/network/"
	neighborNetworkKeyTemplate = NeighborNetworkKeyPrefix + "{{.Network}}/vrf/{{.VrfId}}/server-mode/{{.ServerMode}}"

	SessionPeerKeyPrefix   = "cnf/bgp/session/peer/"
	sessionPeerKeyTemplate = SessionPeerKeyPrefix + "interface/{{.IfName}}/peer-ip/{{.PeerIP}}"

	VrfTableKeyPrefix   = "cnf/bgp/vrf/table/"
	vrfTableKeyTemplate = VrfTableKeyPrefix + "{{.VrfId}}"
)

// NeighborNetworkKey returns derived from Server to represent single neighboring network.
func NeighborNetworkKey(network string, vrf uint32, multiVrfSrv bool) string {
	key := strings.Replace(neighborNetworkKeyTemplate, "{{.Network}}", network, 1)
	key = strings.Replace(key, "{{.VrfId}}", strconv.Itoa(int(vrf)), 1)
	if multiVrfSrv {
		key = strings.Replace(key, "{{.ServerMode}}", "multi-vrf", 1)
	} else {
		key = strings.Replace(key, "{{.ServerMode}}", "per-vrf", 1)
	}
	return key
}

// IsMultiVrfNeighNetKey returns true if key represents neighbor network of a multi-VRF BGP server.
func IsMultiVrfNeighNetKey(key string) bool {
	return strings.HasPrefix(key, NeighborNetworkKeyPrefix) &&
		strings.Contains(key, "/server-mode/multi-vrf")
}

// SessionPeerKey returns key derived from Session to represent Session peer configured in GoBGP.
func SessionPeerKey(ifname, peerIP string) string {
	key := strings.Replace(sessionPeerKeyTemplate, "{{.IfName}}", ifname, 1)
	key = strings.Replace(key, "{{.PeerIP}}", peerIP, 1)
	return key
}

// VrfTableKey returns key derived from VRF to represent VRF table configured in GoBGP.
func VrfTableKey(vrfId uint32) string {
	return strings.Replace(vrfTableKeyTemplate, "{{.VrfId}}", strconv.Itoa(int(vrfId)), 1)
}

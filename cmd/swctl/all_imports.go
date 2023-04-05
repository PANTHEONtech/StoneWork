package main

import (
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/linux"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/netalloc"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/vpp"

	_ "go.pantheon.tech/stonework/proto/abx"
	_ "go.pantheon.tech/stonework/proto/bfd"
	_ "go.pantheon.tech/stonework/proto/isisx"

	_ "go.pantheon.tech/stonework/proto/dhcp4"
)

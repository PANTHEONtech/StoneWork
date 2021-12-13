module go.pantheon.tech/stonework

go 1.14

require (
	git.fd.io/govpp.git v0.3.6-0.20210601140839-da95997338b7
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-errors/errors v1.0.1
	github.com/golang/protobuf v1.4.2
	github.com/namsral/flag v1.7.4-pre
	github.com/onsi/gomega v1.10.3
	github.com/pkg/errors v0.9.1
	github.com/vishvananda/netlink v0.0.0-20180910184128-56b1bd27a9a3
	github.com/vishvananda/netns v0.0.0-20180720170159-13995c7128cc
	go.ligato.io/cn-infra/v2 v2.5.0-alpha.0.20200313154441-b0d4c1b11c73
	go.ligato.io/vpp-agent/v3 v3.3.0-alpha.0.20210716165218-6eac586bfd7d
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
)

// Uncomment to use vpp-agent from sub-module to test changes in the agent before they are merged into the upstream
// (or into any remote fork for that matter).
// TODO: remove the submodule eventually
// replace go.ligato.io/vpp-agent/v3 => ./submodule/vpp-agent

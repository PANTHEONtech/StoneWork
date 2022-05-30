module go.pantheon.tech/stonework

go 1.17

require (
	git.fd.io/govpp.git v0.4.0
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-errors/errors v1.0.1
	github.com/namsral/flag v1.7.4-pre
	github.com/onsi/gomega v1.19.0
	github.com/pkg/errors v0.9.1
	github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df
	go.ligato.io/cn-infra/v2 v2.5.0-alpha.0.20220211111933-3d9ff310b1fa
	go.ligato.io/vpp-agent/v3 v3.3.0
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)

replace (
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.5
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.15-0.20200113171025-3fe6c5262873 // indirect
	github.com/Microsoft/hcsshim v0.8.9 // indirect
	github.com/alecthomas/jsonschema v0.0.0-20200217214135-7152f22193c9 // indirect
	github.com/bennyscetbun/jsongo v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bshuster-repo/logrus-logstash-hook v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/containerd/cgroups v0.0.0-20190919134610-bf292b21730f // indirect
	github.com/containerd/containerd v1.4.13 // indirect
	github.com/coreos/go-iptables v0.4.5 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/docker/docker v20.10.12+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/evalphobia/logrus_fluent v0.4.0 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/fluent/fluent-logger-golang v1.3.0 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/fsouza/go-dockerclient v1.6.6 // indirect
	github.com/ftrvxmtrx/fd v0.0.0-20150925145434-c6d800382fff // indirect
	github.com/goccy/go-graphviz v0.0.9 // indirect
	github.com/goccy/go-yaml v1.8.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0 // indirect
	github.com/lunixbochs/struc v0.0.0-20200521075829-a4cb8d33dbbe // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/moby/sys/mount v0.1.0 // indirect
	github.com/moby/sys/mountinfo v0.4.1 // indirect
	github.com/moby/term v0.0.0-20200429084858-129dac9f73f6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.0.3 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tinylib/msgp v1.0.2 // indirect
	github.com/unrolled/render v0.0.0-20180914162206-b9786414de4d // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.1 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.1 // indirect
	go.etcd.io/etcd/client/v3 v3.5.1 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/image v0.0.0-20220413100746-70e8d0d3baa9 // indirect
	golang.org/x/net v0.0.0-20220526153639-5463443f8c37 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20210629200056-84d6f6074151 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Uncomment to use vpp-agent from the specified directory to test changes in the agent before they are merged into the upstream
// (or into any remote fork for that matter).
// replace go.ligato.io/vpp-agent/v3 => ./path/to/vpp-agent

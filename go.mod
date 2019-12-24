module lrishman

go 1.13

require (
	consul v1.1.2
	datastruck v1.1.2
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878 // indirect
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-rootcerts v1.0.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/hashicorp/memberlist v0.1.5 // indirect
	github.com/hashicorp/serf v0.8.5 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107 // indirect
	google.golang.org/grpc v1.23.0 // indirect
	initconfig v1.1.2 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace consul v1.1.2 => ./pkg/consul

replace initconfig v1.1.2 => ./pkg/initconfig

replace datastruck v1.1.2 => ./pkg/datastruck

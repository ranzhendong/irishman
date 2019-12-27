module etcd

go 1.13

require (
	datastruck v1.1.2
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/etcd-io/etcd v3.3.18+incompatible
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/consul/api v1.3.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/spf13/viper v1.6.1
	initconfig v1.1.2
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	datastruck v1.1.2 => ../datastruck
	initconfig v1.1.2 => ../initconfig
)

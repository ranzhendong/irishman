module consul

go 1.13

require (
	datastruck v1.1.2
	github.com/hashicorp/consul/api v1.3.0
	github.com/spf13/viper v1.6.1
	initconfig v1.1.2 // indirect
)

replace (
	datastruck v1.1.2 => ../datastruck
	initconfig v1.1.2 => ../initconfig
)

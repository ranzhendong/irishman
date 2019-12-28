module upstream

go 1.13

require (
	datastruck v1.1.2
	etcd v1.1.2
)

replace (
	datastruck v1.1.2 => ../datastruck
	etcd v1.1.2 => ../etcd
)

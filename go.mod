module lrishman

go 1.13

require (
	datastruck v1.1.2
	govalidators v1.1.2
	init v1.1.2
	reconstruct v1.1.2
	upstream v1.1.2
)

replace (
	datastruck v1.1.2 => ./pkg/datastruck
	etcd v1.1.2 => ./pkg/etcd
	govalidators v1.1.2 => ./src/govalidators
	init v1.1.2 => ./pkg/init
	reconstruct v1.1.2 => ./pkg/reconstruct
	upstream v1.1.2 => ./pkg/upstream
)

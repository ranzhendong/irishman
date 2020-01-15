module lrishman

go 1.13

require (
	datastruck v1.1.2
	errorhandle v1.1.2
	github.com/xujiajun/nutsdb v0.5.0
	golang.org/x/sys v0.0.0-20200113162924-86b910548bc1 // indirect
	gopkg.in/fatih/set.v0 v0.2.1 // indirect
	healthcheck v1.1.2
	init v1.1.2
)

replace (
	datastruck v1.1.2 => ./pkg/datastruck
	errorhandle v1.1.2 => ./pkg/errorhandle
	etcd v1.1.2 => ./pkg/etcd
	govalidators v1.1.2 => ./src/govalidators
	healthcheck v1.1.2 => ./pkg/healthcheck
	init v1.1.2 => ./pkg/init
	kvnuts v1.1.2 => ./pkg/upstream
	reconstruct v1.1.2 => ./pkg/reconstruct
	upstream v1.1.2 => ./pkg/kvnuts
)

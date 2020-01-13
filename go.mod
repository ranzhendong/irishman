module lrishman

go 1.13

require (
	errorhandle v1.1.2
	gopkg.in/fatih/set.v0 v0.2.1 // indirect
	govalidators v1.1.2 // indirect
	healthcheck v1.1.2
	init v1.1.2
	upstream v1.1.2
)

replace (
	datastruck v1.1.2 => ./pkg/datastruck
	errorhandle v1.1.2 => ./pkg/errorhandle
	etcd v1.1.2 => ./pkg/etcd
	govalidators v1.1.2 => ./src/govalidators
	healthcheck v1.1.2 => ./pkg/healthcheck
	init v1.1.2 => ./pkg/init
	reconstruct v1.1.2 => ./pkg/reconstruct
	upstream v1.1.2 => ./pkg/upstream
)

module lrishman

go 1.13

require (
	errorhandle v1.1.2
	github.com/smokezl/govalidators v0.0.0-20181012065008-5fded539f530 // indirect
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

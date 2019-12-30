module lrishman

go 1.13

require (
	datastruck v1.1.2
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/smokezl/govalidators v0.0.0-20181012065008-5fded539f530 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // indirect
	google.golang.org/api v0.15.0 // indirect
	google.golang.org/grpc v1.23.0 // indirect
	govalidators v1.1.2
	init v1.1.2
	upstream v1.1.2
)

replace (
	datastruck v1.1.2 => ./pkg/datastruck
	etcd v1.1.2 => ./pkg/etcd
	govalidators v1.1.2 => ./src/govalidators
	init v1.1.2 => ./pkg/init
	upstream v1.1.2 => ./pkg/upstream
)

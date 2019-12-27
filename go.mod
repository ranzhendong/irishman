module lrishman

go 1.13

require (
	etcd v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107 // indirect
	google.golang.org/grpc v1.23.0 // indirect
	initconfig v1.1.2
)

replace etcd v1.1.2 => ./pkg/etcd

replace initconfig v1.1.2 => ./pkg/initconfig

replace datastruck v1.1.2 => ./pkg/datastruck

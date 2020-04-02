package gorountinescontroller

import "github.com/ranzhendong/irishman/pkg/etcd"

func etcWatcher() {
	_ = etcd.EtcWatcher()
}

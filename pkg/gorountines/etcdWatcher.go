package gorountines

import (
	"github.com/ranzhendong/irishman/pkg/etcd"
)

//etcWatcher : just for etcd
func etcWatcher(watcherKeys ...string) {
	for _, v := range watcherKeys {
		_ = etcd.EtcWatcher(v)
	}
}

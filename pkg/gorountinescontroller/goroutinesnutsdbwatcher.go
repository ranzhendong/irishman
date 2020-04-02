package gorountinescontroller

import (
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
	"time"
)

//nutsWatcher : Flag NutsDB Upstream watcher
func FlagUpstreamNutsDB() {
	var (
		val []*mvccpb.KeyValue
	)

	for {
		time.Sleep(200 * time.Millisecond)
		if _, _, err := kvnuts.Get("FlagUpstreamNutsDB", "FlagUpstreamNutsDB", "i"); err == nil {
			_ = kvnuts.Del("FlagUpstreamNutsDB", "FlagUpstreamNutsDB")

			//get upstream list key from etcd
			if _, val, err = etcd.EtcGetAll(c.Upstream.EtcdPrefix); err != nil {
				log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
			}

			//set upstream list storage to nutsDB, set flag
			utnf := healthcheck.UpstreamToNutsDBFlag{
				1,
				0}.UpstreamAndHCEtcdToNutsDB
			utnf(val)

			//trigger restart hc
			kvnuts.SetFlagHC()
		}
	}
}

//nutsWatcher : Flag NutsDB Upstream watcher
func FlagHCNutsDB() {
	var (
		val []*mvccpb.KeyValue
	)

	for {
		time.Sleep(200 * time.Millisecond)
		if _, _, err := kvnuts.Get("FlagHCNutsDB", "FlagHCNutsDB", "i"); err == nil {
			_ = kvnuts.Del("FlagHCNutsDB", "FlagHCNutsDB")

			//get upstream list key from etcd
			if _, val, err = etcd.EtcGetAll(c.Upstream.EtcdPrefix); err != nil {
				log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
			}

			//set upstream list storage to nutsDB, set flag
			utnf := healthcheck.UpstreamToNutsDBFlag{
				0,
				1}.UpstreamAndHCEtcdToNutsDB
			utnf(val)

			//trigger restart hc
			kvnuts.SetFlagHC()
		}
	}
}

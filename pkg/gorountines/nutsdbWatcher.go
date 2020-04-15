package gorountines

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
		val  string
		vals []*mvccpb.KeyValue
	)

	for {
		time.Sleep(200 * time.Millisecond)
		if _, _, err := kvnuts.Get("FlagUpstreamNutsDB", "FlagUpstreamNutsDB", "i"); err == nil {
			_ = kvnuts.Del("FlagUpstreamNutsDB", "FlagUpstreamNutsDB")

			WatcherFlag, _, _ := kvnuts.Get("FlagUpstreamNutsDB", "FlagUpstreamNutsDBWatcherFlag", "s")
			log.Println("++++++++++++++", WatcherFlag)

			//set upstream list storage to nutsDB, set flag
			go func() {
				utnf := healthcheck.UpstreamToNutsDBFlag{
					SeparateUpstreamEtcdToNutsForOne: 1,
					HealthCheckEtcdToNuts:            0,
					OneKey:                           WatcherFlag}.UpstreamAndHCFromEtcdToNutsDB
				utnf(vals, val)
			}()

			for {
				time.Sleep(100 * time.Millisecond)
				log.Println("111111111111111111")
				if _, _, err := kvnuts.Get("SetFlagUpstreamReadyTo", "SetFlagUpstreamReadyTo", "i"); err == nil {
					log.Println("time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)")
					_ = kvnuts.Del("SetFlagUpstreamReadyTo", "SetFlagUpstreamReadyTo")

					log.Println("##################")
					//trigger restart hc
					kvnuts.SetFlagHC()
					log.Println("##################")
					goto BREAKFOR
				}
			}
		}
	BREAKFOR:
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

			//get upstream list key from etcd, using upstream prefix
			if _, val, err = etcd.EtcGetAll(c.Upstream.EtcdPrefix); err != nil {
				log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
			}

			//set upstream list storage to nutsDB, set flag
			utnf := healthcheck.UpstreamToNutsDBFlag{
				SeparateUpstreamEtcdToNuts: 0,
				HealthCheckEtcdToNuts:      1}.UpstreamAndHCFromEtcdToNutsDB
			utnf(val)

			//trigger restart hc
			kvnuts.SetFlagHC()
		}
	}
}

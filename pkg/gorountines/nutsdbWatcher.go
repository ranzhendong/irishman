package gorountines

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
	"time"
)

type upstream struct {
	UpstreamName string `json:"upstreamName"`
}

//FlagUpstreamNutsDB : Flag NutsDB Upstream watcher
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
			//log.Println("++++++++++++++", WatcherFlag)

			//set upstream list storage to nutsDB, set flag
			go func() {
				utnf := healthcheck.UpstreamToNutsDBFlag{
					SeparateUpstreamEtcdToNutsForOne: 1,
					HealthCheckEtcdToNuts:            0,
					OneKey:                           WatcherFlag}.UpstreamAndHCFromEtcdToNutsDB
				utnf(vals, val)
			}()

			for {
				time.Sleep(50 * time.Millisecond)
				//log.Println("111111111111111111")
				if _, _, err := kvnuts.Get("SetFlagUpstreamReadyTo", "SetFlagUpstreamReadyTo", "i"); err == nil {
					//log.Println("time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)")
					_ = kvnuts.Del("SetFlagUpstreamReadyTo", "SetFlagUpstreamReadyTo")

					//log.Println("##################")
					//trigger restart hc
					kvnuts.SetFlagHC()
					//log.Println("##################")
					goto BREAKFOR
				}
			}
		}
	BREAKFOR:
	}
}

//FlagHCNutsDB : Flag NutsDB Health check watcher
func FlagHCNutsDB() {

	for {
		time.Sleep(200 * time.Millisecond)
		if _, _, err := kvnuts.Get("FlagHCNutsDB", "FlagHCNutsDB", "i"); err == nil {
			_ = kvnuts.Del("FlagHCNutsDB", "FlagHCNutsDB")

			WatcherFlag, _, _ := kvnuts.Get("FlagHCNutsDB", "FlagHCNutsDBWatcherFlag", "s")
			healthcheck.PostHealthCheckTemplateToNutsDB(WatcherFlag)

			//trigger restart hc
			for {
				time.Sleep(50 * time.Millisecond)
				//log.Println("77777777777777777")
				if _, _, err := kvnuts.Get("SetFlagHealthCheckReadyTo", "SetFlagHealthCheckReadyTo", "i"); err == nil {
					//log.Println("time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)")
					_ = kvnuts.Del("SetFlagHealthCheckReadyTo", "SetFlagHealthCheckReadyTo")
					//log.Println("77777777777777777")
					//trigger restart hc
					_ = kvnuts.Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDBFinishUpstream", 1)
					kvnuts.SetFlagHC()
					//log.Println("77777777777777777")
					goto BREAKFOR
				}
			}
		}
	BREAKFOR:
	}
}

//FlagStartUpstreamNutsDB : Flag NutsDB Start Upstream list recueillir watcher
func FlagStartUpstreamNutsDB() {
	var (
		val []*mvccpb.KeyValue
		u   upstream
	)

	for {
		time.Sleep(200 * time.Millisecond)
		if _, _, err := kvnuts.Get("FlagUpstreamNutsDB", "FlagUpstreamNutsDBStartUpstream", "i"); err == nil {
			_ = kvnuts.Del("FlagUpstreamNutsDB", "FlagUpstreamNutsDBStartUpstream")

			//get upstream list key from etcd, using upstream prefix
			if _, val, err = etcd.EtcGetAll(c.Upstream.EtcdPrefix); err != nil {
				log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
			}

			//get upstream from nutsDB
			UpstreamList, _ := kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

			//remove first, add later
			if len(UpstreamList) > 0 {
				for i := 0; i < len(UpstreamList); i++ {
					//need to use "c.NutsDB.Tag.Up" delete
					_ = kvnuts.SRem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList, UpstreamList[i])
				}
			}

			//add upstream list to nutsDB
			//log.Println("!!!!!1", val)
			for _, v := range val {
				if err := json.Unmarshal(v.Value, &u); err != nil {
					continue
				}
				//set upstream list to nutsDB, as set
				_ = kvnuts.SAdd(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList, u.UpstreamName)
			}

			//UpstreamList, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)
			//for i := 0; i < len(UpstreamList); i++ {
			//	log.Println("*****************************", string(UpstreamList[i]))
			//}

			_ = kvnuts.Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDBFinishUpstream", 1)
		}
	}
}

package healthcheck

import (
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
)

/*
set Separate Upstream To Nuts
set HealthCheck Template To Nuts
just for init

*/
//UpstreamToNutsDB : set upstream list storage to nutsDB
func (f UpstreamToNutsDBFlag) UpstreamAndHCFromEtcdToNutsDB(vals []*mvccpb.KeyValue, val ...string) {
	log.Println("UpstreamAndHCEtcdToNutsDB")
	var tc datastruck.TConfig

	// make sure config loaded
	// set config to tc
	tc = c.TC()

	T := TConfig{
		tc.UpstreamEtcPrefix,
		tc.HealthCheckEtcPrefix,
		tc.TagUp,
		tc.TagDown,
		tc.TagSuccessCode,
		tc.TagFailureCode,
	}

	//Separate Upstream Etcd To Nuts when etcd watcher triggered
	if f.SeparateUpstreamEtcdToNutsForOne != 0 {
		T.SeparateUpstreamFromEtcdToNutsForOne(f.OneKey)
		return
	}

	//set upstream list storage to nutsDB
	for _, v := range vals {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			continue
		}

		//Separate Upstream Etcd To Nuts
		//Multiplexed flag
		if f.SeparateUpstreamEtcdToNuts != 0 {
			T.SeparateUpstreamFromEtcdToNuts(u.UpstreamName)
		}

		//HealthCheck Etcd To Nuts
		if f.HealthCheckEtcdToNuts != 0 {
			//HealthCheck Template To Nuts
			T.HealthCheckEtcdToNuts([]byte(u.UpstreamName))
		}

		//set upstream list to nutsDB, as set
		_ = kvnuts.SAdd(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList, u.UpstreamName)
	}

}

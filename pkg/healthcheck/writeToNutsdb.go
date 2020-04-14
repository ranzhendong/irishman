package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
)

var (
	val string
	err error
)

//repeat remove
func removeRepByMap(slc []map[string]interface{}) (result []map[string]interface{}) {
	tempMap := map[interface{}]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e["ipPort"]] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return
}

/*
All NutsDB bucket, key, and val

for example TheUpstream=vmims
|                       f                      |   set or list   |         bucket        |         key        |   val
|             upstream list recode             |       set       |      UpstreamList     |     UpstreamList   |  ["TheUpstream", "TheUpstream-01", ...]
|                up list recode                |       set       |           Up          |     TheUpstream    |  ["192.168.101.59:8080", "192.168.101.61:8080", ...]
|               down list recode               |       set       |          Down         |     TheUpstream    |  ["192.168.101.59:9000", "192.168.101.61:9000", ...]
|           health check info recode           |      list       |      TheUpstream      |     TheUpstream    |  ["http", "/", 3000, 3, 3000, 4500, 3, 2000]
| health check info success status code recode |       set       |    ScodeTheUpstream   |     TheUpstream    |  [200, 301, 302]
| health check info Failure status code recode |       set       |    FcodeTheUpstream   |     TheUpstream    |  [400, 404, 500, 501, 502, 503, 504, 505]
|          health check status recode          |       k/v       |  TheUpstream+ipPort   |          s         |  times: 1
|          health check status recode          |       k/v       |  TheUpstream+ipPort   |          f         |  times: 1

*/
//SeparateUpstreamToNuts : Separate up&down server from upstream to nutsDB
func (tc *TConfig) SeparateUpstreamFromEtcdToNuts(v string) {
	var u Upstream

	if val, err = etcd.EtcGet(tc.UpstreamEtcPrefix + strFirstToUpper(v)); err != nil {
		log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		log.Println(MyERR.ErrorLog(11005))
	}

	for _, v := range u.Pool {
		if v.Status == "up" {
			_ = kvnuts.SAdd(tc.TagUp, u.UpstreamName, v.IPPort)
		} else {
			_ = kvnuts.SAdd(tc.TagDown, u.UpstreamName, v.IPPort)
		}
	}
}

//SeparateUpstreamFromEtcdToNutsForOne : Separate up&down server from upstream to nutsDB, but just one
func (tc TConfig) SeparateUpstreamFromEtcdToNutsForOne(v string) {
	var (
		u                            Upstream
		c                            datastruck.Config
		EtcdUpIpPort, EtcdDownIpPort []string
	)

	//config loading
	_ = c.Config()

	//get change key from etcd
	if val, err = etcd.EtcGet(v); err != nil {
		log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		log.Println(MyERR.ErrorLog(11005))
	}

	//Separate pool server
	for _, v := range u.Pool {
		if v.Status == "up" {
			EtcdUpIpPort = append(EtcdUpIpPort, v.IPPort)
		} else {
			EtcdDownIpPort = append(EtcdDownIpPort, v.IPPort)
		}
	}

	//get up ip port from nuts
	NutsUpIpPort, _ := kvnuts.SMem(c.NutsDB.Tag.Up, u.UpstreamName)

	//get down ip port from nuts
	//NutsDownIpPort, _ := kvnuts.SMem(c.NutsDB.Tag.Up, u.UpstreamName)

	for _, v := range EtcdUpIpPort {
		log.Println(v)
		//check every ip port
		for i := 0; i < len(NutsUpIpPort); i++ {
			ip := NutsUpIpPort[i]
			log.Println(string(ip))
			if string(ip) == v {
				NutsUpIpPort = append(NutsUpIpPort, []byte(v))
				goto Exit
			}
		}
	Exit:
	}

}

/*
[CheckProtocol, CheckPath, Health.Interval, Health.SuccessTime, Health.SuccessTimeout, UnHealth.Interval, UnHealth.FailuresTime, UnHealth.FailuresTimeout]
[		0, 			  1, 			2, 				3, 						4, 					5, 						6, 						7	     ]
storage health check info as list, but success code and failures code as set.
*/
//HealthCheckToNuts : set health check from etcd to nutsDB
func (tc TConfig) HealthCheckEtcdToNuts(v []byte) {
	var h datastruck.HealthCheck

	//health check to nuts
	if val, err = etcd.EtcGet(tc.HealthCheckEtcPrefix + strFirstToUpper(string(v))); err != nil {
		log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}

	if err = json.Unmarshal([]byte(val), &h); err != nil {
		log.Println(MyERR.ErrorLog(11005))
	}

	_ = kvnuts.LAdd(string(v), v, h.UnHealth.FailuresTimeout)
	_ = kvnuts.LAdd(string(v), v, h.UnHealth.FailuresTime)
	_ = kvnuts.LAdd(string(v), v, h.UnHealth.Interval)
	_ = kvnuts.LAdd(string(v), v, h.Health.SuccessTimeout)
	_ = kvnuts.LAdd(string(v), v, h.Health.SuccessTime)
	_ = kvnuts.LAdd(string(v), v, h.Health.Interval)
	_ = kvnuts.LAdd(string(v), v, h.CheckPath)
	_ = kvnuts.LAdd(string(v), v, h.CheckProtocol)

	//write status code
	for _, t := range h.Health.SuccessStatus {
		//log.Println(t)
		_ = kvnuts.SAdd(tc.TagSuccessCode+string(v), v, t)
	}
	for _, t := range h.UnHealth.FailuresStatus {
		//log.Println(t)
		_ = kvnuts.SAdd(tc.TagFailureCode+string(v), v, t)
	}
}

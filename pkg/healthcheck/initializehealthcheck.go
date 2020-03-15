package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
	"time"
)

type upstream struct {
	UpstreamName string   `json:"upstreamName"`
	Pool         []server `json:"pool"`
}

type server struct {
	IPPort string `json:"ipPort"`
	Status string `json:"status"`
	Weight int    `json:"weight"`
}

//InitHealthCheck : goroutines for Init Health Check
func InitHealthCheck(timeNow time.Time) *MyERR.MyError {
	log.Println("InitHealthCheck")

	var (
		err                            error
		val                            []*mvccpb.KeyValue
		upstreamList, downUpstreamList []string
		healthCheckByte, b             []byte
		upstreamListByte               [][]byte
		u                              upstream
		h                              healthCheck
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(0151), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 0151, TimeStamp: timeNow}
	}

	EtcUpstreamName := c.Upstream.EtcdPrefix
	//get key from etcd
	if _, val, err = etcd.EtcGetAll(EtcUpstreamName); err != nil {
		log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 0104, TimeStamp: timeNow}
	}

	//upstream list storage to nutsDB
	for _, v := range val {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			downUpstreamList = append(downUpstreamList, u.UpstreamName)
			continue
		}
		upstreamList = append(upstreamList, u.UpstreamName)
		//as a number to upstream list
		_ = kvnuts.SAdd(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList, u.UpstreamName)
	}

	for _, v := range upstreamList {
		EtcUpstreamName := c.HealthCheck.EtcdPrefix + strFirstToUpper(v)
		c.HealthCheck.Template.HealthCheckName = v

		//turn struck to json
		if healthCheckByte, err = json.Marshal(c.HealthCheck.Template); err != nil {
			log.Println(MyERR.ErrorLog(0004))
			return &MyERR.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
		}

		// etcd put
		if err = etcd.EtcPut(EtcUpstreamName, string(healthCheckByte)); err != nil {
			log.Printf(MyERR.ErrorLog(0101, fmt.Sprintf("%v", err)))
			return &MyERR.MyError{Error: err.Error(), Code: 0101, TimeStamp: timeNow}
		}
	}

	a := &MyERR.MyError{Code: 000, TimeStamp: timeNow}
	a.Clock()
	if b, err = json.Marshal(a); err != nil {
		log.Println(MyERR.ErrorLog(0004))
		return &MyERR.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
	}

	//split Up, Down from upstream list
	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(11012), fmt.Sprintf("%v", err))
	}

	//get upstream list from nutsDB
	upstreamListByte, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

	for _, v := range upstreamListByte {
		var val string
		if val, err = etcd.EtcGet(c.Upstream.EtcdPrefix + strFirstToUpper(string(v))); err != nil {
			log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
		}
		if err := json.Unmarshal([]byte(val), &u); err != nil {
			log.Println(MyERR.ErrorLog(11005))
		}

		UpDownToNuts(&u)

		//health check to nuts
		if val, err = etcd.EtcGet(c.HealthCheck.EtcdPrefix + strFirstToUpper(string(v))); err != nil {
			log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
		}
		if err := json.Unmarshal([]byte(val), &h); err != nil {
			log.Println(MyERR.ErrorLog(11005))
		}

		TempToNuts(v, &h)

	}

	//ready to hc
	//HC()

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" HealthCheck %v", string(b))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}
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

//UpDownToNuts : get up&down upstream form nutsDB
func UpDownToNuts(u *upstream) {
	for _, v := range u.Pool {
		if v.Status == "up" {
			_ = kvnuts.SAdd(c.NutsDB.Tag.Up, u.UpstreamName, v.IPPort)
		} else {
			_ = kvnuts.SAdd(c.NutsDB.Tag.Down, u.UpstreamName, v.IPPort)
		}
	}
}

/*
[CheckProtocol, CheckPath, Health.Interval, Health.SuccessTime, Health.SuccessTimeout, UnHealth.Interval, UnHealth.FailuresTime, UnHealth.FailuresTimeout]
[		0, 			  1, 			2, 				3, 						4, 					5, 						6, 						7	     ]
storage health check info as list, but success code and failures code as set.
*/

//TempToNuts : set health check template to nutsDB
func TempToNuts(v []byte, h *healthCheck) {
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
		_ = kvnuts.SAdd(c.NutsDB.Tag.SuccessCode+string(v), v, t)
	}
	for _, t := range h.UnHealth.FailuresStatus {
		//log.Println(t)
		_ = kvnuts.SAdd(c.NutsDB.Tag.FailureCode+string(v), v, t)
	}
}

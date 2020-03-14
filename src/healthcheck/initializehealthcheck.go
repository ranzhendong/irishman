package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ranzhendong/irishman/src/datastruck"
	ErrH "github.com/ranzhendong/irishman/src/errorhandle"
	"github.com/ranzhendong/irishman/src/etcd"
	"github.com/ranzhendong/irishman/src/kvnuts"
	"log"
	"time"
)

var c datastruck.Config

type upstream struct {
	UpstreamName string   `json:"upstreamName"`
	Pool         []server `json:"pool"`
}

type server struct {
	IpPort string `json:"ipPort"`
	Status string `json:"status"`
	Weight int    `json:"weight"`
}

func InitHealthCheck(timeNow time.Time) *ErrH.MyError {
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
		log.Println(ErrH.ErrorLog(0151), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 0151, TimeStamp: timeNow}
	}

	EtcUpstreamName := c.Upstream.EtcdPrefix
	//get key from etcd
	if err, _, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
		log.Println(ErrH.ErrorLog(0104), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 0104, TimeStamp: timeNow}
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
			log.Println(ErrH.ErrorLog(0004))
			return &ErrH.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
		}

		// etcd put
		if err = etcd.EtcPut(EtcUpstreamName, string(healthCheckByte)); err != nil {
			log.Printf(ErrH.ErrorLog(0101, fmt.Sprintf("%v", err)))
			return &ErrH.MyError{Error: err.Error(), Code: 0101, TimeStamp: timeNow}
		}
	}

	a := &ErrH.MyError{Code: 000, TimeStamp: timeNow}
	a.Clock()
	if b, err = json.Marshal(a); err != nil {
		log.Println(ErrH.ErrorLog(0004))
		return &ErrH.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
	}

	//split Up, Down from upstream list
	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(11012), fmt.Sprintf("%v", err))
	}

	//get upstream list from nutsDB
	_, upstreamListByte = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

	for _, v := range upstreamListByte {
		var val string
		if err, val = etcd.EtcGet(c.Upstream.EtcdPrefix + strFirstToUpper(string(v))); err != nil {
			log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
		}
		if err := json.Unmarshal([]byte(val), &u); err != nil {
			log.Println(ErrH.ErrorLog(11005))
		}

		UpDownToNuts(&u)

		//health check to nuts
		if err, val = etcd.EtcGet(c.HealthCheck.EtcdPrefix + strFirstToUpper(string(v))); err != nil {
			log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
		}
		if err := json.Unmarshal([]byte(val), &h); err != nil {
			log.Println(ErrH.ErrorLog(11005))
		}

		TempToNuts(v, &h)

	}

	//ready to hc
	//HC()

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" HealthCheck %v", string(b))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
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
func UpDownToNuts(u *upstream) {
	for _, v := range u.Pool {
		if v.Status == "up" {
			_ = kvnuts.SAdd(c.NutsDB.Tag.Up, u.UpstreamName, v.IpPort)
		} else {
			_ = kvnuts.SAdd(c.NutsDB.Tag.Down, u.UpstreamName, v.IpPort)
		}
	}
}

/*
[CheckProtocol, CheckPath, Health.Interval, Health.SuccessTime, Health.SuccessTimeout, UnHealth.Interval, UnHealth.FailuresTime, UnHealth.FailuresTimeout]
[		0, 			  1, 			2, 				3, 						4, 					5, 						6, 						7	     ]
storage health check info as list, but success code and failures code as set.
*/
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

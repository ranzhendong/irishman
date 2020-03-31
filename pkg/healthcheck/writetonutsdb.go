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

//GoroutinesHealthCheck : write HC template to nutsDB
type GoroutinesHealthCheck struct {
	CheckProtocol string `json:"checkProtocol"`
	CheckPath     string `json:"checkPath"`
	Health        struct {
		Interval       int   `json:"interval"`
		SuccessTime    int   `json:"successTime"`
		SuccessTimeout int   `json:"successTimeout"`
		SuccessStatus  []int `json:"successStatus"`
	} `json:"health"`
	UnHealth struct {
		Interval        int   `json:"interval"`
		FailuresTime    int   `json:"failuresTime"`
		FailuresTimeout int   `json:"failuresTimeout"`
		FailuresStatus  []int `json:"failuresStatus"`
	} `json:"unhealth"`
}

//Upstream: write upstream server up&down to nutsDB
type Upstream struct {
	UpstreamName string `json:"upstreamName"`
	Pool         []struct {
		IPPort string `json:"ipPort"`
		Status string `json:"status"`
		Weight int    `json:"weight"`
	} `json:"pool"`
}

var (
	val string
	err error
)

type T datastruck.TConfig

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
func (c T) SeparateUpstreamToNuts(v []byte) {
	var u Upstream

	if val, err = etcd.EtcGet(c.UpstreamEtcPrefix + strFirstToUpper(string(v))); err != nil {
		log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		log.Println(MyERR.ErrorLog(11005))
	}

	for _, v := range u.Pool {
		if v.Status == "up" {
			_ = kvnuts.SAdd(c.TagUp, u.UpstreamName, v.IPPort)
		} else {
			_ = kvnuts.SAdd(c.TagDown, u.UpstreamName, v.IPPort)
		}
	}
}

/*
[CheckProtocol, CheckPath, Health.Interval, Health.SuccessTime, Health.SuccessTimeout, UnHealth.Interval, UnHealth.FailuresTime, UnHealth.FailuresTimeout]
[		0, 			  1, 			2, 				3, 						4, 					5, 						6, 						7	     ]
storage health check info as list, but success code and failures code as set.
*/

//HealthCheckTemplateToNuts : set health check template to nutsDB
func (c T) HealthCheckTemplateToNuts(v []byte) {
	var h datastruck.HealthCheck

	//health check to nuts
	if val, err = etcd.EtcGet(c.HealthCheckEtcPrefix + strFirstToUpper(string(v))); err != nil {
		log.Println(MyERR.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}

	if err := json.Unmarshal([]byte(val), &h); err != nil {
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
		_ = kvnuts.SAdd(c.TagSuccessCode+string(v), v, t)
	}
	for _, t := range h.UnHealth.FailuresStatus {
		//log.Println(t)
		_ = kvnuts.SAdd(c.TagFailureCode+string(v), v, t)
	}
}

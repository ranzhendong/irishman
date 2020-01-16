package healthcheck

import (
	"datastruck"
	"encoding/json"
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"log"
	"lrishman/pkg/kvnuts"
)

type healthCheck struct {
	Health   health   `json:"health"`
	UnHealth unHealth `json:"unhealth"`
}

type health struct {
	Interval       int `json:"interval"`
	SuccessTime    int `json:"successTime"`
	SuccessTimeout int `json:"successTimeout"`
}

//template and put UnHealth
type unHealth struct {
	Interval        int `json:"interval"`
	FailuresTime    int `json:"failuresTime"`
	FailuresTimeout int `json:"failuresTimeout"`
}

func UpDownToNuts() {
	var (
		err          error
		val          string
		upstreamList [][]byte
		u            upstream
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(11012), fmt.Sprintf("%v", err))
	}

	//get upstream list from nutsDB
	_, upstreamList = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

	//split Up, Down and UpstreamList
	for _, v := range upstreamList {
		UpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(string(v))
		if err, val = etcd.EtcGet(UpstreamName); err != nil {
			log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
		}

		if err := json.Unmarshal([]byte(val), &u); err != nil {
			log.Println(ErrH.ErrorLog(11005))
		}

		for _, v := range u.Pool {
			if v.Status == "up" {
				_ = kvnuts.SAdd(c.NutsDB.Tag.Up, u.UpstreamName, v.IpPort)
			} else {
				_ = kvnuts.SAdd(c.NutsDB.Tag.Down, u.UpstreamName, v.IpPort)
			}
		}

	}
	TempToNuts()
}

func TempToNuts() {
	var (
		err          error
		val          string
		upstreamList [][]byte
		h            healthCheck
	)
	//get upstream list from nutsDB
	_, upstreamList = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

	for _, v := range upstreamList {
		UpstreamName := c.HealthCheck.EtcdPrefix + strFirstToUpper(string(v))
		if err, val = etcd.EtcGet(UpstreamName); err != nil {
			log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
		}

		if err := json.Unmarshal([]byte(val), &h); err != nil {
			log.Println(ErrH.ErrorLog(11005))
		}

		_ = kvnuts.LAdd(string(v), v, h.UnHealth.FailuresTimeout)
		_ = kvnuts.LAdd(string(v), v, h.UnHealth.FailuresTime)
		_ = kvnuts.LAdd(string(v), v, h.UnHealth.Interval)
		_ = kvnuts.LAdd(string(v), v, h.Health.SuccessTimeout)
		_ = kvnuts.LAdd(string(v), v, h.Health.SuccessTime)
		_ = kvnuts.LAdd(string(v), v, h.Health.Interval)

	}

	for _, v := range upstreamList {
		log.Println("my string", string(v))
		_, itmes := kvnuts.LIndex(string(v), v, 0, 1)
		for _, v := range itmes {
			log.Println(string(v))
		}
	}
}

func UpHC(a []string, h datastruck.HealthCheck) {
	log.Println(a)
	log.Println(h)
	for _, v := range a {
		log.Println(v)
		log.Println(TCP(v, h.Health.SuccessTimeout))
		//_ = HTTP("vmims.eguagua.cn", h.Health.SuccessTimeout)go get github.com/boltdb/bolt/...
		log.Println(HTTP(v+h.CheckPath, h.Health.SuccessTimeout))
		log.Println(kvnuts.Put(v, v, h.Health.SuccessTime))
		log.Println(kvnuts.Get(v, v, "i"))
	}

}

func DownHC() {

}

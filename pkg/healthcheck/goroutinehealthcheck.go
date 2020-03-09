package healthcheck

import (
	"encoding/json"
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"kvnuts"
	"log"
	"time"
)

type healthCheck struct {
	CheckProtocol string   `json:"checkProtocol"`
	CheckPath     string   `json:"checkPath"`
	Health        health   `json:"health"`
	UnHealth      unHealth `json:"unhealth"`
}

type health struct {
	Interval       int   `json:"interval"`
	SuccessTime    int   `json:"successTime"`
	SuccessTimeout int   `json:"successTimeout"`
	SuccessStatus  []int `json:"successStatus"`
}

//template and put UnHealth
type unHealth struct {
	Interval        int   `json:"interval"`
	FailuresTime    int   `json:"failuresTime"`
	FailuresTimeout int   `json:"failuresTimeout"`
	FailuresStatus  []int `json:"failuresStatus"`
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

/*
[CheckProtocol, CheckPath, Health.Interval, Health.SuccessTime, Health.SuccessTimeout, UnHealth.Interval, UnHealth.FailuresTime, UnHealth.FailuresTimeout]
[		0, 			  1, 			2, 				3, 						4, 					5, 						6, 						7	     ]
storage health check info as list, but success code and failures code as set.
*/
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

	HC()
}

func HC() {
	var (
		upstreamList [][]byte
	)

	_, upstreamList = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

	for _, k := range upstreamList {
		log.Println("my string", string(k))
		//list has eight data, so index[0-7]
		log.Println(kvnuts.LIndex(string(k), k, 0, 7))
		if _, item := kvnuts.LIndex(string(k), k, 0, 7); len(item) != 0 {
			hp := string(item[0])
			hps := string(item[1])
			hi, _ := kvnuts.BytesToInt(item[2], true)
			ht, _ := kvnuts.BytesToInt(item[3], true)
			hto, _ := kvnuts.BytesToInt(item[4], true)
			hfi, _ := kvnuts.BytesToInt(item[5], true)
			hft, _ := kvnuts.BytesToInt(item[6], true)
			hfto, _ := kvnuts.BytesToInt(item[7], true)
			//log.Println(string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
			go UpOneStart(string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
			go DownOneStart(string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
			go test(k)
		}
	}
}

func test(v []byte) {
	var l [][]byte
	for {
		time.Sleep(2 * time.Second)
		_, l = kvnuts.SMem(c.NutsDB.Tag.Up, v)
		for _, s := range l {
			log.Println(string(v), "Success:", string(s))
		}
		_, l = kvnuts.SMem(c.NutsDB.Tag.Down, v)
		for _, s := range l {
			log.Println(string(v), "Failure:", string(s))
		}
	}
}

func UpOneStart(upstreamName, protocal, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
	for {
		time.Sleep(time.Duration(sInterval) * time.Millisecond)
		UpHC(upstreamName, protocal, path, fTimes, fTimeout)
	}
}

func DownOneStart(upstreamName, protocal, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
	for {
		time.Sleep(time.Duration(fInterval) * time.Millisecond)
		DownHC(upstreamName, protocal, path, sTimes, sTimeout)
	}
}

func UpHC(upstreamName, protocal, path string, times, timeout int) {
	// get the upstream up list
	_, ipPort := kvnuts.SMem(c.NutsDB.Tag.Up, upstreamName)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		if protocal == "http" {
			_, statusCode := HTTP(string(ip)+path, timeout)
			log.Println(upstreamName, string(ip), statusCode)

			//the status code can not be in failure, and must be in success code.
			if !kvnuts.SIsMem(c.NutsDB.Tag.FailureCode+upstreamName, upstreamName, statusCode) &&
				kvnuts.SIsMem(c.NutsDB.Tag.SuccessCode+upstreamName, upstreamName, statusCode) {
				continue
			}
		} else {
			if TCP(string(ip), timeout) {
				continue
			}
		}

		if CodeCount(upstreamName+string(ip), "f", times) {
			_ = kvnuts.SRem(c.NutsDB.Tag.Up, upstreamName, ip)
			_ = kvnuts.SAdd(c.NutsDB.Tag.Down, upstreamName, ip)
		}
	}
}

func DownHC(upstreamName, protocal, path string, times, timeout int) {
	// get the upstream down list
	_, ipPort := kvnuts.SMem(c.NutsDB.Tag.Down, upstreamName)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		log.Println(string(ip))
		if protocal == "http" {
			_, statusCode := HTTP(string(ip)+path, timeout)
			log.Println(upstreamName, string(ip), statusCode)

			//the status code must be in success
			if !kvnuts.SIsMem(c.NutsDB.Tag.SuccessCode+upstreamName, upstreamName, statusCode) {
				continue
			}
		} else {
			if !TCP(string(ip), timeout) {
				continue
			}
		}

		if CodeCount(upstreamName+string(ip), "s", times) {
			_ = kvnuts.SRem(c.NutsDB.Tag.Down, upstreamName, ip)
			_ = kvnuts.SAdd(c.NutsDB.Tag.Up, upstreamName, ip)
		}
	}
}

// success && failed counter
func CodeCount(n, key string, times int) bool {
	log.Println(kvnuts.Get(n, key, "i"))
	err, _, nTime := kvnuts.Get(n, key, "i")

	//first be counted
	if err != nil {
		_ = kvnuts.Put(n, key, 1)
		return false
	}

	//counted times less than healthCheck items
	if nTime < times {
		_ = kvnuts.Put(n, key, nTime+1)
		return false
	}

	_ = kvnuts.Del(n, key)
	return true
}

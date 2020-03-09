package healthcheck

import (
	"context"
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

func HC() {
	var (
		upstreamList [][]byte
	)

	_, upstreamList = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)

	//ctx, cancel := context.WithCancel(context.Background())

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
			go UpOneStart(context.TODO(), string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
			go DownOneStart(context.TODO(), string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
			//go UpOneStart(ctx, string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
			//go DownOneStart(ctx, string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
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

func UpOneStart(ctx context.Context, upstreamName, protocal, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
	for {
		time.Sleep(time.Duration(sInterval) * time.Millisecond)
		UpHC(upstreamName, protocal, path, fTimes, fTimeout)
	}
}

func DownOneStart(ctx context.Context, upstreamName, protocal, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
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

package healthcheck

import (
	"context"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
	"time"
)

//healthCheck : template for goroutines
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

type ctxUpstreamList struct {
	upstreamList [][]byte
	ctx          context.Context
	cancel       context.CancelFunc
}

type ctxStart struct {
	upstreamList [][]byte
	ctx          context.Context
}

var (
	c                            datastruck.Config
	ctxCancelChan                = make(chan context.CancelFunc)
	ctxStartChan                 = make(chan ctxStart)
	ctxUpstreamListChan          = make(chan ctxUpstreamList)
	ctxUpstreamListStartFlagChan = make(chan int)
	ctxUpOneStartStartFlagChan   = make(chan int)
	ctxDownOneStartStartFlagChan = make(chan int)
)

//HC : new health check
func HC() {

	// set bit, tell hc controller need to be updated
	go FalgHC()
	for {
		select {
		case Cancels := <-ctxCancelChan:
			log.Println("Cancels := <-ctxCancelChan:")
			Cancels()
		case Start := <-ctxStartChan:
			log.Println("Start := <-ctxStartChan:")
			go upstreamList(Start.ctx, Start.upstreamList)
		case cu := <-ctxUpstreamListChan:
			log.Println("cu := <-ctxUpstreamListChan:")
			go upstreamList(cu.ctx, cu.upstreamList)
			ctxUpstreamListStartFlagChan <- 1
		}
	}
}

func FalgHC() {
	var (
		upstreamList [][]byte
		cul          ctxUpstreamList
		//st           ctxStart
	)

	//first hc
	_ = kvnuts.Del("FalgHC", "FalgHC")
	upstreamList, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("first hc upstreamList", upstreamList)
	cul.upstreamList = upstreamList
	cul.ctx = ctx
	cul.cancel = cancel
	ctxUpstreamListChan <- cul

	for {
		time.Sleep(1 * time.Second)
		if _, _, err := kvnuts.Get("FalgHC", "FalgHC", "i"); err == nil {
			upstreamList, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)
			_ = kvnuts.Del("FalgHC", "FalgHC")
			log.Println("_ = kvnuts.Del....")
			ctxCancelChan <- cul.cancel
			//st.ctx = cul.ctx
			//st.upstreamList = upstreamList
			//ctxStartChan <- st
		} else {
			log.Println("nothing.....")
		}
	}
}

func upstreamList(ctx context.Context, upstreamList [][]byte) {
	ctxs, cancels := context.WithCancel(context.Background())
	for {
		select {
		case <-ctx.Done():
			log.Println("upstreamList退出...", ctx.Err())
			cancels()
			return
		case <-ctxUpstreamListStartFlagChan:
			log.Println(upstreamList, "upstreamList goroutine监控中...")
			time.Sleep(2 * time.Second)
			for _, k := range upstreamList {
				log.Println("my string", string(k))
				//list has eight data, so index[0-7]
				log.Println(kvnuts.LIndex(string(k), k, 0, 7))
				if item, _ := kvnuts.LIndex(string(k), k, 0, 7); len(item) != 0 {
					hp := string(item[0])
					hps := string(item[1])
					hi, _ := kvnuts.BytesToInt(item[2], true)
					ht, _ := kvnuts.BytesToInt(item[3], true)
					hto, _ := kvnuts.BytesToInt(item[4], true)
					hfi, _ := kvnuts.BytesToInt(item[5], true)
					hft, _ := kvnuts.BytesToInt(item[6], true)
					hfto, _ := kvnuts.BytesToInt(item[7], true)
					go UpOneStart(ctxs, string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
					go DownOneStart(ctxs, string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
					go test(k)
					ctxUpOneStartStartFlagChan <- 1
					ctxDownOneStartStartFlagChan <- 1
				}
			}
		}
	}

}

func test(v []byte) {
	var l [][]byte
	for {
		time.Sleep(2 * time.Second)
		l, _ = kvnuts.SMem(c.NutsDB.Tag.Up, v)
		for _, s := range l {
			log.Println(string(v), "Success:", string(s))
		}
		l, _ = kvnuts.SMem(c.NutsDB.Tag.Down, v)
		for _, s := range l {
			log.Println(string(v), "Failure:", string(s))
		}
	}
}

//UpOneStart : up status health check driver
func UpOneStart(ctx context.Context, upstreamName, protocal, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
	for {
		select {
		case <-ctx.Done():
			log.Println("UpOneStart退出...", ctx.Err())
			return
		case <-ctxUpOneStartStartFlagChan:
			log.Println("UpOneStart goroutine监控中...")
			//time.Sleep(2 * time.Second)
			for {
				time.Sleep(time.Duration(sInterval) * time.Millisecond)
				UpHC(upstreamName, protocal, path, fTimes, fTimeout)
			}
		}
	}

}

//DownOneStart : down status health check driver
func DownOneStart(ctx context.Context, upstreamName, protocal, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {

	for {
		select {
		case <-ctx.Done():
			log.Println("DownOneStart监控停止了...", ctx.Err())
			return
		case <-ctxDownOneStartStartFlagChan:
			log.Println("DownOneStart goroutine监控中...")
			//time.Sleep(2 * time.Second)
			for {
				time.Sleep(time.Duration(fInterval) * time.Millisecond)
				DownHC(upstreamName, protocal, path, sTimes, sTimeout)
			}
		}
	}

}

//UpHC : up status ip&port check
func UpHC(upstreamName, protocal, path string, times, timeout int) {
	// get the upstream up list
	ipPort, _ := kvnuts.SMem(c.NutsDB.Tag.Up, upstreamName)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		if protocal == "http" {
			statusCode, _ := HTTP(string(ip)+path, timeout)
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

//DownHC : down status ip&port check
func DownHC(upstreamName, protocal, path string, times, timeout int) {
	// get the upstream down list
	ipPort, _ := kvnuts.SMem(c.NutsDB.Tag.Down, upstreamName)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		log.Println(string(ip))
		if protocal == "http" {
			statusCode, _ := HTTP(string(ip)+path, timeout)
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

//CodeCount : success && failed counter
func CodeCount(n, key string, times int) bool {
	log.Println(kvnuts.Get(n, key, "i"))
	_, nTime, err := kvnuts.Get(n, key, "i")

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

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

//ctxStart: storage upstream list
type ctxStart struct {
	upstreamListBytes [][]byte
}

// upstreamName, protocol, path, sInterval, fTimes, fTimeout
type UpHCS struct {

	//health check upstream name
	un string

	//health check protocol, such as: http,tcp
	ptc string

	// health check successful Interval
	p string

	//successful Interval
	si int

	//failed times
	ft int

	//failed timeout
	fto int
}

type DownHCS struct {
}

var (
	upstreamListBytes   [][]byte
	c                   datastruck.Config
	ctxCancelChan       = make(chan context.CancelFunc, 1)
	ctxRestartHCChan    = make(chan ctxStart)
	ctxStartCancelChan  = make(chan int)
	ctxFirstHCChan      = make(chan int)
	ctxUpstreamListChan = make(chan int)
	ctxUpOneStartChan   = make(chan int)
	ctxDownOneStartChan = make(chan int)
)

//HC : new health check
func HC() {
	var (
		rootCtx context.Context
		cancel  context.CancelFunc
	)

	//first goroutines to hc
	ctx := context.Background()
	rootCtx, cancel = context.WithCancel(ctx)

	//set cancel func to channel
	ctxCancelChan <- cancel

	//check flag if exist
	go FlagHC()

	//if cancel function be triggered
	go ifCancel()

	// start the real hc
	for {
		select {

		//restart hc
		case s := <-ctxRestartHCChan:
			time.Sleep(1 * time.Second)
			rootCtx, cancel = context.WithCancel(ctx)
			ctxCancelChan <- cancel
			log.Println("Start := <-ctxStartChan:")
			go upstreamList(rootCtx, s.upstreamListBytes)
			ctxUpstreamListChan <- 1

		//trigger first HC
		case <-ctxFirstHCChan:
			log.Println("cu := <-ctxUpstreamListChan:")
			go upstreamList(rootCtx, upstreamListBytes)
			ctxUpstreamListChan <- 1
		}
	}
}

//ifCancel function
func ifCancel() {
	for {
		select {
		case <-ctxStartCancelChan:
			cancelFuncs := <-ctxCancelChan
			log.Println("Cancels := <-ctxCancelChan:")
			cancelFuncs()
		}
	}
}

//FlagHC: check flag if exist
func FlagHC() {
	var (
		st ctxStart
	)

	//del flag, avoid key exist
	_ = kvnuts.Del("FalgHC", "FalgHC")

	//get the upstreamListBytes
	upstreamListBytes, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)
	//log.Println("first hc upstreamList", upstreamListBytes)
	//trigger first hc
	ctxFirstHCChan <- 1

	//check flag if exist
	for {
		time.Sleep(1 * time.Second)
		if _, _, err := kvnuts.Get("FalgHC", "FalgHC", "i"); err == nil {
			upstreamListBytes, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)
			_ = kvnuts.Del("FalgHC", "FalgHC")
			log.Println("_ = kvnuts.Del....")

			//if flag exist, trigger ctx cancel function
			ctxStartCancelChan <- 1

			//ready for next hc
			st.upstreamListBytes = upstreamListBytes

			//restart next hc, using the new upstreamListBytes
			ctxRestartHCChan <- st
		} else {
			log.Println("nothing.....")
		}
	}
}

//upstreamList: distribute function about hc
func upstreamList(ctx context.Context, upstreamList [][]byte) {
	for {
		select {

		//if ctx cancel function is triggered, exit
		case <-ctx.Done():
			log.Println("upstreamList......退出...", ctx.Err())
			return

		//start hc
		case <-ctxUpstreamListChan:
			log.Println(upstreamList, "upstreamList goroutine监控中...")
			for _, k := range upstreamList {
				//log.Println("my string", string(k))
				////list has eight data, so index[0-7]
				//log.Println(kvnuts.LIndex(string(k), k, 0, 7))

				//get hc template args
				if item, _ := kvnuts.LIndex(string(k), k, 0, 7); len(item) != 0 {
					hp := string(item[0])
					hps := string(item[1])
					hi, _ := kvnuts.BytesToInt(item[2], true)
					ht, _ := kvnuts.BytesToInt(item[3], true)
					hto, _ := kvnuts.BytesToInt(item[4], true)
					hfi, _ := kvnuts.BytesToInt(item[5], true)
					hft, _ := kvnuts.BytesToInt(item[6], true)
					hfto, _ := kvnuts.BytesToInt(item[7], true)

					//hc for ipPort, who's status is 'down'
					go UpOneStart(ctx, string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
					ctxUpOneStartChan <- 1

					//hc for ipPort, who's status is 'up'
					go DownOneStart(ctx, string(k), hp, hps, hi, ht, hto, hfi, hft, hfto)
					ctxDownOneStartChan <- 1

					go test(k)
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
func UpOneStart(ctx context.Context, upstreamName, protocol, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
	for {
		select {

		//ready for be triggered to hc
		case <-ctxUpOneStartChan:
			log.Println(upstreamName, "UpOneStart goroutine监控中...")

			//periodically round
			for {
				select {

				//if ctx cancel function is triggered, exit
				case <-ctx.Done():
					log.Println(upstreamName, "UpOneStart监控停止了...????", ctx.Err())
					return

				//default to hc
				default:
					time.Sleep(time.Duration(sInterval) * time.Millisecond)
					UpHC(upstreamName, protocol, path, fTimes, fTimeout)
				}

			}
		}
	}

}

//DownOneStart : down status health check driver
func DownOneStart(ctx context.Context, upstreamName, protocol, path string, sInterval, sTimes, sTimeout, fInterval, fTimes, fTimeout int) {
	for {
		select {

		//ready for be triggered to hc
		case <-ctxDownOneStartChan:
			log.Println(upstreamName, "DownOneStart goroutine监控中...")

			//periodically round
			for {
				select {
				//if ctx cancel function is triggered, exit
				case <-ctx.Done():
					log.Println(upstreamName, "DownOneStart监控停止了...????", ctx.Err())
					return
				default:
					time.Sleep(time.Duration(fInterval) * time.Millisecond)
					DownHC(upstreamName, protocol, path, sTimes, sTimeout)
				}
			}
		}
	}
}

//UpHC : up status ip&port check
func UpHC(upstreamName, protocol, path string, times, timeout int) {
	// get the upstream up list
	ipPort, _ := kvnuts.SMem(c.NutsDB.Tag.Up, upstreamName)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		if protocol == "http" {
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
func DownHC(upstreamName, protocol, path string, times, timeout int) {
	// get the upstream down list
	ipPort, _ := kvnuts.SMem(c.NutsDB.Tag.Down, upstreamName)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		log.Println(string(ip))
		if protocol == "http" {
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

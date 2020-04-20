package gorountines

import (
	"context"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
	"time"
)

//ctxStart: storage upstream list
type ctxStart struct {
	upstreamListBytes [][]byte
}

//healthCheckUpstreamName: health check upstream name
type healthCheckUpstreamName string

//healthCheckProtocol: health check protocol, such as: http,tcp
type healthCheckProtocol string

//healthCheckPath: health check successful Interval
type healthCheckPath string

//healthChecksSuccessInterval: successful Interval
type healthCheckSuccessInterval int

//healthChecksSuccessTimes: successful times
type healthCheckSuccessTimes int

//healthChecksSuccessTimeout: successful timeout
type healthCheckSuccessTimeout int

//healthChecksFailureInterval: failed Interval
type healthCheckFailureInterval int

//healthCheckFailureTimes: failed times
type healthCheckFailureTimes int

//healthCheckFailureTimeout: failed timeout
type healthCheckFailureTimeout int

//public: public info
type public struct {
	un  healthCheckUpstreamName
	ptc healthCheckProtocol
	p   healthCheckPath
}

//upHCS: healthCheck for status is 'down'
type upHCS struct {
	pu  public
	si  healthCheckSuccessInterval
	ft  healthCheckFailureTimes
	fto healthCheckFailureTimeout
}

//downHCS: healthCheck for status is 'up'
type downHCS struct {
	pu  public
	fi  healthCheckFailureInterval
	st  healthCheckSuccessTimes
	sto healthCheckSuccessTimeout
}

var (
	upstreamListBytes   [][]byte
	ctxCancelChan       = make(chan context.CancelFunc, 1)
	ctxRestartHCChan    = make(chan ctxStart)
	ctxStartCancelChan  = make(chan int)
	ctxFirstHCChan      = make(chan int)
	ctxUpstreamListChan = make(chan int)
	ctxUpOneStartChan   = make(chan int)
	ctxDownOneStartChan = make(chan int)
)

//HC : new health check
func startHealthCheck() {
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
		if _, _, err := kvnuts.Get("FlagHC", "FlagHC", "i"); err == nil {
			upstreamListBytes, _ = kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList)
			_ = kvnuts.Del("FlagHC", "FlagHC")
			for i := 0; i < len(upstreamListBytes); i++ {
				log.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", string(upstreamListBytes[i]))
			}
			log.Println("kvNuts.Del....")

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

				//get hc template args, list has eight data, so index[0-7]
				if item, _ := kvnuts.LIndex(string(k), k, 0, 7); len(item) != 0 {
					var (
						p   public
						uhc upHCS
						dhc downHCS
					)

					//set public info
					p.un = healthCheckUpstreamName(string(k))
					p.ptc = healthCheckProtocol(string(item[0]))
					p.p = healthCheckPath(string(item[1]))
					uhc.pu = p
					dhc.pu = p

					//shunt upHC
					uhc.si = healthCheckSuccessInterval(kvnuts.BytesToInt(item[2], true))
					uhc.ft = healthCheckFailureTimes(kvnuts.BytesToInt(item[6], true))
					uhc.fto = healthCheckFailureTimeout(kvnuts.BytesToInt(item[7], true))

					//shunt downHC
					dhc.fi = healthCheckFailureInterval(kvnuts.BytesToInt(item[5], true))
					dhc.st = healthCheckSuccessTimes(kvnuts.BytesToInt(item[3], true))
					dhc.sto = healthCheckSuccessTimeout(kvnuts.BytesToInt(item[4], true))

					//hc for ipPort, who's status is 'down'
					go uhc.upOneStart(ctx)
					ctxUpOneStartChan <- 1

					//hc for ipPort, who's status is 'up'
					go dhc.downOneStart(ctx)
					ctxDownOneStartChan <- 1

					go test(ctx, k)
				}
			}
		}
	}
}

//test: print logs about nutsDB all up and down list
func test(ctx context.Context, v []byte) {
	var l [][]byte
	q := 1

	//periodically round
	for {
		select {
		//if ctx cancel function is triggered, exit
		case <-ctx.Done():
			return

		//default to hc
		default:
			time.Sleep(2 * time.Second)
			l, _ = kvnuts.SMem(c.NutsDB.Tag.Up, v)
			for _, s := range l {
				log.Println("[LOG PRINT] Success - Times:", q, " - Upstream:", string(v), " - IP PORT:", string(s))
			}
			l, _ = kvnuts.SMem(c.NutsDB.Tag.Down, v)
			for _, s := range l {
				log.Println("[LOG PRINT] Failure - Times:", q, " - Upstream:", string(v), " - IP PORT:", string(s))
			}
			q += 1
		}
	}

}

//UpOneStart : up status health check driver
func (uhc upHCS) upOneStart(ctx context.Context) {
	for {
		select {

		//ready for be triggered to hc
		case <-ctxUpOneStartChan:
			log.Println(uhc.pu.un, "UpOneStart goroutine监控中...")

			//periodically round
			for {
				select {

				//if ctx cancel function is triggered, exit
				case <-ctx.Done():
					log.Println(uhc.pu.un, "UpOneStart监控停止了...????", ctx.Err())
					return

				//default to hc
				default:
					time.Sleep(time.Duration(uhc.si) * time.Millisecond)
					uhc.upHC()
				}

			}
		}
	}

}

//DownOneStart : down status health check driver
func (dhc downHCS) downOneStart(ctx context.Context) {
	for {
		select {

		//ready for be triggered to hc
		case <-ctxDownOneStartChan:
			log.Println(dhc.pu.un, "DownOneStart goroutine监控中...")

			//periodically round
			for {
				select {
				//if ctx cancel function is triggered, exit
				case <-ctx.Done():
					log.Println(dhc.pu.un, "DownOneStart监控停止了...????", ctx.Err())
					return
				default:
					time.Sleep(time.Duration(dhc.fi) * time.Millisecond)
					dhc.downHC()
				}
			}
		}
	}
}

//UpHC : up status ip&port check
func (uhc upHCS) upHC() {
	//turn upstream name to string name
	name := string(uhc.pu.un)

	// get the upstream up list
	ipPort, _ := kvnuts.SMem(c.NutsDB.Tag.Up, name)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		if uhc.pu.ptc == "http" {
			statusCode, _ := HTTP(string(ip)+string(uhc.pu.p), int(uhc.fto))
			log.Println(name, string(ip), statusCode)

			//the status code can not be in failure, and must be in success code.
			if !kvnuts.SIsMem(c.NutsDB.Tag.FailureCode+name, name, statusCode) &&
				kvnuts.SIsMem(c.NutsDB.Tag.SuccessCode+name, name, statusCode) {
				continue
			}
		} else {
			if TCP(string(ip), int(uhc.fto)) {
				continue
			}
		}

		if CodeCount(name+string(ip), "f", int(uhc.ft)) {
			_ = kvnuts.SRem(c.NutsDB.Tag.Up, name, ip)
			time.Sleep(50 * time.Millisecond)
			_ = kvnuts.SAdd(c.NutsDB.Tag.Down, name, ip)
		}
	}
}

//DownHC : down status ip&port check
func (dhc downHCS) downHC() {
	//turn upstream name to string name
	name := string(dhc.pu.un)

	// get the upstream down list
	ipPort, _ := kvnuts.SMem(c.NutsDB.Tag.Down, name)
	if len(ipPort) == 0 {
		return
	}

	//check every ip port
	for i := 0; i < len(ipPort); i++ {
		ip := ipPort[i]
		log.Println(string(ip))
		if dhc.pu.ptc == "http" {
			statusCode, _ := HTTP(string(ip)+string(dhc.pu.p), int(dhc.sto))
			log.Println(name, string(ip), statusCode)

			//the status code must be in success
			if !kvnuts.SIsMem(c.NutsDB.Tag.SuccessCode+name, name, statusCode) {
				continue
			}
		} else {
			if !TCP(string(ip), int(dhc.sto)) {
				continue
			}
		}

		if CodeCount(name+string(ip), "s", int(dhc.st)) {
			_ = kvnuts.SRem(c.NutsDB.Tag.Down, name, ip)
			_ = kvnuts.SAdd(c.NutsDB.Tag.Up, name, ip)
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

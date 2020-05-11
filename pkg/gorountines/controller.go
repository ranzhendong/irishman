package gorountines

import (
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	"log"
	"time"
)

//func StartHealthCheck() {
//
//}

type GoroutinesMessage struct {
	FlagStartHeathCheck chan int
	StartHealthCheck    func()
	//Error func(*framework.QueuedPodInfo, error)
}

var c datastruck.Config

//Factory: goroutines
func Factory() bool {
	var err error
	log.Println("Factory")

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(0151), fmt.Sprintf("%v", err))
		return false
	}

	//initialize healthCheck
	if err := healthcheck.InitHealthCheck(time.Now()); err.Error != "" {
		return false
	}

	goChan := GoroutinesMessage{
		FlagStartHeathCheck: make(chan int, 1),
	}

	//start health check
	go goChan.startHealthCheck()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			goChan.FlagStartHeathCheck <- 1
		}
	}()

	//create watcher based prefix , c.HealthCheck.EtcdPrefix
	go etcWatcher(c.Upstream.EtcdPrefix, c.HealthCheck.EtcdPrefix)

	//create watcher about monitor nutsDB upstream flag
	go FlagUpstreamNutsDB()

	//create watcher about monitor nutsDB health check flag
	go FlagHCNutsDB()

	//create watcher about monitor nutsDB upstream flag, and generate upstream
	go FlagStartUpstreamNutsDB()

	return true
}

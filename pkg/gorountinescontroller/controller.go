package gorountinescontroller

import (
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	"log"
	"time"
)

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

	//go startHealthCheck()
	go etcdWatcher()
	return true
}

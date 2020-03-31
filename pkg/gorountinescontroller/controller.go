package gorountinescontroller

import (
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	"log"
	"time"
)

//Factory: goroutines
func Factory() {
	//initialize health check
	go healthcheck.InitHealthCheck(time.Now())
	log.Println("Factory")
}

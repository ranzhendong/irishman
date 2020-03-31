package gorountinescontroller

import (
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	"log"
	"time"
)

//Factory: goroutines
func Factory() {
	//initialize healthCheck
	go healthcheck.InitHealthCheck(time.Now())

	//HC()
	log.Println("Factory")
}

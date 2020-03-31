package gorountinescontroller

import (
	"log"
)

//Factory: goroutines
func Factory() {
	//initialize health check
	//go healthcheck.InitHealthCheck(time.Now())
	log.Println("Factory")
}

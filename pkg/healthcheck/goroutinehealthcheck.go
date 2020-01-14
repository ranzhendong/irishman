package healthcheck

import (
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"log"
)

type healthCheck struct {
	HealthCheckName string `json:"healthCheckName"`
}

func SplitUpstreamIpPort() {

	var (
		err error
		val string
		//h                              healthCheck
		//upstreamList, downUpstreamList []string
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(11012), fmt.Sprintf("%v", err))
	}

	//get key from etcd
	if err, val = etcd.EtcGet("UpstreamList"); err != nil {
		log.Println(ErrH.ErrorLog(1102), fmt.Sprintf("; %v", err))
	}
	log.Println(val)

}

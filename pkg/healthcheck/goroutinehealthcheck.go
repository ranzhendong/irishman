package healthcheck

import (
	"datastruck"
	"encoding/json"
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
		err          error
		val          string
		upstreamList []string
		//u                        upstream
		//upListByte, downListByte []byte
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(11012), fmt.Sprintf("%v", err))
	}

	//get key from etcd
	if err, val = etcd.EtcGet("UpstreamList"); err != nil {
		log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}

	if err := json.Unmarshal([]byte(val), &upstreamList); err != nil {
		log.Println(ErrH.ErrorLog(11005))
	}

	//for _, v := range upstreamList {
	//	var upList, downList []string
	//	UpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(v)
	//	if err, val = etcd.EtcGet(UpstreamName); err != nil {
	//		log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
	//	}
	//	//log.Println(val)
	//	if err := json.Unmarshal([]byte(val), &u); err != nil {
	//		log.Println(ErrH.ErrorLog(11005))
	//	}
	//	log.Println(u)
	//	for _, v := range u.Pool {
	//		log.Println(v)
	//		log.Println(v.Status)
	//		if v.Status == "up" {
	//			upList = append(upList, v.IpPort)
	//		} else {
	//			downList = append(downList, v.IpPort)
	//		}
	//	}
	//
	//	log.Println(upList)
	//	log.Println(downList)
	//
	//	if upListByte, err = json.Marshal(upList); err != nil {
	//		log.Println(ErrH.ErrorLog(11004))
	//	}
	//	if downListByte, err = json.Marshal(downList); err != nil {
	//		log.Println(ErrH.ErrorLog(11004))
	//	}
	//
	//	upName := c.Resource.UpList + v
	//	downName := c.Resource.DownList + v
	//	// etcd put
	//	if err = etcd.EtcPut(upName, string(upListByte)); err != nil {
	//		log.Printf(ErrH.ErrorLog(11101, fmt.Sprintf("%v", err)))
	//	}
	//	// etcd put
	//	if err = etcd.EtcPut(downName, string(downListByte)); err != nil {
	//		log.Printf(ErrH.ErrorLog(11101, fmt.Sprintf("%v", err)))
	//	}
	//}

	var testA []string
	var h datastruck.HealthCheck
	testA = append(testA, "172.16.0.51:2379")
	hc := c.HealthCheck.EtcdPrefix + "Vmims"
	if err, val = etcd.EtcGet(hc); err != nil {
		log.Println(ErrH.ErrorLog(11102), fmt.Sprintf("; %v", err))
	}
	if err := json.Unmarshal([]byte(val), &h); err != nil {
		log.Println(ErrH.ErrorLog(11005))
	}
	UpHC(testA, h)

}

func UpHC(a []string, h datastruck.HealthCheck) {
	log.Println(a)
	log.Println(h)
	for _, v := range a {
		log.Println(v)
		log.Println(TCP(v, h.Health.SuccessTimeout))
		//_ = HTTP("vmims.eguagua.cn", h.Health.SuccessTimeout)
		log.Println(HTTP(v+h.CheckPath, h.Health.SuccessTimeout))
	}

}

func DownHC() {

}

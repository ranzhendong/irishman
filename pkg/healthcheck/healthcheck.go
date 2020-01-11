package healthcheck

import (
	"datastruck"
	"encoding/json"
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"time"
)

type upstream struct {
	UpstreamName string `json:"upstreamName"`
}

var c datastruck.Config

//upper the first letter
func strFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

func GetHealthCheck(jsonObj interface{}, timeNow time.Time) (*ErrH.MyError, string) {
	//log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, "ss"
}

//func PutHealthCheck(jsonObj interface{}, timeNow time.Time) *ErrH.MyError {
//	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
//	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
//}
//
//func PatchHealthCheck(jsonObj interface{}, timeNow time.Time) *ErrH.MyError {
//	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
//	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
//}

func InitHealthCheck(timeNow time.Time) *ErrH.MyError {
	log.Println("InitHealthCheck")

	var (
		err                            error
		val                            []*mvccpb.KeyValue
		upstreamList, downUpstreamList []string
		healthCheckByte, b             []byte
		u                              upstream
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(0151), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 0151, TimeStamp: timeNow}
	}

	EtcUpstreamName := c.Upstream.EtcdPrefix
	//get key from etcd
	if err, _, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
		log.Println(ErrH.ErrorLog(0104), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 0104, TimeStamp: timeNow}
	}

	for _, v := range val {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			downUpstreamList = append(downUpstreamList, u.UpstreamName)
			continue
		}
		upstreamList = append(upstreamList, u.UpstreamName)
	}

	for _, v := range upstreamList {
		EtcUpstreamName := c.HealthCheck.EtcdPrefix + strFirstToUpper(v)
		c.HealthCheck.Template.HealthCheckName = v

		//turn struck to json
		if healthCheckByte, err = json.Marshal(c.HealthCheck.Template); err != nil {
			log.Println(ErrH.ErrorLog(0004))
			return &ErrH.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
		}

		// etcd put
		if err = etcd.EtcPut(EtcUpstreamName, string(healthCheckByte)); err != nil {
			log.Printf(ErrH.ErrorLog(0101, fmt.Sprintf("%v", err)))
			return &ErrH.MyError{Error: err.Error(), Code: 0101, TimeStamp: timeNow}
		}
	}

	a := &ErrH.MyError{Code: 000, TimeStamp: timeNow}
	a.Clock()
	if b, err = json.Marshal(a); err != nil {
		log.Println(ErrH.ErrorLog(0004))
		return &ErrH.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
	}
	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" HealthCheck %v", string(b))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

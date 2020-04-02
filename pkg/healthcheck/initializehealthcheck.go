package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"
	"time"
)

var (
	c datastruck.Config
)

//InitHealthCheck : goroutines for Init Health Check
func InitHealthCheck(timeNow time.Time) *MyERR.MyError {
	log.Println("Init HealthCheck")

	var (
		err                            error
		val                            []*mvccpb.KeyValue
		upstreamList, downUpstreamList []string
		healthCheckByte, b             []byte
		u                              Upstream
		tc                             datastruck.TConfig
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(0151), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 0151, TimeStamp: timeNow}
	}

	//set config to tc
	tc = c.TC()
	T := TConfig{
		tc.UpstreamEtcPrefix,
		tc.HealthCheckEtcPrefix,
		tc.TagUp,
		tc.TagDown,
		tc.TagSuccessCode,
		tc.TagFailureCode,
	}

	//get upstream list key from etcd
	if _, val, err = etcd.EtcGetAll(c.Upstream.EtcdPrefix); err != nil {
		log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 0104, TimeStamp: timeNow}
	}

	//set upstream list storage to nutsDB
	for _, v := range val {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			downUpstreamList = append(downUpstreamList, u.UpstreamName)
			continue
		}

		//SeparateUpstreamToNuts
		T.SeparateUpstreamToNuts([]byte(u.UpstreamName))

		//HealthCheckTemplateToNuts
		T.HealthCheckTemplateToNuts([]byte(u.UpstreamName))

		//set upstream list to nutsDB, as set
		_ = kvnuts.SAdd(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList, u.UpstreamName)

		upstreamList = append(upstreamList, u.UpstreamName)
	}

	//healthCheck template storage to etcd
	for _, v := range upstreamList {
		EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(v)

		//HealthCheckName in configure is not be assigned,
		// so need to be set as same as EtcHealthCheckName
		c.HealthCheck.Template.HealthCheckName = v

		//turn struck to json
		if healthCheckByte, err = json.Marshal(c.HealthCheck.Template); err != nil {
			log.Println(MyERR.ErrorLog(0004))
			return &MyERR.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
		}

		// etcd put
		if err = etcd.EtcPut(EtcHealthCheckName, string(healthCheckByte)); err != nil {
			log.Printf(MyERR.ErrorLog(0101, fmt.Sprintf("%v", err)))
			return &MyERR.MyError{Error: err.Error(), Code: 0101, TimeStamp: timeNow}
		}
	}

	a := &MyERR.MyError{Code: 000, TimeStamp: timeNow}
	a.Clock()
	if b, err = json.Marshal(a); err != nil {
		log.Println(MyERR.ErrorLog(0004))
		return &MyERR.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf("Init HealthCheck %v", string(b))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}
}

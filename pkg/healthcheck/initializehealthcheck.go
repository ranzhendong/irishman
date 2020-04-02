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
	u Upstream
)

type UpstreamToNutsDBFlag struct {
	SeparateUpstreamEtcdToNuts int
	HealthCheckEtcdToNuts      int
}

//InitHealthCheck : goroutines for Init Health Check
func InitHealthCheck(timeNow time.Time) *MyERR.MyError {
	log.Println("Init HealthCheck")

	var (
		err error
		val []*mvccpb.KeyValue
		b   []byte
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(0151), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 0151, TimeStamp: timeNow}
	}

	//get upstream list key from etcd
	if _, val, err = etcd.EtcGetAll(c.Upstream.EtcdPrefix); err != nil {
		log.Println(MyERR.ErrorLog(0104), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 0104, TimeStamp: timeNow}
	}

	//healthCheck template storage to etcd
	HealthCheckTemplateToEtcd(val)

	//set upstream list storage to nutsDB, set flag
	utnf := UpstreamToNutsDBFlag{
		1,
		1}.UpstreamAndHCEtcdToNutsDB
	utnf(val)

	a := &MyERR.MyError{Code: 000, TimeStamp: timeNow}
	a.Clock()
	if b, err = json.Marshal(a); err != nil {
		log.Println(MyERR.ErrorLog(0004))
		return &MyERR.MyError{Error: err.Error(), Code: 0004, TimeStamp: timeNow}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf("Init HealthCheck %v", string(b))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}
}

//UpstreamToNutsDB : set upstream list storage to nutsDB
//set Separate Upstream To Nuts
//set HealthCheck Template To Nuts
//just for init
func (f UpstreamToNutsDBFlag) UpstreamAndHCEtcdToNutsDB(val []*mvccpb.KeyValue) {
	log.Println("UpstreamAndHCEtcdToNutsDB")
	var tc datastruck.TConfig

	// make sure config loaded
	// set config to tc
	tc = c.TC()

	T := TConfig{
		tc.UpstreamEtcPrefix,
		tc.HealthCheckEtcPrefix,
		tc.TagUp,
		tc.TagDown,
		tc.TagSuccessCode,
		tc.TagFailureCode,
	}

	//set upstream list storage to nutsDB
	for _, v := range val {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			continue
		}

		//Separate Upstream Etcd To Nuts
		//Multiplexed flag
		if f.SeparateUpstreamEtcdToNuts != 0 {
			T.SeparateUpstreamEtcdToNuts([]byte(u.UpstreamName))
		}

		//HealthCheck Etcd To Nuts
		//Multiplexed flag
		if f.HealthCheckEtcdToNuts != 0 {
			//HealthCheck Template To Nuts
			T.HealthCheckEtcdToNuts([]byte(u.UpstreamName))
		}

		//set upstream list to nutsDB, as set
		_ = kvnuts.SAdd(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList, u.UpstreamName)
	}

}

//HealthCheckTemplateToEtcd : healthCheck template storage to etcd
func HealthCheckTemplateToEtcd(val []*mvccpb.KeyValue) *MyERR.MyError {
	var healthCheckByte []byte
	log.Println("HealthCheckTemplateToEtcd")

	for _, v := range val {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			continue
		}

		EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(u.UpstreamName)

		//if etcd has the value, skip....
		if _, err = etcd.EtcGet(EtcHealthCheckName); err == nil {
			continue
		}

		//HealthCheckName in configure is not be assigned,
		// so need to be set as same as EtcHealthCheckName
		c.HealthCheck.Template.HealthCheckName = u.UpstreamName

		//turn struck to json
		if healthCheckByte, err = json.Marshal(c.HealthCheck.Template); err != nil {
			log.Println(MyERR.ErrorLog(0004))
			return &MyERR.MyError{Error: err.Error(), Code: 0004}
		}

		// etcd put
		if err = etcd.EtcPut(EtcHealthCheckName, string(healthCheckByte)); err != nil {
			log.Printf(MyERR.ErrorLog(0101, fmt.Sprintf("%v", err)))
			return &MyERR.MyError{Error: err.Error(), Code: 0101}
		}
	}
	return &MyERR.MyError{Code: 000}
}

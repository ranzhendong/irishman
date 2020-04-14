package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"log"
	"time"
)

var (
	c datastruck.Config
	u Upstream
)

type UpstreamToNutsDBFlag struct {
	SeparateUpstreamEtcdToNuts       int
	HealthCheckEtcdToNuts            int
	SeparateUpstreamEtcdToNutsForOne int
	HealthCheckEtcdToNutsForOne      int
	OneKey                           string
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
		SeparateUpstreamEtcdToNuts: 1,
		HealthCheckEtcdToNuts:      1}.UpstreamAndHCFromEtcdToNutsDB
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

package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"log"
)

//HealthCheckTemplateToEtcd : healthCheck template storage to etcd
//Even though it's UpstreamName, HealthCheckTemplate is as same as Upstream
func HealthCheckTemplateToEtcd(val []*mvccpb.KeyValue) *MyERR.MyError {
	log.Println("HealthCheckTemplateToEtcd")

	for _, v := range val {
		if err := json.Unmarshal(v.Value, &u); err != nil {
			continue
		}
		PostHealthCheckTemplateToEtcd(u.UpstreamName)
	}
	return &MyERR.MyError{Code: 000}
}

//PostHealthCheckTemplateToEtcd: set template to etcd, but need to config
func PostHealthCheckTemplateToEtcd(name string) {
	var healthCheckByte []byte

	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(name)

	//if etcd has the value, skip....
	if _, err = etcd.EtcGet(EtcHealthCheckName); err == nil {
		return
	}

	//HealthCheckName in configure is not be assigned,
	// so need to be set as same as EtcHealthCheckName
	c.HealthCheck.Template.HealthCheckName = name

	//turn struck to json
	if healthCheckByte, err = json.Marshal(c.HealthCheck.Template); err != nil {
		log.Println(MyERR.ErrorLog(0004))
		return
	}

	// etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(healthCheckByte)); err != nil {
		log.Printf(MyERR.ErrorLog(0101, fmt.Sprintf("%v", err)))
		return
	}
}

//PostHealthCheckTemplateToNutsDB: set template to nutsDB, but need to config
func PostHealthCheckTemplateToNutsDB(name string) {
	var tc datastruck.TConfig

	tc = c.TC()

	T := TConfig{
		tc.UpstreamEtcPrefix,
		tc.HealthCheckEtcPrefix,
		tc.TagUp,
		tc.TagDown,
		tc.TagSuccessCode,
		tc.TagFailureCode,
	}
	T.HealthCheckEtcdToNuts([]byte(name), "the one")
}

//DeleteHealthCheckTemplateToEtcd: Delete Health Check
func DeleteHealthCheckTemplateToEtcd(name string) {
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(name)

	//if etcd has the value, skip....
	if err = etcd.EtcDelete(EtcHealthCheckName); err == nil {
		return
	}
}

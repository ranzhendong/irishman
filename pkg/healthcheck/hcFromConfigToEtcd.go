package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"log"
)

//HealthCheckTemplateToEtcd : healthCheck template storage to etcd
//Even though it's UpstreamName, HealthCheckTemplate is as same as Upstream
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

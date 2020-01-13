package healthcheck

import (
	"datastruck"
	"encoding/json"
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"gopkg.in/fatih/set.v0"
	"log"
	"time"
)

var c datastruck.Config

type upstream struct {
	UpstreamName string `json:"upstreamName"`
}

//repeat remove
func removeRepByMap(slc []int) (result []int) {
	tempMap := map[interface{}]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return
}

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
	var (
		gh  datastruck.GetHealthCheck
		err error
		val string
	)

	////config loading
	//if err = c.Config(); err != nil {
	//	log.Println(ErrH.ErrorLog(4012), fmt.Sprintf("%v", err))
	//	return &ErrH.MyError{Error: err.Error(), Code: 4012, TimeStamp: timeNow}, ""
	//}

	//judge
	if err = gh.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(7003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 7003, TimeStamp: timeNow}, ""
	}

	//get values
	if gh.HealthCheckName == "ALL" {
		EtcHealthCheckName := c.HealthCheck.EtcdPrefix
		//get key from etcd
		if err, val, _ = etcd.EtcGetAll(EtcHealthCheckName); err != nil {
			log.Println(ErrH.ErrorLog(7104), fmt.Sprintf("%v", err))
			return &ErrH.MyError{Error: err.Error(), Code: 7104, TimeStamp: timeNow}, ""
		}
		log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get ALL HealthCheck %v, Values %v", EtcHealthCheckName, val)))
		return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
	}

	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(gh.HealthCheckName)

	//get key from etcd
	if err, val = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Println(err)
		log.Println(ErrH.ErrorLog(7102), fmt.Sprintf("; %v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 7102, TimeStamp: timeNow}, ""
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get HealthCheck: %v, Values %v", EtcHealthCheckName, val)))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
}

func PutHealthCheck(jsonObj interface{}, timeNow time.Time) *ErrH.MyError {
	var (
		h     datastruck.HealthCheck
		err   error
		jsonU []byte
	)

	////config loading
	//if err = c.Config(); err != nil {
	//	log.Println(ErrH.ErrorLog(4012), fmt.Sprintf("%v", err))
	//	return &ErrH.MyError{Error: err.Error(), Code: 4012, TimeStamp: timeNow}
	//}

	//judge
	if err = h.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(8003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 8003, TimeStamp: timeNow}
	}

	//turn to json
	if jsonU, err = json.Marshal(h); err != nil {
		log.Println(ErrH.ErrorLog(8004))
		return &ErrH.MyError{Error: err.Error(), Code: 8004, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(h.HealthCheckName)

	//if exist
	if err, _ = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Printf(ErrH.ErrorLog(8102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 8102, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(8101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 8101, TimeStamp: timeNow}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Put HealthCheck: %v, Values %v", EtcHealthCheckName, string(jsonU))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

func PatchHealthCheck(jsonObj interface{}, timeNow time.Time) (a *ErrH.MyError) {
	var (
		ph, etcdph   datastruck.PatchHealthCheck
		err, errs    error
		middleware   interface{}
		jsonU        []byte
		val          string
		sPool, fPool []int
	)

	//judge
	if err = ph.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(9003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 9003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(ph.HealthCheckName)
	//if exist
	if err, val = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Printf(ErrH.ErrorLog(9102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 9102, TimeStamp: timeNow}
	}

	//turn etcd data to struck, for compare and judge
	if err = json.Unmarshal([]byte(val), &etcdph); err != nil {
		log.Printf(ErrH.ErrorLog(9005), fmt.Sprintf("Etcd PatchHealthCheck String: %v", err))
		return &ErrH.MyError{Code: 9005, TimeStamp: timeNow}
	}

	// first set
	etcdph.Status = ph.Status
	etcdph.CheckProtocol = ph.CheckProtocol
	etcdph.CheckPath = ph.CheckPath
	if ph.Health.Interval != 0 {
		etcdph.Health.Interval = ph.Health.Interval
	}
	if ph.Health.SuccessTime != 0 {
		etcdph.Health.SuccessTime = ph.Health.SuccessTime
	}
	if ph.UnHealth.Interval != 0 {
		etcdph.UnHealth.Interval = ph.UnHealth.Interval
	}
	if ph.UnHealth.FailuresTime != 0 {
		etcdph.UnHealth.FailuresTime = ph.UnHealth.FailuresTime
	}
	if ph.UnHealth.FailuresTimeout != 0 {
		etcdph.UnHealth.FailuresTimeout = ph.UnHealth.FailuresTimeout
	}

	//if Health.SuccessStatus and UnHealth.FailuresStatus all nil
	if ph.Health.SuccessStatus == nil && ph.UnHealth.FailuresStatus == nil {
		middleware = etcdph
		goto READY
	}

	//if ph.Health.SuccessStatus not nil
	if ph.Health.SuccessStatus != nil {
		goto SUCCESS
	}

	//if ph.UnHealth.FailuresStatus not nil
	if ph.UnHealth.FailuresStatus != nil {
		goto FAILURES
	}

SUCCESS:
	for i := 0; i < len(ph.Health.SuccessStatus); i++ {
		sPool = append(sPool, ph.Health.SuccessStatus[i])
	}
	for i := 0; i < len(etcdph.Health.SuccessStatus); i++ {
		sPool = append(sPool, etcdph.Health.SuccessStatus[i])
	}
	etcdph.Health.SuccessStatus = removeRepByMap(sPool)
	middleware = etcdph

	//if ph.UnHealth.FailuresStatus still nil
	if ph.UnHealth.FailuresStatus != nil {
		goto FAILURES
	} else {
		goto READY
	}

FAILURES:
	for i := 0; i < len(ph.UnHealth.FailuresStatus); i++ {
		fPool = append(fPool, ph.UnHealth.FailuresStatus[i])
	}
	for i := 0; i < len(etcdph.UnHealth.FailuresStatus); i++ {
		fPool = append(fPool, etcdph.UnHealth.FailuresStatus[i])
	}
	etcdph.UnHealth.FailuresStatus = removeRepByMap(fPool)
	middleware = etcdph

READY:
	//turn struck or map to json
	if jsonU, err = json.Marshal(middleware); err != nil {
		log.Println(ErrH.ErrorLog(9004))
		return &ErrH.MyError{Error: err.Error(), Code: 9004, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(9101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 9101, TimeStamp: timeNow}
	}

	defer func() {
		_ = recover()
		if errs == nil {
			a = &ErrH.MyError{Code: 000, TimeStamp: timeNow}
		}
	}()
	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Patch HealthCheck %v, Values [%v]", EtcHealthCheckName, string(jsonU))))
	return &ErrH.MyError{Code: 000, Error: errs.Error(), TimeStamp: timeNow}
}

func DeleteHealthCheck(jsonObj interface{}, timeNow time.Time) (a *ErrH.MyError) {
	var (
		dh                               datastruck.DeleteHealthCheck
		etcddh                           datastruck.HealthCheck
		err                              error
		val                              string
		jsonU                            []byte
		sPool, fPool                     []int
		dhset, etcddhSet                 set.Interface
		dhsetf, etcddhSetf               set.Interface
		intersectionSet, differenceSet   set.Interface
		intersectionSetF, differenceSetF set.Interface
	)
	//judge
	if err = dh.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(10003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 10003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(dh.HealthCheckName)
	//if exist
	if err, val = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Printf(ErrH.ErrorLog(10102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 10102, TimeStamp: timeNow}
	}

	//turn etcd data to struck, for compare and judge
	if err = json.Unmarshal([]byte(val), &etcddh); err != nil {
		log.Printf(ErrH.ErrorLog(10005), fmt.Sprintf("Etcd PatchHealthCheck String: %v", err))
		return &ErrH.MyError{Code: 10005, TimeStamp: timeNow}
	}

	////并集
	//set.Union(a, b)
	//
	////交集
	//set.Intersection(a, b)
	//
	////差集
	//set.Difference(a, b)

	//if dh.Health.SuccessStatus not nil
	if dh.Health.SuccessStatus != nil {
		goto SUCCESS
	}
	//if dh.UnHealth.FailuresStatus not nil
	if dh.UnHealth.FailuresStatus != nil {
		goto FAILURES
	}

SUCCESS:
	dhset = set.New(set.ThreadSafe)
	for _, v := range dh.Health.SuccessStatus {
		dhset.Add(v)
	}
	etcddhSet = set.New(set.ThreadSafe)
	for _, v := range etcddh.Health.SuccessStatus {
		etcddhSet.Add(v)
	}

	intersectionSet = set.Intersection(dhset, etcddhSet)
	differenceSet = set.Difference(etcddhSet, intersectionSet)

	//if is least one
	if differenceSet.Size() == 0 {
		log.Printf(ErrH.ErrorLog(10154))
		return &ErrH.MyError{Code: 10154, TimeStamp: timeNow}
	}

	for _, v := range differenceSet.List() {
		sPool = append(sPool, v.(int))
	}
	etcddh.Health.SuccessStatus = sPool

	if dh.UnHealth.FailuresStatus != nil {
		goto FAILURES
	}

FAILURES:
	dhsetf = set.New(set.ThreadSafe)
	for _, v := range dh.UnHealth.FailuresStatus {
		dhsetf.Add(v)
	}
	etcddhSetf = set.New(set.ThreadSafe)
	for _, v := range etcddh.UnHealth.FailuresStatus {
		etcddhSetf.Add(v)
	}

	intersectionSetF = set.Intersection(dhsetf, etcddhSetf)
	differenceSetF = set.Difference(etcddhSetf, intersectionSetF)

	if differenceSetF.Size() == 0 {
		log.Printf(ErrH.ErrorLog(10155))
		return &ErrH.MyError{Code: 10155, TimeStamp: timeNow}
	}

	for _, v := range differenceSetF.List() {
		fPool = append(fPool, v.(int))
	}
	etcddh.UnHealth.FailuresStatus = fPool

	//turn struck or map to json
	if jsonU, err = json.Marshal(etcddh); err != nil {
		log.Println(ErrH.ErrorLog(10004))
		return &ErrH.MyError{Error: err.Error(), Code: 10004, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(10101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 10101, TimeStamp: timeNow}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Dlete HealthCheck %v, New Values [%v]", EtcHealthCheckName, string(jsonU))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

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
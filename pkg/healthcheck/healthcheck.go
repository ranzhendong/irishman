package healthcheck

import (
	"encoding/json"
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"gopkg.in/fatih/set.v0"
	"log"
	"time"
)

type J interface{}

//set interface,timestamp, method to struck
type RHCStruck struct {
	J
	T time.Time
	M string
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

//GetHealthCheck : method for get health check
func (r *RHCStruck) GetHealthCheck() (*MyERR.MyError, string) {
	var (
		gh  datastruck.GetHealthCheck
		err error
		val string
	)

	//judge
	if err = gh.JudgeValidator(r.J); err != nil {
		log.Println(MyERR.ErrorLog(7003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 7003, TimeStamp: r.T}, ""
	}

	//get values
	if gh.HealthCheckName == "ALL" {
		EtcHealthCheckName := c.HealthCheck.EtcdPrefix
		//get key from etcd
		if val, _, err = etcd.EtcGetAll(EtcHealthCheckName); err != nil {
			log.Println(MyERR.ErrorLog(7104), fmt.Sprintf("%v", err))
			return &MyERR.MyError{Error: err.Error(), Code: 7104, TimeStamp: r.T}, ""
		}
		log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Get ALL HealthCheck %v, Values %v", EtcHealthCheckName, val)))
		return &MyERR.MyError{Code: 000, TimeStamp: r.T}, val
	}

	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(gh.HealthCheckName)

	//get key from etcd
	if val, err = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Println(err)
		log.Println(MyERR.ErrorLog(7102), fmt.Sprintf("; %v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 7102, TimeStamp: r.T}, ""
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Get HealthCheck: %v, Values %v", EtcHealthCheckName, val)))
	return &MyERR.MyError{Code: 000, TimeStamp: r.T}, val
}

//PutHealthCheck : method for put health check
func (r *RHCStruck) PutHealthCheck() *MyERR.MyError {
	var (
		h     datastruck.HealthCheck
		err   error
		jsonU []byte
	)

	//judge
	if err = h.JudgeValidator(r.J); err != nil {
		log.Println(MyERR.ErrorLog(8003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 8003, TimeStamp: r.T}
	}

	//turn to json
	if jsonU, err = json.Marshal(h); err != nil {
		log.Println(MyERR.ErrorLog(8004))
		return &MyERR.MyError{Error: err.Error(), Code: 8004, TimeStamp: r.T}
	}

	// Characters joining together
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(h.HealthCheckName)

	//if exist
	if _, err = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Printf(MyERR.ErrorLog(8102), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 8102, TimeStamp: r.T}
	}

	//etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(8101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 8101, TimeStamp: r.T}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Put HealthCheck: %v, Values %v", EtcHealthCheckName, string(jsonU))))
	return &MyERR.MyError{Code: 000, TimeStamp: r.T}
}

//PatchHealthCheck : method for patch health check
func (r *RHCStruck) PatchHealthCheck() (a *MyERR.MyError) {
	var (
		ph, etcdph   datastruck.PatchHealthCheck
		err, errs    error
		middleware   interface{}
		jsonU        []byte
		val          string
		sPool, fPool []int
	)

	//judge
	if err = ph.JudgeValidator(r.J); err != nil {
		log.Println(MyERR.ErrorLog(9003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 9003, TimeStamp: r.T}
	}

	// Characters joining together
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(ph.HealthCheckName)
	//if exist
	if val, err = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Printf(MyERR.ErrorLog(9102), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 9102, TimeStamp: r.T}
	}

	//turn etcd data to struck, for compare and judge
	if err = json.Unmarshal([]byte(val), &etcdph); err != nil {
		log.Printf(MyERR.ErrorLog(9005), fmt.Sprintf("Etcd PatchHealthCheck String: %v", err))
		return &MyERR.MyError{Code: 9005, TimeStamp: r.T}
	}

	// first set
	etcdph.Status = ph.Status
	etcdph.CheckProtocol = ph.CheckProtocol
	etcdph.CheckPath = ph.CheckPath
	if ph.Health.Interval != 0 {
		etcdph.Health.Interval = ph.Health.Interval
	}
	if ph.Health.SuccessTimeout != 0 {
		etcdph.Health.SuccessTimeout = ph.Health.SuccessTimeout
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
		log.Println(MyERR.ErrorLog(9004))
		return &MyERR.MyError{Error: err.Error(), Code: 9004, TimeStamp: r.T}
	}

	//etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(9101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 9101, TimeStamp: r.T}
	}

	defer func() {
		_ = recover()
		if errs == nil {
			a = &MyERR.MyError{Code: 000, TimeStamp: r.T}
		}
	}()
	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Patch HealthCheck %v, Values [%v]", EtcHealthCheckName, string(jsonU))))
	return &MyERR.MyError{Code: 000, Error: errs.Error(), TimeStamp: r.T}
}

//DeleteHealthCheck : method for delete health check
func (r *RHCStruck) DeleteHealthCheck() (a *MyERR.MyError) {
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
	if err = dh.JudgeValidator(r.J); err != nil {
		log.Println(MyERR.ErrorLog(10003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 10003, TimeStamp: r.T}
	}

	// Characters joining together
	EtcHealthCheckName := c.HealthCheck.EtcdPrefix + strFirstToUpper(dh.HealthCheckName)
	//if exist
	if val, err = etcd.EtcGet(EtcHealthCheckName); err != nil {
		log.Printf(MyERR.ErrorLog(10102), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 10102, TimeStamp: r.T}
	}

	//turn etcd data to struck, for compare and judge
	if err = json.Unmarshal([]byte(val), &etcddh); err != nil {
		log.Printf(MyERR.ErrorLog(10005), fmt.Sprintf("Etcd PatchHealthCheck String: %v", err))
		return &MyERR.MyError{Code: 10005, TimeStamp: r.T}
	}
	/*
		使用set包来对数组字符串取交集并集差集，
		使用范围也仅限于同一数组内的元素，可以尝试后续对upstream的patch和delete进行代码重构

		//并集
		set.Union(a, b)

		//交集
		set.Intersection(a, b)

		//差集
		set.Difference(a, b)

	*/

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
		log.Printf(MyERR.ErrorLog(10152))
		return &MyERR.MyError{Code: 10152, TimeStamp: r.T}
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
		log.Printf(MyERR.ErrorLog(10153))
		return &MyERR.MyError{Code: 10153, TimeStamp: r.T}
	}

	for _, v := range differenceSetF.List() {
		fPool = append(fPool, v.(int))
	}
	etcddh.UnHealth.FailuresStatus = fPool

	//turn struck or map to json
	if jsonU, err = json.Marshal(etcddh); err != nil {
		log.Println(MyERR.ErrorLog(10004))
		return &MyERR.MyError{Error: err.Error(), Code: 10004, TimeStamp: r.T}
	}

	//etcd put
	if err = etcd.EtcPut(EtcHealthCheckName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(10101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 10101, TimeStamp: r.T}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Dlete HealthCheck %v, New Values [%v]", EtcHealthCheckName, string(jsonU))))
	return &MyERR.MyError{Code: 000, TimeStamp: r.T}
}

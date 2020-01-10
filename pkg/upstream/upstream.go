package upstream

import (
	"datastruck"
	"encoding/json"
	"errorhandle"
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"log"
	"time"
)

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

//repeat remove
func removeRepByMap(slc []map[string]interface{}) (result []map[string]interface{}) {
	tempMap := map[interface{}]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e["ipPort"]] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return
}

// Get upstream
func GetUpstream(jsonObj interface{}, timeNow time.Time) (*errorhandle.MyError, string) {
	var (
		gu  datastruck.GetUpstream
		err error
		val string
	)

	//judge
	if err = gu.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(1003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 1003, TimeStamp: timeNow}, ""
	}

	if gu.UpstreamName == "ALL" {
		EtcUpstreamName := "Upstream"
		//get key from etcd
		if err, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
			log.Println(ErrH.ErrorLog(1104), fmt.Sprintf("%v", err))
			return &ErrH.MyError{Error: err.Error(), Code: 1104, TimeStamp: timeNow}, ""
		}
		log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get ALL Key [%v], Values [%v]", EtcUpstreamName, val)))
		return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
	}

	EtcUpstreamName := "Upstream" + strFirstToUpper(gu.UpstreamName)
	//get key from etcd
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Println(err)
		log.Println(ErrH.ErrorLog(1102), fmt.Sprintf("; %v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 1102, TimeStamp: timeNow}, ""
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
}

// Full Update upstream
func PutUpstream(jsonObj interface{}, timeNow time.Time) *errorhandle.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(2003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 2003, TimeStamp: timeNow}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(ErrH.ErrorLog(2004))
		return &ErrH.MyError{Error: err.Error(), Code: 2004, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//if exist
	if err, _ = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(2102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 2102, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(2101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 2101, TimeStamp: timeNow}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Put Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

// Create Update upstream
func PostUpstream(jsonObj interface{}, timeNow time.Time) *errorhandle.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(3003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 3003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//if repeat
	if err, _ = etcd.EtcGet(EtcUpstreamName); err == nil {
		log.Printf(ErrH.ErrorLog(3103))
		return &ErrH.MyError{Code: 3103, TimeStamp: timeNow}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(ErrH.ErrorLog(3004))
		return &ErrH.MyError{Error: err.Error(), Code: 3004, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(3101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 3101, TimeStamp: timeNow}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Set Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

// Partial upstream
func PatchUpstream(jsonObj interface{}, timeNow time.Time) *errorhandle.MyError {
	var (
		pu           datastruck.PatchUpstream
		etcdpu       datastruck.PatchUpstream
		puData       map[string]interface{}
		etcdData     map[string]interface{}
		UpstreamPool []map[string]interface{}
		middleware   interface{}
		jsonU        []byte
		puByte       []byte
		err          error
		val          string
	)

	//judge
	if err = pu.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(4003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 4003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(pu.UpstreamName)
	//if exist
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(4102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 4102, TimeStamp: timeNow}
	}

	//turn etcd data to struck, for compare and judge
	if err := json.Unmarshal([]byte(val), &etcdpu); err != nil {
		log.Printf(ErrH.ErrorLog(4005), fmt.Sprintf("Etcd PatchUpstream String: %v", err))
		return &ErrH.MyError{Code: 4005, TimeStamp: timeNow}
	}

	if pu.Pool == nil {
		etcdpu.Algorithms = pu.Algorithms
		middleware = etcdpu
		goto JUST
	}

	if err := json.Unmarshal([]byte(val), &etcdData); err != nil {
		log.Printf(ErrH.ErrorLog(4005), fmt.Sprintf("Etcd PatchUpstream Map: %v", err))
		return &ErrH.MyError{Code: 4005, TimeStamp: timeNow}
	}

	//turn struct to json string，then turn json string to map
	if puByte, err = json.Marshal(pu); err != nil {
		log.Printf(ErrH.ErrorLog(4004), fmt.Sprintf("Request PatchUpstream Struck: %v", err))
		return &ErrH.MyError{Code: 4004, TimeStamp: timeNow}
	}
	if err := json.Unmarshal(puByte, &puData); err != nil {
		log.Printf(ErrH.ErrorLog(4005), fmt.Sprintf("Request PatchUpstream Map: %v", err))
		return &ErrH.MyError{Code: 4005, TimeStamp: timeNow}
	}

	/*
		1. ipport相等
		   1.1 weight存在 （更新status和weight）
		   1.2 weight不存在 （只更新status）
		2. ipport不相等
		   1.1 weight存在 （列表增加，更新三条内容）
		   1.2 weight不存在 （放弃）
	*/
	//replace data
	for k, v := range puData {
		if k != "pool" {
			etcdData[k] = puData[k]
		}
		if k == "pool" && v != nil {
			for ek, ev := range etcdData["pool"].([]interface{}) {
				for _, fv := range v.([]interface{}) {
					for k, sv := range fv.(map[string]interface{}) {

						etcdIpPort := etcdData["pool"].([]interface{})[ek].(map[string]interface{})["ipPort"]
						//etcdStatus := etcdData["pool"].([]interface{})[ek].(map[string]interface{})["status"]
						etcdWeight := etcdData["pool"].([]interface{})[ek].(map[string]interface{})["weight"]
						RequestIpPort := fv.(map[string]interface{})["ipPort"]
						RequestStatus := fv.(map[string]interface{})["status"]
						RequestWeight := fv.(map[string]interface{})["weight"]

						if k == "ipPort" && ev.(map[string]interface{})["ipPort"] == sv {
							if fv.(map[string]interface{})["weight"].(float64) == 0 {
								UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdIpPort, "status": RequestStatus, "weight": etcdWeight})
								//UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdData["pool"].([]interface{})[ek].(map[string]interface{})["ipPort"], "status": fv.(map[string]interface{})["status"], "weight": etcdData["pool"].([]interface{})[ek].(map[string]interface{})["weight"]})
								continue
							} else {
								//etcdData["pool"].([]interface{})[ek].(map[string]interface{})["status"] = fv.(map[string]interface{})["status"]
								//etcdData["pool"].([]interface{})[ek].(map[string]interface{})["weight"] = fv.(map[string]interface{})["weight"]
								//UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdData["pool"].([]interface{})[ek].(map[string]interface{})["ipPort"], "status": fv.(map[string]interface{})["status"], "weight": fv.(map[string]interface{})["weight"]})
								UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdIpPort, "status": RequestStatus, "weight": RequestWeight})
							}
						} else if k == "ipPort" && ev.(map[string]interface{})["ipPort"] != sv {
							if fv.(map[string]interface{})["weight"].(float64) == 0 {
								continue
							} else {
								//UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": fv.(map[string]interface{})["ipPort"], "status": fv.(map[string]interface{})["status"], "weight": fv.(map[string]interface{})["weight"]})
								UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": RequestIpPort, "status": RequestStatus, "weight": RequestWeight})
							}
						}
					}
				}
			}
		}
	}

	//de-weight
	etcdData["pool"] = removeRepByMap(UpstreamPool)
	middleware = etcdData

JUST:
	//turn struck or map to json
	if jsonU, err = json.Marshal(middleware); err != nil {
		log.Println(ErrH.ErrorLog(4004))
		return &ErrH.MyError{Error: err.Error(), Code: 4004, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(4101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 4101, TimeStamp: timeNow}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Patch Key [%v], Values %v", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

// Delete upstream or pool's server
func DeleteUpstream(jsonObj interface{}, timeNow time.Time) *ErrH.MyError {
	var (
		du     datastruck.DeleteUpstream
		etcddu datastruck.PatchUpstream
		//duData       map[string]interface{}
		//etcdData     map[string]interface{}
		//UpstreamPool []map[string]interface{}
		//middleware   interface{}
		//jsonU        []byte
		//puByte       []bytes
		err error
		val string
	)

	//judge
	if err = du.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(5003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 5003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(du.UpstreamName)
	//if exist
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(5102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 5102, TimeStamp: timeNow}
	}

	//if no server, mean to delete server
	if len(val) != 0 {
		goto DELETE
	}

	log.Println()
	//turn etcd data to struck, for compare and judge
	if err := json.Unmarshal([]byte(val), &etcddu); err != nil {
		log.Printf(ErrH.ErrorLog(5005), fmt.Sprintf("Etcd DeleteUpstream String: %v", err))
		return &ErrH.MyError{Code: 5005, TimeStamp: timeNow}
	}

	//if pu.Pool == nil {
	//	etcdpu.Algorithms = pu.Algorithms
	//	middleware = etcdpu
	//	//goto JUST
	//}

DELETE:
	//_ = etcd.EtcDelete(EtcUpstreamName)
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}
}

func toChanTimed(t *time.Timer, ch chan int) {
	t.Reset(1 * time.Second) // No defer, as we don't know which//
	//case will be selected
	select {
	case ch <- 42:
	case <-t.C:
	}
	{
		<-t.C
	}
}

func myTime() {
	// 声明一个退出用的通道
	exit := make(chan int)
	// 打印开始
	fmt.Println("start")
	// 过1秒后, 调用匿名函数
	time.AfterFunc(time.Second, func() {
		// 1秒后, 打印结果
		fmt.Println("one second after")
		// 通知main()的goroutine已经结束
		exit <- 0
	})
	// 等待结束
	<-exit
}

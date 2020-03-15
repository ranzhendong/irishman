package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/ranzhendong/irishman/pkg/etcd"
	"log"
	"time"
)

var (
	c datastruck.Config
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

// GetUpstream : method for get upstream
func GetUpstream(jsonObj interface{}, timeNow time.Time) (*MyERR.MyError, string) {
	var (
		gu  datastruck.GetUpstream
		err error
		val string
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(1012), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 1012, TimeStamp: timeNow}, ""
	}

	//judge
	if err = gu.JudgeValidator(jsonObj); err != nil {
		log.Println(MyERR.ErrorLog(1003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 1003, TimeStamp: timeNow}, ""
	}

	if gu.UpstreamName == "ALL" {
		//EtcUpstreamName := "Upstream"
		EtcUpstreamName := c.Upstream.EtcdPrefix
		//get key from etcd
		if val, _, err = etcd.EtcGetAll(EtcUpstreamName); err != nil {
			log.Println(MyERR.ErrorLog(1104), fmt.Sprintf("%v", err))
			return &MyERR.MyError{Error: err.Error(), Code: 1104, TimeStamp: timeNow}, ""
		}
		log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Get ALL Key [%v], Values [%v]", EtcUpstreamName, val)))
		return &MyERR.MyError{Code: 000, TimeStamp: timeNow}, val
	}

	EtcUpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(gu.UpstreamName)
	//get key from etcd
	if val, err = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Println(MyERR.ErrorLog(1102), fmt.Sprintf("; %v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 1102, TimeStamp: timeNow}, ""
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}, val
}

// PutUpstream : method for full update upstream
func PutUpstream(jsonObj interface{}, timeNow time.Time) *MyERR.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(2012), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 2012, TimeStamp: timeNow}
	}

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(MyERR.ErrorLog(2003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 2003, TimeStamp: timeNow}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(MyERR.ErrorLog(2004))
		return &MyERR.MyError{Error: err.Error(), Code: 2004, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(u.UpstreamName)

	//if exist
	if _, err = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(MyERR.ErrorLog(2102), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 2102, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(2101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 2101, TimeStamp: timeNow}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Put Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}
}

// PostUpstream : method for create update upstream
func PostUpstream(jsonObj interface{}, timeNow time.Time) *MyERR.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(3012), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 3012, TimeStamp: timeNow}
	}

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(MyERR.ErrorLog(3003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 3003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(u.UpstreamName)

	//if repeat
	if _, err = etcd.EtcGet(EtcUpstreamName); err == nil {
		log.Printf(MyERR.ErrorLog(3103))
		return &MyERR.MyError{Code: 3103, TimeStamp: timeNow}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(MyERR.ErrorLog(3004))
		return &MyERR.MyError{Error: err.Error(), Code: 3004, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(3101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 3101, TimeStamp: timeNow}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Set Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}
}

//PatchUpstream : method for partial upstream
func PatchUpstream(jsonObj interface{}, timeNow time.Time) *MyERR.MyError {
	var (
		pu, etcdpu       datastruck.PatchUpstream
		puData, etcdData map[string]interface{}
		UpstreamPool     []map[string]interface{}
		middleware       interface{}
		jsonU, puByte    []byte
		err              error
		val              string
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(4012), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 4012, TimeStamp: timeNow}
	}

	//judge
	if err = pu.JudgeValidator(jsonObj); err != nil {
		log.Println(MyERR.ErrorLog(4003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 4003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(pu.UpstreamName)
	//if exist
	if val, err = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(MyERR.ErrorLog(4102), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 4102, TimeStamp: timeNow}
	}

	//turn etcd data to struck, for compare and judge
	if err := json.Unmarshal([]byte(val), &etcdpu); err != nil {
		log.Printf(MyERR.ErrorLog(4005), fmt.Sprintf("Etcd PatchUpstream String: %v", err))
		return &MyERR.MyError{Code: 4005, TimeStamp: timeNow}
	}

	if pu.Pool == nil {
		etcdpu.Algorithms = pu.Algorithms
		middleware = etcdpu
		goto JUST
	}

	if err := json.Unmarshal([]byte(val), &etcdData); err != nil {
		log.Printf(MyERR.ErrorLog(4005), fmt.Sprintf("Etcd PatchUpstream Map: %v", err))
		return &MyERR.MyError{Code: 4005, TimeStamp: timeNow}
	}

	//turn struct to json string，then turn json string to map
	if puByte, err = json.Marshal(pu); err != nil {
		log.Printf(MyERR.ErrorLog(4004), fmt.Sprintf("Request PatchUpstream Struck: %v", err))
		return &MyERR.MyError{Code: 4004, TimeStamp: timeNow}
	}
	if err := json.Unmarshal(puByte, &puData); err != nil {
		log.Printf(MyERR.ErrorLog(4005), fmt.Sprintf("Request PatchUpstream Map: %v", err))
		return &MyERR.MyError{Code: 4005, TimeStamp: timeNow}
	}

	/*
		1. ipport相等
		   1.1 weight存在 （更新status和weight）
		   1.2 weight不存在 （只更新status）
		2. ipport不相等
		   1.1 weight存在 （列表增加，更新三条内容）
		   1.2 weight不存在 （放弃）
	*/
	//replace algorithms data
	if etcdData["algorithms"].(string) != puData["algorithms"].(string) {
		etcdData["algorithms"] = puData["algorithms"].(string)
	}

	//replace upstream list
	if puData["pool"] != nil {
		for i := 0; i < len(puData["pool"].([]interface{})); i++ {
			RequestIPPort := puData["pool"].([]interface{})[i].(map[string]interface{})["ipPort"]
			RequestStatus := puData["pool"].([]interface{})[i].(map[string]interface{})["status"]
			RequestWeight := puData["pool"].([]interface{})[i].(map[string]interface{})["weight"]

			for e := 0; e < len(etcdData["pool"].([]interface{})); e++ {
				etcdIPPort := etcdData["pool"].([]interface{})[e].(map[string]interface{})["ipPort"]
				etcdStatus := etcdData["pool"].([]interface{})[e].(map[string]interface{})["status"]
				etcdWeight := etcdData["pool"].([]interface{})[e].(map[string]interface{})["weight"]

				if etcdIPPort == RequestIPPort {
					if RequestWeight.(float64) == 0 {
						etcdData["pool"].([]interface{})[e].(map[string]interface{})["status"] = RequestStatus.(string)
						UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdIPPort, "status": RequestStatus, "weight": etcdWeight.(float64)})
					} else {
						etcdData["pool"].([]interface{})[e].(map[string]interface{})["status"] = RequestStatus.(string)
						etcdData["pool"].([]interface{})[e].(map[string]interface{})["weight"] = RequestWeight.(float64)
						UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdIPPort, "status": RequestStatus, "weight": RequestWeight})
					}
					goto breakFor
				} else {
					if RequestWeight.(float64) != 0 {
						UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": RequestIPPort, "status": RequestStatus, "weight": RequestWeight})
					}
				}

				// just for sure if not one server match
				//append the etcd data
				UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdIPPort, "status": etcdStatus, "weight": etcdWeight})
			}
		breakFor:
		}
	}

	//de-weight
	etcdData["pool"] = removeRepByMap(UpstreamPool)
	middleware = etcdData

JUST:
	//turn struck or map to json
	if jsonU, err = json.Marshal(middleware); err != nil {
		log.Println(MyERR.ErrorLog(4004))
		return &MyERR.MyError{Error: err.Error(), Code: 4004, TimeStamp: timeNow}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(4101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 4101, TimeStamp: timeNow}
	}

	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Patch Key [%v], Values %v", EtcUpstreamName, string(jsonU))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}
}

//DeleteUpstream : method for delete upstream or pool's server
func DeleteUpstream(jsonObj interface{}, timeNow time.Time) *MyERR.MyError {
	var (
		du, etcddu       datastruck.DeleteUpstream
		duData, etcdData map[string]interface{}
		UpstreamPool     []map[string]interface{}
		middleware       int
		duByte, jsonU    []byte
		err              error
		val              string
	)

	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(5012), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 5012, TimeStamp: timeNow}
	}

	//judge
	if err = du.JudgeValidator(jsonObj); err != nil {
		log.Println(MyERR.ErrorLog(5003), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 5003, TimeStamp: timeNow}
	}

	// Characters joining together
	EtcUpstreamName := c.Upstream.EtcdPrefix + strFirstToUpper(du.UpstreamName)
	//if exist
	if val, err = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(MyERR.ErrorLog(5102), fmt.Sprintf("%v", err))
		return &MyERR.MyError{Error: err.Error(), Code: 5102, TimeStamp: timeNow}
	}

	//if no server, mean to delete upstream
	if len(du.Pool) == 0 {
		goto DELETEUPSTREAM
	}

	//turn etcd data to struck, for compare and judge
	if err := json.Unmarshal([]byte(val), &etcddu); err != nil {
		log.Printf(MyERR.ErrorLog(5005), fmt.Sprintf("Etcd DeleteUpstream String: %v", err))
		return &MyERR.MyError{Code: 5005, TimeStamp: timeNow}
	}

	//need to least one
	if len(etcddu.Pool) <= 1 {
		middleware = 5106
		goto DELETENOTHING
	}

	if err := json.Unmarshal([]byte(val), &etcdData); err != nil {
		log.Printf(MyERR.ErrorLog(5005), fmt.Sprintf("Etcd DeleteUpstream struck: %v", err))
		return &MyERR.MyError{Code: 5005, TimeStamp: timeNow}
	}

	//turn struct to json string，then turn json string to map
	if duByte, err = json.Marshal(du); err != nil {
		log.Printf(MyERR.ErrorLog(5004), fmt.Sprintf("Request DeleteUpstream Struck: %v", err))
		return &MyERR.MyError{Code: 5004, TimeStamp: timeNow}
	}
	if err := json.Unmarshal(duByte, &duData); err != nil {
		log.Printf(MyERR.ErrorLog(5005), fmt.Sprintf("Request DeleteUpstream Map: %v", err))
		return &MyERR.MyError{Code: 5005, TimeStamp: timeNow}
	}

	//replace data, but need to last one
	for k, duv := range duData {
		if k == "pool" {
			for _, ev := range etcdData["pool"].([]interface{}) {
				for _, v := range duv.([]interface{}) {
					if v.(map[string]interface{})["ipPort"] == ev.(map[string]interface{})["ipPort"] {
						delete(ev.(map[string]interface{}), "ipPort")
						delete(ev.(map[string]interface{}), "status")
						delete(ev.(map[string]interface{}), "weight")
					}
				}
			}
		}
	}

	//filter pool map, not nil
	for _, ev := range etcdData["pool"].([]interface{}) {
		if len(ev.(map[string]interface{})) != 0 {
			UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": ev.(map[string]interface{})["ipPort"], "status": ev.(map[string]interface{})["status"], "weight": ev.(map[string]interface{})["weight"]})
		}
	}

	//can not be delete all,at least one
	if len(UpstreamPool) < 1 {
		middleware = 5107
		goto DELETENOTHING
	}

	//new pool which after delete
	etcdData["pool"] = UpstreamPool

	goto DELETESERVER

DELETESERVER:
	if jsonU, err = json.Marshal(etcdData); err != nil {
		log.Println(MyERR.ErrorLog(5004))
		return &MyERR.MyError{Error: err.Error(), Code: 5004, TimeStamp: timeNow}
	}

	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(MyERR.ErrorLog(5101, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 5101, TimeStamp: timeNow}
	}
	log.Println(MyERR.ErrorLog(000, fmt.Sprintf(" Delete Upstream Key [%v], New Values %v", EtcUpstreamName, string(jsonU))))
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}

DELETEUPSTREAM:
	if err = etcd.EtcDelete(EtcUpstreamName); err != nil {
		log.Printf(MyERR.ErrorLog(5105, fmt.Sprintf("%v", err)))
		return &MyERR.MyError{Error: err.Error(), Code: 5105, TimeStamp: timeNow}
	}
	return &MyERR.MyError{Code: 000, TimeStamp: timeNow}

DELETENOTHING:
	return &MyERR.MyError{Code: middleware, TimeStamp: timeNow}

}

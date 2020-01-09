package upstream

import (
	"datastruck"
	"encoding/json"
	"errorhandle"
	ErrH "errorhandle"
	"etcd"
	"fmt"
	"log"
	"net/http"
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
func GetUpstream(jsonObj interface{}) (*errorhandle.MyError, string) {
	var (
		gu  datastruck.GetUpstream
		err error
		val string
	)

	//judge
	if err = gu.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(1003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 1003}, ""
	}

	if gu.UpstreamName == "ALL" {
		EtcUpstreamName := "Upstream"
		//get key from etcd
		if err, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
			log.Println(ErrH.ErrorLog(1104), fmt.Sprintf("%v", err))
			return &ErrH.MyError{Error: err.Error(), Code: 1104}, ""
		}
		log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get ALL Key [%v], Values [%v]", EtcUpstreamName, val)))
		return &ErrH.MyError{Code: 000}, val
	}

	EtcUpstreamName := "Upstream" + strFirstToUpper(gu.UpstreamName)
	//get key from etcd
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Println(err)
		log.Println(ErrH.ErrorLog(1102), fmt.Sprintf("; %v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 1102}, ""
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
	return &ErrH.MyError{Code: 000}, val
}

// Full Update upstream
func PutUpstream(jsonObj interface{}) *errorhandle.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(2003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 2003}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(ErrH.ErrorLog(2004))
		return &ErrH.MyError{Error: err.Error(), Code: 2004}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//if exist
	if err, _ = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(2102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 2102}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(2101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 2101}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Put Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000}
}

// Create Update upstream
func PostUpstream(jsonObj interface{}) *errorhandle.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(3003), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 3003}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//if repeat
	if err, _ = etcd.EtcGet(EtcUpstreamName); err == nil {
		log.Printf(ErrH.ErrorLog(3103))
		return &ErrH.MyError{Code: 3103}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(ErrH.ErrorLog(3004))
		return &ErrH.MyError{Error: err.Error(), Code: 3004}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(3101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 3101}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Set Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000}
}

// Partial upstream
func PatchUpstream(jsonObj interface{}) *errorhandle.MyError {
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
		return &ErrH.MyError{Error: err.Error(), Code: 4003}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(pu.UpstreamName)
	//if exist
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(4102), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Code: 4102}
	}

	//turn etcd data to struck, for compare and judge
	if err := json.Unmarshal([]byte(val), &etcdpu); err != nil {
		log.Printf(ErrH.ErrorLog(4005), fmt.Sprintf("Etcd PatchUpstream String: %v", err))
		return &ErrH.MyError{Code: 4005}
	}

	if pu.Pool == nil {
		etcdpu.Algorithms = pu.Algorithms
		middleware = etcdpu
		goto JUST
	}

	if err := json.Unmarshal([]byte(val), &etcdData); err != nil {
		log.Printf(ErrH.ErrorLog(4005), fmt.Sprintf("Etcd PatchUpstream Map: %v", err))
		return &ErrH.MyError{Code: 4005}
	}

	//turn struct to json string，then turn json string to map
	if puByte, err = json.Marshal(pu); err != nil {
		log.Printf(ErrH.ErrorLog(4004), fmt.Sprintf("Request PatchUpstream Struck: %v", err))
		return &ErrH.MyError{Code: 4004}
	}
	if err := json.Unmarshal(puByte, &puData); err != nil {
		log.Printf(ErrH.ErrorLog(4005), fmt.Sprintf("Request PatchUpstream Map: %v", err))
		return &ErrH.MyError{Code: 4005}
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
						if k == "ipPort" && ev.(map[string]interface{})["ipPort"] == sv {
							if fv.(map[string]interface{})["weight"].(float64) == 0 {
								UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdData["pool"].([]interface{})[ek].(map[string]interface{})["ipPort"], "status": fv.(map[string]interface{})["status"], "weight": etcdData["pool"].([]interface{})[ek].(map[string]interface{})["weight"]})
								continue
							} else {
								etcdData["pool"].([]interface{})[ek].(map[string]interface{})["status"] = fv.(map[string]interface{})["status"]
								etcdData["pool"].([]interface{})[ek].(map[string]interface{})["weight"] = fv.(map[string]interface{})["weight"]
								UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": etcdData["pool"].([]interface{})[ek].(map[string]interface{})["ipPort"], "status": fv.(map[string]interface{})["status"], "weight": fv.(map[string]interface{})["weight"]})
							}
						} else if k == "ipPort" && ev.(map[string]interface{})["ipPort"] != sv {
							if fv.(map[string]interface{})["weight"].(float64) == 0 {
								continue
							} else {
								UpstreamPool = append(UpstreamPool, map[string]interface{}{"ipPort": fv.(map[string]interface{})["ipPort"], "status": fv.(map[string]interface{})["status"], "weight": fv.(map[string]interface{})["weight"]})
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
		return &ErrH.MyError{Error: err.Error(), Code: 4004}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(4101, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 4101}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Patch Key [%v], Values %v", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000}
}

// Delete upstream
func DeleteUpstream(w http.ResponseWriter, jsonObj interface{}) (err error) {
	var (
		u datastruck.Upstream
	)
	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Printf("[Upstream] JudgeValidator ERR: %v", err)
		return
	}
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)
	_ = etcd.EtcDelete(EtcUpstreamName)
	return
}

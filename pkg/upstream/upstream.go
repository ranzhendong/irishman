package upstream

import (
	"container/list"
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

// Get upstream
func GetUpstream(jsonObj interface{}) (*errorhandle.MyError, string) {
	var (
		gu  datastruck.GetUpstream
		err error
		val string
	)

	//judge
	if err = gu.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(2001), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 2001}, ""
	}

	if gu.UpstreamName == "ALL" {
		EtcUpstreamName := "Upstream"
		//get key from etcd
		if err, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
			log.Println(ErrH.ErrorLog(2005), fmt.Sprintf("%v", err))
			return &ErrH.MyError{Error: err.Error(), Code: 2005}, ""
		}
		log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get ALL Key [%v], Values [%v]", EtcUpstreamName, val)))
		return &ErrH.MyError{Code: 000}, val
	}

	EtcUpstreamName := "Upstream" + strFirstToUpper(gu.UpstreamName)
	//get key from etcd
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Println(err)
		log.Println(ErrH.ErrorLog(2006), fmt.Sprintf("; %v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 2006}, ""
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
		log.Println(ErrH.ErrorLog(3001), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 3001}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(ErrH.ErrorLog(3002))
		return &ErrH.MyError{Error: err.Error(), Code: 3002}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//if exist
	if err, _ = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(3006), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 3006}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(3003, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 3003}
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
		log.Println(ErrH.ErrorLog(4001), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 4001}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//if repeat
	if err, _ = etcd.EtcGet(EtcUpstreamName); err == nil {
		log.Printf(ErrH.ErrorLog(4004))
		return &ErrH.MyError{Code: 4004}
	}

	//turn to json
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(ErrH.ErrorLog(4002))
		return &ErrH.MyError{Error: err.Error(), Code: 4002}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(ErrH.ErrorLog(4003, fmt.Sprintf("%v", err)))
		return &ErrH.MyError{Error: err.Error(), Code: 4003}
	}

	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Set Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &ErrH.MyError{Code: 000}
}

// Partial upstream
func PatchUpstream(jsonObj interface{}) *errorhandle.MyError {
	var (
		pu     datastruck.PatchUpstream
		etcdpu datastruck.PatchUpstream
		//jsonU []byte
		puByte []byte
		err    error
		val    string
	)

	//judge
	if err = pu.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(5001), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Error: err.Error(), Code: 5001}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(pu.UpstreamName)
	//if exist
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		log.Printf(ErrH.ErrorLog(5006), fmt.Sprintf("%v", err))
		return &ErrH.MyError{Code: 5006}
	}

	if err := json.Unmarshal([]byte(val), &etcdpu); err == nil {
		log.Println("epu: ", etcdpu)
		log.Println("pu: ", pu)
	}
	//turn struct to json string，then turn json string to map
	var etcdData map[string]interface{}
	if err := json.Unmarshal([]byte(val), &etcdData); err == nil {
		log.Println("etcdData:", etcdData)
	}
	//struct 到json str
	if puByte, err = json.Marshal(pu); err == nil {
		//log.Println(string(puByte))
	}
	var data map[string]interface{}
	if err := json.Unmarshal(puByte, &data); err == nil {
		log.Println("data:", data)
	}

	l := list.New()
	//var pool []map[string]interface{}
	var service map[string]interface{}
	//pool = make([]map[string]interface{}, 5)
	service = make(map[string]interface{})
	for k, v := range data {
		log.Println("k,v :", k, v)
		if k != "pool" {
			log.Println("etcdData[k] = data[k]")
			etcdData[k] = data[k]
		}

		/*
			1. ipport相等
			   1.1 weight存在 （更新status和weight）
			   1.2 weight不存在 （只更新status）
			2. ipport不相等
			   1.1 weight存在 （列表增加，更新三条内容）
			   1.2 weight不存在 （放弃）
		*/
		if k == "pool" && v != nil {
			for ek, ev := range etcdData["pool"].([]interface{}) {
				for _, fv := range v.([]interface{}) {
					for k, sv := range fv.(map[string]interface{}) {
						log.Println("k,vvv", k, sv)
						log.Println("etcdData[pool]", ev.(map[string]interface{})["ipPort"])
						if k == "ipPort" && ev.(map[string]interface{})["ipPort"] == sv {
							log.Println(1)
							if fv.(map[string]interface{})["weight"].(float64) == 0 {
								log.Println(2)
								etcdData["pool"].([]interface{})[ek].(map[string]interface{})["status"] = fv.(map[string]interface{})["status"]
								continue
							} else {
								log.Println("fv.(map[string]interface{})[weight]", fv.(map[string]interface{})["weight"])
								log.Println(3)
								etcdData["pool"].([]interface{})[ek].(map[string]interface{})["status"] = fv.(map[string]interface{})["status"]
								etcdData["pool"].([]interface{})[ek].(map[string]interface{})["weight"] = fv.(map[string]interface{})["weight"]
							}
						} else if k == "ipPort" && ev.(map[string]interface{})["ipPort"] != sv {
							log.Println(4)
							log.Println("else if k == ipPort && ev.(map[string]interface{})[ipPort] != sv", ev.(map[string]interface{})["ipPort"], sv)
							if fv.(map[string]interface{})["weight"].(float64) == 0 {
								log.Println("fv.(map[string]interface{})[weight]", fv.(map[string]interface{})["weight"])
								log.Println(5)
								continue
							} else {
								log.Println("fv.(map[string]interface{})[weight]", fv.(map[string]interface{})["weight"])
								log.Println(6)
								service["ipPort"] = fv.(map[string]interface{})["ipPort"]
								service["status"] = fv.(map[string]interface{})["status"]
								service["weight"] = fv.(map[string]interface{})["weight"]
								log.Println("ssssssssssssssssssssssss", service)
								l.PushFront(service)
							}
							//for k, v := range pool {
							//	log.Println("ssssssssssssss", k, pool[k])
							//	if len(v) == 0 {
							//		log.Println("len(v) == 0")
							//		pool[k] = service
							//		break
							//	} else {
							//		log.Println("len(v) else")
							//		if pool[k]["ipPort"] == service["ipPort"] {
							//			break
							//		}
							//	}
							//}
							//log.Println("POOL", pool)

						}
					}
				}
			}
		}

	}
	for i := l.Front(); i != nil; i = i.Next() {
		log.Println(i.Value)
	}
	log.Println("END DATA", etcdData)

	//
	////turn to json
	//if jsonU, err = json.Marshal(u); err != nil {
	//	log.Println(ErrH.ErrorLog(4002))
	//	return &ErrH.MyError{Error: err.Error(), Code: 4002}
	//}
	//
	////etcd put
	//if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
	//	log.Printf(ErrH.ErrorLog(4003, fmt.Sprintf("%v", err)))
	//	return &ErrH.MyError{Error: err.Error(), Code: 4003}
	//}
	//
	//log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Set Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
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

package upstream

import (
	"datastruck"
	"encoding/json"
	"errorhandle"
	myErr "errorhandle"
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
func GetUpstream(w http.ResponseWriter, jsonObj interface{}) (err error, val string) {
	var (
		gu datastruck.GetUpstream
	)

	if err = gu.JudgeValidator(jsonObj); err != nil {
		log.Printf("[Upstream] JudgeValidator ERR: %v", err)
		return
	}

	if gu.UpstreamName == "ALL" {
		EtcUpstreamName := "Upstream"
		//get key from etcd
		if err, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
			return
		}
		return
	}
	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(gu.UpstreamName)

	//get key from etcd
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		return
	}

	log.Printf("[GetUpstream]: Get key {%v} Successful! Values %v ", gu.UpstreamName, val)

	return
}

// Full Update upstream, but in this
func PutUpstream(w http.ResponseWriter, jsonObj interface{}) (err error) {
	//var b []byte
	//if err, _ = GetUpstream(w, gu); err != nil {
	//	log.Printf("[PutUpstream]: Get key {%v} Failed ! It Not Exist !", u.UpstreamName)
	//	return
	//}
	//_ = PostUpstream(w, u)
	////return to user
	//
	//if b, err = json.Marshal(u); err == nil {
	//}
	//
	//_, _ = io.WriteString(w, string(b))
	return
}

// Create Update upstream
func PostUpstream(w http.ResponseWriter, jsonObj interface{}) *errorhandle.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(myErr.ErrorLog(4001))
		return &myErr.MyError{Error: err.Error(), Code: 4001}
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)
	if jsonU, err = json.Marshal(u); err != nil {
		log.Println(myErr.ErrorLog(4002))
		return &myErr.MyError{Error: err.Error(), Code: 4002}
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		log.Printf(myErr.ErrorLog(4003, fmt.Sprintf("%v", err)))
		//c := myErr.ErrorLog(4003)
		//log.Println(c)
		return &myErr.MyError{Error: err.Error(), Code: 4003}
	}

	log.Println(myErr.ErrorLog(0000, fmt.Sprintf(" Set Key [%v], Values [%v]", EtcUpstreamName, string(jsonU))))
	return &myErr.MyError{Code: 0000}
}

// Partial upstream
func PatchUpstream(w http.ResponseWriter, jsonObj interface{}) (err error) {
	return
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

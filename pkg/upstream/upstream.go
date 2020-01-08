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

// Get upstream
func GetUpstream(jsonObj interface{}) (*errorhandle.MyError, string) {
	var (
		gu  datastruck.GetUpstream
		err error
		val string
	)

	//judge
	if err = gu.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(2001))
		return &ErrH.MyError{Error: err.Error(), Code: 2001}, ""
	}

	if gu.UpstreamName == "ALL" {
		EtcUpstreamName := "Upstream"
		//get key from etcd
		if err, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
			log.Println(ErrH.ErrorLog(2006), fmt.Sprintf("%v", err))
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
func PostUpstream(jsonObj interface{}) *errorhandle.MyError {
	var (
		u     datastruck.Upstream
		jsonU []byte
		err   error
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Println(ErrH.ErrorLog(4001))
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

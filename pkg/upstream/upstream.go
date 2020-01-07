package upstream

import (
	"datastruck"
	"encoding/json"
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
func PostUpstream(w http.ResponseWriter, jsonObj interface{}) (err error) {
	var (
		u     datastruck.Upstream
		jsonU []byte
	)

	//judge
	if err = u.JudgeValidator(jsonObj); err != nil {
		log.Printf("[Upstream] JudgeValidator ERR: %v", err)
		return
	}

	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)
	if jsonU, err = json.Marshal(u); err != nil {
		log.Printf("[PostUpstream] Json datastruck.Upstream ERR: %v\n", err)
		err = fmt.Errorf("[PostUpstream] Json datastruck.Upstream ERR: %v\n", err)
		return
	}

	//etcd put
	if err = etcd.EtcPut(EtcUpstreamName, string(jsonU)); err != nil {
		return
	}

	log.Printf("[PostUpstream]: Set to etcd Successful!  Key [ %v ], Values [%v] ", EtcUpstreamName, string(jsonU))

	return
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

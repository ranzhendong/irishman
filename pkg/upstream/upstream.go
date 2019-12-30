package upstream

import (
	"datastruck"
	"encoding/json"
	"etcd"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
func GetUpstream(w http.ResponseWriter, u datastruck.Upstream) (err error) {
	var val string
	// Characters joining together
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)

	//get key from etcd
	if err, val = etcd.EtcGet(EtcUpstreamName); err != nil {
		return
	}

	log.Printf("[GetUpstream]: Get key {%v} Successful! Values %v ", u.UpstreamName, val)

	//return to user
	_, _ = io.WriteString(w, val)
	return
}

// Full Update upstream, but in this
func PutUpstream(w http.ResponseWriter, u datastruck.Upstream) (err error) {

	return
}

// Create Update upstream
func PostUpstream(w http.ResponseWriter, u datastruck.Upstream) (err error) {
	var (
		jsonU []byte
	)

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
func PatchUpstream(w http.ResponseWriter, u datastruck.Upstream) (err error) {
	return
}

// Delete upstream
func DeleteUpstream(w http.ResponseWriter, u datastruck.Upstream) (err error) {
	return
}

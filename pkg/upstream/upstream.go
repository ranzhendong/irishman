package upstream

import (
	"datastruck"
	"encoding/json"
	"etcd"
	"fmt"
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
func GetUpstream(r *http.Request, u datastruck.Upstream) {

}

// Full Update upstream
func PutUpstream(r *http.Request, u datastruck.Upstream) {

}

// Create Update upstream
func PostUpstream(r *http.Request, u datastruck.Upstream) (err error) {
	var (
		jsonU []byte
	)
	EtcUpstreamName := "Upstream" + strFirstToUpper(u.UpstreamName)
	log.Println(EtcUpstreamName)
	if jsonU, err = json.Marshal(u); err != nil {
		log.Printf("[PostUpstream] Json datastruck.Upstream ERR: %v\n", err)
		err = fmt.Errorf("[PostUpstream] Json datastruck.Upstream ERR: %v\n", err)
		return
	}
	log.Printf("[PostUpstream] The Request Body: %v", string(jsonU))
	_ = etcd.EtcPut(EtcUpstreamName, string(jsonU))
	return
}

// Partial upstream
func PatchUpstream(r *http.Request, u datastruck.Upstream) {

}

// Delete upstream
func DeleteUpstream(r *http.Request, u datastruck.Upstream) {

}

package kvnuts

import (
	"strings"
	"time"
)

//upper the first letter
func strFirstToLower(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 65 && strArry[0] <= 90 {
		strArry[0] += 32
	}
	return string(strArry)
}

//SetFlagHC : set flag healthCheck to nutsDB, restart goroutine hc
func SetFlagHC() {
	//log.Println("SetFlagHCSetFlagHCSetFlagHCSetFlagHC")
	for {
		time.Sleep(50 * time.Millisecond)
		//log.Println("2222222222222222")
		if _, _, err := Get("FlagUpstreamNutsDB", "FlagUpstreamNutsDBFinishUpstream", "i"); err == nil {
			//log.Println("time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)time.Sleep(100 * time.Millisecond)")
			_ = Del("FlagUpstreamNutsDB", "FlagUpstreamNutsDBFinishUpstream")
			goto BREAKFOR
		}
	}
BREAKFOR:
	_ = Put("FlagHC", "FlagHC", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream and health check to nutsDB
func SetFlagUpstreamAndHCNutsDB(watcherFlag, ectKey string) {
	//log.Println("SetFlagNutsDBSetFlagNutsDBSetFlagNutsDB", strFirstToLower(strings.Split(ectKey, watcherFlag)[1]))
	if watcherFlag == "Upstream" {
		_ = Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDB", 1)
		_ = Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDBWatcherFlag", ectKey)
		_ = Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDBStartUpstream", 1)
	} else {
		_ = Put("FlagHCNutsDB", "FlagHCNutsDB", 1)
		_ = Put("FlagHCNutsDB", "FlagHCNutsDBWatcherFlag", strFirstToLower(strings.Split(ectKey, watcherFlag)[1]))
	}
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagHCNutsDB() {
	//log.Println("FlagHCNutsDBFlagHCNutsDBFlagHCNutsDBFlagHCNutsDB")
	_ = Put("FlagHCNutsDB", "FlagHCNutsDB", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagUpstreamReadyTo() {
	//log.Println("SetFlagUpstreamReadyToSetFlagUpstreamReadyToSetFlagUpstreamReadyTo")
	_ = Put("SetFlagUpstreamReadyTo", "SetFlagUpstreamReadyTo", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagHealthCheckReadyTo() {
	//log.Println("SetFlagHealthCheckReadyToSetFlagHealthCheckReadyToSetFlagHealthCheckReadyTo")
	_ = Put("SetFlagHealthCheckReadyTo", "SetFlagHealthCheckReadyTo", 1)
}

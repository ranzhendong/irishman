package kvnuts

import (
	"log"
	"strings"
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

//SetFlagHC : set flag healthCheck to nutsDB
func SetFlagHC() {
	log.Println("SetFlagHCSetFlagHCSetFlagHCSetFlagHC")
	_ = Put("FlagHC", "FlagHC", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagUpstreamNutsDB(watcherFlag, ectKey string) {
	log.Println("SetFlagNutsDBSetFlagNutsDBSetFlagNutsDB", strFirstToLower(strings.Split(ectKey, watcherFlag)[1]))
	_ = Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDB", 1)
	_ = Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDBWatcherFlag", ectKey)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagHCNutsDB() {
	log.Println("FlagHCNutsDBFlagHCNutsDBFlagHCNutsDBFlagHCNutsDB")
	_ = Put("FlagHCNutsDB", "FlagHCNutsDB", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagUpstreamReadyTo() {
	log.Println("SetFlagUpstreamReadyToSetFlagUpstreamReadyToSetFlagUpstreamReadyTo")
	_ = Put("SetFlagUpstreamReadyTo", "SetFlagUpstreamReadyTo", 1)
}

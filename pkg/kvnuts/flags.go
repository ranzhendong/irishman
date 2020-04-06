package kvnuts

import "log"

//SetFlagHC : set flag healthCheck to nutsDB
func SetFlagHC() {
	log.Println("SetFlagHCSetFlagHCSetFlagHCSetFlagHC")
	_ = Put("FlagHC", "FlagHC", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagUpstreamNutsDB() {
	log.Println("SetFlagNutsDBSetFlagNutsDBSetFlagNutsDB")
	_ = Put("FlagUpstreamNutsDB", "FlagUpstreamNutsDB", 1)
}

//SetFlagUpstreamNutsDB : set flag Upstream to nutsDB
func SetFlagHCNutsDB() {
	log.Println("FlagHCNutsDBFlagHCNutsDBFlagHCNutsDBFlagHCNutsDB")
	_ = Put("FlagHCNutsDB", "FlagHCNutsDB", 1)
}

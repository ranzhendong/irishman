package healthcheck

import (
	ErrH "errorhandle"
	"time"
)

func GetHealthCheck(jsonObj interface{}, timeNow time.Time) (*ErrH.MyError, string) {
	//log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, "ss"
}

//func PutHealthCheck(jsonObj interface{}, timeNow time.Time) *ErrH.MyError {
//	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
//	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
//}
//
//func PatchHealthCheck(jsonObj interface{}, timeNow time.Time) *ErrH.MyError {
//	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get Key [%v], Values [%v]", EtcUpstreamName, val)))
//	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
//}

//func InitHealthCheck() {
//	log.Println("InitHealthCheck")
//
//	var (
//		err error
//		val string
//	)
//
//	EtcUpstreamName := "Upstream"
//	//get key from etcd
//	if err, val = etcd.EtcGetAll(EtcUpstreamName); err != nil {
//		log.Println(ErrH.ErrorLog(1104), fmt.Sprintf("%v", err))
//		return &ErrH.MyError{Error: err.Error(), Code: 1104, TimeStamp: timeNow}, ""
//	}
//	log.Println(ErrH.ErrorLog(000, fmt.Sprintf(" Get ALL Key [%v], Values [%v]", EtcUpstreamName, val)))
//	return &ErrH.MyError{Code: 000, TimeStamp: timeNow}, val
//}

package errorhandle

import (
	"github.com/thinkeridea/go-extend/exstrings"
	"strconv"
	"time"
)

type MyError struct {
	Error        string
	Message      string
	Code         int
	TimeStamp    time.Time
	ExecutorTime string
}

var (
	mux        = make(map[int]string)
	muxS       = make(map[int]string)
	randSlice  = make([]int, 3)
	sRandSlice = make([]int, 2)
)

//registered
/*
000 successful

1-9 method

001-030 system error

140-150 system status

101 - 200 etcd error


*/
func init() {
	muxS[0] = "ServeHTTP: "
	muxS[1] = "Upstream GET: "
	muxS[2] = "Upstream PUT: "
	muxS[3] = "Upstream POST: "
	muxS[4] = "Upstream PATCH: "
	muxS[5] = "Upstream DELETE: "
	muxS[6] = "Viper Watcher: "
	muxS[7] = "HealthCheck GET: "
	muxS[8] = "HealthCheck PUT: "
	muxS[9] = "HealthCheck PATCH: "
	muxS[10] = "HealthCheck DELETE: "
	muxS[11] = "HealthCheck Goroutines: "
	muxS[12] = "NutsDB: "

	mux[000] = "Successful"
	mux[001] = "Upstream: "
	mux[002] = "INIT: Loading Request Body Failed"
	mux[003] = "JudgeValidator Error"
	mux[004] = "Json: Marshal Error"
	mux[005] = "Json: UNMarshal Error"
	mux[006] = "WriteString Error"
	mux[007] = "Not Support Method Error"
	mux[010] = "Url Not Exist"
	mux[011] = "HTTP Server Init Error"
	mux[012] = "Config Json: UNMarshal Error"

	mux[101] = "Etcd Put: Put Key Error"
	mux[102] = "Etcd Get: Key Not Exist Error"
	mux[103] = "Etcd Get: Repeat Key Error"
	mux[104] = "Etcd GetALL: No Key Error"
	mux[105] = "Etcd Delete: Error"
	mux[106] = "Etcd Delete: Etcd Key's Pool Has One ServerList At Least, Delete Canceled !"
	mux[107] = "Etcd Delete: Etcd Key's Pool Has One ServerList At Least, Can Not Delete Them ALL !"

	mux[140] = "Config Change Reloading"
	mux[141] = "IrishMan Is Running With Execute Path"

	mux[151] = "HealthCheck Config Initialize"
	mux[152] = "SuccessStatus Has One Code At Least "
	mux[153] = "FailuresStatus Has One Code At Least "

	mux[161] = "Connect Error"
	mux[162] = "Put Error"
	mux[163] = "Get Error"
}

//register error to message
func (self *MyError) Messages() {
	defer func() {
		_ = recover()
		if self.Message == "" {
			self.Message = "No Error Match"
		} else if self.Error == "" {
			self.Error = self.Message
		} else if self.Error == "" && self.Message == "" {
			self.Error = "No Error Match"
			self.Message = "No Error Match"
		}
	}()
	self.Message = muxS[SCode(self.Code)] + mux[Code(self.Code)]
}

//error log handler
func ErrorLog(code int, content ...string) string {
	if content == nil {
		return muxS[SCode(code)] + mux[Code(code)]
	}
	return muxS[SCode(code)] + mux[Code(code)] + content[0]
}

//timer clock
func (self *MyError) Clock() {
	//if TimeStamp is none
	if len(time.Since(self.TimeStamp).String()) > 20 {
		self.ExecutorTime = time.Since(time.Now()).String()
		return
	}
	self.ExecutorTime = time.Since(self.TimeStamp).String()
}

//code cut out
func Code(e int) (a int) {
	randSlice[0] = e / 100 % 10
	randSlice[1] = e / 10 % 10
	randSlice[2] = e / 1 % 10
	a, _ = strconv.Atoi(exstrings.JoinInts(randSlice, ""))
	return
}

func SCode(e int) (a int) {
	if len(strconv.Itoa(e)) == 4 {
		return e / 1000 % 10
	} else {
		sRandSlice[0] = e / 10000 % 10
		sRandSlice[1] = e / 1000 % 10
	}
	a, _ = strconv.Atoi(exstrings.JoinInts(sRandSlice, ""))
	return
}

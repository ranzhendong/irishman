package errorhandle

import (
	"github.com/thinkeridea/go-extend/exstrings"
	"strconv"
)

type MyError struct {
	Error   string
	Message string
	Code    int
}

var (
	mux       = make(map[int]string)
	muxS      = make(map[int]string)
	randSlice = make([]int, 3)
)

//registered
/*
000 successful

1-9 method

001-030 system

101 - 200 etcd

*/
func init() {
	muxS[1] = "Upstream GET: "
	muxS[2] = "Upstream PUT: "
	muxS[3] = "Upstream POST: "
	muxS[4] = "Upstream PATCH: "
	muxS[5] = "Upstream DELETE: "

	mux[000] = "Successful"
	muxS[001] = "Upstream: "
	mux[002] = "INIT: Loading Body Failed"
	mux[003] = "JudgeValidator Error"
	mux[004] = "Json: Marshal Error"
	mux[005] = "Json: UNMarshal Error"

	mux[101] = "Etcd Put: Put Key Error"
	mux[102] = "Etcd Get: Key Not Exist Error"
	mux[103] = "Etcd Get: Repeat Key Error"
	mux[104] = "Etcd GetALL: No Key Error"
}

//register error to message
func (e *MyError) Messages() {
	defer func() {
		_ = recover()
		if e.Message == "" {
			e.Message = "No Error Match"
		} else if e.Error == "" {
			e.Error = e.Message
		} else if e.Error == "" && e.Message == "" {
			e.Error = "No Error Match"
			e.Message = "No Error Match"
		}
	}()
	e.Message = muxS[e.Code/1000%10] + mux[Code(e.Code)]
}

//error log handler
func ErrorLog(code int, content ...string) string {
	if content == nil {
		return muxS[code/1000%10] + mux[Code(code)]
	}
	return muxS[code/1000%10] + mux[Code(code)] + content[0]
}

func Code(e int) (a int) {
	randSlice[0] = e / 100 % 10
	randSlice[1] = e / 10 % 10
	randSlice[2] = e / 1 % 10
	a, _ = strconv.Atoi(exstrings.JoinInts(randSlice, ""))
	return
}

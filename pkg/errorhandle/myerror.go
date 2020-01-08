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
func init() {
	muxS[2] = "Upstream GET: "
	muxS[3] = "Upstream PUT: "
	muxS[4] = "Upstream POST: "
	mux[000] = "Successful"
	mux[001] = "JudgeValidator Error"
	mux[002] = "Json: Marshal Error"
	mux[003] = "Etcd Put: Put Key Error"
	mux[004] = "Etcd Get: Repeat Key Error"
	mux[005] = "Etcd GetALL: No Key Error"
	mux[006] = "Etcd Get: Key Not Exist Error"
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

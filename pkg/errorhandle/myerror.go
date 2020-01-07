package errorhandle

type MyError struct {
	Error   string
	Message string
	Code    int
}

var (
	mux  = make(map[int]string)
	muxS = make(map[int]string)
)

//registered
func init() {
	muxS[4] = "[Post Upstream] "
	mux[0000] = "Successful"
	mux[4001] = "JudgeValidator"
	mux[4002] = "Json Marshal Error"
	mux[4003] = "Etcd Put Error"
}

//handle the error
func (e *MyError) Messages() {
	e.Message = muxS[e.Code/1000%10] + mux[e.Code]
}

func ErrorLog(code int, content ...string) string {
	defer func() string {
		_ = recover()
		return muxS[code/1000%10] + mux[code]
	}()
	return muxS[code/1000%10] + mux[code] + content[0]
}

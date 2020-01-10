package main

import (
	"encoding/json"
	ErrH "errorhandle"
	"fmt"
	myInit "init"
	"io"
	"log"
	"net/http"
	"time"
	"upstream"
)

//定义map来实现路由转发
var (
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	err error
	n   int
)

type myHandler struct{}

//初始化log函数
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	mux["/upstream"] = myUpstream
}

func main() {
	//configure read
	if err = myInit.Config(); err != nil {
		log.Printf(ErrH.ErrorLog(6142, fmt.Sprintf("%v", err)))
		return
	}

	// server start
	server := http.Server{
		Addr:         ":8080",
		Handler:      &myHandler{},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	if err = server.ListenAndServe(); err != nil {
		log.Printf(ErrH.ErrorLog(0011, fmt.Sprintf("%v", err)))
	}
}

//route handler
func (myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	log.Printf(ErrH.ErrorLog(0010, fmt.Sprintf("%v", r.URL.String())))
	res := &ErrH.MyError{Code: 0010, Error: fmt.Sprintf("%v", r.URL.String())}
	response(w, res)
}

func myUpstream(w http.ResponseWriter, r *http.Request) {
	var (
		jsonObj interface{}
		timeNow = time.Now()
	)

	//loading request body
	if err, jsonObj = myInit.InitializeBody(r.Body); err != nil {
		log.Printf(ErrH.ErrorLog(0002, fmt.Sprintf("%v", err)))
		response(w, &ErrH.MyError{Error: err.Error(), Code: 0002})
		return
	}

	//restful switch
	switch r.Method {
	case "GET":
		if res, val := upstream.GetUpstream(jsonObj, timeNow); res != nil {
			response(w, res, val)
		}

	case "PUT":
		if res := upstream.PutUpstream(jsonObj, timeNow); res != nil {
			response(w, res)
		}

	case "POST":
		if res := upstream.PostUpstream(jsonObj, timeNow); res != nil {
			response(w, res)
		}

	case "PATCH":
		if res := upstream.PatchUpstream(jsonObj, timeNow); res != nil {
			response(w, res)
		}

	case "DELETE":
		if res := upstream.DeleteUpstream(jsonObj, timeNow); res != nil {
			response(w, res)
		}

	default:
		log.Printf(ErrH.ErrorLog(0007), fmt.Sprintf("%v", r.Method))
	}
}

func response(w http.ResponseWriter, res *ErrH.MyError, val ...string) {
	var set int
	defer func() {
		_ = recover()
		if set == 1 {
			return
		}
		res.Messages()
		res.Clock()
		if b, err := json.Marshal(res); err != nil {
			log.Println("Response Json Marshal Error:", err)
		} else {
			n, err = io.WriteString(w, string(b))
			if err != nil {
				log.Printf(ErrH.ErrorLog(0006), fmt.Sprintf("%v", err))
			}
		}
	}()
	if len(val[0]) != 0 {
		set = 1
		n, err = io.WriteString(w, val[0])
		if err != nil {
			log.Printf(ErrH.ErrorLog(0006), fmt.Sprintf("%v", err))
		}
		return
	}
}

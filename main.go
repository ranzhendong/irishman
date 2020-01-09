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
	mux map[string]func(http.ResponseWriter, *http.Request)
	err error
)

type myHandler struct{}

//初始化log函数
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	//configure read
	if err = myInit.Config(); err != nil {
		log.Printf("[MAIN] Init Config filed ! ERR: %v ", err)
		return
	}

	// server start
	server := http.Server{
		Addr:         ":8080",
		Handler:      &myHandler{},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	route(mux)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// route
func route(mux map[string]func(http.ResponseWriter, *http.Request)) {
	//upstream
	mux["/upstream"] = myUpstream
}

//route handler
func (myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	_, _ = io.WriteString(w, "[ServeHTTP] URL:"+r.URL.String()+"IS NOT EXIST")
}

func myUpstream(w http.ResponseWriter, r *http.Request) {
	var (
		jsonObj interface{}
	)

	//loading request body
	if err, jsonObj = myInit.InitializeBody(r.Body); err != nil {
		log.Printf(ErrH.ErrorLog(1002, fmt.Sprintf("%v", err)))
		response(w, &ErrH.MyError{Error: err.Error(), Code: 1002})
		return
	}

	//restful switch
	switch r.Method {
	case "GET":
		if res, val := upstream.GetUpstream(jsonObj); res != nil {
			response(w, res, val)
		}

	case "PUT":
		if res := upstream.PutUpstream(jsonObj); res != nil {
			response(w, res)
		}

	case "POST":
		if res := upstream.PostUpstream(jsonObj); res != nil {
			response(w, res)
		}

	case "PATCH":
		if res := upstream.PatchUpstream(jsonObj); res != nil {
			response(w, res)
		}

	case "DELETE":
		//_ = myUpstream.DeleteUpstream(w, u)
		log.Println("MY DELETE")

	default:
		log.Printf("[ServeHTTP Upstream] Not Support %v", r.Method)
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
		if b, err := json.Marshal(res); err != nil {
			log.Println("Response Json Marshal Error:", err)
		} else {
			_, _ = io.WriteString(w, string(b))
		}
	}()
	if len(val[0]) != 0 {
		set = 1
		_, _ = io.WriteString(w, val[0])
		return
	}
}

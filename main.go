package main

import (
	"encoding/json"
	"errorhandle"
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
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
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
		log.Printf("[Upstream] Can Not Loading body %v", r.Body)
		return
	}

	//restful switch
	switch r.Method {
	case "GET":

		_, val := upstream.GetUpstream(w, jsonObj)
		//return to user
		_, _ = io.WriteString(w, val)
		log.Println("MY GET")

	case "PUT":
		//_ = myUpstream.PutUpstream(w, u)
		//if b, err = json.Marshal(u); err == nil {
		//	_, _ = io.WriteString(w, string(b))
		//}
		log.Println("MY PUT")

	case "POST":

		if res := upstream.PostUpstream(w, jsonObj); res != nil {
			response(w, res)
		}
		log.Println("MY POST")

	case "PATCH":
		//_ = myUpstream.PatchUpstream(w, u)
		log.Println("MY PATCH")

	case "DELETE":
		//_ = myUpstream.DeleteUpstream(w, u)
		log.Println("MY DELETE")

	default:
		log.Printf("[ServeHTTP Upstream] Not Support %v", r.Method)
	}

}

func response(w http.ResponseWriter, res *errorhandle.MyError) {

	res.Messages()

	if b, err := json.Marshal(res); err != nil {
		fmt.Println(err)

	} else {
		_, _ = io.WriteString(w, string(b))
	}
}

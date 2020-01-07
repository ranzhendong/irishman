package main

import (
	myInit "init"
	"io"
	"log"
	"net/http"
	"time"
	myUpstream "upstream"
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
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
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
	mux["/upstream"] = upstream
}

//route handler
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	_, _ = io.WriteString(w, "[ServeHTTP] URL:"+r.URL.String()+"IS NOT EXIST")
}

func upstream(w http.ResponseWriter, r *http.Request) {
	var (
		//b       []byte
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

		_, val := myUpstream.GetUpstream(w, jsonObj)
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

		_ = myUpstream.PostUpstream(w, jsonObj)

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

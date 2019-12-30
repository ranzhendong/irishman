package main

import (
	"datastruck"
	"encoding/json"
	"fmt"
	"govalidators"
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

	if err = myInit.Config(); err != nil {
		log.Printf("[MAIN] Init Config filed ! ERR: %v ", err)
		return
	}

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

// 路由
func route(mux map[string]func(http.ResponseWriter, *http.Request)) {
	//镜像更新
	mux["/lua"] = lua
	mux["/upstream"] = upstream
}

//路由的转发
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		//用这个handler实现路由转发，相应的路由调用相应func
		h(w, r)
		return
	}
	//log.Println("[ServeHTTP] URL:" + r.URL.String() + "IS NOT EXIST")
	//log.Println(r.Method)
	//log.Println(r.Header)
	//var body []byte
	//if body, err = ioutil.ReadAll(r.Body); err != nil {
	//	log.Printf("[InitCheck] Read Body ERR: %v\n", err)
	//	err = fmt.Errorf("[InitCheck] Read Body ERR: %v\n", err)
	//	return
	//}
	//log.Println(string(body))
	//log.Println(r)
	_, _ = io.WriteString(w, "[ServeHTTP] URL:"+r.URL.String()+"IS NOT EXIST")
}

func lua(w http.ResponseWriter, r *http.Request) {

	fmt.Println("lua")
}

func upstream(w http.ResponseWriter, r *http.Request) {

	//loading request body
	var (
		u datastruck.Upstream
		b []byte
	)
	if err, u = myInit.InitializeBody(r.Body); err != nil {
		log.Printf("[Upstream] Can Not Loading body %v", u)
		return
	}

	switch r.Method {
	case "GET":
		_, val := myUpstream.GetUpstream(w, u)
		//return to user
		_, _ = io.WriteString(w, val)
		log.Println("MY GET")
	case "PUT":
		_ = myUpstream.PutUpstream(w, u)
		if b, err = json.Marshal(u); err == nil {
			_, _ = io.WriteString(w, string(b))
		}
		log.Println("MY PUT")
	case "POST":
		validator := govalidators.New()
		if err := validator.Validate(u); err != nil {
			fmt.Println(err)
			return
		}
		_ = myUpstream.PostUpstream(w, u)
		if b, err = json.Marshal(u); err == nil {
			_, _ = io.WriteString(w, string(b))
		}
		log.Println("MY POST")
	case "PATCH":
		_ = myUpstream.PatchUpstream(w, u)
		log.Println("MY PATCH")
	case "DELETE":
		_ = myUpstream.DeleteUpstream(w, u)
		log.Println("MY DELETE")
	default:
		log.Printf("[ServeHTTP Upstream] Not Support %v", r.Method)
	}
}

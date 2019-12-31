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
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	myUpstream "upstream"
)

//定义map来实现路由转发
var (
	mux map[string]func(http.ResponseWriter, *http.Request)
	err error
)

type myHandler struct{}

//type ValidatorF func(params map[string]interface{}, val reflect.Value, args ...string) (bool, error)
//type Validator interface {
//	Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error)
//}

type IpPortValidator struct {
	EMsg string
}

func (self *IpPortValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {
	const (
		IP = "^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)$"
	)

	defer func() {
		_ = recover()
		if err != nil {
			log.Printf("[Validate]: IP or Port Is Valid")
			err = fmt.Errorf("[Validate]: IP or Port Is Valid")
		}
	}()

	sep := ":"
	arr := strings.Split(val.String(), sep)

	if !regexp.MustCompile(IP).MatchString(arr[0]) {
		err = fmt.Errorf("IP illegal")
		return false, err
	}

	a, err := strconv.Atoi(arr[1])
	if int(1024) >= a || int(65535) <= a {
		err = fmt.Errorf("PORT illegal")
		return false, err
	}

	return true, nil
}

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

		validator.SetValidators(map[string]interface{}{
			"ipPort": &IpPortValidator{},
		})

		if err := validator.Validate(u); err != nil {
			log.Println(err)
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

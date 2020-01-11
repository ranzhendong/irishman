package main

import (
	"datastruck"
	"encoding/json"
	ErrH "errorhandle"
	"fmt"
	"healthcheck"
	myInit "init"
	"io"
	"log"
	"net/http"
	"time"
	"upstream"
)

var (
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	c   datastruck.Config
	err error
	n   int
)

type myHandler struct{}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	mux["/upstream"] = myUpstream
	mux["/healthcheck"] = healthCheck
}

func main() {
	//configure read
	if err = myInit.Config(); err != nil {
		log.Printf(ErrH.ErrorLog(6142, fmt.Sprintf("%v", err)))
		return
	}

	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(0012), fmt.Sprintf("%v", err))
		return
	}

	//initialize health check
	healthcheck.InitHealthCheck(time.Now())

	// server start
	server := http.Server{
		Addr:         c.Server.Bind,
		Handler:      &myHandler{},
		ReadTimeout:  time.Duration(c.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Server.WriteTimeout) * time.Second,
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

func healthCheck(w http.ResponseWriter, r *http.Request) {
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
		if res, val := healthcheck.GetHealthCheck(jsonObj, timeNow); res != nil {
			response(w, res, val)
		}

	//case "PUT":
	//	if res := healthcheck.PutHealthCheck(jsonObj, timeNow); res != nil {
	//		response(w, res)
	//	}
	//
	//case "PATCH":
	//	if res := healthcheck.PatchHealthCheck(jsonObj, timeNow); res != nil {
	//		response(w, res)
	//	}

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

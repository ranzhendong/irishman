package main

import (
	"datastruck"
	"encoding/json"
	ErrH "errorhandle"
	"fmt"
	"healthcheck"
	myInit "init"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"upstream"
)

var (
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	c   datastruck.Config
	err error
)

type myHandler struct{}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//configure read
	if err = myInit.Config(); err != nil {
		log.Printf(ErrH.ErrorLog(6142, fmt.Sprintf("%v", err)))
		return
	}
	mux["/upstream"] = myUpstream
	mux["/healthcheck"] = healthCheck
}

func main() {
	//config loading
	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(0012), fmt.Sprintf("%v", err))
		return
	}

	//remove nutsDB
	if err = os.RemoveAll(c.NutsDB.Path); err != nil {
		log.Println(ErrH.ErrorLog(12123), fmt.Sprintf("; %v", err))
	}

	//judge if remove nutsDB successful
	if files, err := ioutil.ReadDir(c.NutsDB.Path); err == nil {
		var f os.FileInfo
		for _, f = range files {
			log.Println(ErrH.ErrorLog(12124), fmt.Sprintf(";The file:%v", f.Name()))
		}
		return
	}

	//initialize health check
	go healthcheck.InitHealthCheck(time.Now())

	//config about server
	server := http.Server{
		Addr:         c.Server.Bind,
		Handler:      &myHandler{},
		ReadTimeout:  time.Duration(c.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Server.WriteTimeout) * time.Second,
	}

	// server start
	log.Println(ErrH.ErrorLog(142))
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

	case "PUT":
		if res := healthcheck.PutHealthCheck(jsonObj, timeNow); res != nil {
			response(w, res)
		}

	case "PATCH":
		if res := healthcheck.PatchHealthCheck(jsonObj, timeNow); res != nil {
			response(w, res)
		}

	case "DELETE":
		if res := healthcheck.DeleteHealthCheck(jsonObj, timeNow); res != nil {
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
			_, err = io.WriteString(w, string(b))
			if err != nil {
				log.Printf(ErrH.ErrorLog(0006), fmt.Sprintf("%v", err))
			}
		}
	}()
	if len(val[0]) != 0 {
		set = 1
		_, err = io.WriteString(w, val[0])
		if err != nil {
			log.Printf(ErrH.ErrorLog(0006), fmt.Sprintf("%v", err))
		}
		return
	}
}

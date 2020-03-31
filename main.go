package main

import (
	"encoding/json"
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	gc "github.com/ranzhendong/irishman/pkg/gorountinescontroller"
	"github.com/ranzhendong/irishman/pkg/healthcheck"
	MyInit "github.com/ranzhendong/irishman/pkg/init"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"github.com/ranzhendong/irishman/pkg/upstream"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
	if err = MyInit.Config(); err != nil {
		log.Printf(MyERR.ErrorLog(6142, fmt.Sprintf("%v", err)))
		return
	}

	//set route
	mux["/upstream"] = myUpstream
	mux["/healthcheck"] = healthCheck
	mux["/nutsdb"] = nutsDB
}

func main() {
	//config loading
	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(0012), fmt.Sprintf("%v", err))
		return
	}

	//remove nutsDB
	if err = os.RemoveAll(c.NutsDB.Path); err != nil {
		log.Println(MyERR.ErrorLog(12123), fmt.Sprintf("; %v", err))
	}

	//judge if remove nutsDB successful
	if files, err := ioutil.ReadDir(c.NutsDB.Path); err == nil {
		var f os.FileInfo
		for _, f = range files {
			log.Println(MyERR.ErrorLog(12124), fmt.Sprintf(";The file:%v", f.Name()))
		}
		return
	}

	//goroutines controller: hc, etcd watcher
	gc.Factory()

	//config about server
	server := http.Server{
		Addr:         c.Server.Bind,
		Handler:      &myHandler{},
		ReadTimeout:  time.Duration(c.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Server.WriteTimeout) * time.Second,
	}

	// server start
	log.Println(MyERR.ErrorLog(142))
	if err = server.ListenAndServe(); err != nil {
		log.Printf(MyERR.ErrorLog(0011, fmt.Sprintf("%v", err)))
	}
}

//route handler
func (myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	log.Printf(MyERR.ErrorLog(0010, fmt.Sprintf("%v", r.URL.String())))
	res := &MyERR.MyError{Code: 0010, Error: fmt.Sprintf("%v", r.URL.String())}
	response(w, res)
}

func nutsDB(w http.ResponseWriter, r *http.Request) {
	log.Println("nutsDB.....")
	_ = kvnuts.Put("FalgHC", "FalgHC", 1)
}

func myUpstream(w http.ResponseWriter, r *http.Request) {
	var (
		rs upstream.RStruck
	)

	//set timestamp
	rs.T = time.Now()

	//loading request body
	if rs.J, err = MyInit.InitializeBody(r.Body); err != nil {
		log.Printf(MyERR.ErrorLog(0002, fmt.Sprintf("%v", err)))
		response(w, &MyERR.MyError{Error: err.Error(), Code: 0002})
		return
	}

	//restful switch method
	switch r.Method {
	case "GET":
		if res, val := rs.GetUpstream(); res != nil {
			response(w, res, val)
		}

	case "PUT":
		if res := rs.PutUpstream(); res != nil {
			response(w, res)
		}

	case "POST":
		if res := rs.PostUpstream(); res != nil {
			response(w, res)
		}

	case "PATCH":
		if res := rs.PatchUpstream(); res != nil {
			response(w, res)
		}

	case "DELETE":
		if res := rs.DeleteUpstream(); res != nil {
			response(w, res)
		}

	default:
		log.Printf(MyERR.ErrorLog(0007), fmt.Sprintf("%v", r.Method))
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	var (
		rhos healthcheck.RHCStruck
	)

	//set timestamp
	rhos.T = time.Now()

	//loading request body
	if rhos.J, err = MyInit.InitializeBody(r.Body); err != nil {
		log.Printf(MyERR.ErrorLog(0002, fmt.Sprintf("%v", err)))
		response(w, &MyERR.MyError{Error: err.Error(), Code: 0002})
		return
	}

	//restful switch method
	switch r.Method {
	case "GET":
		if res, val := rhos.GetHealthCheck(); res != nil {
			response(w, res, val)
		}

	case "PUT":
		if res := rhos.PutHealthCheck(); res != nil {
			response(w, res)
		}

	case "PATCH":
		if res := rhos.PatchHealthCheck(); res != nil {
			response(w, res)
		}

	case "DELETE":
		if res := rhos.DeleteHealthCheck(); res != nil {
			response(w, res)
		}

	default:
		log.Printf(MyERR.ErrorLog(0007), fmt.Sprintf("%v", r.Method))
	}

}

func response(w http.ResponseWriter, res *MyERR.MyError, val ...string) {
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
				log.Printf(MyERR.ErrorLog(0006), fmt.Sprintf("%v", err))
			}
		}
	}()
	if len(val[0]) != 0 {
		set = 1
		_, err = io.WriteString(w, val[0])
		if err != nil {
			log.Printf(MyERR.ErrorLog(0006), fmt.Sprintf("%v", err))
		}
		return
	}
}

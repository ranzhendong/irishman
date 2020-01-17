package healthcheck

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

func TCP(ip string, pingTimeout int) bool {
	var (
		conn net.Conn
		err  error
	)
	if conn, err = net.DialTimeout("tcp", ip, time.Duration(pingTimeout)*time.Millisecond); err != nil {
		return false
	}
	conn.Close()
	return true
}

func HTTP(request string, t int) (err error, code int) {
	var (
		myRequest  *http.Request
		myResponse *http.Response
	)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Millisecond * time.Duration(t),
	}

	//request
	request = "http://" + request
	//log.Println(request, t)
	myRequest, _ = http.NewRequest("PUT", request, nil)

	//set my request
	if myResponse, err = client.Do(myRequest); err != nil {
		log.Println(err)
		code = 504
		return
	}

	defer func() {
		_ = recover()
		if err == nil {
			myResponse.Body.Close()
		}
	}()

	return nil, myResponse.StatusCode
}

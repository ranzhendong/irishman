package healthcheck

import (
	"crypto/tls"
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

func HTTP(request string, t int) (error, int) {
	var (
		err        error
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
	//log.Println(request)
	request = "http://" + request
	myRequest, _ = http.NewRequest("PUT", request, nil)

	//if response exist
	if myResponse, err = client.Do(myRequest); err != nil {
		return err, 0
	}

	defer myResponse.Body.Close()

	//log.Println(myResponse.StatusCode)
	return nil, myResponse.StatusCode
}

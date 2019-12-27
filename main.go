package main

import (
	"etcd"
	"initconfig"
	"log"
	"time"
)

type Test struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Info     Info   `json:"info"`
}

type Info struct {
	Address string `json:"address"`
}

//初始化log函数
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	var (
		err error
	)

	if err = initconfig.ConfigAnalysis(); err != nil {
		return
	}

	if err = etcd.EtcPut("info", "zhendong"); err != nil {
		return
	}

	if err = etcd.EtcGet("info"); err != nil {
		return
	}

	if err = etcd.EtcWatcher(); err != nil {
		return
	}

	var count int
	count = 1
	for {
		//log.Println("检查次数:[", count, "]")
		count = count + 1
		time.Sleep(time.Duration(2) * time.Second)
	}
}

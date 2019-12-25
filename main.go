package main

import (
	"consul"
	"encoding/json"
	"log"
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
		t   Test
	)

	// set key to consul
	t = Test{
		"zhendong",
		"0922",
		Info{
			"hangzhou",
		},
	}

	v, err := json.Marshal(t)
	//fmt.Println(string(v))
	if err = consul.SetKey("info", v); err != nil {
		return
	}

	// get key to consul
	if err = consul.GetKey("info"); err != nil {
		return
	}

}

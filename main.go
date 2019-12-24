package main

import (
	"consul"
	"encoding/json"
	"fmt"
	"log"
)

type Test struct {
	Name string `json:"name"`
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
	t.Name = "zhendong"
	k, err := json.Marshal(t)
	fmt.Println(string(k))
	if err = consul.SetKey("info", k); err != nil {
		log.Println(err)
	}

	// get key to consul
	if err = consul.GetKey("info"); err != nil {
		log.Println(err)
	}

}

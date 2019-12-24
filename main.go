package main

import (
	"consul"
	"datastruck"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"initconfig"
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
		c   datastruck.Config
	)

	if err = initconfig.Config(); err != nil {
		log.Println(err)
	}

	if err = viper.Unmarshal(&c); err != nil {
		log.Printf("[User] Unable To Decode Into Config Struct, %v", err)
		err = fmt.Errorf("[User] Unable To Decode Into Config Struct, %v", err)
		return
	}

	fmt.Println(c)

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

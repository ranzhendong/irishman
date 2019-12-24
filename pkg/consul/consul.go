package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
)

func GetKey(key string) (err error) {
	if err = viper.AddRemoteProvider("consul", "172.16.0.51:8500", key); err != nil {
		log.Println(err)
		return
	}
	viper.SetConfigType("json")
	if err = viper.ReadRemoteConfig(); err != nil {
		log.Println(err)
		return
	}
	log.Println(viper.Get("port"))
	log.Println(viper.Get("name"))
	return
}

func SetKey(key string, value []byte) (err error) {

	// Get a new client
	config := api.DefaultConfig()
	config.Address = "172.16.0.51:8500"
	client, err := api.NewClient(config)
	if err != nil {
		return
	}

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	if _, err = kv.Put(&api.KVPair{Key: key, Value: value}, nil); err != nil {
		return
	}
	return
}

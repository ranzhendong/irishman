package consul

import (
	"datastruck"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"initconfig"
	"log"
)

var (
	c datastruck.Config
)

func GetKey(key string) (err error) {

	if err, c = initconfig.ConfigAnalysis(); err != nil {
		return
	}

	if err = viper.AddRemoteProvider("consul", c.Consul.Host, key); err != nil {
		log.Printf("[GetKey] Can Not Connect Consul, ERR: %v", err)
		return
	}

	viper.SetConfigType("json")
	if err = viper.ReadRemoteConfig(); err != nil {
		log.Printf("[GetKey] Can Not Read ConsulKey, ERR: %v", err)
		return
	}

	log.Printf("[GetKey] Get Key:%v, Key Gather:%v", key, viper.AllKeys())

	return
}

func SetKey(key string, value []byte) (err error) {

	if err, c = initconfig.ConfigAnalysis(); err != nil {
		return
	}

	// Get a new client
	config := api.DefaultConfig()
	config.Address = c.Consul.Host
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("[SetKey] Can Not Connect Consul, ERR: %v", err)
		return
	}

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	if _, err = kv.Put(&api.KVPair{Key: key, Value: value}, nil); err != nil {
		log.Printf("[SetKey] Can Not Set Key to Consul, ERR: %v", err)
		return
	}

	log.Printf("[SetKey] Set Key:%v, Value:%v", key, string(value))
	return
}

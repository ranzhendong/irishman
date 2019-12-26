package consul

import (
	"datastruck"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"initconfig"
	"log"
	"time"
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

func ConsulWatcher() {
	var err error

	if err, c = initconfig.ConfigAnalysis(); err != nil {
		return
	}

	// alternatively, you can create a new viper instance.
	var runtimeViper = viper.New()

	if err = runtimeViper.AddRemoteProvider("consul", c.Consul.Host, "info"); err != nil {
		log.Printf("[GetKey] Can Not Connect Consul, ERR: %v", err)
		return
	}

	runtimeViper.SetConfigType("json") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	err = runtimeViper.ReadRemoteConfig()

	fmt.Println(runtimeViper.AllKeys())
	//// unmarshal config
	//runtimeViper.Unmarshal(&runtime_conf)

	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtimeViper.WatchRemoteConfig()
			if err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}

			// unmarshal new config into our runtime config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			//runtimeViper.Unmarshal(&runtime_conf)

			fmt.Println(runtimeViper.AllKeys())
		}
	}()

}

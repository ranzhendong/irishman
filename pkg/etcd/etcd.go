package etcd

import (
	"context"
	"datastruck"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"time"
)

var (
	c datastruck.Config
)

//func GetKey(key string) (err error) {
//
//	if err, c = initconfig.ConfigAnalysis(); err != nil {
//		return
//	}
//
//	if err = viper.AddRemoteProvider("consul", c.Consul.Host, key); err != nil {
//		log.Printf("[GetKey] Can Not Connect Consul, ERR: %v", err)
//		return
//	}
//
//	viper.SetConfigType("json")
//	if err = viper.ReadRemoteConfig(); err != nil {
//		log.Printf("[GetKey] Can Not Read ConsulKey, ERR: %v", err)
//		return
//	}
//
//	log.Printf("[GetKey] Get Key:%v, Key Gather:%v", key, viper.AllKeys())
//
//	return
//}
//
//func SetKey(key string, value []byte) (err error) {
//
//	if err, c = initconfig.ConfigAnalysis(); err != nil {
//		return
//	}
//
//	// Get a new client
//	config := api.DefaultConfig()
//	config.Address = c.Consul.Host
//	client, err := api.NewClient(config)
//	if err != nil {
//		log.Printf("[SetKey] Can Not Connect Consul, ERR: %v", err)
//		return
//	}
//
//	// Get a handle to the KV API
//	kv := client.KV()
//
//	// PUT a new KV pair
//	if _, err = kv.Put(&api.KVPair{Key: key, Value: value}, nil); err != nil {
//		log.Printf("[SetKey] Can Not Set Key to Consul, ERR: %v", err)
//		return
//	}
//
//	log.Printf("[SetKey] Set Key:%v, Value:%v", key, string(value))
//	return
//}
//
//func ConsulWatcher() {
//	var err error
//
//	if err, c = initconfig.ConfigAnalysis(); err != nil {
//		return
//	}
//
//	// alternatively, you can create a new viper instance.
//	var runtimeViper = viper.New()
//
//	if err = runtimeViper.AddRemoteProvider("consul", c.Consul.Host, "info"); err != nil {
//		log.Printf("[GetKey] Can Not Connect Consul, ERR: %v", err)
//		return
//	}
//
//	runtimeViper.SetConfigType("json") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
//
//	// read from remote config the first time.
//	err = runtimeViper.ReadRemoteConfig()
//
//	fmt.Println(runtimeViper.AllKeys())
//	//// unmarshal config
//	//runtimeViper.Unmarshal(&runtime_conf)
//
//	// open a goroutine to watch remote changes forever
//	go func() {
//		for {
//			time.Sleep(time.Second * 5) // delay after each request
//
//			// currently, only tested with etcd support
//			err := runtimeViper.WatchRemoteConfig()
//			if err != nil {
//				log.Printf("unable to read remote config: %v", err)
//				continue
//			}
//
//			// unmarshal new config into our runtime config struct. you can also use channel
//			// to implement a signal to notify the system of the changes
//			//runtimeViper.Unmarshal(&runtime_conf)
//
//			fmt.Println(runtimeViper.AllKeys())
//		}
//	}()
//
//}

func EtcGet(key string) (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		//getResp *clientv3.GetResponse
		getOp  clientv3.Op
		opResp clientv3.OpResponse
	)

	// Unmarshal to struck
	if err = viper.Unmarshal(&c); err != nil {
		log.Printf("[Main] Unable To Decode Into Config Struct, %v", err)
		return
	}

	//set config
	config = clientv3.Config{
		Endpoints:   []string{c.Etcd.Host},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	// 创建Op
	getOp = clientv3.OpGet(key)

	// 执行Op
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("数据Revision", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value", string(opResp.Get().Kvs[0].Value))

	return
}

func EtcPut(key, val string) (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
	)

	// Unmarshal to struck
	if err = viper.Unmarshal(&c); err != nil {
		log.Printf("[Main] Unable To Decode Into Config Struct, %v", err)
		return
	}

	//set config
	config = clientv3.Config{
		Endpoints:   []string{c.Etcd.Host},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	// 写入
	if _, err = kv.Put(context.TODO(), key, val); err != nil {
		fmt.Println(err)
		return
	}

	return
}

func EtcWatcher() (err error) {

	var (
		config             clientv3.Config
		client             *clientv3.Client
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
	)

	// Unmarshal to struck
	if err = viper.Unmarshal(&c); err != nil {
		log.Printf("[Main] Unable To Decode Into Config Struct, %v", err)
		return
	}

	//set config
	config = clientv3.Config{
		Endpoints:   []string{c.Etcd.Host},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个监听器
	watcher = clientv3.NewWatcher(client)
	// 启动监听 5秒后关闭

	ctx, _ := context.WithCancel(context.TODO())
	//time.AfterFunc(10*time.Second, func() {
	//	cancelFunc()
	//})
	watchRespChan = watcher.Watch(ctx, "name", clientv3.WithRev(watchStartRevision))

	go func() { // 处理kv变化事件
		for {

			for watchResp = range watchRespChan {
				for _, event = range watchResp.Events {
					switch event.Type {
					case mvccpb.PUT:
						fmt.Println("key patch", string(event.Kv.Value))
					case mvccpb.DELETE:
						fmt.Println("key delete", string(event.Kv.Key))
					}
				}
			}
		}
	}()

	return
}

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

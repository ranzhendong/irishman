package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"log"

	//github.com/spf13/viper/remote need to be here
	_ "github.com/spf13/viper/remote"
	"time"
)

//etcConnect : connect etcd
func etcConnect() (client *clientv3.Client, err error) {
	var (
		c         datastruck.Config
		config    clientv3.Config
		statusRes *clientv3.StatusResponse
	)

	//config loading
	if err = c.Config(); err != nil {
		return
	}

	//set config
	config = clientv3.Config{
		Endpoints:   []string{c.Etcd.Host},
		DialTimeout: time.Duration(c.Etcd.Timeout) * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}
	//timeout control
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Etcd.Timeout)*time.Second)
	defer cancel()
	statusRes, err = client.Status(timeoutCtx, c.Etcd.Host)
	if err != nil || statusRes == nil {
		return
	}
	return
}

//EtcGet : get key
func EtcGet(key string) (val string, err error) {
	var (
		client *clientv3.Client
		kv     clientv3.KV
		getOp  clientv3.Op
		opResp clientv3.OpResponse
	)

	if client, err = etcConnect(); err != nil {
		err = fmt.Errorf(" Etcd Client Initialize Failed")
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	// 创建Op
	getOp = clientv3.OpGet(key)

	// 执行Op
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		err = fmt.Errorf("KV.DO Failed")
		return
	}

	defer func() {
		_ = recover()
		if val == "" {
			err = fmt.Errorf("No Key [%v]", key)
		}
	}()

	val = string(opResp.Get().Kvs[0].Value)
	return
}

//EtcGetAll : get prefix key
func EtcGetAll(key string) (val string, rVal []*mvccpb.KeyValue, err error) {
	var (
		client  *clientv3.Client
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)

	if client, err = etcConnect(); err != nil {
		err = fmt.Errorf(" Etcd Client Initialize Failed")
		return
	}

	kv = clientv3.NewKV(client)

	if getResp, err = kv.Get(context.TODO(), key, clientv3.WithPrefix()); err != nil {
		err = fmt.Errorf("KV.DO Failed")
		return
	}

	for _, v := range getResp.Kvs {
		val = val + string(v.Value) + "\n\n"
	}
	rVal = getResp.Kvs
	return
}

//EtcPut : put key
func EtcPut(key, val string) (err error) {
	var (
		client *clientv3.Client
		kv     clientv3.KV
	)
	if client, err = etcConnect(); err != nil {
		err = fmt.Errorf(" Etcd Client Initialize Failed")
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	// 写入
	if _, err = kv.Put(context.TODO(), key, val); err != nil {
		err = fmt.Errorf("KV.DO Failed")
		return
	}

	return
}

//EtcDelete : delete key
func EtcDelete(key string) (err error) {
	var (
		client *clientv3.Client
		kv     clientv3.KV
	)
	if client, err = etcConnect(); err != nil {
		err = fmt.Errorf(" Etcd Client Initialize Failed")
		return
	}

	kv = clientv3.NewKV(client)

	if _, err = kv.Delete(context.TODO(), key); err != nil {
		err = fmt.Errorf("KV.DO Failed")
		return
	}

	return
}

//EtcWatcher : watcher key if change
func EtcWatcher(key string) (err error) {

	var (
		client             *clientv3.Client
		watchStartRevision int64
		watcher            clientv3.Watcher
	)
	if client, err = etcConnect(); err != nil {
		err = fmt.Errorf(" Etcd Client Initialize Failed")
		return
	}

	// 创建一个监听器
	watcher = clientv3.NewWatcher(client)

	ctxRoot := context.WithValue(context.Background(), "watcherFlag", key)
	watchRespChan := watcher.Watch(ctxRoot, key, clientv3.WithPrefix(), clientv3.WithRev(watchStartRevision))
	log.Println("EtcWatcher KEYS", key)
	go Watcher(ctxRoot, watchRespChan)

	return
}

//Watcher : goroutines if watcher
func Watcher(ctx context.Context, watchRespChan <-chan clientv3.WatchResponse) {

	var (
		watchResp clientv3.WatchResponse
		event     *clientv3.Event
	)

	for {
		for watchResp = range watchRespChan {
			for _, event = range watchResp.Events {
				switch event.Type {

				//be triggered when method is put, post, patch
				case mvccpb.PUT:
					log.Println("EtcWatcher PUT", string(event.Kv.Key), string(event.Kv.Value))
					//set flag SetFlagNutsDB, nutsDB watcher is triggered
					kvnuts.SetFlagUpstreamNutsDB(ctx.Value("watcherFlag").(interface{}).(string), string(event.Kv.Key))

				//be triggered when method is delete
				case mvccpb.DELETE:
					log.Println("EtcWatcher DELETE", string(event.Kv.Key))

				default:
					log.Println("EtcWatcher DEFAULT", string(event.Kv.Key))
				}
			}
		}
	}
}

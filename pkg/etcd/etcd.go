package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	_ "github.com/spf13/viper/remote"
	"time"
)

//etcd connect function
func etcConnect() (err error, client *clientv3.Client) {
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

func EtcGet(key string) (err error, val string) {
	var (
		client *clientv3.Client
		kv     clientv3.KV
		getOp  clientv3.Op
		opResp clientv3.OpResponse
	)

	if err, client = etcConnect(); err != nil {
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

func EtcGetAll(key string) (err error, val string, rVal []*mvccpb.KeyValue) {
	var (
		client  *clientv3.Client
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)

	if err, client = etcConnect(); err != nil {
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

func EtcPut(key, val string) (err error) {
	var (
		client *clientv3.Client
		kv     clientv3.KV
	)
	if err, client = etcConnect(); err != nil {
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

func EtcDelete(key string) (err error) {
	var (
		client *clientv3.Client
		kv     clientv3.KV
	)
	if err, client = etcConnect(); err != nil {
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

func EtcWatcher() (err error) {

	var (
		client             *clientv3.Client
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
	)
	if err, client = etcConnect(); err != nil {
		err = fmt.Errorf(" Etcd Client Initialize Failed")
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

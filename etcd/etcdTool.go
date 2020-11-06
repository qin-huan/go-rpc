package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
)

func Connect() (*clientv3.Client) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatalln("error.connect: ", err.Error())
	}
	log.Println("succeed to connect to etcd.")
	return cli
}

func Get(cli *clientv3.Client, key string) ([]*mvccpb.KeyValue ,error) {
	res, err := cli.Get(context.Background(), key, clientv3.WithPrefix())

	if err != nil {
		log.Fatalf("error.etcdGet: ", err.Error())
		return nil, err
	}
	// log.Println("kvs length is : ", len(res.Kvs))
	return res.Kvs, err
}

func Put(cli *clientv3.Client, key string, val string) error {
	_, err := cli.Put(context.Background(), key, val)
	if err != nil {
		log.Fatalln("error.etcdPut : ", err.Error())
		return err
	}
	return nil
}

func Watch(cli *clientv3.Client, key string)  {
	log.Println("watching...")
	watcher := clientv3.NewWatcher(cli)
	ch := watcher.Watch(context.Background(), key, clientv3.WithPrefix())
	log.Println(ch)
}
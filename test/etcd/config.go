package main

import (
	"context"
	"fmt"
	"ngo2/pkgs/Etcd"
)

func main() {

	config := Etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}
	cli, err := Etcd.NewClient(&config)
	if err != nil {
		panic(err.Error())
		return
	}

	_, err = cli.Put(context.Background(), "test", "test123")
	if err != nil {
		panic(err.Error())
		return
	}

	get, err := cli.Get(context.Background(), "test")
	if err != nil {
		panic(err.Error())
		return
	}

	fmt.Printf("key: %v, value: %v", string(get.Kvs[0].Key), string(get.Kvs[0].Value))

}

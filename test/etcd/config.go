package main

import (
	"context"
	"fmt"
	"github.com/kevinu2/ngo2/pkgs/Etcd"
	"ngo.new/pkgs/Log"
)

func main() {

	config := Etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}
	cli, err := Etcd.NewClient(&config)
	if err != nil {
		Log.Logger().Panic(err.Error())
		return
	}

	_, err = cli.Put(context.Background(), "test", "test123")
	if err != nil {
		Log.Logger().Panic(err.Error())
		return
	}

	get, err := cli.Get(context.Background(), "test")
	if err != nil {
		Log.Logger().Panic(err.Error())
		return
	}

	fmt.Println(get)

}

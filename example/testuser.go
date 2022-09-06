package main

import (
	"context"
	"films/utils"
	"fmt"
	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	jsoniter "github.com/json-iterator/go"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
)

func main() {
	req := map[string]string{"auth": "9503611fd03cd9235b1a24bf3837c1f9"}
	//req := map[string]interface{}{"mobile": "18111111111", "type": 1, "code": "1111"}
	var rsp map[string][]byte //[respone]=> json 二进制
	//var rs map[string]interface{}
	var rs string

	conf := utils.Instance()
	reg := consul.NewRegistry(registry.Addrs(conf.Get("hosts", "registry-consul", "address").String(""))) //连接在线服务注册中心
	cl := grpc.NewClient(client.Registry(reg))
	request := cl.NewRequest("user", "UserService.Authorize", req, client.WithContentType("application/json"))

	if err := cl.Call(context.TODO(), request, &rsp); err != nil {

		fmt.Println("ss", err)
	} else {

		jsoniter.Unmarshal(rsp["respone"], &rs)
		var rrs map[string]interface{}
		jsoniter.Unmarshal([]byte(rs), &rrs)
		fmt.Println(rrs)
	}
}

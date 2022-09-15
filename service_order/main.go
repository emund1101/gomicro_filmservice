package main

import (
	"films/service_order/handler"
	pb "films/service_order/proto"
	"films/utils"
	plugin_micro "films/utils/go2sky_micro"
	"films/utils/orm"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
	"gorm.io/gorm"
)

var db *gorm.DB
var conf config.Config
var rdb *redis.Client

//初始化前
func init() {
	conf = utils.Instance()
	//初始化链路跟踪
	utils.StartTracer(conf.Get("hosts", "reporter", "address").String(""), conf.Get("services", "service_order", "name").String(""))
	db = orm.Initconf("database_order", conf.Get("services", "service_order", "name").String(""))
	rdb = orm.InitRedis("redis")

}

func main() {

	reg := consul.NewRegistry(registry.Addrs(conf.Get("hosts", "registry-consul", "address").String("")))
	srv := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name(conf.Get("services", "service_order", "name").String("")),
		micro.Version(conf.Get("services", "service_order", "version").String("")),
		micro.Registry(reg),
	)
	srv.Init(
		micro.WrapHandler(plugin_micro.NewHandlerWrapper(utils.GetTracer(), conf.Get("services", "service_order", "name").String(""))),
	)

	odersrv := new(handler.OrderService)
	odersrv.Db = db
	odersrv.Redis = rdb

	pb.RegisterOrderServiceHandler(srv.Server(), odersrv)
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"films/service_user/handler"
	pb "films/service_user/proto"
	"films/utils"
	plugin_micro "films/utils/go2sky_micro"
	"films/utils/orm"
	"fmt"
	"github.com/SkyAPM/go2sky"
	consul "github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client
var conf config.Config
var tracer go2sky.Reporter

//初始化前
func init() {
	conf = utils.Instance()
	//初始化链路跟踪
	utils.StartTracer(conf.Get("hosts", "reporter", "address").String(""), conf.Get("services", "service_user", "name").String(""))
	db = orm.Initconf("database_user", conf.Get("services", "service_user", "name").String(""))
	rdb = orm.InitRedis("redis")
}

func main() {

	//注册服务
	reg := consul.NewRegistry(registry.Addrs(conf.Get("hosts", "registry-consul", "address").String("")))
	srv := micro.NewService(
		micro.Server(grpc.NewServer()), //使用grpc作为服务器
		micro.Name(conf.Get("services", "service_user", "name").String("")),
		micro.Version(conf.Get("services", "service_user", "version").String("")),
		//micro.Address(":8002"), //对外访问端口地址
		micro.Registry(reg), //注册服务发现
		//micro.WrapHandler()//注册包装器,熔断，限流等的中间件

	)

	srv.Init(
		micro.WrapHandler(plugin_micro.NewHandlerWrapper(utils.GetTracer(), conf.Get("services", "service_user", "name").String(""))),
	)

	user := new(handler.UserService)
	user.Db = db   //设置gorm
	user.Rdb = rdb //设置redis
	// Register handler  注册rpc的服务
	pb.RegisterUserServiceHandler(srv.Server(), user)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"films/service_film/handler"
	pb "films/service_film/proto"
	"films/utils"
	plugin_micro "films/utils/go2sky_micro"
	"films/utils/orm"
	grpc_client "github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/config"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"gorm.io/gorm"
)

var db *gorm.DB
var conf config.Config
var rdb *redis.Client
var cl client.Client

func init() {
	conf = utils.Instance()
	//初始化链路跟踪
	utils.StartTracer(conf.Get("hosts", "reporter", "address").String(""), conf.Get("services", "service_film", "name").String(""))
	db = orm.Initconf("database_films", conf.Get("services", "service_film", "name").String(""))
	rdb = orm.InitRedis("redis")

}

func main() {

	reg := consul.NewRegistry(registry.Addrs(conf.Get("hosts", "registry-consul", "address").String("")))
	cl = grpc_client.NewClient(client.Registry(reg))
	srv := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name(conf.Get("services", "service_film", "name").String("")),
		micro.Version(conf.Get("services", "service_film", "version").String("")),
		micro.Metadata(map[string]string{"protocol": "tcp/grpc", "header": ""}), //定义
		micro.Registry(reg),
		micro.Client(cl),
		//micro.Config(conf),
	)

	//初始化
	srv.Init(
		micro.WrapHandler(
			//wservice.NewHandlerWrapper(srv),
			plugin_micro.NewHandlerWrapper(utils.GetTracer(), conf.Get("services", "service_user", "name").String("")),
		), //将服务上下文传入到处理器使用.
	//	micro.WrapClient(plugin_micro.NewClientWrapper(utils.GetTracer(), plugin_micro.WithClientWrapperReportTags(conf.Get("services", "service_user", "name").String("")))),
	)

	// Register handler  注册rpc的服务
	filmsrv := new(handler.FilmService)
	filmsrv.Db = db
	filmsrv.Rdb = rdb
	filmsrv.GRPClient = cl

	pb.RegisterFilmServiceHandler(srv.Server(), filmsrv)
	// 运行服务
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"films/utils"
	plugin_micro "films/utils/go2sky_micro"
	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	fh_router "github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
	"strings"
	"time"
	"unicode"
)

var cl client.Client
var conf config.Config

//grpc网关
func main() {

	conf = utils.Instance()
	utils.StartTracer(conf.Get("hosts", "reporter", "address").String(""), "gateway")
	reg := consul.NewRegistry(registry.Addrs(conf.Get("hosts", "registry-consul", "address").String(""))) //连接在线服务注册中心
	cl = grpc.NewClient(client.Registry(reg))

	router := fh_router.New()
	InitRouter(router, reg) //注册路由器
	default_http := conf.Get("api", "host").String("") + ":" + conf.Get("api", "port").String("8899")
	fasthttp.ListenAndServe(default_http, router.Handler)
}

//分发路由器
func InitRouter(router *fh_router.Router, reg registry.Registry) {
	//获取consul
	routers_post := conf.Get("routers", "post").StringSlice([]string{})
	routers_get := conf.Get("routers", "get").StringSlice([]string{})

	list, err := reg.ListServices()
	if err == nil {
		if len(routers_post) > 0 {
			for _, path := range routers_post {
				for _, item := range list {
					addRouter("POST", path, item, router)
				}
			}
		}
		if len(routers_get) > 0 {
			for _, path := range routers_get {
				for _, item := range list {
					addRouter("GET", path, item, router)
				}
			}
		}

	}
}

func addRouter(rtype string, url string, list *registry.Service, router *fh_router.Router) {
	if strings.Contains(url, "/"+list.Name+"/") {
		switch rtype {
		case "POST":
			router.POST(url, fasthttpdeal)
		case "GET":
			router.GET(url, fasthttpdeal)
		}
	}
}

//转发请求grpc服务
func fasthttpdeal(ctx *fasthttp.RequestCtx) {
	var rsp map[string][]byte
	var req map[string]interface{}
	path := string(ctx.Path())

	url := lo.Substring(path, 1, uint(len(path)-1))
	path_slice := strings.Split(url, "/")
	a := []rune(path_slice[0])
	a[0] = unicode.ToUpper(a[0]) //将第一字符转大写
	//context_todo := context.TODO()
	context_bg, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := jsoniter.Unmarshal(ctx.Request.Body(), &req); err != nil { //json转map
		//使用传统模式获取参数
		ctx.Request.PostArgs().VisitAll(func(key, value []byte) {
			req[string(key)] = string(value) //转为map格式
		})
		//加上自定义的
	}

	request := cl.NewRequest(path_slice[0], string(a)+"Service."+path_slice[1], req, client.WithContentType("application/json"))
	//注入链路context
	plugin_micro.NewGateCall(utils.GetTracer(), context_bg, request, path)
	nctx := utils.GetContext()

	if err := cl.Call(nctx, request, &rsp); err == nil {
		ctx.Write(rsp["respone"]) //输出结果
	} else {
		ctx.WriteString(err.Error())
	}

}

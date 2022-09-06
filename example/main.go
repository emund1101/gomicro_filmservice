package main

import (
	fh_router "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fh_router.New()
	router.GET("/", fasthttp_deal)
	fasthttp.ListenAndServe("0.0.0.0:8899", router.Handler)
}

func fasthttp_deal(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("你好")

}

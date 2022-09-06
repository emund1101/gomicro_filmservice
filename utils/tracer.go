package utils

import (
	"context"
	"fmt"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"os"
)

var tracer *go2sky.Tracer
var ctx context.Context

func StartTracer(serviceAddr, serviceName string) go2sky.Reporter {
	re, err := reporter.NewGRPCReporter(serviceAddr)
	if err != nil {
		fmt.Println("创建gosky reporter失败", err)
		os.Exit(0)
	}
	//	defer re.Close()

	//初始化tracer
	tracer, err = go2sky.NewTracer("service_"+serviceName, go2sky.WithReporter(re))
	if err != nil {
		fmt.Println("tracer 失败", err)
	}
	return re
}

func GetTracer() *go2sky.Tracer {
	return tracer
}

func GetContext() context.Context {
	return ctx
}

func SetTContext(ctx2 context.Context) {
	ctx = ctx2
}

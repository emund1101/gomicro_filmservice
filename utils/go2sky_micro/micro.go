package go2sky_micro

import (
	"context"
	"errors"
	"films/utils"
	"fmt"
	"github.com/SkyAPM/go2sky"
	jsoniter "github.com/json-iterator/go"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strings"
	"time"
)

const (
	componentIDGoMicroClient = 5008
	componentIDGoMicroServer = 5009
	componentIDGOHttpClient  = 5005
)

var errTracerIsNil = errors.New("tracer is nil")

type clientWrapper struct {
	client.Client

	sw         *go2sky.Tracer
	reportTags []string
}

// ClientOption allow optional configuration of Client
type ClientOption func(*clientWrapper)

// WithClientWrapperReportTags customize span tags
func WithClientWrapperReportTags(reportTags ...string) ClientOption {
	return func(c *clientWrapper) {
		c.reportTags = append(c.reportTags, reportTags...)
	}
}

// Call is used for client calls
func (s *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	span, err := s.sw.CreateExitSpan(ctx, name, req.Service(), func(key, value string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[key] = value
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return err
	}

	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

	defer span.End()

	msg, _ := jsoniter.Marshal(req.Body())
	span.Tag("请求内容", string(msg))
	span.Tag("请求类型", req.ContentType())

	for _, k := range s.reportTags {
		if v, ok := metadata.Get(ctx, k); ok {
			span.Tag(go2sky.Tag(k), v)
		}
	}
	if err = s.Client.Call(ctx, req, rsp, opts...); err != nil {
		span.Error(time.Now(), err.Error())
	}
	return err
}

// Stream is used streaming
func (s *clientWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	span, err := s.sw.CreateExitSpan(ctx, name, req.Service(), func(key, value string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[key] = value
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return nil, err
	}

	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

	defer span.End()
	for _, k := range s.reportTags {
		if v, ok := metadata.Get(ctx, k); ok {
			span.Tag(go2sky.Tag(k), v)
		}
	}
	stream, err := s.Client.Stream(ctx, req, opts...)
	if err != nil {
		span.Error(time.Now(), err.Error())
	}
	return stream, err
}

// Publish is used publish message to subscriber
func (s *clientWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	name := fmt.Sprintf("Pub to %s", p.Topic())
	span, err := s.sw.CreateExitSpan(ctx, name, p.ContentType(), func(key, value string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[key] = value
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return err
	}

	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

	defer span.End()
	for _, k := range s.reportTags {
		if v, ok := metadata.Get(ctx, k); ok {
			span.Tag(go2sky.Tag(k), v)
		}
	}
	if err = s.Client.Publish(ctx, p, opts...); err != nil {
		span.Error(time.Now(), err.Error())
	}
	return err
}

// NewClientWrapper accepts a go2sky Tracer and returns a Client Wrapper
func NewClientWrapper(sw *go2sky.Tracer, options ...ClientOption) client.Wrapper {
	return func(c client.Client) client.Client {
		co := &clientWrapper{
			sw:     sw,
			Client: c,
		}
		for _, option := range options {
			option(co)
		}
		return co
	}
}

// NewCallWrapper accepts an go2sky Tracer and returns a Call Wrapper
func NewCallWrapper(sw *go2sky.Tracer, reportTags ...string) client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			if sw == nil {
				return errTracerIsNil
			}

			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			span, err := sw.CreateExitSpan(ctx, name, req.Service(), func(key, value string) error {
				mda, _ := metadata.FromContext(ctx)
				md := metadata.Copy(mda)
				md[key] = value
				ctx = metadata.NewContext(ctx, md)
				return nil
			})
			if err != nil {
				return err
			}

			span.SetComponent(componentIDGoMicroClient)
			span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

			defer span.End()
			for _, k := range reportTags {
				if v, ok := metadata.Get(ctx, k); ok {
					span.Tag(go2sky.Tag(k), v)
				}
			}
			if err = cf(ctx, node, req, rsp, opts); err != nil {
				span.Error(time.Now(), err.Error())
			}
			return err
		}
	}
}

// NewSubscriberWrapper accepts a go2sky Tracer and returns a Handler Wrapper
func NewSubscriberWrapper(sw *go2sky.Tracer, reportTags ...string) server.SubscriberWrapper {
	return func(next server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			if sw == nil {
				return errTracerIsNil
			}

			name := "Sub from " + msg.Topic()
			span, err := sw.CreateExitSpan(ctx, name, msg.ContentType(), func(key, value string) error {
				mda, _ := metadata.FromContext(ctx)
				md := metadata.Copy(mda)
				md[key] = value
				ctx = metadata.NewContext(ctx, md)
				return nil
			})
			if err != nil {
				return err
			}

			span.SetComponent(componentIDGoMicroClient)
			span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

			defer span.End()
			for _, k := range reportTags {
				if v, ok := metadata.Get(ctx, k); ok {
					span.Tag(go2sky.Tag(k), v)
				}
			}
			if err = next(ctx, msg); err != nil {
				span.Error(time.Now(), err.Error())
			}
			return err
		}
	}
}

// NewHandlerWrapper accepts a go2sky Tracer and returns a Subscriber Wrapper
//拦截器微服务触发的方法
func NewHandlerWrapper(sw *go2sky.Tracer, reportTags ...string) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			if sw == nil {
				return errTracerIsNil
			}

			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			span, ctx, err := sw.CreateEntrySpan(ctx, name, func(key string) (string, error) {
				str, _ := metadata.Get(ctx, strings.Title(key))

				return str, nil
			})
			if err != nil {
				return err
			}

			span.SetComponent(componentIDGoMicroServer)
			span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

			defer span.End()
			for _, k := range reportTags {
				if v, ok := metadata.Get(ctx, k); ok {
					span.Tag(go2sky.Tag(k), v)
				}
			}
			utils.SetTContext(ctx)
			utils.SetSpanName(name)
			if err = fn(ctx, req, rsp); err != nil {
				span.Error(time.Now(), err.Error())
			}
			//记录数据
				 
			span.Tag("响应内容", fmt.Sprintf("%s",rsp))

			return err
		}
	}
}

func NewGateCall(sw *go2sky.Tracer, ctx context.Context, req client.Request, path string) error {
	//	span, ctx, err := sw.CreateEntrySpan(ctx, path, func(key string) (string, error) {
	//		str, _ := metadata.Get(ctx, strings.Title(key))
	//		return str, nil
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	span.SetComponent(componentIDGOHttpClient)
	//	span.SetSpanLayer(agentv3.SpanLayer_Http)
	//	defer span.End()

	//name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	span, err := sw.CreateExitSpan(ctx, path, req.Service(), func(key, value string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[key] = value
		ctx = metadata.NewContext(ctx, md)

		return nil
	})
	if err != nil {
		return err
	}

	span.SetComponent(componentIDGOHttpClient)
	//span1.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(agentv3.SpanLayer_Http)
	utils.SetTContext(ctx)

	defer span.End()
	//记录request 数据

	msg, _ := jsoniter.Marshal(req.Body())
	span.Tag("请求内容", string(msg))
	span.Tag("请求类型", req.ContentType())
	span.Tag("请求方法", "post")
	return err

}

func NewCall(sw *go2sky.Tracer, ctx context.Context, req client.Request) error {
	name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	span, err := sw.CreateExitSpan(ctx, name, req.Service(), func(key, value string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[key] = value
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return err
	}
	utils.SetTContext(ctx)
	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

	defer span.End()
	msg, _ := jsoniter.Marshal(req.Body())
	span.Tag("请求内容", string(msg))
	span.Tag("请求类型", req.ContentType())

	return err
}

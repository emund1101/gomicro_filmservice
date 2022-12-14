// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/film.proto

package film

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for FilmService service

func NewFilmServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for FilmService service

type FilmService interface {
	GetAllList(ctx context.Context, in *FilmListRequest, opts ...client.CallOption) (*Respone, error)
	GetInfoList(ctx context.Context, in *FilmListRequest, opts ...client.CallOption) (*Respone, error)
	GetHallInfo(ctx context.Context, in *FilmInfoRequest, opts ...client.CallOption) (*Respone, error)
}

type filmService struct {
	c    client.Client
	name string
}

func NewFilmService(name string, c client.Client) FilmService {
	return &filmService{
		c:    c,
		name: name,
	}
}

func (c *filmService) GetAllList(ctx context.Context, in *FilmListRequest, opts ...client.CallOption) (*Respone, error) {
	req := c.c.NewRequest(c.name, "FilmService.GetAllList", in)
	out := new(Respone)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filmService) GetInfoList(ctx context.Context, in *FilmListRequest, opts ...client.CallOption) (*Respone, error) {
	req := c.c.NewRequest(c.name, "FilmService.GetInfoList", in)
	out := new(Respone)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filmService) GetHallInfo(ctx context.Context, in *FilmInfoRequest, opts ...client.CallOption) (*Respone, error) {
	req := c.c.NewRequest(c.name, "FilmService.GetHallInfo", in)
	out := new(Respone)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for FilmService service

type FilmServiceHandler interface {
	GetAllList(context.Context, *FilmListRequest, *Respone) error
	GetInfoList(context.Context, *FilmListRequest, *Respone) error
	GetHallInfo(context.Context, *FilmInfoRequest, *Respone) error
}

func RegisterFilmServiceHandler(s server.Server, hdlr FilmServiceHandler, opts ...server.HandlerOption) error {
	type filmService interface {
		GetAllList(ctx context.Context, in *FilmListRequest, out *Respone) error
		GetInfoList(ctx context.Context, in *FilmListRequest, out *Respone) error
		GetHallInfo(ctx context.Context, in *FilmInfoRequest, out *Respone) error
	}
	type FilmService struct {
		filmService
	}
	h := &filmServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&FilmService{h}, opts...))
}

type filmServiceHandler struct {
	FilmServiceHandler
}

func (h *filmServiceHandler) GetAllList(ctx context.Context, in *FilmListRequest, out *Respone) error {
	return h.FilmServiceHandler.GetAllList(ctx, in, out)
}

func (h *filmServiceHandler) GetInfoList(ctx context.Context, in *FilmListRequest, out *Respone) error {
	return h.FilmServiceHandler.GetInfoList(ctx, in, out)
}

func (h *filmServiceHandler) GetHallInfo(ctx context.Context, in *FilmInfoRequest, out *Respone) error {
	return h.FilmServiceHandler.GetHallInfo(ctx, in, out)
}

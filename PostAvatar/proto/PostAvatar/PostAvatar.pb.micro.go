// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/PostAvatar/PostAvatar.proto

package go_micro_srv_PostAvatar

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for PostAvatar service

type PostAvatarService interface {
	PostAvatar(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type postAvatarService struct {
	c    client.Client
	name string
}

func NewPostAvatarService(name string, c client.Client) PostAvatarService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.PostAvatar"
	}
	return &postAvatarService{
		c:    c,
		name: name,
	}
}

func (c *postAvatarService) PostAvatar(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "PostAvatar.PostAvatar", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PostAvatar service

type PostAvatarHandler interface {
	PostAvatar(context.Context, *Request, *Response) error
}

func RegisterPostAvatarHandler(s server.Server, hdlr PostAvatarHandler, opts ...server.HandlerOption) error {
	type postAvatar interface {
		PostAvatar(ctx context.Context, in *Request, out *Response) error
	}
	type PostAvatar struct {
		postAvatar
	}
	h := &postAvatarHandler{hdlr}
	return s.Handle(s.NewHandler(&PostAvatar{h}, opts...))
}

type postAvatarHandler struct {
	PostAvatarHandler
}

func (h *postAvatarHandler) PostAvatar(ctx context.Context, in *Request, out *Response) error {
	return h.PostAvatarHandler.PostAvatar(ctx, in, out)
}

package main

import (
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"house/PostRet/handler"
	"house/PostRet/subscriber"

	PostRet "house/PostRet/proto/PostRet"
)

func main() {
	// 创建grpc服务
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostRet"),
		micro.Version("latest"),
	)

	// 初始化服务
	service.Init()

	// 注册核心服务
	PostRet.RegisterPostRetHandler(service.Server(), new(handler.PostRet))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostRet", service.Server(), new(subscriber.PostRet))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostRet", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

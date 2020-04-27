package main

import (
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"house/PostLogin/handler"
	"house/PostLogin/subscriber"

	PostLogin "house/PostLogin/proto/PostLogin"
)

func main() {
	// 创建服务
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostLogin"),
		micro.Version("latest"),
	)

	// 初始化服务
	service.Init()

	// 注册相应的核心服务
	PostLogin.RegisterPostLoginHandler(service.Server(), new(handler.PostLogin))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostLogin", service.Server(), new(subscriber.PostLogin))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostLogin", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

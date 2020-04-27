package main

import (
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"house/GetSession/handler"
	"house/GetSession/subscriber"

	GetSession "house/GetSession/proto/GetSession"
)

func main() {
	// 创建服务
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetSession"),
		micro.Version("latest"),
	)

	//初始化服务
	service.Init()

	// 注册服务
	GetSession.RegisterGetSessionHandler(service.Server(), new(handler.GetSession))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSession", service.Server(), new(subscriber.GetSession))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSession", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

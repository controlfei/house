package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"house/GetArea/handler"
	"house/GetArea/subscriber"

	GetArea "house/GetArea/proto/GetArea"
)

func main() {
	//创建grpc服务
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetArea"),
		micro.Version("latest"),
	)

	// 初始化服务
	service.Init()

	// 注册服务
	GetArea.RegisterGetAreaHandler(service.Server(), new(handler.GetArea))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetArea", service.Server(), new(subscriber.GetArea))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetArea", service.Server(), subscriber.Handler)

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

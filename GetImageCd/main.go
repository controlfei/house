package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"house/GetImageCd/handler"
	"house/GetImageCd/subscriber"

	GetImageCd "house/GetImageCd/proto/GetImageCd"
)

func main() {
	// 创建服务
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetImageCd"),
		micro.Version("latest"),
	)

	// 初始化服务
	service.Init()

	// 注册句柄
	GetImageCd.RegisterGetImageCdHandler(service.Server(), new(handler.GetImageCd))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetImageCd", service.Server(), new(subscriber.GetImageCd))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetImageCd", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

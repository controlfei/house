package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"house/GetHouseInfo/handler"
	"house/GetHouseInfo/subscriber"

	GetHouseInfo "house/GetHouseInfo/proto/GetHouseInfo"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.GetHouseInfo"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetHouseInfo.RegisterGetHouseInfoHandler(service.Server(), new(handler.GetHouseInfo))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetHouseInfo", service.Server(), new(subscriber.GetHouseInfo))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetHouseInfo", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

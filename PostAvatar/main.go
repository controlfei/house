package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"house/PostAvatar/handler"
	PostAvatar "house/PostAvatar/proto/PostAvatar"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostAvatar"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	PostAvatar.RegisterPostAvatarHandler(service.Server(), new(handler.PostAvatar))

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.PostAvatar", service.Server(), new(subscriber.PostAvatar))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.PostAvatar", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

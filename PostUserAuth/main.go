package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"house/PostUserAuth/handler"
	PostUserAuth "house/PostUserAuth/proto/PostUserAuth"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostUserAuth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	PostUserAuth.RegisterPostUserAuthHandler(service.Server(),new(handler.PostUserAuth))

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.PostUserAuth", service.Server(), new(subscriber.PostUserAuth))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.PostUserAuth", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

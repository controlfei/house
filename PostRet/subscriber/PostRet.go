package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	PostRET "house/PostRet/proto/PostRet"
)

type PostRet struct{}

func (e *PostRet) Handle(ctx context.Context, msg *PostRET.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *PostRET.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}

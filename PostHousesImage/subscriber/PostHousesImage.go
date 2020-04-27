package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	PostHousesImage "house/PostHousesImage/proto/PostHousesImage"
)

type PostHousesImage struct{}

func (e *PostHousesImage) Handle(ctx context.Context, msg *PostHousesImage.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *PostHousesImage.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}

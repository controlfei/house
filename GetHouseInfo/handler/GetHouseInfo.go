package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	GetHouseInfo "house/GetHouseInfo/proto/GetHouseInfo"
)

type GetHouseInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetHouseInfo) Call(ctx context.Context, req *GetHouseInfo.Request, rsp *GetHouseInfo.Response) error {
	log.Log("Received GetHouseInfo.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetHouseInfo) Stream(ctx context.Context, req *GetHouseInfo.StreamingRequest, stream GetHouseInfo.GetHouseInfo_StreamStream) error {
	log.Logf("Received GetHouseInfo.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetHouseInfo.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetHouseInfo) PingPong(ctx context.Context, stream GetHouseInfo.GetHouseInfo_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetHouseInfo.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

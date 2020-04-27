package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"house/homeweb/utils"
	//"github.com/micro/go-micro/util/log"

	DeleteSESSION "house/DeleteSession/proto/DeleteSession"
)

type DeleteSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *DeleteSession) DeleteSession(ctx context.Context, req *DeleteSESSION.Request, rsp *DeleteSESSION.Response) error {
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//链接redis创建句柄
	redis_config_map := map[string]string{
		"key" : utils.G_server_name,
		"conn" : utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbnum" : utils.G_redis_dbnum,
	}
	//确定了解信息
	beego.Info(redis_config_map)
	//将map转化为json
	redis_config, _ := json.Marshal(redis_config_map)
	//链接redis
	bm,err := cache.NewCache("redis",string(redis_config))
	if err != nil {
		beego.Info("创建缓存失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	sesssionid := req.Sessionid
	//创建缓存的key
	sesssionuser_id := sesssionid + "user_id"
	bm.Delete(sesssionuser_id)
	sesssionname:= sesssionid + "name"
	bm.Delete(sesssionname)
	sesssionmobile := sesssionid + "mobile"
	bm.Delete(sesssionmobile)

	return nil
}
//
//// Stream is a server side stream handler called via client.Stream or the generated client code
//func (e *DeleteSession) Stream(ctx context.Context, req *DeleteSESSION.StreamingRequest, stream DeleteSESSION.DeleteSession_StreamStream) error {
//	log.Logf("Received DeleteSession.Stream request with count: %d", req.Count)
//
//	for i := 0; i < int(req.Count); i++ {
//		log.Logf("Responding: %d", i)
//		if err := stream.Send(&DeleteSession.StreamingResponse{
//			Count: int64(i),
//		}); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
//func (e *DeleteSession) PingPong(ctx context.Context, stream DeleteSession.DeleteSession_PingPongStream) error {
//	for {
//		req, err := stream.Recv()
//		if err != nil {
//			return err
//		}
//		log.Logf("Got ping %v", req.Stroke)
//		if err := stream.Send(&DeleteSession.Pong{Stroke: req.Stroke}); err != nil {
//			return err
//		}
//	}
//}

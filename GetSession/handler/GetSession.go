package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"house/homeweb/utils"

	"github.com/micro/go-micro/util/log"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"encoding/json"

	GetSESSION "house/GetSession/proto/GetSession"
)

type GetSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSession) GetSession(ctx context.Context, req *GetSESSION.Request, rsp *GetSESSION.Response) error {
	beego.Info("getsession is begign")
	//初始化返回
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	sessionid := req.Sessionid
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
	//通过sessionid + name(type: string)获取username
	namec := sessionid + "name"
	username  := bm.Get(namec)
	if username ==nil{
		beego.Info("没找到相应的用户",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	rsp.UserName,_ = redis.String(username,nil)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSession) Stream(ctx context.Context, req *GetSESSION.StreamingRequest, stream GetSESSION.GetSession_StreamStream) error {
	log.Logf("Received GetSession.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetSESSION.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetSession) PingPong(ctx context.Context, stream GetSESSION.GetSession_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetSESSION.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

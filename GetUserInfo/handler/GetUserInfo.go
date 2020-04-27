package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/garyburd/redigo/redis"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"house/homeweb/models"
	"house/homeweb/utils"
	"strconv"

	GetUserINFO "house/GetUserInfo/proto/GetUserInfo"
)

type GetUserInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserInfo) GetUserInfo(ctx context.Context, req *GetUserINFO.Request, rsp *GetUserINFO.Response) error {
	beego.Info("this is getuserinfo")
	//初始化返回
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//获取用户id
	sessionid := req.Sessionid
	sessioniduser_id := sessionid + "user_id"
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
	userid := bm.Get(sessioniduser_id)
	userid_str ,err := redis.String(userid,nil)
	beego.Info("this is ",userid_str)
	id,_ := strconv.Atoi(userid_str)
	beego.Info("this is ",id)
	//查询数据库用户信息
	user := models.User{Id: id}
	o := orm.NewOrm()
	err = o.Read(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_USERERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.UserId = strconv.Itoa(user.Id)
	rsp.Name = user.Name
	rsp.Mobile = user.Mobile
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card
	rsp.AvatarUrl = user.Avatar_url


	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
//func (e *GetUserInfo) Stream(ctx context.Context, req *GetUserINFO.StreamingRequest, stream GetUserINFO.GetUserInfo_StreamStream) error {
//	log.Logf("Received GetUserInfo.Stream request with count: %d", req.Count)
//
//	for i := 0; i < int(req.Count); i++ {
//		log.Logf("Responding: %d", i)
//		if err := stream.Send(&GetUserINFO.StreamingResponse{
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
//func (e *GetUserInfo) PingPong(ctx context.Context, stream GetUserInfo.GetUserInfo_PingPongStream) error {
//	for {
//		req, err := stream.Recv()
//		if err != nil {
//			return err
//		}
//		log.Logf("Got ping %v", req.Stroke)
//		if err := stream.Send(&GetUserInfo.Pong{Stroke: req.Stroke}); err != nil {
//			return err
//		}
//	}
//}

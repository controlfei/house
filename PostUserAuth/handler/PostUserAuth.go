package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gomodule/redigo/redis"
	"house/homeweb/models"
	"house/homeweb/utils"
	"strconv"
	"time"

	//"github.com/micro/go-micro/util/log"

	PostUserAUTH "house/PostUserAuth/proto/PostUserAuth"
)

type PostUserAuth struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostUserAuth) PostUserAuth(ctx context.Context, req *PostUserAUTH.Request, rsp *PostUserAUTH.Response) error {
	beego.Info("this is postuserauth")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//获取sessionid
	sesssionid := req.Sessionid
	//链接redis
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
	//查找userid
	sessionuser_id := sesssionid + "user_id"
	id_str,_  := redis.String(bm.Get(sessionuser_id),nil)
	id,err := strconv.Atoi(id_str)

	//更新数据库表
	user := models.User{Id: id,Real_name: req.RealName,Id_card: req.IdCard}
	o := orm.NewOrm()
	_,err = o.Update(&user,"id_card","real_name")
	if err != nil {
		beego.Info("更新数据库表失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	bm.Put(sessionuser_id,id_str,time.Second*600)

	return nil
}

//// Stream is a server side stream handler called via client.Stream or the generated client code
//func (e *PostUserAuth) Stream(ctx context.Context, req *PostUserAuth.StreamingRequest, stream PostUserAuth.PostUserAuth_StreamStream) error {
//	log.Logf("Received PostUserAuth.Stream request with count: %d", req.Count)
//
//	for i := 0; i < int(req.Count); i++ {
//		log.Logf("Responding: %d", i)
//		if err := stream.Send(&PostUserAuth.StreamingResponse{
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
//func (e *PostUserAuth) PingPong(ctx context.Context, stream PostUserAuth.PostUserAuth_PingPongStream) error {
//	for {
//		req, err := stream.Recv()
//		if err != nil {
//			return err
//		}
//		log.Logf("Got ping %v", req.Stroke)
//		if err := stream.Send(&PostUserAuth.Pong{Stroke: req.Stroke}); err != nil {
//			return err
//		}
//	}
//}

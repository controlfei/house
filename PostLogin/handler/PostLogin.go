package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"github.com/micro/go-micro/util/log"
	"house/homeweb/models"
	"house/homeweb/utils"
	"time"

	PostLOGIN "house/PostLogin/proto/PostLogin"
)



type PostLogin struct{}

// 用户登录的核心服务
func (e *PostLogin) PostLogin(ctx context.Context, req *PostLOGIN.Request, rsp *PostLOGIN.Response) error {
	beego.Info("this is postloagin begin")
	///初始化返回参数
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//获取参数
	mobile := req.Mobile
	pwd := req.Password

	//从数据库查找是否有该用户
	o := orm.NewOrm()
	user := models.User{}
	//创建查询句柄
	qs := o.QueryTable("user")
	qs.Filter("mobile",mobile).One(&user)
	if utils.Md5String(pwd) != user.Password_hash{
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//创建sessionid
	sesssionid := utils.Md5String(req.Mobile + req.Password)
	rsp.Sessionid = sesssionid



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
	//创建缓存的key
	sesssionuser_id := sesssionid + "user_id"
	bm.Put(sesssionuser_id,user.Id,time.Second*3600)
	sesssionname:= sesssionid + "name"
	bm.Put(sesssionname,user.Name,time.Second*3600)
	sesssionmobile := sesssionid + "mobile"
	bm.Put(sesssionmobile,user.Mobile,time.Second*3600)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostLogin) Stream(ctx context.Context, req *PostLOGIN.StreamingRequest, stream PostLOGIN.PostLogin_StreamStream) error {
	log.Logf("Received PostLogin.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PostLOGIN.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostLogin) PingPong(ctx context.Context, stream PostLOGIN.PostLogin_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PostLOGIN.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

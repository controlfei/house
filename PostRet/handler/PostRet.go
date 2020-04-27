package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/util/log"
	PostRET "house/PostRet/proto/PostRet"
	"house/homeweb/models"
	"house/homeweb/utils"
	"time"
)
//明文加密
func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type PostRet struct{}

//注册服务代码的核心
func (e *PostRet) PostRet(ctx context.Context, req *PostRET.Request, rsp *PostRET.Response) error {
	beego.Info("this is  postret")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	//验证手机验证码是否正确
	//首先链接redis数据库
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
	//从redis服务器中读取手机短信验证码
	val := bm.Get(req.Mobile)
	if val == nil{
		beego.Info("this is q")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	//将val转化成str
	val_str, _ := redis.String(val, nil)

	//对比验证短信验证码是否正确
	if  val_str != req.SmsCode {
		beego.Info("this is data err")
		rsp.Errno = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
		return nil
	}
	//完成数据验证，将数据添加到数据库
	o := orm.NewOrm()
	user := models.User{
		Mobile: req.Mobile,
		Password_hash: Md5String(req.Password),
		Name: req.Mobile,
	}
	id ,err := o.Insert(&user)
	if err != nil {
		beego.Info("数据加入数据库失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	beego.Info(id)
	//创建sessionid
	sessionid := Md5String(req.Mobile + req.Password)
	rsp.SessionId = sessionid
	//以sessionid为key的一部分来创建session
	bm.Put(sessionid + "name",req.Mobile,time.Second * 3600)
	bm.Put(sessionid + "id",id,time.Second * 3600)
	bm.Put(sessionid + "mobile",req.Mobile,time.Second * 3600)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostRet) Stream(ctx context.Context, req *PostRET.StreamingRequest, stream PostRET.PostRet_StreamStream) error {
	log.Logf("Received PostRet.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PostRET.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostRet) PingPong(ctx context.Context, stream PostRET.PostRet_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PostRET.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

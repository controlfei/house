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
	"path"
	"strconv"

	PostAVATAR "house/PostAvatar/proto/PostAvatar"
)

type PostAvatar struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostAvatar) PostAvatar(ctx context.Context, req *PostAVATAR.Request, rsp *PostAVATAR.Response) error {
	beego.Info("this is postavatat ")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//验证图片上传是否完整
	size := len(req.Avatar)
	if size != int(req.Filesize) {
		beego.Info("数据缺少")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//获取文件后缀
	ext := path.Ext(req.Fileext)  //返回的是 .jpg有个点
	//beego.Info(req.Avatar)
	beego.Info(ext)
	//调用fdfs上传到图片服务器
	fileid ,err := utils.UploadByBuffer(req.Avatar,ext[1:])
	if err != nil {
		beego.Info(err)
		beego.Info("上传失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	beego.Info(fileid)

	//获取sessionid
	sessionid := req.SessionId

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
	//拼接key获取userid
	sessionuser_id :=  sessionid + "user_id"
	userid ,_:= redis.String(bm.Get(sessionuser_id),nil)
	id,_  := strconv.Atoi(userid)

	//将图片的链接地址存储到user表中
	user := models.User{Id: id,Avatar_url: fileid}
	//链接数据
	o := orm.NewOrm()
	//更新数据
	_,err = o.Update(&user,"avatar_url")
	if  err != nil {
		beego.Info("跟新数据失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.AvatarUrl = fileid
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
//func (e *PostAvatar) Stream(ctx context.Context, req *PostAVATAR.StreamingRequest, stream PostAVATAR.PostAvatar_StreamStream) error {
//	log.Logf("Received PostAvatar.Stream request with count: %d", req.Count)
//
//	for i := 0; i < int(req.Count); i++ {
//		log.Logf("Responding: %d", i)
//		if err := stream.Send(&PostAvatar.StreamingResponse{
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
//func (e *PostAvatar) PingPong(ctx context.Context, stream PostAvatar.PostAvatar_PingPongStream) error {
//	for {
//		req, err := stream.Recv()
//		if err != nil {
//			return err
//		}
//		log.Logf("Got ping %v", req.Stroke)
//		if err := stream.Send(&PostAvatar.Pong{Stroke: req.Stroke}); err != nil {
//			return err
//		}
//	}
//}

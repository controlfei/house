package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	GetUserHOUSE "house/GetUserHouses/proto/GetUserHouses"
	"house/homeweb/models"
	"house/homeweb/utils"
	"strconv"
)

type GetUserHouses struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserHouses) GetUserHouses(ctx context.Context, req *GetUserHOUSE.Request, rsp *GetUserHOUSE.Response) error {
	beego.Info("this is getuserhouses")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//获取sessionid
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
	//拼接key
	sessionuser_id := sessionid + "user_id"
	id_str,_ := redis.String(bm.Get(sessionuser_id),nil)
	id,_ := strconv.Atoi(id_str)
	//链接数据库获取数据
	houses := []models.House{}
	o := orm.NewOrm()
	qs := o.QueryTable("house")
	_,err = qs.Filter("user_id",id).All(&houses)
	if err != nil {
		beego.Info("查找房屋失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	house_list,_ := json.Marshal(houses)
	rsp.Mix = house_list

	return nil
}

//// Stream is a server side stream handler called via client.Stream or the generated client code
//func (e *GetUserHouses) Stream(ctx context.Context, req *GetUserHouses.StreamingRequest, stream GetUserHouses.GetUserHouses_StreamStream) error {
//	log.Logf("Received GetUserHouses.Stream request with count: %d", req.Count)
//
//	for i := 0; i < int(req.Count); i++ {
//		log.Logf("Responding: %d", i)
//		if err := stream.Send(&GetUserHouses.StreamingResponse{
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
//func (e *GetUserHouses) PingPong(ctx context.Context, stream GetUserHouses.GetUserHouses_PingPongStream) error {
//	for {
//		req, err := stream.Recv()
//		if err != nil {
//			return err
//		}
//		log.Logf("Got ping %v", req.Stroke)
//		if err := stream.Send(&GetUserHouses.Pong{Stroke: req.Stroke}); err != nil {
//			return err
//		}
//	}
//}

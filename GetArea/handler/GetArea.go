package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"github.com/micro/go-micro/util/log"
	GetAREA "house/GetArea/proto/GetArea"
	"house/homeweb/models"
	"house/homeweb/utils"
	"time"
)

type GetArea struct{}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetArea) Stream(ctx context.Context, req *GetAREA.StreamingRequest, stream GetAREA.GetArea_StreamStream) error {
	log.Logf("Received GetArea.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetAREA.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetArea) PingPong(ctx context.Context, stream GetAREA.GetArea_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetAREA.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

// 获取地区信息服务核心
func (e *GetArea) GetArea(ctx context.Context, req *GetAREA.Request, rsp *GetAREA.Response) error {
	beego.Info(" getarea api/v1.0/areas !!!")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

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
	//获取缓存数据
	areas_info_value := bm.Get("areas_info")
	//如果areas_info_value 不为空则是有缓存的
	if areas_info_value != nil {
		//获取到缓存发送给前端
		ares_info := []map[string]interface{}{}
		//解码
		err = json.Unmarshal(areas_info_value.([]byte),&ares_info)
		//进行循环赋值
		for key, value := range ares_info {
			beego.Info(key,value)
			//创建对的数据类型并进行赋值
			area := GetAREA.Response_Address{Aid: int32(value["aid"].(float64)),Aname:value["aname"].(string) }

			//将数据保存到切片
			rsp.Data = append(rsp.Data,&area)
		}
		return nil
	}
	fmt.Println("testing")
	//如果没有缓存就从数据库里取，并进行缓存
	o := orm.NewOrm()
	//接受地区的信息切片
	var areas []models.Area
	//创建查询条件
	qs := o.QueryTable("area")
	//查询所有地区
	num,err := qs.All(&areas)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	if num == 0{
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(utils.RECODE_NODATA)
		return nil
	}
	beego.Info("写入缓存")
	//将查到的数据转换成json并写入缓存
	ares_info_str,_ := json.Marshal(areas)
	//Put(key string, val interface{}, timeout time.Duration) error
	//存入缓存
	err = bm.Put("areas_info",ares_info_str,time.Second*3600)
	if err != nil {
		beego.Info("数据信息存入缓存失败",err)
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//返回地区信息
	for key, value := range areas {
		beego.Info(key,value)

		area := GetAREA.Response_Address{Aid: int32(value.Id),Aname: value.Name}
		rsp.Data = append(rsp.Data,&area)
	}

	return nil
}

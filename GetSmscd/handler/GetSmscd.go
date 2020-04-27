package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"house/homeweb/models"
	"house/homeweb/utils"
	"github.com/micro/go-micro/util/log"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	GetSmsCD "house/GetSmscd/proto/GetSmscd"
	"math/rand"
	"strconv"
	"time"
)

type GetSmscd struct{}

// 获取手机验证码的核心服务
func (e *GetSmscd) GetSmscd(ctx context.Context, req *GetSmsCD.Request, rsp *GetSmsCD.Response) error {
	beego.Info("getphonecode is begign")
	//初始化返回
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//验证手机号是否已经被注册
	//创建数据库句柄
	o := orm.NewOrm()
	//使用手机号进行查询
	user := models.User{Mobile: req.Mobile}
	err := o.Read(&user)
	if err == nil {
		rsp.Errno = utils.RECODE_MOBILEERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//验证图片验证码是否正确
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
	//通过uuid进行查找
	val := bm.Get(req.Uuid)
	beego.Info("this is val",val)
	if val == nil {
		beego.Info("查找失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	//格式转换
	val_str,err := redis.String(val,nil)
	if val_str != req.Imagestr {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}
	//生成随机数
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(8998)+1001
	beego.Info("this is ",size)
	//发送短信
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "youralikey", "youralisecret")
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = req.Mobile
	request.SignName = "yoursignname"
	request.TemplateCode = "youralicode"
	request.TemplateParam = "{\"code\":\""+strconv.Itoa(size) +"\"}"
	_, err = client.SendSms(request)
	if err != nil {
		fmt.Println("failed")
		fmt.Print(err.Error())
	}
	beego.Info(req.Mobile)
	//将手机验证码放入缓存
	err = bm.Put(req.Mobile,size,time.Second*60)
	if err != nil {
		beego.Info("faield is this ")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSmscd) Stream(ctx context.Context, req *GetSmsCD.StreamingRequest, stream GetSmsCD.GetSmscd_StreamStream) error {
	log.Logf("Received GetSmscd.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetSmsCD.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetSmscd) PingPong(ctx context.Context, stream GetSmsCD.GetSmscd_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetSmsCD.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

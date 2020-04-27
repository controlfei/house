package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"house/homeweb/utils"
	"image/color"
	"time"

	"github.com/micro/go-micro/util/log"

	GetImageCD "house/GetImageCd/proto/GetImageCd"
)

type GetImageCd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetImageCd) Call(ctx context.Context, req *GetImageCD.Request, rsp *GetImageCD.Response) error {
	log.Log("Received GetImageCd.Call request")
	rsp.Errno = "Hello " + req.Uuid
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetImageCd) Stream(ctx context.Context, req *GetImageCD.StreamingRequest, stream GetImageCD.GetImageCd_StreamStream) error {
	log.Logf("Received GetImageCd.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetImageCD.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetImageCd) PingPong(ctx context.Context, stream GetImageCD.GetImageCd_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetImageCD.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

//生成验证码
func (e *GetImageCd) GetImageCd(ctx context.Context, req *GetImageCD.Request, rsp *GetImageCD.Response) error {
	//生成图片
	cap := captcha.New()
	//通过句柄调用 字体文件
	if err := cap.SetFont("/home/linpengfei/go/src/house/homeweb/static/fzzy.ttf"); err != nil {
		panic(err.Error())
	}
	//设置图片大小
	cap.SetSize(90,41)
	//设置感染强度
	cap.SetDisturbance(captcha.NORMAL)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	image,str := cap.Create(4,captcha.NUM)

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
	//将验证码和uuid存入缓存
	err = bm.Put(req.Uuid,str,time.Second*300)
	//将图片分解成我们proto定义的数据
	img1 := *image
	img2 := *img1.RGBA

	//返回信息】
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//开始拆分图片
	rsp.Pix = []byte(img2.Pix)
	rsp.Stride = int64(img2.Stride)
	rsp.MIN = &GetImageCD.Response_Point{X: int64(img2.Rect.Min.X),Y: int64(img2.Rect.Min.Y)}
	rsp.MAX = &GetImageCD.Response_Point{X: int64(img2.Rect.Max.X),Y: int64(img2.Rect.Max.Y)}

	return nil
}

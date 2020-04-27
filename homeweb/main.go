package main

import (
        "github.com/julienschmidt/httprouter"
        "github.com/micro/go-micro/util/log"
        "github.com/micro/go-micro/web"
    "house/homeweb/handler"
    _ "house/homeweb/models"
        "net/http"
)


func main() {
	// 创建web服务
    service :=  web.NewService(
                web.Name("go.micro.web.homeweb"),
                web.Version("latest"),
                web.Address(":8999"),
    )

	// 初始化服务
	if err := service.Init(); err != nil {
	    log.Fatal(err)
	}
	rou := httprouter.New()
	//映射静态页面
	rou.NotFound = http.FileServer(http.Dir("html"))
	//后续陆续添加服务所以这个文件的这个地方会一直添加内容
	rou.GET("/api/v1.0/areas",handler.GetArea)
	//下面两个目前并不实现服务
	//获取session
	rou.GET("/api/v1.0/session",handler.GetSession)
	//获取index
	rou.GET("/api/v1.0/house/index",handler.GetIndex)
	//获取验证码
	rou.GET("/api/v1.0/imagecode/:uuid",handler.GetImageCd)
	//获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile",handler.GetSmscd)
	//用户注册
	rou.POST("/api/v1.0/users",handler.PostRet)
	//用户登录
	rou.POST("/api/v1.0/sessions",handler.PostLogin)
	//退出登录
	rou.DELETE("/api/v1.0/session",handler.DeleteSession)
	//获取用户信息
	rou.GET("/api/v1.0/user",handler.GetUserInfo)

	//检查用户实名认证
	rou.GET("/api/v1.0/user/auth",handler.GetUserAuth)
	//上传头像
	rou.POST("/api/v1.0/user/avatar",handler.PostAvatar)
	//上传头像
	rou.POST("/api/v1.0/user/auth",handler.PostUserAuth)
	//获取用户发布的房源信息
	rou.GET("/api/v1.0/user/houses",handler.GetUserHouses)
	//上传房屋图片
	rou.POST("/api/v1.0/houses/:id/images",handler.PostHousesImage)
	//发布房屋信息
	rou.POST("/api/v1.0/houses",handler.PostHouses)

	//获取房屋信息
	rou.GET("/api/v1.0/houses/:id",handler.GetHouseInfo)

	// 注册服务
    service.Handle("/", rou)

    //运行服务
    if err := service.Run(); err != nil {
        log.Fatal(err)
    }
}

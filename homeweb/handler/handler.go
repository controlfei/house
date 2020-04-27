package handler

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/service/grpc"
	DeleteSESSION "house/DeleteSession/proto/DeleteSession"
	GetAREA "house/GetArea/proto/GetArea"
	GetImageCD "house/GetImageCd/proto/GetImageCd"
	GetSESSION "house/GetSession/proto/GetSession"
	GetSmsCD "house/GetSmscd/proto/GetSmscd"
	GetUserHOUSES "house/GetUserHouses/proto/GetUserHouses"
	GetUSERINFO "house/GetUserInfo/proto/GetUserInfo"
	PostAVATAR "house/PostAvatar/proto/PostAvatar"
	PostHOUSES "house/PostHouses/proto/PostHouses"
	PostLOGIN "house/PostLogin/proto/PostLogin"
	PostRET "house/PostRet/proto/PostRet"
	PostUserAUTH "house/PostUserAuth/proto/PostUserAuth"
	GetHouseINFO  "house/GetHouseInfo/proto/GetHouseInfo"
	PostHousesIMAGE "house/PostHousesImage/proto/PostHousesImage"
	"house/homeweb/models"
	"house/homeweb/utils"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

//func HomewebCall(w http.ResponseWriter, r *http.Request) {
//	// decode the incoming request as json
//	var request map[string]interface{}
//	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//
//	// call the backend service
//	homewebClient := homeweb.NewHomewebService("go.micro.srv.homeweb", client.DefaultClient)
//	rsp, err := homewebClient.Call(context.TODO(), &homeweb.Request{
//		Name: request["name"].(string),
//	})
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//
//	// we want to augment the response
//	response := map[string]interface{}{
//		"msg": rsp.Msg,
//		"ref": time.Now().UnixNano(),
//	}
//
//	// encode and write the response as json
//	if err := json.NewEncoder(w).Encode(response); err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//}

//获取地区服务
func GetArea(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	//创建服务获取句柄
	server := grpc.NewService()
	//初始化服务
	server.Init()
	//创建获取地区的服务并且返回句柄
	areaClient := GetAREA.NewGetAreaService("go.micro.srv.GetArea", server.Client())
	rsp, err := areaClient.GetArea(context.TODO(), &GetAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 创建返回类型的切片
	area_list := []models.Area{}
	//循环读取服务返回的数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id:int(value.Aid),Name: value.Aname,Houses: nil}
		area_list = append(area_list,tmp)
	}
	//创建返回数据map
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data" : area_list,
	}

	w.Header().Set("Content-Type","application/json")

	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

//获取session
func GetSession(w http.ResponseWriter,r *http.Request,_ httprouter.Params)  {
	//通过cookie获取sessionid  name是userlogin之前就定义好的
	cookie,err :=r.Cookie("userlogin")
	if  err != nil || cookie.Value == ""{
		// 接受返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_SESSIONERR,
			"Errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-type","application/json")

		//将返回数据map发送到前段
		if err := json.NewEncoder(w).Encode(response);err != nil {
			http.Error(w,err.Error(),503)
			return
		}
		return
	}
	//创建服务获取句柄
	server := grpc.NewService()
	//初始化服务
	server.Init()
	// 链接服务并传递参数
	SessionClient := GetSESSION.NewGetSessionService("go.micro.srv.GetSession", server.Client())
	rsp, err := SessionClient.GetSession(context.TODO(), &GetSESSION.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data:=make(map[string]string,2)
	data["name"] = rsp.UserName
	fmt.Println(data)
	// 接受返回的数据
	response := map[string]interface{}{
		"Errno": rsp.Errno,
		"Errmsg": rsp.Errmsg,
		"data" : data,
		//"Errno" : utils.RECODE_OK,
		//"Errmsg" : utils.RecodeText(utils.RECODE_OK),
	}
	w.Header().Set("Content-type","application/json")

	//将返回数据map发送到前段
	if err := json.NewEncoder(w).Encode(response);err != nil {
		http.Error(w,err.Error(),503)
		return
	}
}

//获取首页轮播
func GetIndex(w http.ResponseWriter,r *http.Request,_ httprouter.Params)  {
	//创建返回数据
	response := map[string]interface{}{
		"errno" : utils.RECODE_OK,
		"errmsg" : utils.RecodeText(utils.RECODE_OK),
	}
	w.Header().Set("Content-type","application/json")

	//将返回数据map发送到前段
	if err := json.NewEncoder(w).Encode(response);err != nil {
		http.Error(w,err.Error(),503)
		return
	}
}

//获取验证码
func GetImageCd(w http.ResponseWriter, r *http.Request,p httprouter.Params) {
	//获取服务器句柄
	service := grpc.NewService()
	//初始化服务
	service.Init()
	// 链接服务端并且传递参数
	codeClient := GetImageCD.NewGetImageCdService("go.micro.srv.GetImageCd", service.Client())
	rsp, err := codeClient.GetImageCd(context.TODO(), &GetImageCD.Request{
		Uuid: p.ByName("uuid"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var img image.RGBA
	img.Pix = []uint8(rsp.Pix)
	img.Stride = int(rsp.Stride)
	img.Rect.Min.X = int(rsp.MIN.X)
	img.Rect.Min.Y= int(rsp.MIN.Y)
	img.Rect.Max.X = int(rsp.MAX.X)
	img.Rect.Max.Y = int(rsp.MAX.Y)

	var image captcha.Image
	image.RGBA = &img
	//返回图片
	png.Encode(w,image)
}

//获取短信验证码
func GetSmscd(w http.ResponseWriter, r *http.Request,p httprouter.Params) {
	//获取参数
	text := r.URL.Query()["text"][0]
	id := r.URL.Query()["id"][0]
	mobile := p.ByName("mobile")
	//对参数进行判断
	mobile_reg := regexp.MustCompile("0?(13|14|15|17|18|19)[0-9]{9}")
	bl := mobile_reg.MatchString(mobile)

	if bl == false{
		// 接受返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_MOBILEERR,
			"Errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
		}
		//设置返回头格式
		w.Header().Set("Content-Type","application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	//创建服务获取句柄
	server := grpc.NewService()
	//初始化服务
	server.Init()
	// 链接服务并传递参数
	SmsClient := GetSmsCD.NewGetSmscdService("go.micro.srv.GetSmscd", server.Client())
	rsp, err := SmsClient.GetSmscd(context.TODO(), &GetSmsCD.Request{
		Mobile: mobile,
		Imagestr: text,
		Uuid: id,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 接受返回的数据
	response := map[string]interface{}{
		"Errno": rsp.Errno,
		"Errmsg": rsp.Errmsg,
	}
	//设置返回头格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}


//用户注册
func PostRet(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// 解析并获取参数
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mobile := request["mobile"]
	pwd :=  request["password"]
	sms_code := request["sms_code"]

	//对参数进行验证
	if mobile.(string) == "" || pwd.(string) == "" || sms_code.(string) == "" {
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_DATAERR,
			"Errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	//创建并获取服务句柄
	server := grpc.NewService()
	//初始化服务
	server.Init()
	//链接服务并传递参数
	homewebClient := PostRET.NewPostRetService("go.micro.srv.PostRet",server.Client())
	rsp, err := homewebClient.PostRet(context.TODO(), &PostRET.Request{
		Mobile:  mobile.(string),
		Password: pwd.(string),
		SmsCode:  sms_code.(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//给浏览器设置cookie
	//首先获取cookie  统一cookie名称为 userlogin
	cook,err := r.Cookie("userlogin")
	if err != nil  || cook.Value == ""{
		cookie := http.Cookie{
			Name:       "userlogin",
			Value:      rsp.SessionId,
			Path:       "/",
			MaxAge:     3600,
		}
		//给浏览器设置cookie
		http.SetCookie(w,&cookie)
	}

	// 得到服务返回的数据
	response := map[string]interface{}{
		"Errno": rsp.Errno,
		"Errmsg": rsp.Errmsg,
	}
	//设置返回头格式
	w.Header().Set("Content-Type","application/json")

	// 将返回数据封装成json格式
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//用户登录
func PostLogin(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// 解析获取参数
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mobile := request["mobile"].(string)
	pwd := request["password"].(string)

	if mobile == "" || pwd == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_DATAERR,
			"Errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	//创建grpc客户端 、
	server := grpc.NewService()
	//初始化服务
	server.Init()

	// 链接服务并请求
	LoginClient := PostLOGIN.NewPostLoginService("go.micro.srv.PostLogin", server.Client())
	rsp, err := LoginClient.PostLogin(context.TODO(), &PostLOGIN.Request{
		Mobile: mobile,
		Password: pwd,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 请求后台服务成功之后，将cookie写入浏览器
	//首先判断是否存在cookie,存在说明已经登陆了，直接返回
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		//设置cookie
		cooke := http.Cookie{
			Name:   "userlogin",
			Value:  rsp.Sessionid,
			Path:   "/",
			MaxAge: 600,
		}
		http.SetCookie(w,&cooke)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// 将数据封装成json并返回相应
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//用户退出用户
func DeleteSession(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	/// 获取sessionid
	cookie,err  := r.Cookie("userlogin")
	if err != nil  || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	sessionid := cookie.Value
	//创建grpc客户端
	server := grpc.NewService()
	//服务初始化
	server.Init()
	// 链接服务
	DeleteSessionClient := DeleteSESSION.NewDeleteSessionService("go.micro.srv.DeleteSession",server.Client())
	rsp, err := DeleteSessionClient.DeleteSession(context.TODO(), &DeleteSESSION.Request{
		Sessionid: sessionid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//删除sessionid
	cookie ,err  = r.Cookie("userlogin")
	if err == nil  || cookie.Value != ""{
		cookie := http.Cookie{Name: "userlogin",Value: "",Path: "/",MaxAge: -1}
		http.SetCookie(w,&cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno" : rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	//创建grpc客户端
	server := grpc.NewService()
	//初始化服务
	server.Init()
	// call the backend service
	UserInfoClient := GetUSERINFO.NewGetUserInfoService("go.micro.srv.GetUserInfo", server.Client())
	rsp, err := UserInfoClient.GetUserInfo(context.TODO(), &GetUSERINFO.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//封装用户的信息
	data := make(map[string]interface{})
	data["user_id"]  = rsp.UserId
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2url(rsp.AvatarUrl)

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data" :data,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//上传用户头像
func PostAvatar(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// 解析参数
	File,FileHeader,err := r.FormFile("avatar")
	if err != nil {
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	beego.Info("文件大小,",FileHeader.Size)
	beego.Info("文件名,",FileHeader.Filename)
	//创建一个文件大小的切片
	filebuf := make([]byte,FileHeader.Size)
	//将file的数据读取到filebuf中
	_,err = File.Read(filebuf)
	beego.Info(filebuf)
	if err != nil {
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	/// 获取sessionid
	cookie,err  := r.Cookie("userlogin")
	if err != nil  || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	//创建服务
	server := grpc.NewService()
	//初始化服务
	server.Init()
	// call the backend service
	AvatarClient := PostAVATAR.NewPostAvatarService("go.micro.srv.PostAvatar", server.Client())
	rsp, err := AvatarClient.PostAvatar(context.TODO(), &PostAVATAR.Request{
		SessionId: cookie.Value,
		Fileext: FileHeader.Filename,
		Filesize: FileHeader.Size,
		Avatar: filebuf,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := make(map[string]string)
	data["avatar_url"] = utils.AddDomain2url(rsp.AvatarUrl)

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data" :data,
	}

	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//用户信息检查
func GetUserAuth(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	//创建grpc客户端
	server := grpc.NewService()
	//初始化服务
	server.Init()
	// call the backend service
	UserInfoClient := GetUSERINFO.NewGetUserInfoService("go.micro.srv.GetUserInfo", server.Client())
	rsp, err := UserInfoClient.GetUserInfo(context.TODO(), &GetUSERINFO.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//封装用户的信息
	data := make(map[string]interface{})
	data["user_id"]  = rsp.UserId
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2url(rsp.AvatarUrl)

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data" :data,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
// 进行实名认证
func PostUserAuth(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// 获取参数 idcard  realname
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//验证参数

	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	//创建服务句柄
	server := grpc.NewService()
	//初始化服务
	server.Client()
	// 请求服务
	UserAuthClient := PostUserAUTH.NewPostUserAuthService("go.micro.srv.PostUserAuth", server.Client())
	rsp, err := UserAuthClient.PostUserAuth(context.TODO(), &PostUserAUTH.Request{
		Sessionid: cookie.Value,
		RealName: request["real_name"].(string),
		IdCard: request["id_card"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 接受返回值
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取用户发布的房源信息
func GetUserHouses(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	//创建服务句柄
	server := grpc.NewService()
	//初始化服务
	server.Client()
	// call the backend service
	GetHouseClient := GetUserHOUSES.NewGetUserHousesService("go.micro.srv.GetUserHouses", server.Client())
	rsp, err := GetHouseClient.GetUserHouses(context.TODO(), &GetUserHOUSES.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//接受返回的参数
	house_list := []models.House{}  //查找所有相关数据的所有字段翻进去
	json.Unmarshal(rsp.Mix,&house_list)

	var houses []interface{}
	for _, val := range house_list {
		//返回我们需要的字段  所以在models进行封装返回函数
		houses = append(houses,val.To_house_info())
	}
	data := make(map[string]interface{})
	data["houses"] = houses

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data": data,
	}

	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//发布房屋信息
func PostHouses(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	//获取前端传来的json二进制留
	body ,_ := ioutil.ReadAll(r.Body)

	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	//创建服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	PhousesClient := PostHOUSES.NewPostHousesService("go.micro.srv.PostHouses", client.DefaultClient)
	rsp, err := PhousesClient.PostHouses(context.TODO(), &PostHOUSES.Request{
		Sessionid: cookie.Value,
		Body: body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data:= make(map[string]interface{})
	data["house_id"] = rsp.HousesId

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Error,
		"errmsg": rsp.Errmsg,
		"data": data,

	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//上传房屋的图片
func PostHousesImage(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	//创建服务
	server := grpc.NewService()
	server.Init()

	//获取房屋ID
	houseid := ps.ByName("id")

	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	file,handler,err := r.FormFile("house_image")

	if err != nil {
		beego.Info("formfile faiels")
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	filebuf := make([]byte,handler.Size)
	//读取文件到filebuf
	_,err=file.Read(filebuf)
	if err != nil {
		beego.Info(" 读取到buf失败")
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	// call the backend service
	houseImageClient := PostHousesIMAGE.NewPostHousesImageService("go.micro.srv.PostHousesImage", server.Client())
	rsp, err := houseImageClient.PostHousesImage(context.TODO(), &PostHousesIMAGE.Request{
		Sessionid: cookie.Value,
		Id: houseid,
		Image: filebuf,
		Filesize: handler.Size,
		Filename: handler.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := make(map[string]interface{})
	data["url"] = utils.AddDomain2url(rsp.Url)

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data": data,

	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取房屋信息
func GetHouseInfo(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	//获取房屋id
	id := ps.ByName("id")
	//获取sessionid
	cookie,err := r.Cookie("userlogin")
	if err != nil || cookie.Value == ""{
		// 得到服务返回的数据
		response := map[string]interface{}{
			"Errno": utils.RECODE_USERERR,
			"Errmsg": utils.RecodeText(utils.RECODE_USERERR),
		}

		beego.Info("this is one")
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		// 将返回数据封装成json格式
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	server := grpc.NewService()
	server.Init()
	// call the backend service
	homewebClient := GetHouseINFO.NewGetHouseInfoService("go.micro.srv.GetHouseInfo", server.Client())
	rsp, err := homewebClient.GetHouseInfo(context.TODO(), &GetHouseINFO.Request{
		Id: id,
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	house := models.House{}
	json.Unmarshal(rsp.Housedata,&house)
	data := make(map[string]interface{})
	data["user_id"] = int(rsp.Userid)
	data["house"] = house.To_one_house_desc()

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data" :data,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
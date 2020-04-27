package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"house/homeweb/models"
	"house/homeweb/utils"
	"strconv"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	PostHOUSES "house/PostHouses/proto/PostHouses"
	"encoding/json"
)

type PostHouses struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostHouses) PostHouses(ctx context.Context, req *PostHOUSES.Request, rsp *PostHOUSES.Response) error {
	beego.Info("this is post house")
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

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
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	//拼接key
	sessionuser_id := sessionid + "user_id"
	id_str,_ := redis.String(bm.Get(sessionuser_id),nil)
	id,_ := strconv.Atoi(id_str)

	//解析发送过来的body参数
	var request = make(map[string]interface{})
	house := models.House{}
	house.Title = request["title"].(string)
	price , _ := strconv.Atoi(request["price"].(string))
	house.Price = price * 100
	house.Address = request["address"].(string)
	house.Room_count,_ = strconv.Atoi(request["room_count"].(string))
	house.Acreage,_ = strconv.Atoi(request["acrege"].(string))
	house.Unit = request["unit"].(string)
	house.Capacity,_ = strconv.Atoi(request["capacity"].(string))
	house.Beds = request["beds"].(string)
	deposit,_ := strconv.Atoi(request["deposit"].(string))
	house.Deposit = deposit
	house.Min_days,_ = strconv.Atoi(request["min_days"].(string))
	house.Max_days,_ = strconv.Atoi(request["max_days"].(string))
	facility := []*models.Facility{}
	for _, val := range request["facility"].([]interface{}) {
		//将设施编号转换成为对应的类型
		fid , _ := strconv.Atoi(val.(string))
		//创建临时变量 使用设施编号创建的设施表变量的指针
		ftmp := &models.Facility{Id: fid}
		facility = append(facility,ftmp)
	}

	area_id ,_ :=strconv.Atoi(request["area_id"].(string))
	area := models.Area{Id: area_id}

	user := models.User{Id: id}
	house.Area = &area
	house.User = &user

	//插入数据库
	o := orm.NewOrm()
	house_id,err  := o.Insert(&house)
	if err != nil {
		beego.Info("插入数据库失败",err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	//多对多插入房屋设施
	m2m := o.QueryM2M(&house,"Facilities")
	_,err = m2m.Add(facility)
	if err != nil {
		beego.Info("插入设施表失败",err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	rsp.HousesId = strconv.Itoa(int(house_id))

	return nil
}

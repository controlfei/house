package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"house/homeweb/models"
	"house/homeweb/utils"
	"path"
	"strconv"

	PostHousesIMAGE "house/PostHousesImage/proto/PostHousesImage"
)

type PostHousesImage struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostHousesImage) PostHousesImage(ctx context.Context, req *PostHousesIMAGE.Request, rsp *PostHousesIMAGE.Response) error {
	beego.Info("this is post house image")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	fileext := path.Ext(req.Filename)

	fileid ,err := utils.UploadByBuffer(req.Image,fileext[1:])
	if err != nil {
		beego.Info("上传失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	house_id ,_:= strconv.Atoi(req.Id)

	house := models.House{Id: house_id}

	o := orm.NewOrm()
	err = o.Read(&house)

	if  err != nil {
		beego.Info("插入数据库事变")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	if house.Index_image_url == "" {
		house.Index_image_url = fileid
	}

	houseimage := models.HouseImage{House: &house,Url: fileid}
	house.Images = append(house.Images,&houseimage)

	_,err = o.Insert(&houseimage)
	if  err != nil {
		beego.Info("插入数据库事变2")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	_,err = o.Update(&house)
	if  err != nil {
		beego.Info("插入数据库事变e")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.Url = fileid



	return nil
}

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	fc  "github.com/tedcy/fdfs_client"
)

//字符串拼接

// 将url加上 http:// ip:port/ 前缀
func AddDomain2url(url string) string  {
	domain_url := "http://" + G_fastdfs_addr + G_fastdfs_port +"/" +url
	return domain_url
}

//明文加密
func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func UploadByBuffer(filebuffer []byte,fileExt string) (fileid string,err error) {
	//返回fdfs客户端句柄
	client,err  := fc.NewClientWithConfig("/home/linpengfei/go/src/house/homeweb/conf/client.conf")
	if err != nil {
		fmt.Println("创建客户端句柄失败")
		fmt.Println(err)
		fileid = ""
		return
	}
	fileid,err = client.UploadByBuffer(filebuffer,fileExt)
	if err != nil {
		fmt.Println("上传文件失败")
		fileid = ""
		return
	}
	return fileid,nil
}
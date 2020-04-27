package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"math/rand"
	"strconv"
	"time"
)

func main()  {
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(8998)+1001
	//发送短信
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAIJssLApJ2Ola3", "Im8YNliA3bNJq4oiydK4NqZ34W7smw")
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = "13121331024"
	request.SignName = "飞哥商场"
	request.TemplateCode = "SMS_181864283"

	request.TemplateParam = "{\"code\":\""+strconv.Itoa(size) +"\"}"
	fmt.Println(request.TemplateParam)
	res, err := client.SendSms(request)
	if err != nil {
		fmt.Println("failed")
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", res)
	//isv.INVALID_PARAMETERS
}
package main

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

var aliyunClient *dysmsapi20170525.Client

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func AliyunInit(AccessKeyId, AccessKeySecret string) {
	var _err error
	aliyunClient, _err = CreateClient(&AccessKeyId, &AccessKeySecret)
	if _err != nil {
		panic(fmt.Sprintf("aliyunClient init fail! err=%v", _err))
	}
}

func AliyunSendSms(code, phone string) error {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}
	runtime := &util.RuntimeOptions{}
	rsp, err := aliyunClient.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		return err
	}
	fmt.Println(rsp)
	if *rsp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("AliyunSendSms fail, statusCode %d", *rsp.StatusCode))
	}
	if *rsp.Body.Code != "OK" {
		return errors.New(fmt.Sprintf("AliyunSendSms fail, bodyCode=%s msg=%s", *rsp.Body.Code, *rsp.Body.Message))
	}
	return nil
}

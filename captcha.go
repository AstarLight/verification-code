package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
	"time"
)

// 图形验证码
//验证码属性: https://captcha.mojotv.cn/

var store = base64Captcha.DefaultMemStore

//获取验证码
func MakeCaptcha(codeLen int) (string, string, error) {
	//定义一个driver
	var driver base64Captcha.Driver
	driverDigit := &base64Captcha.DriverDigit{
		Height:   80,  //高度
		Width:    240, //宽度
		MaxSkew:  0.7,
		Length:   codeLen, //数字个数
		DotCount: 80,
	}
	driver = driverDigit
	//生成验证码
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	code := GetCodeById(id)
	return code, b64s, err
}

// tag:唯一标记，如phone username等
// from: 标记是哪个业务申请的验证码
func GenCaptchaCodeKey(tag, from string) string {
	return "CAPTCHA-CODE-" + from + "-" + tag
}

func CreateCaptcha(tag, from string, ttl, codeLen int) (string, string, error) {
	key := GenCaptchaCodeKey(tag, "login")
	code, b64s, err := MakeCaptcha(codeLen)
	if err != nil {
		return "", "", err
	}

	// 如果code没有过期，是不允许再发送的
	err = redisDb.Set(ctx, key, code, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		fmt.Printf("SendSmsCode AliyunSendSms fail %v\n", err)
		return "", "", err
	}

	return code, b64s, err
}

func GetCodeById(id string) string {
	return store.Get(id, true)
}

func VerifyCaptchaCode(tag, inputCode, from string) error {
	key := GenCaptchaCodeKey(tag, from)
	code, err := redisDb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errCodeNotExist
		}
		return err
	}

	// 对比后马上删除
	err = redisDb.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("redis del fail %v\n", err)
		return err
	}

	if inputCode != code {
		return errCodeNotMatch
	}

	return nil
}

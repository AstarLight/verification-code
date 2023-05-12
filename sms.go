package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// 手机短信验证码
var (
	errValidateReachLimit = errors.New("validate times reach limit")
	errCodeNotMatch       = errors.New("code not match")
	errCodeNotExpire      = errors.New("code not expire")
	errCodeNotExist       = errors.New("code not exist")
)

// from: 标记是哪个业务申请的验证码
func GenSmsCodeKey(phone, from string) string {
	return "SMS-CODE-" + from + "-" + phone
}

// from: 标记是哪个业务申请的验证码
func ValidateSmsCode(phone, code, from string) error {
	key := GenSmsCodeKey(phone, from)

	smsCode, err := redisDb.Get(ctx, key).Result()
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

	if smsCode != code {
		return errCodeNotMatch
	}

	return nil
}

func SendSmsCode(phone, code, from string, ttl int) error {
	// 先插入redis，再发短信
	key := GenSmsCodeKey(phone, from)
	// 如果code没有过期，是不允许再发送的
	success, err := redisDb.SetNX(ctx, key, code, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		fmt.Printf("SendSmsCode AliyunSendSms fail %v\n", err)
		return err
	}

	if !success {
		return errCodeNotExpire
	}

	// 发短信，调用第三方接口
	err = AliyunSendSms(code, phone)
	if err != nil {
		fmt.Printf("SendSmsCode AliyunSendSms fail %v\n", err)
		return nil
	}
	return nil
}

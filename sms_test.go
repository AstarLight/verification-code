package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	redisAddr := "10.10.40.231:6379"
	password := ""
	db := 0
	RedisInit(redisAddr, password, db)

	AccessKeyId := "" // 阿里云accesskey
	AccessKeySecret := "" // 阿里云accessSecret
	AliyunInit(AccessKeyId, AccessKeySecret)
	retCode := m.Run() // 执行测试
	os.Exit(retCode)   // 退出测试
}

func TestSendSmsCode(t *testing.T) {
	code := "658574"
	phone := "15913137599"
	from := "login"
	smsTTL := 60
	err := SendSmsCode(phone, code, from, smsTTL)
	if err != nil {
		t.Errorf("expected:%v, got:%v", nil, err)
	}

	// 再次发送，期望是不能发送，因为验证码有60秒的Cd
	err = SendSmsCode(phone, code, from, smsTTL)
	if err != errCodeNotExpire {
		t.Errorf("expected:%v, got:%v", errCodeNotExpire, err)
	}

	// 验证
	err = ValidateSmsCode(phone, code, from)
	if err != nil {
		t.Errorf("expected:%v, got:%v", nil, err)
	}

	// 再次验证，此时验证码理应删除
	err = ValidateSmsCode(phone, code, from)
	if err != errCodeNotExist {
		t.Errorf("expected:%v, got:%v", errCodeNotExist, err)
	}
}

func TestSendSmsCodeWithWrongCode(t *testing.T) {
	code := "658574"
	phone := "15913137599"
	from := "login"
	smsTTL := 60
	
	err := SendSmsCode(phone, code, from, smsTTL)
	if err != nil {
		t.Errorf("expected:%v, got:%v", nil, err)
	}

	// 用错误的code验证
	err = ValidateSmsCode(phone, "789789", from)
	if err != errCodeNotMatch {
		t.Errorf("expected:%v, got:%v", errCodeNotMatch, err)
	}

	// 用错误的phone验证
	err = ValidateSmsCode("15935978964", code, from)
	if err != errCodeNotExist {
		t.Errorf("expected:%v, got:%v", errCodeNotExist, err)
	}
}

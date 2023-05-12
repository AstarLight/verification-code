package main

import (
	"testing"
	"fmt"
)

func TestSendAndVerify(t *testing.T) {
	mailUser := "lijunshi_mhxy7@163.com"
	mailPass := "" // 授权码
	mailTo := "lijunshi2015@163.com"
	codeLen := 4
	codeTTL := 5*60
	from := "login" // 业务场景标记

	code := GenRandomCode(codeLen)

	subjectText := "登录验证码" // 邮件主题
	bodyText := fmt.Sprintf("您的验证码是:%s， 有效期为5分钟", code) // 邮件正文

	options := &MailOptions{
		MailUser: mailUser,
		MailPass: mailPass,
		MailTo:   mailTo,
		Subject:  subjectText,
		Body:     bodyText,
	}

	err := SendMailCode(options, code, from, codeTTL)
	if err != nil {
		t.Error("SendMailCode error", err)
	}

	err = ValidateMailCode(mailTo, code, from)
	if err != nil {
		t.Error("ValidateMailCode error", err)
	}
}

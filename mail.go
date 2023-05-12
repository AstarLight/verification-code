package main

// 邮件验证码
import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
	"strings"
	"time"
)

type MailOptions struct {
	MailHost string
	MailPort int
	MailUser string // 发件人
	MailPass string // 发件人密码
	MailTo   string // 收件人 多个用,分割
	Subject  string // 邮件主题
	Body     string // 邮件内容
}

func MailSend(o *MailOptions) error {

	m := gomail.NewMessage()

	//设置发件人
	m.SetHeader("From", o.MailUser)

	//设置发送给多个用户
	mailArrTo := strings.Split(o.MailTo, ",")
	m.SetHeader("To", mailArrTo...)

	//设置邮件主题
	m.SetHeader("Subject", o.Subject)

	//设置邮件正文
	m.SetBody("text/html", o.Body)

	d := gomail.NewDialer(o.MailHost, o.MailPort, o.MailUser, o.MailPass)

	return d.DialAndSend(m)
}

// from: 标记是哪个业务申请的验证码
func GenMailCodeKey(mailAddr, from string) string {
	return "MAIL-CODE-" + from + "-" + mailAddr
}

func SendMailCode(mailOption *MailOptions, code, from string, ttl int) error {
	// 先插入redis，再发邮件
	key := GenMailCodeKey(mailOption.MailTo, from)
	// 如果code没有过期，是不允许再发送的
	success, err := redisDb.SetNX(ctx, key, code, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		fmt.Printf("SendMailCode redis fail %v\n", err)
		return err
	}

	if !success {
		return errCodeNotExpire
	}

	// 发邮件
	options := &MailOptions{
		MailHost: "smtp.163.com",
		MailPort: 465,
		MailUser: mailOption.MailUser,
		MailPass: mailOption.MailPass,
		MailTo:   mailOption.MailTo,
		Subject:  mailOption.Subject,
		Body:     mailOption.Body,
	}
	err = MailSend(options)
	if err != nil {
		return err
	}
	return nil
}

// from: 标记是哪个业务申请的验证码
func ValidateMailCode(mailAddr, code, from string) error {
	key := GenMailCodeKey(mailAddr, from)

	mailCode, err := redisDb.Get(ctx, key).Result()
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

	if mailCode != code {
		return errCodeNotMatch
	}

	return nil
}

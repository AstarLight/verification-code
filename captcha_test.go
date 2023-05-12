package main

import (
	"fmt"
	"testing"
)

func TestMakeCaptcha(t *testing.T) {
	codeLen := 6
	code, b64s, err := MakeCaptcha(codeLen)
	if err != nil {
		t.Errorf("expected:%v, got:%v", nil, err)
	}

	fmt.Println("code=", code)
	fmt.Println("b64s=", b64s)
}

func TestCaptchaCreateAndVerify(t *testing.T) {
	phone := "15913137596"
	codeLen := 6
	captchaTTL := 300
	from := "login" // 业务场景标记
	code, b64image, err := CreateCaptcha(phone, from, captchaTTL, codeLen)
	if err != nil {
		t.Errorf("expected:%v, got:%v", nil, err)
	}

	err = VerifyCaptchaCode(phone, code, from)
	if err != nil {
		t.Errorf("expected:%v, got:%v", nil, err)
	}

	fmt.Println("imageLen=", len(b64image))

}

package main

import (
	"fmt"
	"math/rand"

)

// 返回随机的n位 code
func GenRandomCode(n int) string {
	var code string
	for i := 0; i < n; i++ {
		r := rand.Intn(10)
		code += fmt.Sprintf("%d", r)
	}

	return code
}

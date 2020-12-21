package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

//用户名+Md5Key+当前时间后，用md5加密得到了cookie
func CookieGenerator(name string) string{
	s:=name+Md5Key+time.Now().String()
	ctx:=md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}
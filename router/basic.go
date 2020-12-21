package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	d "iot_admin/util/decode"
	"net/http"
	"strings"
	"time"
)

var loginCache *cache.Cache

const (
	KEY="123456789"
	COOKIEKEY="iot_admin_cookie"
)

//跨域中间件
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)                                    // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, If-Modified-Since, Cache-Control, Content-Type, Pragma,Cookie,signature")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Content-Disposition,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                   //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

type wrap struct{
	Data []byte `json:"data"`
}


//统一解密rsa
func unwrapAndDecodeRsa(c *gin.Context) ([]byte,error){
	var w wrap
	err:=c.BindJSON(&w)
	if err!=nil{
		return nil,err
	}
	if w.Data==nil{
		return nil,errors.New("empty data")
	}
	data,err:= d.BodyDecodeRSA(w.Data)
	return data,err
}

//检查签名，只限于get
func checkSignature(c *gin.Context) error{
	if c.GetHeader("signature")!= d.Signature {
		return errors.New("wrong signature")
	}
	return nil
}

//检查登录，适用所有非登录接口
func checkCookie(c *gin.Context) (*loginState,error){
	cookie,err:=c.Cookie(COOKIEKEY)
	if err!=nil{
		return nil,err
	}
	ls,ok:=loginCache.Get(cookie)
	if !ok{
		return nil,errors.New("wrong cookie")
	}
	return ls.(*loginState),nil
}

func NewRouter() *gin.Engine{
	r:=gin.Default()
	r.Use(cors())
	loginCache=cache.New(20*time.Hour,24*time.Hour)
	LoginRouter(r,"")
	DeviceRouter(r,"/operations/device")
	return r
}
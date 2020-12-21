package router

import (
	"encoding/json"
	g "github.com/GramYang/gylog"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	sc "iot_admin/sqlx_client"
	"iot_admin/util"
)

func LoginRouter(r *gin.Engine, base string){
	if base !=""{
		rg:=r.Group("/"+base)
		rg.POST("/login",Login)
	}else{
		r.POST("/login",Login)
	}
}

type login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type loginState struct{
	Name string
	Cookie string
	Other map[string]string
}

//返回状态码：成功、失败
func Login(c *gin.Context){
	data,err:=unwrapAndDecodeRsa(c)
	if err!=nil{
		g.Errorln(err)
		c.JSON(400,gin.H{
			"message":err.Error(),
		})
		return
	}
	var l login
	err=json.Unmarshal(data,&l)
	if err!=nil{
		g.Errorln(err)
		c.JSON(400,gin.H{
			"message":err.Error(),
		})
		return
	}
	a,err:=sc.GetAdminByName(l.UserName)
	if err!=nil{
		g.Errorln(err)
		c.JSON(500,gin.H{
			"message":err.Error(),
		})
		return
	}
	if a!=nil && a.Password==l.Password{
		cookie:=util.CookieGenerator(a.UserName)
		//保存loginState，期限为1天
		loginCache.Set(cookie,
			&loginState{Name:a.UserName,Cookie:cookie,Other:make(map[string]string)},
			cache.DefaultExpiration)
		//返回用户名
		c.JSON(200,gin.H{
			"userName":a.UserName,
			"cookie":cookie,
		})
		return
	}
	g.Errorln("password wrong")
	c.JSON(400,gin.H{
		"message":"password wrong",
	})
}
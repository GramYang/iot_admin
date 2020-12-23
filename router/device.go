package router

import (
	"encoding/base64"
	"encoding/json"
	g "github.com/GramYang/gylog"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	ct "iot_admin/ctwing_client"
	sc "iot_admin/sqlx_client"
	"iot_admin/util"
	"path"
	"strconv"
	"time"
)

const fileName = "114514"

func DeviceRouter(r *gin.Engine, base string) {
	if base != "" {
		rg := r.Group("/" + base)
		rg.POST("/create", DeviceCreate)
		rg.POST("/update", DeviceUpdate)
		rg.GET("/delete", DeviceDelete)
		rg.GET("/query", DeviceQuery)
		rg.GET("/queryList", DeviceListQuery)
		rg.GET("/queryLocalList", DeviceListLocal)
		rg.GET("/templateDownload", TemplateDownload)
		rg.POST("/templateParse", DevicesTemplateCreate)
		rg.POST("/templateUpload", TemplateUpload)
	} else {
		r.POST("/create", DeviceCreate)
		r.POST("/update", DeviceUpdate)
		r.GET("/delete", DeviceDelete)
		r.GET("/query", DeviceQuery)
		r.GET("/queryList", DeviceListQuery)
		r.GET("/queryLocalList", DeviceListLocal)
		r.GET("/templateDownload", TemplateDownload)
		r.POST("/templateParse", DevicesTemplateCreate)
		r.POST("/templateUpload", TemplateUpload)
	}
}

//创建设备
//返回状态码：成功、失败、未登录
func DeviceCreate(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	data, err := unwrapAndDecodeRsa(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	resp, err := ct.CreateDevice(string(data))
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode != 200 {
		c.JSON(500, gin.H{
			"message": "iot error code:" + strconv.Itoa(resp.StatusCode),
		})
		return
	}
	bodyStr, _ := ioutil.ReadAll(resp.Body)
	var mapResult map[string]interface{}
	err = json.Unmarshal(bodyStr, &mapResult)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": "iot response body error",
		})
		return
	}
	//直接返回结果，200并不代表成功
	c.JSON(200, gin.H{
		"code":    mapResult["code"],
		"message": mapResult["msg"],
	})
	//创建设备返回的信息并不完整，需要查询
	pid := strconv.Itoa(int(mapResult["result"].(map[string]interface{})["productId"].(float64)))
	did := mapResult["result"].(map[string]interface{})["deviceId"].(string)
	//只保存pid和did
	err = sc.SaveDevice(pid, did)
	if err != nil {
		g.Errorln(err)
		//出错了就写入本地文件，总之数据不能丢失
		util.WriteDeviceToLocal(pid, did)
	}
	resp1, err := ct.QuerySingleDevice(pid, did)
	if err != nil {
		g.Errorln(err)
		return
	}
	if resp1.StatusCode == 200 {
		bodyStr1, _ := ioutil.ReadAll(resp1.Body)
		//这里失败了没关系，还可以补充
		err = sc.SaveDeviceDetail(bodyStr1)
		if err != nil {
			g.Errorln(err.Error())
		}
	}
}

//模板批量创建设备
//返回状态码：成功、失败、未登录
func DevicesTemplateCreate(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	data, err := unwrapAndDecodeRsa(c)
	if err != nil {
		g.Errorln("unwrapAndDecodeRsa", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	g.Debugln("template raw json: ", string(data))
	//拆分json数组，获取json请求体和其int列号
	jsonMap, err := util.SplitJSONArrayAndCol(data)
	if err != nil {
		g.Errorln("SplitJSONArrayAndCol", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	//并发请求
	m := &util.ConMap{M: map[int]*util.ResponseWrap{}}
	g.Debugln("template parse map: ", jsonMap)
	util.ConcurrentOpt1(m, jsonMap, func(body string) *util.ResponseWrap {
		rw := &util.ResponseWrap{}
		resp, err := ct.CreateDevice(body)
		if err != nil {
			rw.Err = err
		}
		if resp != nil {
			rw.Status = resp.StatusCode
			bodyStr, _ := ioutil.ReadAll(resp.Body)
			rw.Body = string(bodyStr)
		}
		return rw
	})
	//查看并发请求结果
	errResult := map[int]string{}
	var successArr []int
	//结果分类
	for k, v := range m.M {
		g.Debugln("template create: ", k, v)
		//网络请求出错
		if v.Err != nil {
			errResult[k] = v.Err.Error()
		}
		//网络请求返回非200
		if v.Status != 200 {
			errResult[k] = v.Body
		}
		//网络请求成功，但是设备创建失败
		tmpMap := map[string]interface{}{}
		_ = json.Unmarshal([]byte(v.Body), &tmpMap)
		if tmpMap["code"] != 0 {
			errResult[k] = v.Body
		}
		successArr = append(successArr, k)
	}
	//如果出错了，就返回错误明细
	if len(errResult) > 0 {
		data, _ := json.Marshal(errResult)
		g.Debugln("创建失败" + string(data))
		c.JSON(200, gin.H{
			"message": "创建失败" + string(data),
		})
	} else {
		g.Debugln("创建成功" + strconv.Itoa(len(successArr)) + "个")
		c.JSON(200, gin.H{
			"message": "创建成功" + strconv.Itoa(len(successArr)) + "个",
		})
	}
	//保存创建成功的设备信息
	if len(successArr) == 0 {
		return
	}
	util.ConcurrentOpt2(m, successArr, func(col int, params map[string]interface{}) {
		if int(params["code"].(float64)) != 0 {
			g.Debugln("col " + strconv.Itoa(col) + " create error")
			return
		}
		g.Debugln("col " + strconv.Itoa(col) + " create ok")
		//创建设备返回的信息并不完整，需要查询
		pid := strconv.Itoa(int(params["result"].(map[string]interface{})["productId"].(float64)))
		did := params["result"].(map[string]interface{})["deviceId"].(string)
		g.Debugln("SaveDevice: ", pid, did)
		//只保存pid和did
		err = sc.SaveDevice(pid, did)
		if err != nil {
			g.Errorln(err)
			//出错了就写入本地文件，总之数据不能丢失
			util.WriteDeviceToLocal(pid, did)
		}
		resp1, err := ct.QuerySingleDevice(pid, did)
		g.Debugln("QuerySingleDevice: ", resp1)
		if err != nil {
			g.Errorln(err)
			return
		}
		if resp1.StatusCode == 200 {
			bodyStr1, _ := ioutil.ReadAll(resp1.Body)
			//这里失败了没关系，还可以补充
			g.Debugln("SaveDeviceDetail: ", string(bodyStr1))
			err = sc.SaveDeviceDetail(bodyStr1)
			if err != nil {
				g.Errorln("SaveDeviceDetail: ", err.Error())
			}
		}
	})
}

//更新设备
//返回状态码：成功、失败、未登录
func DeviceUpdate(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	data, err := unwrapAndDecodeRsa(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	resp, err := ct.UpdateDevice(data)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode != 200 {
		c.JSON(500, gin.H{
			"message": "iot error code:" + strconv.Itoa(resp.StatusCode),
		})
		return
	}
	bodyStr, _ := ioutil.ReadAll(resp.Body)
	var mapResult map[string]interface{}
	err = json.Unmarshal(bodyStr, &mapResult)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": "iot response body error",
		})
		return
	}
	//直接返回结果，200并不代表成功
	c.JSON(200, gin.H{
		"code":    mapResult["code"],
		"message": mapResult["msg"],
	})
}

//删除设备
//返回状态码：成功、失败、未登录
func DeviceDelete(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = checkSignature(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	params := c.Query("params")
	if params == "" {
		g.Errorln("missing params")
		c.JSON(400, gin.H{
			"message": "missing params",
		})
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(params)
	var mapResult1 map[string]interface{}
	_ = json.Unmarshal(decoded, &mapResult1)
	resp, err := ct.DeleteDevice(mapResult1["productId"].(string), mapResult1["deviceIds"].(string))
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode != 200 {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": "iot error code:" + strconv.Itoa(resp.StatusCode),
		})
		return
	}
	bodyStr, _ := ioutil.ReadAll(resp.Body)
	var mapResult map[string]interface{}
	err = json.Unmarshal(bodyStr, &mapResult)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": "iot response body error",
		})
		return
	}
	//直接返回结果，200并不代表成功
	c.JSON(200, gin.H{
		"code":    mapResult["code"],
		"message": mapResult["msg"],
	})
}

//查询单个设备
//返回状态码：成功、失败、未登录
func DeviceQuery(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = checkSignature(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	params := c.Query("params")
	if params == "" {
		g.Errorln("missing params")
		c.JSON(400, gin.H{
			"message": "missing params",
		})
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(params)
	var mapResult1 map[string]interface{}
	_ = json.Unmarshal(decoded, &mapResult1)
	resp, err := ct.QuerySingleDevice(mapResult1["productId"].(string), mapResult1["deviceId"].(string))
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode != 200 {
		g.Errorln("iot error code:" + strconv.Itoa(resp.StatusCode))
		c.JSON(500, gin.H{
			"message": "iot error code:" + strconv.Itoa(resp.StatusCode),
		})
		return
	}
	bodyStr, _ := ioutil.ReadAll(resp.Body)
	var mapResult2 map[string]interface{}
	err = json.Unmarshal(bodyStr, &mapResult2)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": "iot response body error",
		})
		return
	}
	//直接返回结果，200并不代表成功
	c.JSON(200, gin.H{
		"code":    mapResult2["code"],
		"message": mapResult2["msg"],
	})
}

//查询设备列表
//返回状态码：成功、失败、未登录
func DeviceListQuery(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = checkSignature(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	params := c.Query("params")
	if params == "" {
		g.Errorln("missing params")
		c.JSON(400, gin.H{
			"message": "missing params",
		})
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(params)
	var mapResult1 map[string]interface{}
	_ = json.Unmarshal(decoded, &mapResult1)
	resp, err := ct.QueryDevices(mapResult1["productId"].(string), mapResult1["searchValue"].(string),
		mapResult1["page"].(string))
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode != 200 {
		g.Errorln("iot error code:" + strconv.Itoa(resp.StatusCode))
		c.JSON(500, gin.H{
			"message": "iot error code:" + strconv.Itoa(resp.StatusCode),
		})
		return
	}
	bodyStr, _ := ioutil.ReadAll(resp.Body)
	var mapResult2 map[string]interface{}
	err = json.Unmarshal(bodyStr, &mapResult2)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": "iot response body error",
		})
		return
	}
	//直接返回结果，200并不代表成功
	c.JSON(200, gin.H{
		"code":    mapResult2["code"],
		"message": mapResult2["msg"],
		"result":  mapResult2["result"],
	})
}

//查询本地设备列表
//返回状态码：成功、失败、未登录
func DeviceListLocal(c *gin.Context) {
	_, err := checkCookie(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = checkSignature(c)
	if err != nil {
		g.Errorln(err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	params := c.Query("params")
	if params == "" {
		g.Errorln("missing params")
		c.JSON(400, gin.H{
			"message": "missing params",
		})
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(params)
	var mapResult1 map[string]interface{}
	_ = json.Unmarshal(decoded, &mapResult1)
	if mapResult1["deviceName"] == nil {
		mapResult1["deviceName"] = ""
	}
	if mapResult1["page"] == nil {
		c.JSON(400, gin.H{
			"message": "page invalid",
		})
		return
	}
	res, err := sc.QueryDeviceByName(mapResult1["deviceName"].(string), mapResult1["page"].(string))
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"result":  res,
		"pageNum": mapResult1["page"].(string),
	})
}

//下载创建设备模板xlsx，这是一个静态文件服务接口，get请求没有参数，也没必要认证，不登录也能下载。
func TemplateDownload(c *gin.Context) {
	g.Debugln("ctwing.xlsx download")
	fileName := "ctwing.xlsx"
	c.Header("Content-Disposition", "attachment;filename="+fileName)
	c.Header("Content-Type", "application/.nmsl")
	c.File("assets/" + fileName)
}

//上传模板xlsx
//返回状态码：成功、失败、未登录
func TemplateUpload(c *gin.Context) {
	file, _ := c.FormFile(fileName)
	if file != nil {
		g.Debugln("receive template: " + file.Filename)
		dst := path.Join("template/", file.Filename+time.Now().Format(".20060102.150405"))
		_ = c.SaveUploadedFile(file, dst)
	}
}

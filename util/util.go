package util

import (
	"encoding/json"
	"errors"
	g "github.com/GramYang/gylog"
	"sync"
)

//返回的map，key是deviceId，value是创建设备的json body
func SplitJSONArrayAndCol(data []byte) (map[int]string, map[string]interface{}, error) {
	var res = map[int]string{}
	var mapArr []map[string]interface{}
	_ = json.Unmarshal(data, &mapArr)
	g.Debugln("xlsx map", mapArr)
	//检查template的数据是否有整col空填的情况
	for i := 1; i < len(mapArr); i++ {
		if len(mapArr) <= 1 {
			return nil, nil, errors.New("missing col params")
		}
	}
	//保存额外的参数，比如秘钥
	var res1 = map[string]interface{}{}
	//公共参数提出来，要写入ctwing的请求体
	var operator, imsi, pskValue string
	var autoObserver, productId int
	t := mapArr[0]
	c := int(t["col"].(float64))
	//解析公共变量，在col 1中
	if c == 1 {
		if t["operator"] != nil {
			tmp, ok := t["operator"].(string)
			if ok {
				operator = tmp
			}
		}
		if t["imsi"] != nil {
			tmp, ok := t["imsi"].(string)
			if ok {
				imsi = tmp
			}
		}
		if t["pskValue"] != nil {
			tmp, ok := t["pskValue"].(string)
			if ok {
				pskValue = tmp
			}
		}
		if t["autoObserver"] != nil {
			tmp, ok := t["autoObserver"].(int)
			if ok {
				autoObserver = tmp
			}
		}
		if t["productId"] != nil {
			tmp, ok := t["productId"].(float64)
			if ok {
				productId = int(tmp)
			}
		}
		//secretKey存入res1
		if t["secretKey"] != nil {
			tmp, ok := t["secretKey"].(string)
			if ok {
				res1["secretKey"] = tmp
				delete(t, "secretKey")
			}
		}
	}
	for k, v := range mapArr {
		c := int(v["col"].(float64))
		delete(mapArr[k], "col")
		if operator != "" {
			v["operator"] = operator
		}
		if imsi != "" {
			v["imsi"] = imsi
		}
		if pskValue != "" {
			v["pskValue"] = pskValue
		}
		if autoObserver != 0 {
			v["autoObserver"] = autoObserver
		}
		if productId != 0 {
			v["productId"] = productId
		}
		data, _ := json.Marshal(v)
		res[c] = string(data)
	}
	g.Debugln("xlsx result", res)
	g.Debugln("xlsx result1", res1)
	return res, res1, nil
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

type ResponseWrap struct {
	Status int
	Body   string
	Err    error
}

//标记创建设备的是其excel记录中的列号
type ConMap struct {
	M  map[int]*ResponseWrap
	Mu sync.Mutex
}

func (c *ConMap) Set(k int, v *ResponseWrap) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.M[k] = v
}

func (c *ConMap) Get(k int) *ResponseWrap {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return c.M[k]
}

//请求的参数容器是一个map，回调函数的参数是一个string
func ConcurrentOpt1(m *ConMap, params map[int]string, f func(body string) *ResponseWrap) {
	var wg sync.WaitGroup
	for k, v := range params {
		k1 := k
		v1 := v
		wg.Add(1)
		go func() {
			rw := f(v1)
			m.Set(k1, rw)
			wg.Done()
		}()
	}
	wg.Wait()
}

//m1是并发请求结果，params是创建成功的设备信息col号数组，paraMap是额外的数据（比如要本地保存的设备秘钥，可以为空）
func ConcurrentOpt2(m1 *ConMap, params []int, paraMap map[string]interface{}, f func(col int, params map[string]interface{})) {
	var wg sync.WaitGroup
	add := false
	if len(paraMap) > 0 {
		add = true
	}
	for _, v := range params {
		v1 := v
		rw := m1.Get(v1)
		var mapRes map[string]interface{}
		_ = json.Unmarshal([]byte(rw.Body), &mapRes)
		if add {
			for k, v := range paraMap {
				mapRes[k] = v
			}
		}
		wg.Add(1)
		go func() {
			f(v1, mapRes)
			wg.Done()
		}()
	}
	wg.Wait()
}

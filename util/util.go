package util

import (
	"encoding/json"
	"errors"
	"sync"
)

//返回的map，key是deviceId，value是创建设备的json body
func SplitJSONArrayAndCol(data []byte) (map[int]string, error) {
	var res = map[int]string{}
	var mapArr []map[string]interface{}
	_ = json.Unmarshal(data, &mapArr)
	//检查template的数据是否有整列空填的情况，是否为矩阵
	rowLength := len(mapArr[0])
	var tmp int
	for i := 1; i < len(mapArr); i++ {
		tmp = max(tmp, len(mapArr[i]))
	}
	if tmp != rowLength {
		return nil, errors.New("it's not maxtrix")
	}
	for k, v := range mapArr {
		c := int(v["col"].(float64))
		delete(mapArr[k], "col")
		data, _ := json.Marshal(v)
		res[c] = string(data)
	}
	return res, nil
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

//有两个ConMap，从m1中抽取信息f处理后去写入m2，请求参数容器是一个slice，回调函数参数是一个map
func ConcurrentOpt2(m1 *ConMap, params []int, f func(col int, params map[string]interface{})) {
	var wg sync.WaitGroup
	for _, v := range params {
		v1 := v
		rw := m1.Get(v1)
		var mapRes map[string]interface{}
		_ = json.Unmarshal([]byte(rw.Body), &mapRes)
		wg.Add(1)
		go func() {
			f(v1, mapRes)
			wg.Done()
		}()
	}
	wg.Wait()
}

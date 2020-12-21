package aepapi

import (
    "apis/core"
    "net/http"
)


//参数body: 类型json, 参数不可以为空
//  描述:body,具体参考平台api说明
func BatchCreateNBDevice(appKey string, appSecret string, body string) (*http.Response, error) {
    path := "/aep_nb_device_management/batchNBDevice"
    
    var headers map[string]string = nil
    var param map[string]string = nil
    version := "20200828140355"

    application := appKey
    key := appSecret

    return aepsdkcore.SendAepHttpRequest(path, headers, param, body, version, application, key, "POST")
}


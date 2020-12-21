package aepapi

import (
    "apis/core"
    "net/http"
)


//参数body: 类型json, 参数不可以为空
//  描述:body,具体参考平台api说明
func QueryPublicByPublicProductId(appKey string, appSecret string, body string) (*http.Response, error) {
    path := "/aep_public_product_management/publicProducts"
    
    var headers map[string]string = nil
    var param map[string]string = nil
    version := "20190507003930"

    application := appKey
    key := appSecret

    return aepsdkcore.SendAepHttpRequest(path, headers, param, body, version, application, key, "POST")
}

//参数body: 类型json, 参数不可以为空
//  描述:body,具体参考平台api说明
func QueryPublicByProductId(appKey string, appSecret string, body string) (*http.Response, error) {
    path := "/aep_public_product_management/publicProductList"
    
    var headers map[string]string = nil
    var param map[string]string = nil
    version := "20190507004139"

    application := appKey
    key := appSecret

    return aepsdkcore.SendAepHttpRequest(path, headers, param, body, version, application, key, "POST")
}

//参数body: 类型json, 参数不可以为空
//  描述:body,具体参考平台api说明
func InstantiateProduct(appKey string, appSecret string, body string) (*http.Response, error) {
    path := "/aep_public_product_management/instantiateProduct"
    
    var headers map[string]string = nil
    var param map[string]string = nil
    version := "20200801233037"

    application := appKey
    key := appSecret

    return aepsdkcore.SendAepHttpRequest(path, headers, param, body, version, application, key, "POST")
}

//参数searchValue: 类型String, 参数可以为空
//  描述:设备类型、厂商ID、厂商名称
//参数pageNow: 类型long, 参数可以为空
//  描述:当前页数
//参数pageSize: 类型long, 参数可以为空
//  描述:每页记录数
func QueryAllPublicProductList(appKey string, appSecret string, searchValue string , pageNow string , pageSize string ) (*http.Response, error) {
    path := "/aep_public_product_management/allPublicProductList"
    
    var headers map[string]string = nil
    var param map[string]string = make(map[string]string)
    param["searchValue"] = searchValue
    param["pageNow"] = pageNow
    param["pageSize"] = pageSize

    version := "20200829005548"

    application := appKey
    key := appSecret

    return aepsdkcore.SendAepHttpRequest(path, headers, param, "", version, application, key, "GET")
}

//参数searchValue: 类型String, 参数可以为空
//  描述:产品名称
//参数pageNow: 类型long, 参数可以为空
//  描述:当前页数
//参数pageSize: 类型long, 参数可以为空
//  描述:每页记录数
//参数productIds: 类型String, 参数可以为空
//  描述:私有产品idList
func QueryMyPublicProductList(appKey string, appSecret string, searchValue string , pageNow string , pageSize string , productIds string ) (*http.Response, error) {
    path := "/aep_public_product_management/myPublicProductList"
    
    var headers map[string]string = nil
    var param map[string]string = make(map[string]string)
    param["searchValue"] = searchValue
    param["pageNow"] = pageNow
    param["pageSize"] = pageSize
    param["productIds"] = productIds

    version := "20200829005359"

    application := appKey
    key := appSecret

    return aepsdkcore.SendAepHttpRequest(path, headers, param, "", version, application, key, "GET")
}


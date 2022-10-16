package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

//检查自建应用获取token参数
func checkSelfToken(tokenInfo map[string]string){
	if tokenInfo["appKey"] == "" {
		panic("请填充token参数中的appKey信息")
	}
	if tokenInfo["appSecret"] == "" {
		panic("请填充token参数中的appSecret信息")
	}
}

//检查生态应用获取token参数
func checkSuiteToken(tokenInfo map[string]string){
	if tokenInfo["suiteKey"] == "" {
		panic("请填充token参数中的appKey信息")
	}
	if tokenInfo["suiteSecret"] == "" {
		panic("请填充token参数中的appSecret信息")
	}
	if tokenInfo["tenantId"] == "" {
		panic("请填充token参数中的tenantId信息")
	}
}
//检查常规参数
func checkSplitParam(httpMethod string, hostUrl string,businessUrl string){
	if httpMethod == "" {
		panic("请填充token参数中的appKey信息")
	}
	if !strings.HasPrefix(hostUrl,"http"){
		panic("请使用以http协议开头的hostUrl")
	}
	if businessUrl== "" {
		panic("请填充businessUrl")
	}
}

// OptSelfRequest 使用token信息获取accessToken后执行自建应用请求
func OptSelfRequest(httpMethod string, hostUrl string,businessUrl string, params map[string]string, header map[string]string, data map[string]interface{}, tokenInfo map[string]string) string {
	checkSelfToken(tokenInfo)
	checkSplitParam(httpMethod , hostUrl ,businessUrl )
	accessToken := OptSelfTokenWithCache(tokenInfo["appKey"], tokenInfo["appSecret"], hostUrl)
	requestUrl := hostUrl +"/iuap-api-gateway" + businessUrl
	method := strings.ToUpper(httpMethod)
	if method == "POST" {
		return OptRequestPostWithAccessToken(requestUrl, params, header, data, accessToken)
	}else if method == "GET" {
		return OptRequestGetWithAccessToken(requestUrl, params, header, accessToken)
	}else{
		panic("该sdk只支持post 和get请求方式")
	}
}

// OptSuiteRequest 使用token信息获取accessToken后执行生态应用请求
func OptSuiteRequest(httpMethod string,hostUrl string,businessUrl string,  params map[string]string, header map[string]string, data map[string]interface{}, tokenInfo map[string]string) string {
	checkSuiteToken(tokenInfo)
	checkSplitParam(httpMethod , hostUrl ,businessUrl )
	accessToken := OptSuiteTokenWithCache(tokenInfo["suiteKey"], tokenInfo["tenantId"],tokenInfo["suiteSecret"], hostUrl)
    requestUrl := hostUrl +"/iuap-api-gateway" + businessUrl
	method := strings.ToUpper(httpMethod)
	if method == "POST" {
		return OptRequestPostWithAccessToken(requestUrl, params, header, data, accessToken)
	}else if method == "GET" {
		return OptRequestGetWithAccessToken(requestUrl, params, header, accessToken)
	}else{
		panic("该sdk只支持post 和get请求方式")
	}
}

// OptRequestGetWithAccessToken 执行Get操作
func OptRequestGetWithAccessToken(requestUrl string, params map[string]string, header map[string]string, accessToken string) string {
	urlParam := ""
	if params == nil {
		params = make(map[string]string)
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["access_token"] = accessToken
	client := &http.Client{}
	req, err := http.NewRequest("get", requestUrl, strings.NewReader(urlParam))
	if err != nil {
		panic(err)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k,v :=range params{
		q.Add(k, v)
	}
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	return string(body)
}

// OptRequestPostWithAccessToken 执行Post操作
func OptRequestPostWithAccessToken(requestUrl string, params map[string]string, header map[string]string, data map[string]interface{}, accessToken string) string {
	if params == nil {
		params = make(map[string]string)
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["access_token"] = accessToken
	jsonData, _ := json.Marshal(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	fmt.Println(req.Header.Get("Content-Type"))
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k,v :=range params{
		q.Add(k, v)
	}
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

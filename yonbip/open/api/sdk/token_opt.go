package sdk

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"github.com/tmzjlyw/yonyou-openapi-sdk/yonbip/open/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var cache = freecache.NewCache(10 * 1024 * 1024)

type tokenReqInfo struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    tokenInfo `json:"data"`
}

type tokenInfo struct {
	AccessToken string `json:"access_token"`
	Expire      int    `json:"expire"`
}

//获取
func OptSelfToken(appKey string, appSecret string, hostUrl string) string {
	tokenParamDic := make(map[string]string)
	tokenParamDic["appKey"] = appKey
	tokenParamDic["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	return string(OptTokenRequest(hostUrl+"/iuap-api-auth/open-auth/selfAppAuth/getAccessToken", tokenParamDic, appSecret))
}

func OptSuiteToken(suiteKey string, tenantId string, appSecret string, hostUrl string) string {
	tokenParamDic := make(map[string]string)
	tokenParamDic["suiteKey"] = suiteKey
	tokenParamDic["tenantId"] = tenantId
	tokenParamDic["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	return string(OptTokenRequest(hostUrl+"/iuap-api-auth/open-auth/suiteApp/getAccessToken", tokenParamDic, appSecret))

}
func OptSelfTokenWithCache(appKey string, appSecret string, hostUrl string) string {
	key := appKey
	value, _ := cache.Get([]byte(key))
	if value != nil {
		return string(value)
	}
	tokenParamDic := make(map[string]string)
	tokenParamDic["appKey"] = appKey
	tokenParamDic["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	result := OptTokenRequest(hostUrl+"/iuap-api-auth/open-auth/selfAppAuth/getAccessToken", tokenParamDic, appSecret)
	tokenResult := tokenReqInfo{}
	err := json.Unmarshal(result, &tokenResult)
	if err == nil {
		accessToken := tokenResult.Data.AccessToken
		cache.Set([]byte(key), []byte(accessToken), tokenResult.Data.Expire/1000-10)
		return accessToken
	}
	panic(err)
}

func OptSuiteTokenWithCache(suiteKey string, tenantId string, suiteSecret string, hostUrl string) string {
	key := suiteKey + tenantId
	value, _ := cache.Get([]byte(key))
	if value != nil {
		return string(value)
	}
	tokenParamDic := make(map[string]string)
	tokenParamDic["suiteKey"] = suiteKey
	tokenParamDic["tenantId"] = tenantId
	tokenParamDic["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	result := OptTokenRequest(hostUrl+"/iuap-api-auth/open-auth/suiteApp/getAccessToken", tokenParamDic, suiteSecret)
	tokenResult := tokenReqInfo{}
	err := json.Unmarshal(result, &tokenResult)
	if err == nil {
		accessToken := tokenResult.Data.AccessToken
		cache.Set([]byte(key), []byte(accessToken), tokenResult.Data.Expire/1000-10)
		return accessToken
	}
	panic(err)
}

func OptTokenRequest(requestUrl string, tokenParamDic map[string]string, appSecret string) []byte {
	sortData := utils.SortTokenParam(tokenParamDic)
	sign := utils.EncoderSha256(appSecret, sortData)
	tokenParamDic["signature"] = sign
	uri, err := url.Parse(requestUrl)
	if err != nil {
		panic(err)
	}
	values := url.Values{} //拼接query参数
	for k, v := range tokenParamDic {
		values.Add(k, v)
	}
	uri.RawQuery = values.Encode()
	urlPath := uri.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"sort"
)

func EncoderSha256(key string, data string) string {
	m := hmac.New(sha256.New, []byte(key))
	_,_ =m.Write([]byte(data))
	sum := m.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}


func SortListParam(data []string) string {
	sort.Strings(data)
	result :=""
	for  _, key := range  data {
		result += key
	}
	return result
}


func SortTokenParam(tokenInfo map[string]string) string {
	var keys []string
	result := ""
	for key, _ := range tokenInfo {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for  _, key := range  keys {
		result += key + tokenInfo[key]
	}
	return result
}

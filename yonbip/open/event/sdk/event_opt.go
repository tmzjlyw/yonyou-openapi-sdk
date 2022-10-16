package sdk

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"strings"
	"yonyou-openapi-sdk/yonbip/open/utils"
)

func DecryptEventEncrypt(appSecret string, holder map[string]string) string {
	holderInfoList := []string{holder["nonce"], holder["encrypt"], holder["timestamp"]}
	sortedData := utils.SortListParam(holderInfoList)
	sha256 := utils.EncoderSha256(appSecret, sortedData)
	secret := buildAesKeyFromSecret(appSecret)+"="
	secretBytes, _ :=  base64.StdEncoding.DecodeString(secret)
	encrypt, _ :=  base64.StdEncoding.DecodeString(holder["encrypt"])
	if sha256 == holder["signature"] {
		return AesDecryptCBC(encrypt, secretBytes)
	}
	panic("验签失败")

}

func AesDecryptCBC(encrypted []byte, key []byte) string {
	block, er := aes.NewCipher(key)        // 分组秘钥
	if er != nil {
		panic(er)
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	//encrypted = MakeBlocksFull(encrypted,blockSize)
	decrypted := make([]byte, len(encrypted))                   // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	padding := pkcs5UnPadding(decrypted)
	dataMapJson , _ := json.Marshal(padding)
	dataString := string(dataMapJson)
	return dataString

}

func pkcs5UnPadding(origData []byte) map[string]string {
    le := len(origData)
	padIndex := int(string(origData)[le-2])
	if padIndex < 1 || padIndex > 32{
		padIndex = 0
	}

	origData = origData[0:le -padIndex]
	// 分离16位随机字符串, 网络字节序和corpId
	networkOrder := origData[16:20]
	xmlLength := recoverNetworkBytesOrder(networkOrder)
	plainInfo := origData[20:20 + xmlLength]
	suitKey := origData[20 + xmlLength:]
	return map[string]string {"event_info": string(plainInfo), "suit_key": string(suitKey)}
}

func buildAesKeyFromSecret(appSecret string) string {
	encodingAesKey := strings.Replace(appSecret, "-", "", -1)
	if len(encodingAesKey) == 43 {
		return encodingAesKey
	}
	if len(encodingAesKey) > 43 {
		return encodingAesKey[:43]
	}
	for {
		if len(encodingAesKey) >= 43 {
			break
		}
		encodingAesKey = encodingAesKey + "0"
	}
	return encodingAesKey
}

func recoverNetworkBytesOrder(order_bytes []byte) int {
	source_number := 0

	j := 0
	for ; j < 4; {
		source_number <<= 8
		source_number |= int(order_bytes[j]) & 0xff
		j++
	}
	return source_number

}


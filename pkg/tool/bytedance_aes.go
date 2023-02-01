package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strings"
)

func byteDanceDecrypt(key, iv []byte, text string) (string, error) {
	decodeData, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", nil
	}
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//输出到[]byte数组
	originData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)
	//去除填充,并返回
	return string(originData[:16]), nil
}

// FastAesKey 更加keyToken获取加密key然后解密
// https://api.juejin.cn/user_api/v1/video/key_token?aid=2608&uuid=6898099958165407246&spider=0
// https://kds.bytedance.com/kds/api/v3/keys?source=jarvis&ak=6320731a6c39fc14ab44e6c50102c65b&token=HMAC-SHA1%3A1.0%3A1675240485%3AAKLTNmEwYWEzZmJhMDE0NDUyYTk1MThiYTk2NjQ4MmY1ZTk%3Aoku5E9Q%2BQ6Tsf5DGJWjDSDxuy9U%3D&_=1675236886731
func FastAesKey(data string) (string, error) {
	ds := strings.Split(data, ":")
	kv := ds[0]
	key := kv[16:] + kv[:16]
	iv := kv[16:]
	encData := ds[1]
	return byteDanceDecrypt([]byte(key), []byte(iv), encData)
}

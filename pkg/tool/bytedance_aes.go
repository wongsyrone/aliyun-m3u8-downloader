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

// PlayAuthDecrypt 火山引擎视频云 playAuth 解密
// 参考地址：https://www.52pojie.cn/thread-1726219-1-1.html
func PlayAuthDecrypt(playAuth string) string {
	a1, _ := base64.StdEncoding.DecodeString(playAuth)
	a2 := len(a1)
	var v6, v7, v8, v9, v10, v11 int
	v6 = 0
	v7 = 0
	v8 = 0
	v9 = 0
	v10 = 0
	v11 = 0

	if a2 >= 3 {
		v6 = 1
		v7 = int(a1[0] ^ a1[1] ^ a1[2])
		v9 = a2 - v7
		v11 = v9 + 47
		if v7-48 < 1 {
			v10 = 1
		}
		if v11 < 1 {
			v8 = 1
		}
	}
	if v8|v10 == 0 {
		v12 := make([]byte, v11)
		v13 := v7 - 47
		for i := 0; i < v13; i++ {
		}
		for i := 0; i < v11; i++ {
			v12[i] = a1[i+v6]
		}
		v15 := 0
		v16 := -6
		v17 := byte(85)
		for v15 != v11 {
			v18 := v12[v15]
			v19 := 0
			v20 := v15
			v21 := v18
			if v15&1 == 0 {
				v21 = v17
			}
			for v20 > 0 {
				v19++
				v20 &= v20 - 1
			}
			if v15&1 == 0 {
				v17 = byte(v16)
			}
			v22 := (v18 ^ v17) - byte(v19)
			v17 = v21
			v12[v15] = byte(v22 - 21)
			if v15&1 == 0 {
				v16 = int(v18)
			}
			v15++
		}
		return string(v12[1 : v11-1])
	}
	return ""
}

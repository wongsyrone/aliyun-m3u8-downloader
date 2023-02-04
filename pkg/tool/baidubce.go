package tool

import (
	"crypto/aes"
	"encoding/hex"
)

// BaiduDecrypt 百度智能云key解密
// 参考： http://aqxbk.com/archives/security/securitysite/2021/08_02/62815.html
func BaiduDecrypt(key, src string) (string, error) {
	encrypted, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	return string(decrypted), nil
}

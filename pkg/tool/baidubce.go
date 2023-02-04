package tool

import (
	"crypto/aes"
	"encoding/hex"
)

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

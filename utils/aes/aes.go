package aes

import (
	"bytes"
	cryptoAes "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Aes struct {
	key string
	iv  string
}

func New(key, iv string) *Aes {
	return &Aes{
		key: key,
		iv:  iv,
	}
}

func (a *Aes) i() {}

func (a *Aes) Encrypt(encryptStr string) (string, error) {
	encryptBytes := []byte(encryptStr)
	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes = pkcs5Padding(encryptBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(a.iv))
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func (a *Aes) Decrypt(decryptStr string) (string, error) {
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)
	if err != nil {
		return "", err
	}

	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(a.iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted = pkcs5UnPadding(decrypted)
	return string(decrypted), nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}

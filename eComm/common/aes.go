package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
)

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func MakeIv() []byte {
	iv := make([]byte, aes.BlockSize)
	if n, err := rand.Read(iv); n != aes.BlockSize || err != nil {
		log.Fatal(n, err)
	}
	return iv
}

func Aes256(plaintext []byte, key string, bIV []byte, blockSize int) []byte {
	bKey := []byte(key)
	bPlaintext := PKCS5Padding(plaintext, blockSize, len(plaintext))
	block, _ := aes.NewCipher(bKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return ciphertext
}

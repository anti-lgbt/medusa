package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncryptPassword(data []byte) string {
	block, _ := aes.NewCipher([]byte(createHash("ngulamdjtmemay")))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	base64_str := base64.StdEncoding.EncodeToString(ciphertext)

	return base64_str
}

func DecryptPassword(password string) []byte {
	data, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		panic(err.Error())
	}

	key := []byte(createHash("ngulamdjtmemay"))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func main() {
	password := "ha09031192"

	encrypt := EncryptPassword([]byte(password))

	base64_encrypt := base64.StdEncoding.EncodeToString(encrypt)

	log.Println(base64_encrypt)
	nn, _ := base64.StdEncoding.DecodeString(base64_encrypt)
	decrypt := DecryptPassword(nn)
	log.Println(string(decrypt))
}

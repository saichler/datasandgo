package securityutil

import (
	"crypto/aes"
	"io"
	"crypto/cipher"
	"encoding/base64"
	"crypto/rand"
	"log"
	"errors"
)

type SecurityKey struct {
	key []byte
}

func (key *SecurityKey)Enc(data []byte) ([]byte, error) {
	//@TODO manage keys per network node uuid
	key.key = []byte("12345678901234567890123456789012")
	block, err := aes.NewCipher(key.key)
	if err !=nil{
		log.Fatal("Failed to load encryption!", err)
	}

	b := base64.StdEncoding.EncodeToString(data)
	cipherdata := make([]byte, aes.BlockSize+len(b))

	iv := cipherdata[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err!=nil {
		log.Fatal("Failed to encrypt data")
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherdata[aes.BlockSize:], []byte(b))
	return cipherdata, nil
}

func (key *SecurityKey) Dec(encData []byte) ([]byte, error) {
	if len(encData) < aes.BlockSize {
		return nil, errors.New("Not an encrypted data")
	}
	key.key = []byte("12345678901234567890123456789012")
	block, err := aes.NewCipher(key.key)
	if err !=nil{
		log.Fatal("Failed to load encryption!", err)
	}
	iv := encData[:aes.BlockSize]
	encData = encData[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(encData, encData)
	data, err := base64.StdEncoding.DecodeString(string(encData))
	if err != nil {
		return nil, err
	}
	return data, nil
}

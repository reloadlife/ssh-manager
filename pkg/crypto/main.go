package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"io"
)

type Crypto struct {
	privateKey []byte
	publicKey  []byte
}

func NewCrypto(username, password string) *Crypto {
	sha := sha512.New()
	password = base64.StdEncoding.EncodeToString(sha.Sum([]byte(password)))
	c := &Crypto{
		privateKey: []byte(password),
		publicKey:  []byte(base64.StdEncoding.EncodeToString([]byte(username))),
	}
	return c
}

func (c *Crypto) cypherKey() []byte {
	var b []byte

	for i, by := range c.privateKey {
		b = append(b, by)
		for j, pb := range c.publicKey {
			b = append(b, pb)
			if j > i {
				break
			}
		}
	}
	return to32Bytes(b)
}

func to32Bytes(key []byte) []byte {
	if len(key) < 32 {
		return append(key, to32Bytes(key)...)
	}
	return key[:32]
}

func (c *Crypto) Encode(msg []byte) []byte {
	block, err := aes.NewCipher(c.cypherKey())
	if err != nil {
		panic(err)
	}
	cipherText := make([]byte, aes.BlockSize+len(msg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], msg)
	return []byte(base64.RawStdEncoding.EncodeToString(cipherText))
}

func (c *Crypto) Decode(msg []byte) []byte {
	cipherText, err := base64.RawStdEncoding.DecodeString(string(msg))
	if err != nil {
		return nil
	}
	block, err := aes.NewCipher(c.cypherKey())
	if err != nil {
		return nil
	}
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return nil
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return cipherText
}

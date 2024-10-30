package cicd

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

const (
	AES_KEY = "12345678abcdefgh"
	AES_IV  = "abcdefgh12345678"
)

// CTR 128bit no padding
func AesEny(plaintext []byte) []byte {
	var (
		block cipher.Block
		err   error
	)
	if block, err = aes.NewCipher([]byte(AES_KEY)); err != nil {
		log.Fatal(err)
	}
	stream := cipher.NewCTR(block, []byte(AES_IV))
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext
}

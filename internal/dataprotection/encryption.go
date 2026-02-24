package dataprotection

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func EncryptField(plain, key []byte) (string, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(blk)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)
	cipherText := gcm.Seal(nonce, nonce, plain, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptField(encoded string, key []byte) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blk)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	nonce, c := blob[:ns], blob[ns:]
	return gcm.Open(nil, nonce, c, nil)
}

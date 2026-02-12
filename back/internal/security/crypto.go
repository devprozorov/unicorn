package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type AESCrypt struct {
	key []byte
}

func NewAESCrypt(keyB64 string) (*AESCrypt, error) {
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("TOTP_ENC_KEY_B64 must decode to 32 bytes")
	}
	return &AESCrypt{key: key}, nil
}

func (c *AESCrypt) EncryptToB64(plain []byte) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil { return "", err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return "", err }
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return "", err }
	ct := gcm.Seal(nil, nonce, plain, nil)
	out := append(nonce, ct...)
	return base64.StdEncoding.EncodeToString(out), nil
}

func (c *AESCrypt) DecryptFromB64(encB64 string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encB64)
	if err != nil { return nil, err }
	block, err := aes.NewCipher(c.key)
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	ct := data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

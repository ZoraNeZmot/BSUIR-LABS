package algos

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"oop/lab5/internal/funcplugins"
)

// aesCFB is a real symmetric cipher: the user-supplied passphrase is
// stretched through SHA-256 to a 32-byte key, the 16-byte IV is
// generated fresh on every Encode and prepended to the payload before
// the result is base64-encoded.
type aesCFB struct{}

func (aesCFB) ID() string          { return "aes-cfb" }
func (aesCFB) DisplayName() string { return "AES-256-CFB" }
func (aesCFB) Description() string {
	return "AES-256 in CFB mode. Key is SHA-256 of passphrase; IV is random and prepended."
}

func (aesCFB) Parameters() []funcplugins.ParamSpec {
	return []funcplugins.ParamSpec{
		{Name: "passphrase", Label: "Passphrase", Kind: funcplugins.ParamSecret, Default: "change me"},
	}
}

func (aesCFB) Encode(data []byte, p map[string]string) ([]byte, error) {
	pass := p["passphrase"]
	if pass == "" {
		return nil, errors.New("aes: passphrase must not be empty")
	}
	key := sha256.Sum256([]byte(pass))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	out := make([]byte, aes.BlockSize+len(data))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], data)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(out)))
	base64.StdEncoding.Encode(encoded, out)
	return encoded, nil
}

func (aesCFB) Decode(data []byte, p map[string]string) ([]byte, error) {
	pass := p["passphrase"]
	if pass == "" {
		return nil, errors.New("aes: passphrase must not be empty")
	}
	raw := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(raw, data)
	if err != nil {
		return nil, err
	}
	raw = raw[:n]
	if len(raw) < aes.BlockSize {
		return nil, errors.New("aes: ciphertext too short")
	}
	key := sha256.Sum256([]byte(pass))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	iv := raw[:aes.BlockSize]
	cipherText := raw[aes.BlockSize:]
	plain := make([]byte, len(cipherText))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plain, cipherText)
	return plain, nil
}

func init() { funcplugins.RegisterAlgorithm(aesCFB{}) }

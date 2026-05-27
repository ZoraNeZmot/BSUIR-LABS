// Package algos hosts the built-in encryption algorithms used by the
// lab-5 functional-plugin system.
package algos

import (
	"encoding/base64"
	"errors"

	"oop/lab6/internal/funcplugins"
)

// xorCipher implements a straightforward repeating-key XOR cipher.
// Output is base64-encoded so that the (otherwise binary) ciphertext
// can still be stored next to the human-readable text format.
type xorCipher struct{}

func (xorCipher) ID() string          { return "xor" }
func (xorCipher) DisplayName() string { return "XOR Cipher" }
func (xorCipher) Description() string {
	return "Repeating-key XOR with base64 framing. Fast, symmetric, weak вЂ” fine for a lab assignment."
}

func (xorCipher) Parameters() []funcplugins.ParamSpec {
	return []funcplugins.ParamSpec{
		{Name: "key", Label: "Key", Kind: funcplugins.ParamSecret, Default: "secret"},
	}
}

func (xorCipher) Encode(data []byte, p map[string]string) ([]byte, error) {
	key := p["key"]
	if key == "" {
		return nil, errors.New("xor: key must not be empty")
	}
	cipher := xorBytes(data, []byte(key))
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(cipher)))
	base64.StdEncoding.Encode(dst, cipher)
	return dst, nil
}

func (xorCipher) Decode(data []byte, p map[string]string) ([]byte, error) {
	key := p["key"]
	if key == "" {
		return nil, errors.New("xor: key must not be empty")
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(dst, data)
	if err != nil {
		return nil, err
	}
	return xorBytes(dst[:n], []byte(key)), nil
}

// xorBytes returns src XOR-ed with a repeating key.
func xorBytes(src, key []byte) []byte {
	out := make([]byte, len(src))
	for i := range src {
		out[i] = src[i] ^ key[i%len(key)]
	}
	return out
}

func init() { funcplugins.RegisterAlgorithm(xorCipher{}) }

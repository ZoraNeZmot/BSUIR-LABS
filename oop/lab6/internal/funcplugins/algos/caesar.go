package algos

import (
	"errors"
	"strconv"

	"oop/lab6/internal/funcplugins"
)

// caesarCipher shifts every byte of the input by a configured amount.
// Unlike the classic letter-only Caesar this one operates on the full
// byte range so it is reversible regardless of the input alphabet.
type caesarCipher struct{}

func (caesarCipher) ID() string          { return "caesar" }
func (caesarCipher) DisplayName() string { return "Caesar (byte) Cipher" }
func (caesarCipher) Description() string {
	return "Adds the configured shift modulo 256 to every byte. Trivially reversible."
}

func (caesarCipher) Parameters() []funcplugins.ParamSpec {
	return []funcplugins.ParamSpec{
		{Name: "shift", Label: "Shift", Kind: funcplugins.ParamInt, Default: "7"},
	}
}

func (caesarCipher) Encode(data []byte, p map[string]string) ([]byte, error) {
	shift, err := parseShift(p["shift"])
	if err != nil {
		return nil, err
	}
	return shiftBytes(data, shift), nil
}

func (caesarCipher) Decode(data []byte, p map[string]string) ([]byte, error) {
	shift, err := parseShift(p["shift"])
	if err != nil {
		return nil, err
	}
	return shiftBytes(data, -shift), nil
}

func parseShift(s string) (int, error) {
	if s == "" {
		return 0, errors.New("caesar: shift must not be empty")
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func shiftBytes(src []byte, shift int) []byte {
	mod := ((shift % 256) + 256) % 256
	out := make([]byte, len(src))
	for i, b := range src {
		out[i] = byte((int(b) + mod) % 256)
	}
	return out
}

func init() { funcplugins.RegisterAlgorithm(caesarCipher{}) }

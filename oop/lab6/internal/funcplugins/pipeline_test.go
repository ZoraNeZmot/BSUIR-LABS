package funcplugins_test

import (
	"bytes"
	"testing"

	"oop/lab6/internal/funcplugins"
	_ "oop/lab6/internal/funcplugins/algos"
)

// TestPipelineRoundTrip stacks every built-in algorithm (XOR, Caesar,
// AES, friend-checksum adapter) and checks that Decode is the inverse
// of Encode. It also exercises the Adapter pattern end-to-end because
// the friend-checksum adapter is part of the chain.
func TestPipelineRoundTrip(t *testing.T) {
	plugins := []*funcplugins.FuncPlugin{
		{ID: "p-caesar", Name: "caesar", AlgorithmID: "caesar", Parameters: map[string]string{"shift": "11"}, Enabled: true},
		{ID: "p-xor", Name: "xor", AlgorithmID: "xor", Parameters: map[string]string{"key": "lab6"}, Enabled: true},
		{ID: "p-aes", Name: "aes", AlgorithmID: "aes-cfb", Parameters: map[string]string{"passphrase": "demo"}, Enabled: true},
		{ID: "p-chk", Name: "friend-checksum", AlgorithmID: "friend-checksum", Parameters: map[string]string{"algorithm": "sha256"}, Enabled: true},
	}
	for _, p := range plugins {
		if err := funcplugins.AddPlugin(p); err != nil {
			t.Fatalf("AddPlugin %s: %v", p.ID, err)
		}
	}

	original := []byte("[BEGIN Car]\nID=v-1\nManufacturer=Toyota\n[END]\n")
	encoded, err := funcplugins.Encode(original)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if bytes.Equal(encoded, original) {
		t.Fatalf("Encode produced unchanged bytes; pipeline is a no-op")
	}
	decoded, err := funcplugins.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !bytes.Equal(decoded, original) {
		t.Fatalf("round-trip mismatch:\n  original=%q\n  decoded =%q", original, decoded)
	}

	// Sanity: tampering with the encoded payload must be caught by the
	// checksum adapter on the way back.
	if len(encoded) > 0 {
		encoded[len(encoded)-1] ^= 0xFF
	}
	if _, err := funcplugins.Decode(encoded); err == nil {
		t.Fatalf("expected Decode to reject tampered payload")
	}
}

package bencode

import "testing"

func TestMarshalDictionaryOrder(t *testing.T) {
	got, err := Marshal(map[string]any{
		"z": int64(1),
		"a": "x",
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	want := "d1:a1:x1:zi1ee"
	if string(got) != want {
		t.Fatalf("Marshal() = %q, want %q", string(got), want)
	}
}

func TestMarshalList(t *testing.T) {
	got, err := Marshal([]any{"abc", int64(42)})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(got) != "l3:abci42ee" {
		t.Fatalf("Marshal() = %q", string(got))
	}
}

package storage_test

import (
	"bytes"
	"testing"

	"oop/lab3/internal/storage"
	"oop/lab3/internal/vehicle"
	_ "oop/lab3/internal/vehicles"
)

// TestTextRoundTrip verifies that Marshal/Unmarshal are inverse on a
// representative sample of every concrete class in the hierarchy.
func TestTextRoundTrip(t *testing.T) {
	names := vehicle.Names()
	if len(names) == 0 {
		t.Fatal("registry is empty -- vehicles package not imported?")
	}
	src := make([]vehicle.Vehicle, 0, len(names))
	for i, name := range names {
		v, err := vehicle.Create(name)
		if err != nil {
			t.Fatalf("Create(%q): %v", name, err)
		}
		// Set a couple of fields so the test exercises the value
		// encoder for every primitive kind.
		for _, f := range v.Fields() {
			switch f.Name {
			case "ID":
				_ = f.Set("v-" + itoa(i))
			case "Manufacturer":
				_ = f.Set("Acme")
			case "Year":
				_ = f.Set("2024")
			}
		}
		src = append(src, v)
	}

	var buf bytes.Buffer
	if err := storage.Marshal(&buf, src); err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	dst, err := storage.Unmarshal(&buf)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(dst) != len(src) {
		t.Fatalf("count mismatch: got %d want %d", len(dst), len(src))
	}
	for i := range src {
		if dst[i].TypeName() != src[i].TypeName() {
			t.Fatalf("item #%d type: got %s want %s", i, dst[i].TypeName(), src[i].TypeName())
		}
		gotFields := dst[i].Fields()
		wantFields := src[i].Fields()
		if len(gotFields) != len(wantFields) {
			t.Fatalf("item #%d field count: got %d want %d", i, len(gotFields), len(wantFields))
		}
		for j := range wantFields {
			if gotFields[j].Get() != wantFields[j].Get() {
				t.Fatalf("item #%d field %s: got %q want %q",
					i, wantFields[j].Name, gotFields[j].Get(), wantFields[j].Get())
			}
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

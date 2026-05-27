package bencode

import "testing"

func TestUnmarshalInt(t *testing.T) {
	v, n, err := Unmarshal([]byte("i42e"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatalf("consumed %d", n)
	}
	if v.(int64) != 42 {
		t.Fatalf("value %v", v)
	}
}

func TestUnmarshalDict(t *testing.T) {
	data := []byte("d3:foo3:bar3:numi99ee")
	v, _, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	d := v.(map[string]any)
	if string(d["foo"].([]byte)) != "bar" {
		t.Fatalf("foo=%v", d["foo"])
	}
	if d["num"].(int64) != 99 {
		t.Fatalf("num=%v", d["num"])
	}
}

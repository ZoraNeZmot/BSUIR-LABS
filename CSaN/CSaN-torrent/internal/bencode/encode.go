package bencode

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MustMarshal(v any) []byte {
	out, err := Marshal(v)
	if err != nil {
		panic(err)
	}
	return out
}

func encode(buf *bytes.Buffer, v any) error {
	switch x := v.(type) {
	case string:
		writeString(buf, []byte(x))
	case []byte:
		writeString(buf, x)
	case int:
		writeInt(buf, int64(x))
	case int64:
		writeInt(buf, x)
	case uint64:
		writeInt(buf, int64(x))
	case []any:
		buf.WriteByte('l')
		for _, item := range x {
			if err := encode(buf, item); err != nil {
				return err
			}
		}
		buf.WriteByte('e')
	case map[string]any:
		buf.WriteByte('d')
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			writeString(buf, []byte(k))
			if err := encode(buf, x[k]); err != nil {
				return err
			}
		}
		buf.WriteByte('e')
	default:
		return fmt.Errorf("unsupported bencode type %T", v)
	}
	return nil
}

func writeString(buf *bytes.Buffer, value []byte) {
	buf.WriteString(strconv.Itoa(len(value)))
	buf.WriteByte(':')
	buf.Write(value)
}

func writeInt(buf *bytes.Buffer, value int64) {
	buf.WriteByte('i')
	buf.WriteString(strconv.FormatInt(value, 10))
	buf.WriteByte('e')
}

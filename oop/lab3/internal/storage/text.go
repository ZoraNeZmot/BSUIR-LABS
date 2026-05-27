// Package storage holds Marshal/Unmarshal helpers that serialise an
// ordered slice of vehicle.Vehicle into the text format chosen by the
// current variant (variant 3 — plain text, key=value).
//
// The encoder/decoder never branch on concrete Vehicle types: they only
// read from FieldDescriptor closures, so brand new classes added through
// the registry are supported automatically.
package storage

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"oop/lab3/internal/vehicle"
)

// File format
//
//	[BEGIN <TypeName>]
//	FieldName=encoded value
//	...
//	[END]
//
// One block per object, blocks are separated by a single blank line.
// Values are encoded by replacing CR/LF/`\` with their `\n`/`\r`/`\\`
// counterparts so that any string can be safely round-tripped.

const (
	tagBegin = "[BEGIN "
	tagEnd   = "[END]"
)

// encodeValue escapes characters that would otherwise break the
// line-oriented format.
func encodeValue(s string) string {
	r := strings.NewReplacer(`\`, `\\`, "\n", `\n`, "\r", `\r`)
	return r.Replace(s)
}

// decodeValue reverses encodeValue.
func decodeValue(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case '\\':
				b.WriteByte('\\')
			default:
				b.WriteByte(s[i+1])
			}
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// Marshal writes the entire vehicles slice to w in the lab-3 text
// format. Field iteration is delegated to each Vehicle's Fields().
func Marshal(w io.Writer, items []vehicle.Vehicle) error {
	bw := bufio.NewWriter(w)
	for i, v := range items {
		if i > 0 {
			if _, err := fmt.Fprintln(bw); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(bw, "%s%s]\n", tagBegin, v.TypeName()); err != nil {
			return err
		}
		for _, f := range v.Fields() {
			if _, err := fmt.Fprintf(bw, "%s=%s\n", f.Name, encodeValue(f.Get())); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(bw, tagEnd); err != nil {
			return err
		}
	}
	return bw.Flush()
}

// Unmarshal reads vehicles previously written by Marshal. Unknown type
// tags or unknown field names are reported as errors so the user can
// fix bad files instead of silently losing data.
func Unmarshal(r io.Reader) ([]vehicle.Vehicle, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var result []vehicle.Vehicle
	var current vehicle.Vehicle
	var currentFields map[string]vehicle.FieldDescriptor
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimRight(scanner.Text(), "\r")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, tagBegin) && strings.HasSuffix(line, "]") {
			typeName := strings.TrimSuffix(strings.TrimPrefix(line, tagBegin), "]")
			v, err := vehicle.Create(typeName)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			current = v
			currentFields = make(map[string]vehicle.FieldDescriptor)
			for _, f := range v.Fields() {
				currentFields[f.Name] = f
			}
			continue
		}
		if line == tagEnd {
			if current == nil {
				return nil, fmt.Errorf("line %d: unexpected %s", lineNo, tagEnd)
			}
			result = append(result, current)
			current = nil
			currentFields = nil
			continue
		}
		if current == nil {
			return nil, fmt.Errorf("line %d: data outside object block", lineNo)
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNo, line)
		}
		name := line[:eq]
		raw := decodeValue(line[eq+1:])
		f, ok := currentFields[name]
		if !ok {
			return nil, fmt.Errorf("line %d: %s has no field %q", lineNo, current.TypeName(), name)
		}
		if err := f.Set(raw); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if current != nil {
		return nil, fmt.Errorf("missing %s tag at end of file", tagEnd)
	}
	return result, nil
}

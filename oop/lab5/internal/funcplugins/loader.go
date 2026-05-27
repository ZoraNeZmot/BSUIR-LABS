package funcplugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadDir loads every *.json file in dir as a functional plugin. Files
// that fail to parse are reported in errs but do not abort the scan.
func LoadDir(dir string) (int, []error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, []error{fmt.Errorf("read dir: %w", err)}
	}
	loaded := 0
	var errs []error
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if err := LoadFile(path); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", e.Name(), err))
			continue
		}
		loaded++
	}
	return loaded, errs
}

// LoadFile loads a single descriptor file.
func LoadFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var p FuncPlugin
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}
	if p.Parameters == nil {
		p.Parameters = make(map[string]string)
	}
	algo, ok := LookupAlgorithm(p.AlgorithmID)
	if !ok {
		return fmt.Errorf("unknown algorithm %q", p.AlgorithmID)
	}
	// Merge defaults so the settings dialog always has every key
	// present even if the JSON omitted some entries.
	for _, ps := range algo.Parameters() {
		if _, present := p.Parameters[ps.Name]; !present {
			p.Parameters[ps.Name] = ps.Default
		}
	}
	p.source = path
	return AddPlugin(&p)
}

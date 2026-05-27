// Package plugins implements dynamic loading of new Vehicle classes
// from JSON descriptor files. Dropping a fresh `*.json` file into the
// plugins directory and re-running the program (or pressing the
// "Reload plugins" button) is enough to extend the hierarchy: no source
// file inside the host program has to be modified.
//
// Plugin descriptor schema (lab 4)
// --------------------------------
//
//	{
//	  "typeName":      "Bicycle",
//	  "category":      "Land",
//	  "summary":       "[%s] %s %s (%s)",
//	  "summaryFields": ["TypeName", "Manufacturer", "Model", "Year"],
//	  "fields": [
//	    {"name": "ID",            "label": "Identifier",      "kind": "string"},
//	    {"name": "Manufacturer",  "label": "Manufacturer",    "kind": "string"},
//	    {"name": "Model",         "label": "Model",           "kind": "string"},
//	    {"name": "Year",          "label": "Year",            "kind": "int",   "default": "2024"},
//	    {"name": "GearCount",     "label": "Number of gears", "kind": "int"},
//	    {"name": "FrameMaterial", "label": "Frame material",  "kind": "string","default": "Aluminum"},
//	    {"name": "IsElectric",    "label": "Electric",        "kind": "bool"}
//	  ]
//	}
package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"oop/lab4/internal/vehicle"
)

// FieldSpec describes one editable field declared by a plugin.
type FieldSpec struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Kind    string `json:"kind"`
	Default string `json:"default"`
}

// Spec is the in-memory representation of a single plugin descriptor.
type Spec struct {
	TypeName      string      `json:"typeName"`
	Category      string      `json:"category"`
	Summary       string      `json:"summary"`
	SummaryFields []string    `json:"summaryFields"`
	Fields        []FieldSpec `json:"fields"`

	source string // file path, kept for diagnostics
}

// kindMap maps the string kind tag found in JSON to the closed
// FieldKind enum understood by the rest of the codebase. A new field
// kind would only require an extra entry here, not a switch in every
// consumer.
var kindMap = map[string]vehicle.FieldKind{
	"string": vehicle.KindString,
	"int":    vehicle.KindInt,
	"float":  vehicle.KindFloat,
	"bool":   vehicle.KindBool,
}

// validators check that a string value is acceptable for a given kind
// before it is stored. They are looked up by kind tag, not by class.
var validators = map[string]func(string) error{
	"string": func(string) error { return nil },
	"int":    func(s string) error { _, err := strconv.Atoi(s); return err },
	"float":  func(s string) error { _, err := strconv.ParseFloat(s, 64); return err },
	"bool":   func(s string) error { _, err := strconv.ParseBool(s); return err },
}

// LoadDir scans dir for *.json files and registers every well-formed
// descriptor with the global vehicle registry. The number of plugins
// loaded and the list of errors (per file) are returned so the UI can
// display them.
func LoadDir(dir string) (int, []error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, []error{fmt.Errorf("read plugin dir: %w", err)}
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

// LoadFile reads one descriptor file and registers it.
func LoadFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var spec Spec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}
	spec.source = path
	if err := validateSpec(&spec); err != nil {
		return err
	}
	vehicle.Register(spec.TypeName, func() vehicle.Vehicle {
		return newPluginVehicle(&spec)
	})
	return nil
}

// validateSpec catches malformed plugins early.
func validateSpec(s *Spec) error {
	if strings.TrimSpace(s.TypeName) == "" {
		return fmt.Errorf("typeName is required")
	}
	if strings.TrimSpace(s.Category) == "" {
		s.Category = "Other"
	}
	if len(s.Fields) == 0 {
		return fmt.Errorf("fields list is empty")
	}
	seen := map[string]struct{}{}
	for i := range s.Fields {
		f := &s.Fields[i]
		if f.Name == "" {
			return fmt.Errorf("field #%d: empty name", i)
		}
		if _, dup := seen[f.Name]; dup {
			return fmt.Errorf("field %q is declared twice", f.Name)
		}
		seen[f.Name] = struct{}{}
		if _, ok := kindMap[f.Kind]; !ok {
			return fmt.Errorf("field %q: unknown kind %q", f.Name, f.Kind)
		}
		if f.Label == "" {
			f.Label = f.Name
		}
	}
	return nil
}

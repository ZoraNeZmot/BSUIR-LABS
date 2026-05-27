package plugins

import (
	"fmt"
	"strings"

	"oop/lab4/internal/vehicle"
)

// pluginVehicle is the generic Vehicle implementation backing every
// plugin-defined class. All field values are stored as their raw string
// representation: kind-aware validators inside the Set closure ensure
// the string can later be parsed back without surprises.
type pluginVehicle struct {
	spec   *Spec
	values map[string]string
}

// newPluginVehicle creates a fresh instance with default values from
// the descriptor.
func newPluginVehicle(spec *Spec) *pluginVehicle {
	v := &pluginVehicle{spec: spec, values: make(map[string]string, len(spec.Fields))}
	for _, f := range spec.Fields {
		v.values[f.Name] = f.Default
	}
	return v
}

func (p *pluginVehicle) TypeName() string { return p.spec.TypeName }
func (p *pluginVehicle) Category() string { return p.spec.Category }

// Summary builds the master-list label according to the descriptor's
// summary template. Missing template falls back to a sensible default.
func (p *pluginVehicle) Summary() string {
	if p.spec.Summary == "" || len(p.spec.SummaryFields) == 0 {
		return fmt.Sprintf("[%s] %s", p.spec.TypeName, p.values["Model"])
	}
	args := make([]any, 0, len(p.spec.SummaryFields))
	for _, name := range p.spec.SummaryFields {
		if name == "TypeName" {
			args = append(args, p.spec.TypeName)
			continue
		}
		args = append(args, p.values[name])
	}
	return fmt.Sprintf(p.spec.Summary, args...)
}

// Fields returns generic FieldDescriptors, one per declared field. The
// closures capture the FieldSpec by value to avoid the classic loop
// variable trap.
func (p *pluginVehicle) Fields() []vehicle.FieldDescriptor {
	out := make([]vehicle.FieldDescriptor, 0, len(p.spec.Fields))
	for _, fs := range p.spec.Fields {
		fs := fs
		validate := validators[fs.Kind]
		out = append(out, vehicle.FieldDescriptor{
			Name:  fs.Name,
			Label: fs.Label,
			Kind:  kindMap[fs.Kind],
			Get:   func() string { return p.values[fs.Name] },
			Set: func(s string) error {
				s = strings.TrimSpace(s)
				if validate != nil {
					if err := validate(s); err != nil {
						return fmt.Errorf("field %q: %w", fs.Name, err)
					}
				}
				p.values[fs.Name] = s
				return nil
			},
		})
	}
	return out
}

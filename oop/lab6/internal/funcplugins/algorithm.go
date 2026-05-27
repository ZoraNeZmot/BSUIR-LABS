// Package funcplugins implements the "variant 3 вЂ” encryption" lab-5
// requirement. It exposes:
//
//   - a small set of built-in cryptographic algorithms (XOR, Caesar,
//     AES-CFB) registered through their package's init();
//   - a JSON-descriptor format identical in spirit to lab 4, used to
//     declare *plugin instances* (a particular algorithm together with
//     its concrete parameters) and dropped in a folder;
//   - a pipeline that applies every enabled plugin in registration
//     order before saving and in reverse order when loading.
//
// The algorithm/plugin split makes the 10-point bonus trivial: each
// algorithm advertises its own parameter list and the settings dialog
// builds an editor for it generically.
package funcplugins

// ParamKind enumerates the primitive parameter types accepted by an
// algorithm. The settings dialog reads it to choose between a regular
// entry, a password entry and a numeric entry.
type ParamKind int

const (
	ParamString ParamKind = iota
	ParamSecret
	ParamInt
)

// ParamSpec describes one configuration knob an algorithm exposes.
type ParamSpec struct {
	Name    string
	Label   string
	Kind    ParamKind
	Default string
}

// Algorithm is the strategy interface for any data-transformation step
// the host can apply before saving / after loading. Implementations are
// stateless: all configuration travels through the params map so the
// same algorithm object can back several plugin instances.
type Algorithm interface {
	ID() string
	DisplayName() string
	Description() string
	Parameters() []ParamSpec
	Encode(data []byte, params map[string]string) ([]byte, error)
	Decode(data []byte, params map[string]string) ([]byte, error)
}

package funcplugins

import "fmt"

// Encode runs every enabled plugin's Encode in registration order.
// `data` is the freshly serialized payload; the result must be written
// to disk verbatim.
func Encode(data []byte) ([]byte, error) {
	out := data
	for _, p := range EnabledPlugins() {
		algo, ok := LookupAlgorithm(p.AlgorithmID)
		if !ok {
			return nil, fmt.Errorf("plugin %q: algorithm %q vanished", p.ID, p.AlgorithmID)
		}
		next, err := algo.Encode(out, p.Parameters)
		if err != nil {
			return nil, fmt.Errorf("plugin %q: %w", p.ID, err)
		}
		out = next
	}
	return out, nil
}

// Decode runs every enabled plugin's Decode in REVERSE registration
// order вЂ” the inverse of Encode.
func Decode(data []byte) ([]byte, error) {
	out := data
	enabled := EnabledPlugins()
	for i := len(enabled) - 1; i >= 0; i-- {
		p := enabled[i]
		algo, ok := LookupAlgorithm(p.AlgorithmID)
		if !ok {
			return nil, fmt.Errorf("plugin %q: algorithm %q vanished", p.ID, p.AlgorithmID)
		}
		next, err := algo.Decode(out, p.Parameters)
		if err != nil {
			return nil, fmt.Errorf("plugin %q: %w", p.ID, err)
		}
		out = next
	}
	return out, nil
}

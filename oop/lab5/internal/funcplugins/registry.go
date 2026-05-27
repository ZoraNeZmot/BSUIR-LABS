package funcplugins

import (
	"fmt"
	"sort"
	"sync"
)

var (
	algoLock     sync.RWMutex
	algorithms   = map[string]Algorithm{}
	algoOrder    []string
	pluginsLock  sync.RWMutex
	pluginsList  []*FuncPlugin
	pluginsByID  = map[string]*FuncPlugin{}
)

// RegisterAlgorithm makes an Algorithm available to the host. Each
// built-in algorithm calls this from its own init().
func RegisterAlgorithm(a Algorithm) {
	algoLock.Lock()
	defer algoLock.Unlock()
	if _, ok := algorithms[a.ID()]; !ok {
		algoOrder = append(algoOrder, a.ID())
	}
	algorithms[a.ID()] = a
}

// LookupAlgorithm fetches an algorithm by ID.
func LookupAlgorithm(id string) (Algorithm, bool) {
	algoLock.RLock()
	defer algoLock.RUnlock()
	a, ok := algorithms[id]
	return a, ok
}

// AlgorithmIDs returns the IDs of all registered algorithms.
func AlgorithmIDs() []string {
	algoLock.RLock()
	defer algoLock.RUnlock()
	out := make([]string, len(algoOrder))
	copy(out, algoOrder)
	sort.Strings(out)
	return out
}

// FuncPlugin is the runtime representation of a JSON descriptor.
type FuncPlugin struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	AlgorithmID string            `json:"algorithm"`
	Parameters  map[string]string `json:"parameters"`
	Enabled     bool              `json:"enabled"`

	source string // path on disk, kept for diagnostics
}

// AddPlugin registers a plugin instance. Re-using an existing ID
// silently replaces the previous record.
func AddPlugin(p *FuncPlugin) error {
	if p.ID == "" {
		return fmt.Errorf("plugin id is required")
	}
	if _, ok := LookupAlgorithm(p.AlgorithmID); !ok {
		return fmt.Errorf("plugin %q references unknown algorithm %q", p.ID, p.AlgorithmID)
	}
	pluginsLock.Lock()
	defer pluginsLock.Unlock()
	if existing, ok := pluginsByID[p.ID]; ok {
		*existing = *p
		return nil
	}
	pluginsByID[p.ID] = p
	pluginsList = append(pluginsList, p)
	return nil
}

// Plugins returns a snapshot of registered plugins in the order they
// were added (which is also the encoding order).
func Plugins() []*FuncPlugin {
	pluginsLock.RLock()
	defer pluginsLock.RUnlock()
	out := make([]*FuncPlugin, len(pluginsList))
	copy(out, pluginsList)
	return out
}

// EnabledPlugins returns plugins with Enabled==true in their declared
// order. The encoding pipeline iterates this slice forward; the
// decoding pipeline walks it in reverse.
func EnabledPlugins() []*FuncPlugin {
	pluginsLock.RLock()
	defer pluginsLock.RUnlock()
	out := make([]*FuncPlugin, 0, len(pluginsList))
	for _, p := range pluginsList {
		if p.Enabled {
			out = append(out, p)
		}
	}
	return out
}

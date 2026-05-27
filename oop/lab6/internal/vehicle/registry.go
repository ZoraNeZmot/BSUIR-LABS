package vehicle

import (
	"fmt"
	"sort"
	"sync"
)

// Factory is the constructor signature used to instantiate a Vehicle of
// a specific concrete type. Each class provides one and registers it
// from its own init() function.
type Factory func() Vehicle

// registryLock guards both maps so plugins loaded from goroutines stay
// safe even though all current registration happens on init().
var registryLock sync.RWMutex
var factories = make(map[string]Factory)
var orderedNames []string

// Register adds a factory under the given type name. Re-registering an
// existing name is allowed (last definition wins) so that plugins can
// override built-in classes if they really need to.
func Register(typeName string, f Factory) {
	registryLock.Lock()
	defer registryLock.Unlock()
	if _, exists := factories[typeName]; !exists {
		orderedNames = append(orderedNames, typeName)
	}
	factories[typeName] = f
}

// Create instantiates a fresh Vehicle of the requested type. It returns
// an error if the type has not been registered.
func Create(typeName string) (Vehicle, error) {
	registryLock.RLock()
	f, ok := factories[typeName]
	registryLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("vehicle type %q is not registered", typeName)
	}
	return f(), nil
}

// Names returns all registered type names sorted alphabetically.
func Names() []string {
	registryLock.RLock()
	defer registryLock.RUnlock()
	out := make([]string, len(orderedNames))
	copy(out, orderedNames)
	sort.Strings(out)
	return out
}

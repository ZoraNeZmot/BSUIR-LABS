// Package appsettings is the lab-6 demonstration of the SINGLETON
// design pattern.
//
// Why a singleton? The application has a small set of cross-cutting
// preferences (last opened file, default plugin directories, the user
// name shown in the title bar) that have no natural owner -- the UI,
// the storage layer and the plugin loader all need to read them. A
// process-wide single instance avoids threading the same struct
// through every constructor while keeping access serialized.
package appsettings

import "sync"

// Settings holds the cross-cutting preferences. New fields can be added
// without touching the consumers thanks to the typed accessors.
type Settings struct {
	mu             sync.RWMutex
	UserName       string
	LastOpenedFile string
	PluginDir      string
	FuncPluginDir  string
}

var (
	once     sync.Once
	instance *Settings
)

// Instance returns the lazily initialised singleton. The sync.Once
// guarantee gives us a safe, idiomatic implementation without leaking
// the constructor.
func Instance() *Settings {
	once.Do(func() {
		instance = &Settings{
			UserName:      "operator",
			PluginDir:     "plugins",
			FuncPluginDir: "funcplugins",
		}
	})
	return instance
}

// Get returns a snapshot of the settings. Exposing a copy keeps the
// internal mutex private to this package.
func (s *Settings) Get() Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Settings{
		UserName:       s.UserName,
		LastOpenedFile: s.LastOpenedFile,
		PluginDir:      s.PluginDir,
		FuncPluginDir:  s.FuncPluginDir,
	}
}

// Update applies fn under the write lock.
func (s *Settings) Update(fn func(*Settings)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s)
}

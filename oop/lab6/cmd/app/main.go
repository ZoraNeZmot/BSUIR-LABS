// Lab 6 -- "Patterns".
//
// Lab 6 carries forward every feature of lab 5 (hierarchy plugins,
// encryption pipeline, settings) and on top adds:
//
//   - the Adapter pattern, used to bridge a friend's library
//     (internal/friendlib) with our Algorithm contract through
//     internal/funcplugins/algos/checksum_adapter.go;
//   - the Singleton pattern, demonstrated by internal/appsettings;
//   - the Observer pattern, demonstrated by internal/events.
//
// Built-in classes are wired through their package's init() functions
// (blank imports below). Hierarchy plugins live as *.json descriptors
// inside plugins/, functional plugins live as *.json descriptors
// inside funcplugins/.
package main

import (
	"flag"
	"fmt"
	"os"

	"oop/lab6/internal/funcplugins"
	_ "oop/lab6/internal/funcplugins/algos"
	"oop/lab6/internal/plugins"
	"oop/lab6/internal/ui"
	_ "oop/lab6/internal/vehicles"
)

func main() {
	pluginDir := flag.String("plugins", "plugins", "directory with hierarchy plugin descriptors")
	funcDir := flag.String("funcplugins", "funcplugins", "directory with functional plugin descriptors")
	flag.Parse()
	ui.PluginDir = *pluginDir
	ui.FuncPluginDir = *funcDir

	n, errs := plugins.LoadDir(*pluginDir)
	fmt.Fprintf(os.Stderr, "Loaded %d hierarchy plugin(s) from %q\n", n, *pluginDir)
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "  plugin error: %v\n", e)
	}
	fn, ferrs := funcplugins.LoadDir(*funcDir)
	fmt.Fprintf(os.Stderr, "Loaded %d functional plugin(s) from %q\n", fn, *funcDir)
	for _, e := range ferrs {
		fmt.Fprintf(os.Stderr, "  funcplugin error: %v\n", e)
	}
	ui.Run()
}

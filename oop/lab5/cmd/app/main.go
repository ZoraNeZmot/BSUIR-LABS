// Lab 5 -- "Plugins (functionality)", variant 3 (encryption/decryption).
//
// Built-in classes are wired through their package's init() functions
// (blank import below). Additional classes live as JSON descriptor
// files inside the plugins/ directory; functional plugins (encryption
// instances) live as JSON descriptors inside the funcplugins/ directory
// and reference one of the algorithms registered through algos/*.go.
package main

import (
	"flag"
	"fmt"
	"os"

	"oop/lab5/internal/funcplugins"
	_ "oop/lab5/internal/funcplugins/algos"
	"oop/lab5/internal/plugins"
	"oop/lab5/internal/ui"
	_ "oop/lab5/internal/vehicles"
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

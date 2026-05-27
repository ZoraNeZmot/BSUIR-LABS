// Lab 4 — "Plugins (hierarchy)".
//
// Built-in classes are wired through their package's init() functions
// (blank import below). Additional classes live as JSON descriptor
// files inside the plugins/ directory; LoadDir registers them with the
// same vehicle registry used by the built-in classes, so the rest of
// the program treats them identically.
package main

import (
	"flag"
	"fmt"
	"os"

	"oop/lab4/internal/plugins"
	"oop/lab4/internal/ui"
	_ "oop/lab4/internal/vehicles"
)

func main() {
	dir := flag.String("plugins", "plugins", "directory with plugin descriptor files")
	flag.Parse()
	ui.PluginDir = *dir

	n, errs := plugins.LoadDir(*dir)
	fmt.Fprintf(os.Stderr, "Loaded %d plugin(s) from %q\n", n, *dir)
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "  plugin error: %v\n", e)
	}
	ui.Run()
}

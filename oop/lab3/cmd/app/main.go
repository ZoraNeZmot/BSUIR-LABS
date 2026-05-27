// Lab 3 — "Object serialization", variant 3 (Text format).
//
// The blank import of internal/vehicles wires every concrete class into
// the global registry through their init() functions: that is the only
// place where the main module knows that such classes exist.
package main

import (
	_ "oop/lab3/internal/vehicles"
	"oop/lab3/internal/ui"
)

func main() {
	ui.Run()
}

package startup

import (
	"flag"
	"fmt"
)

const (
	DefaultName    = "Anonim"
	DefaultPort    = 56789
	DefaultPortStr = "56789"
)

var (
	FlagVerbose *bool  = flag.Bool("v", false, "Verbose mode")
	FlagPort    *int   = flag.Int("p", DefaultPort, "Port for creating a connection")
	Name        string = DefaultName
)

func init() {
	flag.Parse()
	if flag.NArg() > 0 {
		Name = flag.Arg(0)
	}

	if *FlagVerbose {
		fmt.Printf("User: %s registered.\tPort used: %d\nBroadcasting INVITING message.\n", Name, *FlagPort)
	}
}

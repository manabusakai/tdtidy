package tdtidy

import (
	"fmt"
	"os"
)

var debug = debugger((os.Getenv("DEBUG") == "true"))

type debugger bool

func (d debugger) Printf(format string, args ...interface{}) {
	if d {
		fmt.Printf("[debug] %s\n", fmt.Sprintf(format, args...))
	}
}

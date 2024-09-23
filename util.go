package tdtidy

import (
	"fmt"
	"os"
	"time"
)

var debug = debugger((os.Getenv("DEBUG") == "true"))

type debugger bool

func (d debugger) Printf(format string, args ...interface{}) {
	if d {
		fmt.Printf("[debug] %s\n", fmt.Sprintf(format, args...))
	}
}

func sleep() {
	// Refill rate of API actions per second.
	// https://docs.aws.amazon.com/AmazonECS/latest/APIReference/request-throttling.html
	const refillRate = 1

	time.Sleep(refillRate * time.Second)
}

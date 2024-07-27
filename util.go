package tdtidy

import (
	"time"
)

func chunk(items []string, chunkSize int) (chunks [][]string) {
	if len(items) == 0 {
		return [][]string{}
	}
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[:chunkSize])
	}

	return append(chunks, items)
}

func sleep() {
	// Refill rate of API actions per second.
	// https://docs.aws.amazon.com/AmazonECS/latest/APIReference/request-throttling.html
	const refillRate = 1

	time.Sleep(refillRate * time.Second)
}

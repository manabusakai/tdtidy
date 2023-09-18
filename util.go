package tdtidy

func chunk(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[:chunkSize])
	}

	return append(chunks, items)
}

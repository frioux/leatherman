package fuzz

func Fuzz(data []byte) int {
	if len(data) < 1 {
		return 0
	}

	return 1
}

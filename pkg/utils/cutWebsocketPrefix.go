package utils

func CutMessagePrefix(bytes []byte) []byte {
	var pos int
	for i, b := range bytes {
		if b == '{' {
			pos = i
			break
		}
	}

	return bytes[pos:]
}

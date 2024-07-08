package parse

const stringDelimiter = byte(':')

// ParseInt parses a decimal number in p
// the index of p which finished parsing is returned.
func ParseInt(p []byte) (int64, int) {
	var number int64

	var i int

	var magnitude int64 = 1

	if p[0] == '-' {
		magnitude = -1
		i++
	}

	for ; i < len(p); i++ {
		cSubZ := p[i] - '0'

		if !(cSubZ < 10) {
			break
		}

		digit := int64(cSubZ)
		number = (number * 10) + digit
	}

	return magnitude * number, i
}

// ParseString parses a bencoded string
// the index of p which finished parsing is returned.
func ParseString(p []byte) ([]byte, int) {
	slen, i := ParseInt(p)

	if p[i] != stringDelimiter {
		return nil, -1
	}

	if int(slen) > len(p) {
		return nil, -2
	}

	idx := int(slen) + i
	return p[i+1 : idx+1], idx + 1
}

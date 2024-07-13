package tracker

// Escape20 is a specialised form of net/url.escape
// optimised for info_hash/peer_id.
//
// In the worst case, if every character of src
// needed escaping as hex, then we would need
// three bytes per byte of src.
// The caller must ensure dst can hold up to
// 60 bytes.
func Escape20(dst []byte, src *[20]byte) int {
	const upperhex = "0123456789ABCDEF"

	j := 0
	for i := 0; i < 20; i++ {
		switch c := src[i]; {
		case c == ' ':
			dst[j] = '+'
			j++
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' ||
			c == '-' || c == '_' || c == '.' || c == '~':
			dst[j] = src[i]
			j++
		default:
			dst[j] = '%'
			dst[j+1] = upperhex[c>>4]
			dst[j+2] = upperhex[c&15]
			j += 3
		}
	}

	return j
}

package tracker

import (
	"math/rand"
	"net/url"
	"testing"
)

func TestEscape20(t *testing.T) {
	// we will test, exhaustively
	// one byte of the input
	// and fuzz the rest
	t.Run("OneByte", func(t *testing.T) {
		var dst [60]byte
		var src [20]byte

		// set every character to an alphanumeric char
		// that won't be expanded
		for i := 0; i < 20; i++ {
			src[i] = 'a'
		}

		var i int
		for ; i < 256; i++ {
			src[0] = byte(i)

			n := Escape20(dst[:], &src)

			got := dst[:n]
			want := url.QueryEscape(string(src[:]))

			if string(got) != want {
				t.Fatalf("%s/%d got: %s, want: %s", t.Name(), i, got, want)
			}
		}
	})

	t.Run("Fuzz", func(t *testing.T) {
		var dst [60]byte
		var src [20]byte

		for i := 0; i < 100000; i++ {
			_, err := rand.Read(src[:])
			if err != nil {
				t.Fatalf("rand.Read: %v", err)
			}

			n := Escape20(dst[:], &src)

			got := dst[:n]
			want := url.QueryEscape(string(src[:]))

			if string(got) != want {
				t.Fatalf("input: %x: got: %s, want: %s", string(src[:]), got, want)
			}

		}
	})
}

func BenchmarkEscape20(b *testing.B) {
	const corpusSize = 8192

	corpus := make([][20]byte, corpusSize)

	for i := 0; i < len(corpus); i++ {
		_, err := rand.Read(corpus[i][:])
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst [60]byte

		_ = Escape20(dst[:], &corpus[i%corpusSize])
	}
}

package parse

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		wantNum  int64
		wantByte byte
	}{
		{
			name:     "Correct",
			s:        "1829e",
			wantNum:  1829,
			wantByte: 'e',
		},
		{
			name:     "Correct",
			s:        "1e",
			wantNum:  1,
			wantByte: 'e',
		},
		{
			name:     "LengthOne",
			s:        "1",
			wantNum:  1,
			wantByte: '1',
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotNum, _ := ParseInt([]byte(tc.s))

			if gotNum != tc.wantNum {
				t.Fatalf("%s: got: %d, want: %d", tc.name, gotNum, tc.wantNum)
			}

			// if gotByte != tc.wantByte {
			// 	t.Fatalf("%s: got: %d, want: %d", tc.name, gotByte, tc.wantByte)
			// }
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name      string
		p         string
		want      string
		wantError bool
		wantChar  byte
	}{
		{
			name:     "Correct",
			p:        "3:dogZYX",
			want:     "dog",
			wantChar: 'Z',
		},
		{
			name:     "Empty",
			p:        "0:ZYX",
			want:     "",
			wantChar: 'Z',
		},
		{
			name:     "LengthOne",
			p:        "1:aZYX",
			want:     "a",
			wantChar: 'Z',
		},
		{
			name:      "BadLengthGT",
			p:         "99:cheese",
			wantError: true,
		},
		{
			name:      "NoDelim",
			p:         "3cat",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, gotIdx := ParseString([]byte(tc.p))

			got := string(gotBytes)
			gotError := gotIdx < 0

			if gotError != tc.wantError {
				t.Fatalf("%s: got error: %v, want: %v", tc.name, gotError, tc.wantError)
			}

			if got != tc.want {
				t.Fatalf("%s: got: %s, want: %s", tc.name, got, tc.want)
			}

			if tc.wantChar != byte(0) {
				if gotIdx >= len(tc.p) {
					t.Fatalf("idx greater than length")
				}

				gotChar := tc.p[gotIdx]
				if gotChar != tc.wantChar {
					t.Fatalf("got char: %s, want: %s (idx: %d)",
						[]byte{gotChar},
						[]byte{tc.wantChar},
						gotIdx)
				}

			}

		})
	}
}

// BENCH

func Benchmark_parseNum(b *testing.B) {
	const numCases = 128
	casesBytes := make([][]byte, 0)
	casesStrings := make([]string, 0)

	for i := 0; i < numCases; i++ {
		n := rand.Int()
		nstr := strconv.Itoa(n)
		casesBytes = append(casesBytes, []byte(nstr))
		casesStrings = append(casesStrings, nstr)
	}

	b.Run("strconv_positive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = strconv.Atoi(casesStrings[i%numCases])
		}
	})

	b.Run("ParseNum_positive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ParseInt(casesBytes[i%numCases])
		}
	})

	for i := 0; i < len(casesBytes); i++ {
		n := make([]byte, 0, len(casesBytes))
		n = append(n, '-')
		n = append(n, casesBytes[i]...)

		casesBytes[i] = n
		casesStrings[i] = "-" + casesStrings[i]
	}

	b.Run("strconv_negative", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = strconv.Atoi(casesStrings[i%numCases])
		}
	})

	b.Run("ParseNum_negative", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ParseInt(casesBytes[i%numCases])
		}
	})
}

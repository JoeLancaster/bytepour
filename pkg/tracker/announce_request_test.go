package tracker

import (
	"math/rand"
	"net/url"
	"strconv"
	"testing"
)

func TestBuild(t *testing.T) {
	var b URLBuffer

	url := Build(&b, []byte("bbc.co.uk:9000"), &AnnounceRequest{
		InfoHash: [20]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		PeerId: [20]byte{'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		Downloaded: 21944892,
	})

	if url == nil {
		t.Fatal("did not build")
	}

	t.Log(string(url))
}

func Test_uintLen(t *testing.T) {
	for i := 0; i < 100000; i++ {
		n := rand.Uint64()

		strlen := len(strconv.FormatUint(n, 10))
		ulen := uintLen(n)

		if strlen != ulen {
			t.Fatalf("lengths dont match; strconv says: %d, uintLen says: %d", strlen, ulen)
		}

	}
}

func BenchmarkBytePourBuild(b *testing.B) {
	host := "http://bbc.co.uk:9000"

	req := AnnounceRequest{
		InfoHash: [20]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		PeerId: [20]byte{'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		Port:       6007,
		Uploaded:   192000,
		Downloaded: 198498929,
		Left:       234848334,
		NumWant:    50,
		Event:      EventCompleted,
	}

	for i := 0; i < b.N; i++ {
		var b URLBuffer
		x := Build(&b, []byte(host), &req)
		_ = x
	}
}

func BenchmarkNetBuild(b *testing.B) {
	host := "bbc.co.uk:9000"

	req := AnnounceRequest{
		InfoHash: [20]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		PeerId: [20]byte{'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		Port:       6007,
		Uploaded:   192000,
		Downloaded: 198498929,
		Left:       234848334,
		NumWant:    50,
	}
	for i := 0; i < b.N; i++ {
		var u url.URL

		u.Scheme = "http"
		u.Host = host
		q := make(url.Values)
		q.Set("info_hash", string(req.InfoHash[:]))
		q.Set("peer_id", string(req.PeerId[:]))
		q.Set("port", strconv.FormatUint(req.Port, 10))
		q.Set("uploaded", strconv.FormatUint(req.Uploaded, 10))
		q.Set("downloaded", strconv.FormatUint(req.Downloaded, 10))
		q.Set("left", strconv.FormatUint(req.Left, 10))
		q.Set("numwant", strconv.FormatUint(req.NumWant, 10))
		q.Set("event", "completed")
		u.RawQuery = q.Encode()

		_ = u.String()
	}
}

func BenchmarkNetEscape(b *testing.B) {
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
		_ = url.QueryEscape(string(corpus[i%corpusSize][:]))
	}
}

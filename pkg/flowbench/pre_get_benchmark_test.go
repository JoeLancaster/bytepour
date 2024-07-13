package flowbench

import (
	"crypto/sha1"
	_ "embed"
	"testing"

	"github.com/joelancaster/bytepour/pkg/bencode/aot"
	"github.com/joelancaster/bytepour/pkg/bittorrent"
	"github.com/joelancaster/bytepour/pkg/metainfo"
	"github.com/joelancaster/bytepour/pkg/tracker"
)

//go:embed testdata/debian.torrent
var debian []byte

/*
Benchmark parsing a .torrent file into
a metainfo struct
then constructing a query url
*/
func Benchmark_Announce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var URL tracker.URLBuffer
		var mi metainfo.MetaInfoPreCompute

		err := aot.DecodeMetaInfoFile(&mi, debian)
		if err.IsError() {
			b.Fatal()
		}

		req := tracker.AnnounceRequest{
			PeerId:     bittorrent.IdBP,
			InfoHash:   sha1.Sum(mi.InfoDict),
			Port:       6007,
			Uploaded:   0,
			Downloaded: 0,
			Left:       mi.Info.Length,
			NumWant:    50,
		}

		tracker.Build(&URL, mi.Announce, &req)

	}
}

package aot

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/joelancaster/bytepour/bencode/aot/testdata"
	"github.com/joelancaster/bytepour/metainfo"

	jackpal "github.com/jackpal/bencode-go"
)

//go:embed testdata/debian.torrent
var debian []byte

func TestAOTMetaInfo(t *testing.T) {
	mi := metainfo.MetaInfoPreCompute{}

	err := DecodeMetaInfoFile(&mi, debian)

	if err.IsError() {
		t.Logf("error: %s", err.Error())
		sourceStr, caret := whereDebug(debian, int(err.Where()))
		t.Log(sourceStr)
		t.Log(caret)
		t.Fail()
	}

	want := &testdata.WantDebianMetaInfo

	if !mi.Eq(want) {
		t.Fatalf("decoded metainfo not equal")
	}
}

func whereDebug(source []byte, where int) (string, string) {
	var from, to int
	to = where + 10
	from = where - 10

	if from < 0 {
		from = 0
	}

	if to+1 > len(source) {
		to = len(source)
	}

	sourceSnippet := source[from:to]

	srcLen := len(sourceSnippet)

	caretPos := srcLen - (to - where)

	caretLine := strings.Repeat(" ", caretPos) + "^" + strings.Repeat(" ", to)

	return string(sourceSnippet), caretLine
}

func BenchmarkBytePourDebian(b *testing.B) {
	var mi metainfo.MetaInfoPreCompute
	for i := 0; i < b.N; i++ {
		_ = DecodeMetaInfoFile(&mi, debian)
	}
}

func BenchmarkJackpalDebian(b *testing.B) {
	rdr := bytes.NewReader(debian)
	for i := 0; i < b.N; i++ {
		_, _ = jackpal.Decode(rdr)
		b.StopTimer()
		rdr.Reset(debian)
		b.StartTimer()
	}
}

package metainfo

import (
	"bytes"
	"encoding/json"
)

// MetaInfoPreCompute is the top-level
// dictionary of a metainfo file.
type MetaInfoPreCompute struct {
	// The URL of the tracker.
	Announce []byte `bencode:"announce"`
	// Optional free-form comment field.
	Comment []byte `bencode:"comment"`
	// Substring of the input that is the info dict.
	InfoDict []byte `bencode:"-" json:"-"`
	// The info dictionary, containing file info.
	Info Info `bencode:"info"`
}

// Info is the info dictionary in a metainfo file.
type Info struct {
	// Length of the file, in bytes.
	Length uint64 `bencode:"length"`
	// The name of the file.
	Name []byte `bencode:"name"`
	// The pieces of a file, kept as a single
	// string.
	Pieces []byte `bencode:"pieces" json:"-"`
	// The length of each piece.
	PieceLength uint64 `bencode:"piece length"`
}

// Eq compares a MetaInfoPreCompute for equality.
func (a *MetaInfoPreCompute) Eq(b *MetaInfoPreCompute) bool {
	if a == b {
		return true
	}

	return a.Info.Eq(&b.Info) &&
		bytes.Equal(a.InfoDict, b.InfoDict) &&
		bytes.Equal(a.Announce, b.Announce) &&
		bytes.Equal(a.Comment, b.Comment)

}

// Eq compares an Info for equality.
func (a *Info) Eq(b *Info) bool {
	if a == b {
		return true
	}

	return a.Length == b.Length &&
		a.PieceLength == b.PieceLength &&
		bytes.Equal(a.Name, b.Name) &&
		bytes.Equal(a.Pieces, b.Pieces)
}

// String implements the stringer interface for
// MetaInfoPreCompute. Debug use only.
func (m *MetaInfoPreCompute) String() string {
	s, _ := json.Marshal(m)

	return string(s)
}

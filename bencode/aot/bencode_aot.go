package aot

import (
	"github.com/joelancaster/bytepour/bencode/parse"
	"github.com/joelancaster/bytepour/metainfo"
)

// DecodeMetaInfoFile parses a bencode representation of a meta info file
// a.k.a. .torrent files.
func DecodeMetaInfoFile(mi *metainfo.MetaInfoPreCompute, p []byte) parse.Error {
	integers := map[string]*uint64{
		"length":       &mi.Info.Length,
		"piece length": &mi.Info.PieceLength,
	}

	bstrings := map[string]*[]byte{
		"announce": &mi.Announce,
		"comment":  &mi.Comment,
		"name":     &mi.Info.Name,
		"pieces":   &mi.Info.Pieces,
	}

	const maxLength = 0x7FFFFFFE

	var i uint32

	if len(p) >= maxLength {
		return parse.MakeError(parse.ErrInputTooLong, 0, 0)
	}

	// Valid bencodings should have a dictionary at the top level.
	if p[0] != 'd' {
		return parse.MakeError(parse.ErrNoTopLevelDict, 0, 0)
	}

	var s stack
	s.push(parse.Dict)

	var nextStr *[]byte
	var nextInt *uint64

	for i = 1; i < uint32(len(p)); {
		numeric := (p[i] - '0') < 10

		switch {
		case p[i] == parse.OpenList:
			i++
			s.push(parse.List)
		case p[i] == parse.OpenDict:
			i++
			s.push(parse.Dict)
		case p[i] == parse.OpenInt:
			n, j := parse.ParseInt(p[i+1:])

			if nextInt != nil {
				*nextInt = uint64(n)
				nextInt = nil
			}

			i += uint32(j) + 1
			if p[i] != 'e' {
				return parse.MakeError(parse.ErrUnexpectedEndOfTerm, i, parse.Int)
			}
			i++
		case numeric:
			bs, j := parse.ParseString(p[i:])
			if j < 0 {
				return parse.MakeError(parse.ErrUnexpectedEndOfTerm, i, parse.String)
			}

			i += uint32(j)

			// We only care about dictionaries for
			// a metainfo file.
			if s.topType() != parse.Dict {
				break
			}

			// This is a string value of a key in a
			// dictionary.
			if nextStr != nil {
				*nextStr = bs
				nextStr = nil

				break
			}

			intElem, ok := integers[string(bs)]
			if ok {
				nextInt = intElem
				break
			}

			// this is a key
			// possibly for an element we care about
			strElem, ok := bstrings[string(bs)]
			if ok {
				nextStr = strElem

				break
			}

			// Not something we care about, ignore it
		case p[i] == 'e':
			s.pop()
			// it would be impossible to be in the info dict after popping
			// from anything.
			i++
		default:
			// The input is malformed, or the parser is confused
			// at any state of the parse, we expect one of the valid characters
			// that begins a term.
			return parse.MakeError(parse.ErrConfusion, i, 0)
		}
	}

	return parse.ErrOk
}

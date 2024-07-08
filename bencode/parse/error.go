package parse

import (
	"strconv"
	"strings"
)

const (
	// ErrOk is the zero-value of Error
	// and indicates a successful parse.
	ErrOk = Error(0)
)

const (
	// While parsing a term, we reached the end when
	// we should not have.
	ErrUnexpectedEndOfTerm = uint64(1)
	// We are parsing nested terms too deep.
	ErrTermDepthLimit = uint64(2)
	// The parser has entered a state
	// that it cannot deal with.
	// This should never happen.
	ErrConfusion = uint64(3)
	// The input is too long.
	ErrInputTooLong = uint64(4)
	// The input's top level term is not a dictionary.
	ErrNoTopLevelDict = uint64(5)
	// The top level dictionary does not have
	// and "announce" key.
	ErrNoAnnounce = uint64(6)
)

// errorStrings is a lookup table of
// error codes to their string representation.
var errorStrings [7]string

// termStrings is a lookup table of
// term types to their string representation.
var termStrings [5]string

func init() {
	errorStrings[ErrUnexpectedEndOfTerm] = "unexpected end of term"
	errorStrings[ErrTermDepthLimit] = "reached maximum term depth limit"
	errorStrings[ErrConfusion] = "confusion"
	errorStrings[ErrNoTopLevelDict] = "bencode object does not have top-level dict"
	errorStrings[ErrNoAnnounce] = "no announce key"

	termStrings[List] = "list"
	termStrings[Dict] = "dict"
	termStrings[Int] = "int"
	termStrings[String] = "string"
	termStrings[StringHeader] = "string_header"
}

// Error is a custom error type for parsing errors.
type Error uint64

// IsError checks if e is an error,
// or if it is ok.
func (e Error) IsError() bool {
	return e != ErrOk
}

// Error implements the error interface
// for Error.
func (e Error) Error() string {
	return e.String()
}

// String implements the stringer interface
// for Error.
func (e Error) String() string {
	const (
		whatPart  = "error decoding bencode object: "
		whenPart  = " when parsing a "
		wherePart = " at character "
	)

	var (
		what  = e.what()
		where = e.Where()
		when  = e.term()
	)

	reason := errorStrings[what]
	term := termStrings[when]

	switch what {
	case ErrInputTooLong, ErrConfusion, ErrNoTopLevelDict, ErrNoAnnounce:
		return whatPart + reason
	}

	var sb strings.Builder

	sb.Grow(len(whatPart) + len(whenPart) + len(wherePart) + len(reason) + len(term))
	sb.WriteString(whatPart)
	sb.WriteString(reason)
	sb.WriteString(whenPart)
	sb.WriteString(term)
	sb.WriteString(wherePart)
	sb.WriteString(strconv.Itoa(int(where)))

	return sb.String()
}

// what yields the error code from e.
func (e Error) what() uint64 {
	return (uint64(e) & 0x00000000_0000_FF_00) >> 8
}

// Where yields the character position of the
// input that caused e.
func (e Error) Where() uint32 {
	s := uint64(e) & 0xFFFFFFFF_00_00_00_00 >> 32

	return uint32(s)
}

// term yields the term type that caused e.
func (e Error) term() Term {
	return Term(uint64(e) & 0x00000000_0000_00_FF)
}

// MakeError constructs an Error from the error code,
// character position, and term type.
func MakeError(code uint64, where uint32, term Term) Error {
	/*
	 * The highest 32 bits of an Error is the character index
	 * of the input that caused the error.
	 *
	 * The next 32 bits are as such
	 * 00000000 00000000 00000000 00000000
	 * Reserved Reserved  Code     Term
	 *
	 *
	 */

	var e uint64

	e |= (uint64(where) << 32)

	e |= (code << 8)

	e |= uint64(term)

	return Error(e)
}

package parse

import "testing"

func TestMakeError(t *testing.T) {
	var (
		wantErr  = ErrUnexpectedEndOfTerm
		wantPos  = uint32(184299)
		wantTerm = List
	)
	e := MakeError(wantErr, wantPos, wantTerm)
	if got := e.what(); wantErr != got {
		t.Fatalf("want error: %d, got: %d, (%x)", wantErr, got, e)
	}

	if got := e.Where(); wantPos != got {
		t.Fatalf("want pos: %d, got: %d", wantPos, got)
	}

	if got := e.term(); wantTerm != got {
		t.Fatalf("want term: %d, got: %d", wantTerm, got)
	}

}

func TestIsError(t *testing.T) {
	expectErrors := []uint64{
		ErrUnexpectedEndOfTerm, ErrTermDepthLimit, ErrConfusion, ErrInputTooLong,
	}

	for i := 0; i < len(expectErrors); i++ {
		e := MakeError(expectErrors[i], 0, 0)

		if !e.IsError() {
			t.Fatalf("expected %d to be error", expectErrors[i])
		}
	}

	e := ErrOk
	if e.IsError() {
		t.Fatalf("ErrOk should not be an error")
	}
}

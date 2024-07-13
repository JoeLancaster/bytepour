package aot

import "github.com/joelancaster/bytepour/pkg/bencode/parse"

const maxDepth = 256

type stack struct {
	st [maxDepth]parse.Term
	sp int
}

func (s *stack) topType() parse.Term {
	return s.st[s.sp]
}

func (s *stack) depth() int {
	return s.sp
}

func (s *stack) pop() {
	s.sp--

	if s.sp < 0 {
		panic("negative stack pointer")
	}

	return
}

func (s *stack) push(t parse.Term) {
	s.sp++

	if s.sp > maxDepth {
		panic("term limit reached")
	}

	s.st[s.sp] = t
}

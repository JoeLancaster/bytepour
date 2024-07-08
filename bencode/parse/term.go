package parse

type Term byte

const (
	String       = Term(0)
	Int          = Term(1)
	List         = Term(2)
	Dict         = Term(3)
	StringHeader = Term(4)
)

const (
	OpenList = byte('l')
	OpenDict = byte('d')
	OpenInt  = byte('i')
	EndTerm  = byte('e')
)

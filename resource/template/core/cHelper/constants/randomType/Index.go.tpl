package randomType

type Value int

const (
	All Value = iota
	Num
	Char
	LowerChar
	UpperChar
	Symbol
	NumAndChar
	CharAndSymbol
)

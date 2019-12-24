package token

const (
	ILLEGAL Type = "ILLEGAL"
	EOF          = "EOF"

	// IDENTIFIERS, LITERALS
	IDENTIFIER = "IDENTIFIER" // "subtract", "foo", "bar"..
	INTEGER    = "INTEGER"    // 1, 5, 1231...

	// OPERATORS
	ASSIGN = "="
	PLUS   = "+"

	// DELIMITERS
	COMMA            = ","
	SEMICOLON        = ";"
	LEFTPARENTHESIS  = "("
	RIGHTPARENTHESIS = ")"
	LEFTBRACKET      = "{"
	RIGHTBRACKET     = "}"

	// KEYWORDS
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var reservedKeywords = map[string]Type{
	"fn" : FUNCTION,
	"let" : LET,
}

type Type string

type Token struct {
	Typ     Type
	Literal string
}

func LookupIdentified(s string) Type {
	if typ, ok := reservedKeywords[s]; ok {
		return typ
	}
	return IDENTIFIER
}

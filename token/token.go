package token

const (
	// Meta
	ILLEGAL Type = "ILLEGAL"
	EOF          = "EOF"

	// Idettifiers, literals
	IDENTIFIER = "IDENTIFIER" // "subtract", "foo", "bar"..
	INTEGER    = "INTEGER"    // 1, 5, 1231...

	// OPERATORS
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LESST    = "<"
	GREATERT = ">"
	EQUAL    = "=="
	NEQUAL   = "!="

	// Delimiters
	COMMA            = ","
	SEMICOLON        = ";"
	LEFTPARENTHESIS  = "("
	RIGHTPARENTHESIS = ")"
	LEFTBRACKET      = "{"
	RIGHTBRACKET     = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var reservedKeywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Type represents the type of the token.
type Type string

// Token aggregates the Type and its Literal
// to make the process of parsing easier.
type Token struct {
	Typ     Type
	Literal string
}

// LookupIdentifier checks whether s is a reserved keyword
// otherwise returns that it is a IDENTIFIER.
func LookupIdentifier(s string) Type {
	if typ, ok := reservedKeywords[s]; ok {
		return typ
	}
	return IDENTIFIER
}

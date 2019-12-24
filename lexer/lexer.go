package lexer

import (
	"github.com/despire/interpreter/token"
	"unicode"
)

const (
	// NULL char represents EOF in the input string
	NULL = 0
)

// Lexer is used to parse the input into individual tokens.
type Lexer struct {
	input        string
	position     int  // current reading position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	char         byte // current character
}

// New returns an initialized Lexer on the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	// init fields
	l.readChar()

	return l
}

// readChar advances the pointers in the input buffer to the next character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// skipWhitespace advances the pointers in the input buffer to the next non-whitespace character.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.char)) || l.char == '\n' {
		l.readChar()
	}
}

// NextToken returns the next token in the input buffer.
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	// if the current pointer is on a whitespace
	// skip it.
	l.skipWhitespace()

	switch l.char {
	case charFromToken(token.LEFTPARENTHESIS):
		t = token.Token{Typ: token.LEFTPARENTHESIS, Literal: string(l.char)}
	case charFromToken(token.RIGHTPARENTHESIS):
		t = token.Token{Typ: token.RIGHTPARENTHESIS, Literal: string(l.char)}
	case charFromToken(token.LEFTBRACKET):
		t = token.Token{Typ: token.LEFTBRACKET, Literal: string(l.char)}
	case charFromToken(token.RIGHTBRACKET):
		t = token.Token{Typ: token.RIGHTBRACKET, Literal: string(l.char)}
	case charFromToken(token.SEMICOLON):
		t = token.Token{Typ: token.SEMICOLON, Literal: string(l.char)}
	case charFromToken(token.COMMA):
		t = token.Token{Typ: token.COMMA, Literal: string(l.char)}
	case charFromToken(token.PLUS):
		t = token.Token{Typ: token.PLUS, Literal: string(l.char)}
	case charFromToken(token.ASSIGN):
		t = token.Token{Typ: token.ASSIGN, Literal: string(l.char)}
	case NULL:
		t = token.Token{Typ: token.EOF, Literal: string(l.char)}
	default:
		switch {
		case isLetter(l.char):
			curr := l.position

			for isLetter(l.char) {
				l.readChar()
			}

			literal := l.input[curr:l.position]

			t = token.Token{Typ: token.LookupIdentified(literal), Literal: literal}

			// the pointer in the buffer is set to the first non ascii character
			// so we just return the token.
			return t
		case unicode.IsDigit(rune(l.char)):
			curr := l.position

			for unicode.IsDigit(rune(l.char)) {
				l.readChar()
			}

			literal := l.input[curr:l.position]

			t = token.Token{Typ: token.INTEGER, Literal: literal}

			// same as above.
			return t
		default:
			t = token.Token{Typ: token.ILLEGAL, Literal: string(l.char)}
		}
	}

	// advance in the buffer
	l.readChar()

	return t
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func charFromToken(typ string) byte {
	return typ[0]
}

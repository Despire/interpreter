package lexer

import (
	"github.com/despire/interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
 x + y;
};

let result = add(five, ten);`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFTPARENTHESIS, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHTPARENTHESIS, ")"},
		{token.LEFTBRACKET, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHTBRACKET, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFTPARENTHESIS, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHTPARENTHESIS, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, "\x00"},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		if tok.Typ != tt.expectedType {
			t.Errorf("token type mismatch, have=%q, want=%q", tok.Typ, tt.expectedType)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Errorf("token literal mismatch, have=%q, want=%q", tok.Literal, tt.expectedLiteral)
		}
	}
}

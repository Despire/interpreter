package parser

import (
	"github.com/despire/interpreter/ast"
	"github.com/despire/interpreter/lexer"
	"github.com/despire/interpreter/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statement: []ast.Statement {
			&ast.LetStatement{
				Token: token.Token{
					Typ: token.LET,
					Literal: "let",
				},
				Identifier: &ast.Identifier{
					Token: token.Token{
						Typ: token.IDENTIFIER,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Expression: &ast.Identifier{
					Token: token.Token{
						Typ: token.IDENTIFIER,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("String() = %v, want  %v", program.String(), "let myVar = anotherVar;")
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. have=%d", len(program.Statement))
	}

	for _, statement := range program.Statement {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, have=%T", statement)
			continue
		}

		if returnStatement.Literal() != "return" {
			t.Errorf("returnStatement.Literal() not 'return', have %q", returnStatement.Literal())
		}
	}

}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() failed, returned nil")
	}

	if len(program.Statement) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. have=%d", len(program.Statement))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statement[i]

		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.Literal() != "let" {
		t.Errorf("s.Literal() = %v, want %v", s.Literal(), "let")
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("node is not of type *ast.LetStatement. have=%T", s)
		return false
	}

	if letStatement.Identifier.Value != name {
		t.Errorf("letStatement.Identifier.Value not %q. have=%s", name, letStatement.Identifier.Value)
		return false
	}

	if letStatement.Identifier.Literal() != name {
		t.Errorf("letStatement.Identifier.Literal not %q, have=%s", name, letStatement.Identifier.Literal())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

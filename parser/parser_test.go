package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Despire/interpreter/ast"
	"github.com/Despire/interpreter/lexer"
	"github.com/Despire/interpreter/token"
)

func TestOperatorPrecedenceParse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"add(a, b, 1, 2* 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		in    string
		left  int
		op    string
		right int
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for i, tt := range tests {
		t.Run("infix-expressions "+strconv.Itoa(i), func(t *testing.T) {
			l := lexer.New(tt.in)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statement) != 1 {
				t.Fatalf("program.Statements does not contain %d statements, have %d", 1, len(program.Statement))
			}

			statement, ok := program.Statement[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, have = %T", program.Statement[0])
			}

			expression, ok := statement.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("expression is not *ast.InfixExpression, have = %T", statement.Expression)
			}

			if !testIntegerLiteral(t, expression.Left, tt.left) {
				return
			}

			if expression.Operator != tt.op {
				t.Fatalf("expression.Operator is not %q, have = %s", tt.op, expression.Operator)
			}

			if !testIntegerLiteral(t, expression.Right, tt.right) {
				return
			}
		})
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		in  string
		op  string
		val int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for i, tt := range tests {
		t.Run("prefix-tests-"+strconv.Itoa(i), func(t *testing.T) {
			l := lexer.New(tt.in)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statement) != 1 {
				t.Fatalf("program.Statements does not contain %d statemtns, have %d", 1, len(program.Statement))
			}

			statement, ok := program.Statement[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not *ast.Expression, have = %T", program.Statement[0])
			}

			expression, ok := statement.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("statement is not ast.PrefixExpression have = %T", statement.Expression)
			}

			if expression.Operator != tt.op {
				t.Fatalf("expression.Operator is not %q, have %q", tt.op, expression.Operator)
			}

			if !testIntegerLiteral(t, expression.Right, int(tt.val)) {
				return
			}
		})
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 1 {
		t.Fatalf("program statements mismatch, have %d, want %d", len(program.Statement), 1)
	}

	statement, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, have = %T", program.Statement[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. have = %T", statement.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("identifier.Value not %d, have %d", 5, literal.Value)
	}
	if literal.Literal() != "5" {
		t.Errorf("identifier.Literal() not %s, have = %s", "5", literal.Literal())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 1 {
		t.Fatalf("program statements mismatch, have %d, want %d", len(program.Statement), 1)
	}

	statement, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, have = %T", program.Statement[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. have = %T", statement.Expression)
	}
	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value not %s, have %s", "foobar", identifier.Value)
	}
	if identifier.Literal() != "foobar" {
		t.Errorf("identifier.Literal() not %s, have = %s", "foobar", identifier.Literal())
	}
}

func TestString(t *testing.T) {
	program := &ast.Program{
		Statement: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Typ:     token.LET,
					Literal: "let",
				},
				Identifier: &ast.Identifier{
					Token: token.Token{
						Typ:     token.IDENTIFIER,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Expression: &ast.Identifier{
					Token: token.Token{
						Typ:     token.IDENTIFIER,
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

func testIdentifier(t *testing.T, expression ast.Expression, val string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not *ast.Identifier have = %T", expression)
		return false
	}

	if identifier.Value != val {
		t.Errorf("identifier.Value is not %q, have = %q", val, identifier.Value)
		return false
	}

	if identifier.Literal() != val {
		t.Errorf("identifier.Literal() is not %q, have = %q", val, identifier.Literal())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int) bool {
	v, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integer literal not *ast.IntegerLiteral, have = %T", il)
		return false
	}

	if v.Value != val {
		t.Errorf("integer.Value not %d, have %d", val, v.Value)
		return false
	}

	if v.Literal() != fmt.Sprintf("%d", val) {
		t.Errorf("integer.Literal not %d, have %s", val, v.Literal())
		return false
	}

	return true
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

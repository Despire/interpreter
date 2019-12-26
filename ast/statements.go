package ast

import (
	"github.com/despire/interpreter/token"
	"strings"
)

type (
	// Identifier represents a value that is binded
	// to a name.
	Identifier struct {
		Token token.Token
		Value string
	}

	// IntegerLiteral represents and integer expression.
	IntegerLiteral struct {
		Token token.Token
		Value int
	}

	// LetStatement consists of
	// a identified (the LHS of the statement)
	// and an expression (the RHS of the statement).
	// In the case of the LetStatement the Identifier will have
	// no value, since the value will be assigned after the
	// evaluation of the statement.
	LetStatement struct {
		Identifier *Identifier
		Expression Expression
		Token      token.Token
	}

	// ReturnStatement consists of the
	// expression it should return.
	ReturnStatement struct {
		Token      token.Token
		Expression Expression
	}

	// ExpressionStatement represents a single line
	// that consist only of a single expression (e.g y + 2;)
	ExpressionStatement struct {
		Token      token.Token
		Expression Expression
	}
)

// implement Statement interface for type checking.
func (s *LetStatement) statement()      {}
func (s *LetStatement) Literal() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	buff := new(strings.Builder)

	buff.WriteString(s.Literal() + " ")
	buff.WriteString(s.Identifier.String())
	buff.WriteString(" = ")

	if s.Expression != nil {
		buff.WriteString(s.Expression.String())
	}

	buff.WriteString(";")

	return buff.String()
}

// implement Expression interface for type checking.
func (i *Identifier) expression()     {}
func (i *Identifier) Literal() string { return i.Token.Literal }
func (i *Identifier) String() string  { return i.Value }

// implement Expression interface for type checking.
func (il *IntegerLiteral) expression()     {}
func (il *IntegerLiteral) Literal() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string  { return il.Token.Literal }

// implement Statement interface for type checking.
func (r *ReturnStatement) statement()      {}
func (r *ReturnStatement) Literal() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	buff := new(strings.Builder)

	buff.WriteString(r.Literal() + " ")
	if r.Expression != nil {
		buff.WriteString(r.Expression.String())
	}

	buff.WriteString(";")

	return buff.String()
}

// implement Statement interface for type checking.
func (e *ExpressionStatement) statement()      {}
func (e *ExpressionStatement) Literal() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

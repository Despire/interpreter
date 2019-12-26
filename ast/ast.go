package ast

import (
	"strings"
)

type (
	// PrefixParseHandler handles the prefix operators.
	PrefixParseHandler func() Expression

	// InfixParseHandler handles infix operators.
	InfixParseHandler func(expression Expression) Expression

	// Node represents top level node in ast.
	Node interface {
		// Literal returns the literal value associated with the token.
		Literal() string
		// String returns the string representation of the node.
		String() string
	}

	// Statement represents statements in the program.
	Statement interface {
		Node

		// Used to distinguish between statements and expressions.
		statement()
	}

	// Expression represents expressions in the program.
	Expression interface {
		Node

		// Used to distinguish between statements and expressions.
		expression()
	}

	// Program represents the root node of the ast.
	Program struct {
		Statement []Statement // statement nodes.
	}
)

func (p *Program) Literal() string {
	if len(p.Statement) > 0 {
		return p.Statement[0].Literal()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	buff := new(strings.Builder)

	for _, s := range p.Statement {
		buff.WriteString(s.String())
	}

	return buff.String()
}

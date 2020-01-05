package ast

import (
	"strings"

	"github.com/despire/interpreter/token"
)

type (
	// FunctionLiteral represents a function
	// expression.
	FunctionLiteral struct {
		Token      token.Token
		Parameters []*Identifier
		Body       *BlockStatement
	}

	// Block statement represents a series
	// of statements wrapped in a '{}'
	BlockStatement struct {
		Token      token.Token
		Statements []Statement
	}

	// Identifier represents a value that is binded
	// to a name.
	Identifier struct {
		Token token.Token
		Value string
	}

	// IntegerLiteral represents an integer expression.
	IntegerLiteral struct {
		Token token.Token
		Value int
	}

	// BooleanLiteral represents an boolean expression.
	BooleanLiteral struct {
		Token token.Token
		Value bool
	}

	// IfExpression represents and if/else expression.
	IfExpression struct {
		Token       token.Token
		Condition   Expression
		Consequence *BlockStatement
		Alternative *BlockStatement
	}

	// CallExpression represents a function
	// call expresion.
	CallExpression struct {
		Token     token.Token
		Function  Expression
		Arguments []Expression
	}

	// PrefixExpression represents an operator
	// that is allowed to prefix an expression.
	PrefixExpression struct {
		Token    token.Token
		Operator string
		Right    Expression
	}

	// InfixExpression represents an binary
	// operator that contains a left, right expression.
	InfixExpression struct {
		Token    token.Token
		Left     Expression
		Operator string
		Right    Expression
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

// implement the Expression interface for type checking.
func (fl *CallExpression) expression()     {}
func (fl *CallExpression) Literal() string { return fl.Token.Literal }
func (fl *CallExpression) String() string {
	buff := new(strings.Builder)

	args := []string{}
	for _, a := range fl.Arguments {
		args = append(args, a.String())
	}

	buff.WriteString(fl.Function.String())
	buff.WriteString("(")
	buff.WriteString(strings.Join(args, ", "))
	buff.WriteString(")")

	return buff.String()
}

// implement the Expression interface for type checking.
func (fl *FunctionLiteral) expression()     {}
func (fl *FunctionLiteral) Literal() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	buff := new(strings.Builder)

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	buff.WriteString(fl.Literal())
	buff.WriteString("(")
	buff.WriteString(strings.Join(params, ", "))
	buff.WriteString(") ")
	buff.WriteString(fl.Body.String())
}

// implement the Statement interface for type checking.
func (bs *BlockStatement) statement()      {}
func (bs *BlockStatement) Literal() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	buff := new(strings.Builder)

	for _, s := range bs.Statements {
		buff.WriteString(s.String())
	}

	return buff.String()
}

// implement Expression interface for type checking.
func (ie *IfExpression) expression()     {}
func (ie *IfExpression) Literal() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	buff := new(strings.Builder)

	buff.WriteString("if")
	buff.WriteString(ie.Condition.String())
	buff.WriteString(" ")
	buff.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		buff.WriteString(" ")
		buff.WriteString(ie.Alternative.String())
	}

	return buff.String()
}

// implement Expression interface for type checking.
func (bl *BooleanLiteral) expression()     {}
func (bl *BooleanLiteral) Literal() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string  { return bl.Token.Literal }

// implement Expression interface for type checking.
func (ie *InfixExpression) expression()     {}
func (ie *InfixExpression) Literal() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	buff := new(strings.Builder)

	buff.WriteString("(")
	buff.WriteString(ie.Left.String())
	buff.WriteString(" " + ie.Operator + " ")
	buff.WriteString(ie.Right.String())
	buff.WriteString(")")

	return buff.String()
}

// implement expression interface for type checking.
func (pe *PrefixExpression) expression()     {}
func (pe *PrefixExpression) Literal() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	buff := new(strings.Builder)

	buff.WriteString("(")
	buff.WriteString(pe.Operator)
	buff.WriteString(pe.Right.String())
	buff.WriteString(")")

	return buff.String()
}

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

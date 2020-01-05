package eval

import (
	"github.com/despire/interpreter/ast"
	"github.com/despire/interpreter/objects"
	"github.com/despire/interpreter/token"
)

var (
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
	NULL  = &objects.Null{}
)

func Eval(node ast.Node) objects.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statement)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &objects.Integer{
			Value: int64(node.Value),
		}
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		return evalPrefix(node.Operator, Eval(node.Right))
	case *ast.InfixExpression:
		return evalInfix(node.Operator, Eval(node.Left), Eval(node.Right))
	}

	return nil
}

func evalIntegerInfix(op string, left objects.Object, right objects.Object) objects.Object {
	lVal := left.(*objects.Integer).Value
	rVal := right.(*objects.Integer).Value

	switch op {
	case token.PLUS:
		return &objects.Integer{
			Value: int64(lVal) + int64(rVal),
		}
	case token.MINUS:
		return &objects.Integer{
			Value: int64(lVal) - int64(rVal),
		}
	case token.ASTERISK:
		return &objects.Integer{
			Value: int64(lVal) * int64(rVal),
		}
	case token.SLASH:
		return &objects.Integer{
			Value: int64(lVal) / int64(rVal),
		}
	default:
		return NULL
	}
}

func evalInfix(op string, left objects.Object, right objects.Object) objects.Object {
	switch {
	case left.Type() == objects.INTEGER && right.Type() == objects.INTEGER:
		return evalIntegerInfix(op, left, right)
	default:
		return NULL
	}
}

func evalBang(exp objects.Object) objects.Object {
	switch exp {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinus(exp objects.Object) objects.Object {
	if exp.Type() != objects.INTEGER {
		return NULL
	}

	return &objects.Integer{
		Value: -(exp.(*objects.Integer).Value),
	}
}

func evalPrefix(op string, exp objects.Object) objects.Object {
	switch op {
	case token.BANG:
		return evalBang(exp)
	case token.MINUS:
		return evalMinus(exp)
	default:
		return NULL
	}
}

func evalStatements(statements []ast.Statement) objects.Object {
	var result objects.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

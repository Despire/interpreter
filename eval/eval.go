package eval

import (
	"fmt"

	"github.com/Despire/interpreter/ast"
	"github.com/Despire/interpreter/objects"
	"github.com/Despire/interpreter/token"
)

var (
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
	NULL  = &objects.Null{}
)

func Eval(node ast.Node, env *objects.Environment) objects.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statement, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.Expression, env)
		if isError(val) {
			return val
		}
		return &objects.Return{
			Value: val,
		}
	case *ast.LetStatement:
		val := Eval(node.Expression, env)
		if isError(val) {
			return val
		}
		env.Set(node.Identifier.Value, val)
	case *ast.IntegerLiteral:
		return &objects.Integer{
			Value: int64(node.Value),
		}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.FunctionLiteral:
		return &objects.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.BlockStatement:
		return evalBlock(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.PrefixExpression:
		exp := Eval(node.Right, env)
		if isError(exp) {
			return exp
		}
		return evalPrefix(node.Operator, exp)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfix(node.Operator, left, right)
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressionList(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)
	}

	return nil
}

func unwrapreturnValue(o objects.Object) objects.Object {
	if ret, ok := o.(*objects.Return); ok {
		return ret.Value
	}

	return o
}

func extendFunctionEnv(fn *objects.Function, args []objects.Object) *objects.Environment {
	env := objects.NewEnclosedEnvironment(fn.Env)

	for i, p := range fn.Parameters {
		env.Set(p.Value, args[i])
	}

	return env
}

func applyFunction(fn objects.Object, args []objects.Object) objects.Object {
	function, ok := fn.(*objects.Function)
	if !ok {
		return newError(fmt.Sprintf("not a function: %s", fn.Type()))
	}

	eenv := extendFunctionEnv(function, args)
	eval := Eval(function.Body, eenv)
	return unwrapreturnValue(eval)
}

func evalExpressionList(exp []ast.Expression, env *objects.Environment) []objects.Object {
	var result []objects.Object

	for _, e := range exp {
		eval := Eval(e, env)
		if isError(eval) {
			return []objects.Object{eval}
		}

		result = append(result, eval)
	}

	return result
}

func isError(o objects.Object) bool {
	if o != nil {
		return o.Type() == objects.ERROR
	}
	return false
}

func isOk(o objects.Object) bool {
	switch o {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIdentifier(node *ast.Identifier, env *objects.Environment) objects.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError(fmt.Sprintf("identifier not found: " + node.Value))
	}

	return val
}

func evalIfExpression(exp *ast.IfExpression, env *objects.Environment) objects.Object {
	condition := Eval(exp.Condition, env)
	if isError(condition) {
		return condition
	}

	if isOk(condition) {
		return Eval(exp.Consequence, env)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	} else {
		return NULL
	}
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
	case token.LESST:
		if lVal < rVal {
			return TRUE
		}
		return FALSE
	case token.GREATERT:
		if lVal > rVal {
			return TRUE
		}
		return FALSE
	case token.EQUAL:
		if lVal == rVal {
			return TRUE
		}
		return FALSE
	case token.NEQUAL:
		if lVal != rVal {
			return TRUE
		}
		return FALSE
	default:
		return newError(fmt.Sprintf("unknown operator: %s %s %s", left.Type(), op, right.Type()))
	}
}

func evalInfix(op string, left objects.Object, right objects.Object) objects.Object {
	switch {
	case left.Type() == objects.INTEGER && right.Type() == objects.INTEGER:
		return evalIntegerInfix(op, left, right)
	case op == token.EQUAL:
		if left == right {
			return TRUE
		}
		return FALSE
	case op == token.NEQUAL:
		if left != right {
			return TRUE
		}
		return FALSE
	case left.Type() != right.Type():
		return newError(fmt.Sprintf("type mismatch: %s %s %s", left.Type(), op, right.Type()))
	default:
		return newError(fmt.Sprintf("unknown operator: %s %s %s", left.Type(), op, right.Type()))
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
		return newError(fmt.Sprintf("unknown operator: -%s", exp.Type()))
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
		return newError(fmt.Sprintf("unknown operator: %s%s", op, exp.Type()))
	}
}

func evalProgram(statements []ast.Statement, env *objects.Environment) objects.Object {
	var result objects.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *objects.Return:
			return result.Value
		case *objects.Error:
			return result
		}
	}

	return result
}

func evalBlock(block *ast.BlockStatement, env *objects.Environment) objects.Object {
	var result objects.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil && (result.Type() == objects.RETURN || result.Type() == objects.ERROR) {
			return result
		}
	}

	return result
}

func newError(s string) *objects.Error {
	return &objects.Error{
		Value: s,
	}
}

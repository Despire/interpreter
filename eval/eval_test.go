package eval

import (
	"github.com/despire/interpreter/lexer"
	"github.com/despire/interpreter/objects"
	"github.com/despire/interpreter/parser"
	"testing"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		have := testEval(tt.input)
		testBooleanObject(t, have, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		have := testEval(tt.input)
		testBooleanObject(t, have, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, tt := range tests {
		have := testEval(tt.input)
		testIntegerObject(t, have, tt.expected)
	}
}

func testEval(input string) objects.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj objects.Object, expected int64) bool {
	result, ok := obj.(*objects.Integer)
	if !ok {
		t.Errorf("object is not Integer. have=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. have=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj objects.Object, expected bool) bool {
	result, ok := obj.(*objects.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. have=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. have=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

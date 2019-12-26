package parser

import (
	"fmt"
	"github.com/despire/interpreter/ast"
	"github.com/despire/interpreter/lexer"
	"github.com/despire/interpreter/token"
	"strconv"
)

type precedence int

const (
	LOWEST precedence = iota
	EQUALS
	LTGT
	SUM
	PRODUCT
	PREFIX
	FNCALL
)

// Parser parses the token from the lexer,
// to create a data structure (ast) to represent
// the source code.
type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	token     token.Token
	peekToken token.Token

	prefixParseHandlers map[token.Type]ast.PrefixParseHandler
	infixParseHandlers  map[token.Type]ast.InfixParseHandler
}

// New returns a initialized parser.
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:               lexer,
		errors:              []string{},
		prefixParseHandlers: map[token.Type]ast.PrefixParseHandler{},
	}

	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)

	// read two token to set 'token', 'peekToken' fields.
	p.nextToken()
	p.nextToken()

	return p
}

// Errors return any errors encountered while parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken advances to the next token from the lexer.
func (p *Parser) nextToken() {
	p.token = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseProgram parses the source code,
// and returns the data structure representing it.
func (p *Parser) ParseProgram() *ast.Program {
	program := new(ast.Program)

	for !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			program.Statement = append(program.Statement, statement)
		}

		p.nextToken()
	}

	return program
}

// parseStatement parses the next statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.token.Typ {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: p.token,
	}

	val, err := strconv.ParseInt(literal.Token.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", literal.Token.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = int(val)

	return literal
}

func (p *Parser) parseExpression(pr precedence) ast.Expression {
	prefix, ok := p.prefixParseHandlers[p.token.Typ]
	if !ok {
		return nil
	}

	leftSide := prefix()

	return leftSide
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token:      p.token,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekToken.Typ == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.token,
	}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: p.token,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Identifier = &ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) registerPrefix(typ token.Type, fn ast.PrefixParseHandler) {
	p.prefixParseHandlers[typ] = fn
}

func (p *Parser) registerInfix(typ token.Type, fn ast.InfixParseHandler) {
	p.infixParseHandlers[typ] = fn
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.token.Typ == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekToken.Typ == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Typ)
	p.errors = append(p.errors, msg)
}

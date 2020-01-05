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

var precedences = map[token.Type]precedence{
	token.EQUAL:           EQUALS,
	token.NEQUAL:          EQUALS,
	token.LESST:           LTGT,
	token.GREATERT:        LTGT,
	token.PLUS:            SUM,
	token.MINUS:           SUM,
	token.SLASH:           PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.LEFTPARENTHESIS: FNCALL,
}

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
		infixParseHandlers:  map[token.Type]ast.InfixParseHandler{},
	}

	p.registerPrefix(token.TRUE, p.parseBool)
	p.registerPrefix(token.FALSE, p.parseBool)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LEFTPARENTHESIS, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NEQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESST, p.parseInfixExpression)
	p.registerInfix(token.GREATERT, p.parseInfixExpression)
	p.registerInfix(token.LEFTPARENTHESIS, p.parseCallExpression)

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

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Typ == token.RIGHTPARENTHESIS {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Typ == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RIGHTPARENTHESIS) {
		return nil
	}

	return args
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:     p.token,
		Function:  fn,
		Arguments: p.parseCallArguments(),
	}
	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Typ == token.RIGHTPARENTHESIS {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifiers = append(identifiers, &ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	})

	for p.peekToken.Typ == token.COMMA {
		p.nextToken()
		p.nextToken()

		identifiers = append(identifiers, &ast.Identifier{
			Token: p.token,
			Value: p.token.Literal,
		})
	}

	if !p.expectPeek(token.RIGHTPARENTHESIS) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{
		Token: p.token,
	}

	if !p.expectPeek(token.LEFTPARENTHESIS) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LEFTBRACKET) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.token,
	}

	p.nextToken()

	for !p.curTokenIs(token.RIGHTBRACKET) && !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.token,
	}

	if !p.expectPeek(token.LEFTPARENTHESIS) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHTPARENTHESIS) {
		return nil
	}

	if !p.expectPeek(token.LEFTBRACKET) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Typ == token.ELSE {
		p.nextToken()

		if !p.expectPeek(token.LEFTBRACKET) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHTPARENTHESIS) {
		return nil
	}

	return expression
}

func (p *Parser) parseBool() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.token,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
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
		p.noPrefixParseError(p.token.Typ)
		return nil
	}

	leftSide := prefix()

	for !(p.peekToken.Typ == token.SEMICOLON) && pr < p.peekPrecedence() {
		infix, ok := p.infixParseHandlers[p.peekToken.Typ]
		if !ok {
			break
		}

		p.nextToken()

		leftSide = infix(leftSide)
	}

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

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Typ == token.SEMICOLON {
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

	p.nextToken()

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Typ == token.SEMICOLON {
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

func (p *Parser) currentPrecedence() precedence {
	if p, ok := precedences[p.token.Typ]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() precedence {
	if p, ok := precedences[p.peekToken.Typ]; ok {
		return p
	}

	return LOWEST
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

func (p *Parser) noPrefixParseError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

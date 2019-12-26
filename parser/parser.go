package parser

import (
	"fmt"
	"github.com/despire/interpreter/ast"
	"github.com/despire/interpreter/lexer"
	"github.com/despire/interpreter/token"
)

// Parser parses the token from the lexer,
// to create a data structure (ast) to represent
// the source code.
type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	token     token.Token
	peekToken token.Token
}

// New returns a initialized parser.
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

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
		return nil
	}
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

func (p *Parser) curTokenIs(t token.Type) bool { return p.token.Typ == t }

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

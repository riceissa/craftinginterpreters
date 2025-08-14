package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			fmt.Printf("error in parse(): %v", err)
		}
		statements = append(statements, stmt)
	}
	return statements
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value := p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return Print{value}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr := p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return Expression{expr}, nil
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(VAR) {
		result, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		} else {
			return result, nil
		}
	}
	result, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, err
	} else {
		return result, nil
	}
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return Var{name, initializer}, nil
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{false}
	}
	if p.match(TRUE) {
		return Literal{true}
	}
	if p.match(NIL) {
		return Literal{nil}
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().literal}
	}

	if p.match(IDENTIFIER) {
		return Variable{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}

	parse_error(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) match(token_types ...TokenType) bool {
	for _, token_type := range token_types {
		if p.check(token_type) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(token_type TokenType, message string) (Token, error) {
	if p.check(token_type) {
		return p.advance(), nil
	}
	parse_error(p.peek(), message)
	err := errors.New("error inside of consume()")
	return Token{}, err
}

func parse_error(token Token, message string) {
	if token.token_type == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().token_type == SEMICOLON {
			return
		}

		switch p.peek().token_type {
		case CLASS:
			return
		case FUN:
			return
		case VAR:
			return
		case FOR:
			return
		case IF:
			return
		case WHILE:
			return
		case PRINT:
			return
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) check(token_type TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().token_type == token_type
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().token_type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

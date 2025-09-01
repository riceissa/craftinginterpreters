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
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		b, err := p.block()
		if err != nil {
			return nil, err
		}
		return Block{b}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr
	var err error
	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}
	return Return{keyword, value}, nil
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt = nil
	if p.match(SEMICOLON) {
		// do nothing
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr = nil
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = Block{
			[]Stmt{body, Expression{increment}},
		}
	}

	if condition == nil {
		condition = &Literal{true}
	}
	body = While{condition, body}

	if initializer != nil {
		body = Block{
			[]Stmt{initializer, body},
		}
	}

	return body, err
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return While{condition, body}, nil
}

func (p *Parser) ifStatement() (If, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return If{}, err
	}
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return If{}, err
	}
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return If{}, err
		}
	}

	return If{condition, thenBranch, elseBranch}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}

	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return Print{value}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return Expression{expr}, nil
}

func (p *Parser) function(kind string) (Function, error) {
	name, err := p.consume(IDENTIFIER, fmt.Sprintf("Expect %v name.", kind))
	if err != nil {
		return Function{}, err
	}
	_, err = p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %v name.", kind))
	if err != nil {
		return Function{}, err
	}
	parameters := []Token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				log_parse_error(p.peek(), "Can't have more than 255 parameters.")
			}
			ident, err := p.consume(IDENTIFIER, "Expect parameters name.")
			if err != nil {
				return Function{}, err
			}
			parameters = append(parameters, ident)
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %v body.", kind))
	body, err := p.block()
	if err != nil {
		return Function{}, err
	}
	return Function{name, parameters, body}, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &Logical{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &Logical{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		_ = p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		switch v := expr.(type) {
		case *Variable:
			var name Token = v.name
			return &Assign{name, value}, nil
		case *Get:
			return &Set{v.object, v.name, value}, nil
		default:
			fmt.Printf("Invalid assignment target: %q\n", stringify(v))
		}
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
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

func (p *Parser) classDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *Variable = nil
	if p.match(LESS) {
		_, err = p.consume(IDENTIFIER, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = &Variable{p.previous()}
		// Kind of weird how the book does it here... Why not just
		// consume and then take the name from the consume call?
		// Why consume without assigning to a name, and then make
		// a separate call to previous?
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []Function{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		f, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, f)
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return Class{name, superclass, methods}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return Var{name, initializer}, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{operator, right}, nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &Get{expr, name}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := []Expr{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				log_parse_error(p.peek(), "Can't have more than 255 arguments.")
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &Call{callee, paren, arguments}, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{false}, nil
	}
	if p.match(TRUE) {
		return &Literal{true}, nil
	}
	if p.match(NIL) {
		return &Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}, nil
	}

	if p.match(THIS) {
		return &This{p.previous()}, nil
	}

	if p.match(IDENTIFIER) {
		return &Variable{p.previous()}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping{expr}, nil
	}

	log_parse_error(p.peek(), "Expect expression.")
	return nil, nil
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
	log_parse_error(p.peek(), message)
	err := errors.New("error inside of consume()")
	return Token{}, err
}

func log_parse_error(token Token, message string) {
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

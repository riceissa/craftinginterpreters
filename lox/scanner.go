package main

import (
	"log"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewScanner(source string) Scanner {
	scanner := Scanner{}
	scanner.source = source
	scanner.line = 1
	return scanner
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	result := s.source[s.current]
	s.current += 1
	return result
}

func (s *Scanner) addSimpleToken(token_type TokenType) {
	s.addToken(token_type, nil)
}

func (s *Scanner) addToken(token_type TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{token_type, text, literal, s.line})
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addSimpleToken(LEFT_PAREN)
	case ')':
		s.addSimpleToken(RIGHT_PAREN)
	case '{':
		s.addSimpleToken(LEFT_BRACE)
	case '}':
		s.addSimpleToken(RIGHT_BRACE)
	case ',':
		s.addSimpleToken(COMMA)
	case '.':
		s.addSimpleToken(DOT)
	case '-':
		s.addSimpleToken(MINUS)
	case '+':
		s.addSimpleToken(PLUS)
	case ';':
		s.addSimpleToken(SEMICOLON)
	case '*':
		s.addSimpleToken(STAR)
	case '!':
		if s.match('=') {
			s.addSimpleToken(BANG_EQUAL)
		} else {
			s.addSimpleToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addSimpleToken(EQUAL_EQUAL)
		} else {
			s.addSimpleToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addSimpleToken(LESS_EQUAL)
		} else {
			s.addSimpleToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addSimpleToken(GREATER_EQUAL)
		} else {
			s.addSimpleToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addSimpleToken(SLASH)
		}
	case ' ':
		// do nothing
	case '\r':
		// do nothing
	case '\t':
		// do nothing
	case '\n':
		s.line++
	case '"':
		s.scan_string()
	default:
		if isDigit(c) {
			s.scan_number()
		} else if isAlpha(c) {
			s.scan_identifier()
		} else {
			log_error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) scan_number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	f, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		log.Fatal(err)
	}
	s.addToken(NUMBER, f)
}

func (s *Scanner) scan_string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		log_error(s.line, "Unterminated string.")
		return
	}

	s.advance() // The closing "

	// Trim the surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) scan_identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	token_type, ok := keywords[text]
	if !ok {
		token_type = IDENTIFIER
	}
	s.addSimpleToken(token_type)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

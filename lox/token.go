package main

import (
	"fmt"
)

type Token struct {
	token_type TokenType
	lexeme     string
	literal    any
	line       int
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", t.token_type, t.lexeme, t.literal)
}

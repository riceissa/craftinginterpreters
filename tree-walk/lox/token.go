package main

import (
	"fmt"
)

type Token struct {
	tokenType TokenType
	lexeme     string
	literal    any
	line       int
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}

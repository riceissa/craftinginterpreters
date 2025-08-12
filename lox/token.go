package main

type Token struct {
	token_type TokenType // can we just call this "type" instead of token_type? or does golang disallow this?
	lexeme     string
	literal    any
	line       int
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", t.token_type, t.lexeme, t.literal)
}

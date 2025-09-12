enum TokenType {
    // Single-character tokens.
    LeftParen, RightParen,
    LeftBrace, RightBrace,
    Comma, Dot, Minus, Plus,
    Semicolon, Slash, Star,
    // One or two character tokens.
    Bang, BangEqual,
    Equal, EqualEqual,
    Greater, GreaterEqual,
    Less, LessEqual,
    // Literals.
    Identifier, String, Number,
    // Keywords.
    And, Class, Else, False,
    For, Fun, If, Nil, Or,
    Print, Return, Super, This,
    True, Var, While,

    Error, Eof
}

struct Token {
    token_type: TokenType,
    content: &str,
    line: isize,
}

struct Scanner {
    source: &str,
    current: &str,
    line: isize,
}

impl Scanner {
    pub fn new(source: &str) -> Self {
        Scanner {
            source: &source[..],
            current: &source[0..1],
            line: 1,
        }
    }

    pub fn is_at_end(&self) -> bool {
        return self.current == '\0'; // TODO: i think i should instead check that the current == the source string's length
    }

    fn make_token(&self, token_type: TokenType) -> Token {
        Token {
            token_type: token_type,
            start: scanner.start,
            length: scanner.current - scanner.start,
            line: scanner.line
        }
    }

    fn error_token(&self, message: &str) -> Token {
        Token {
            token_type: token_type,
            start: scanner.start,
            length: scanner.current - scanner.start,
            line: scanner.line
        }
    }

    pub fn scan_token(&self) -> Token {
        self.current = ;

        if self.is_at_end() { return self.make_token(TokenType::Eof); }

        return self.error_token("Unexpected character.");

    }
}

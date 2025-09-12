pub fn compile(source: &str) {
    let scanner = Scanner::new(source);

    let mut line = -1;
    loop {
        let token: Token = scanner.scan_token();
        if token.line != line {
            print!("{:4} ", token.line);
            line = token.line;
        } else {
            print!("   | ");
        }
        println!("{:2} '{}'", token.token_type, token.content);

        if token.token_type == TokenType::Eof { break; }
    }
}

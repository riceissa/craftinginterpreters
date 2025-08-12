jlox: lox/lox.go lox/scanner.go lox/token.go lox/token_type.go lox/ast_printer.go lox/expr.go lox/parser.go
	go build -o $@ ./lox

.PHONY: clean
clean:
	rm ./jlox

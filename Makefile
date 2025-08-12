jlox: lox/lox.go lox/scanner.go lox/token.go lox/token_type.go
	go build -o $@ ./lox

generate_ast: tool/generate_ast.go
	go build -o $@ ./tool && ./generate_ast ./lox

.PHONY: clean
clean:
	rm ./jlox

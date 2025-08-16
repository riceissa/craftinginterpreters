SRCS=lox/lox.go lox/scanner.go lox/token.go lox/token_type.go lox/ast_printer.go lox/expr.go lox/parser.go lox/interpreter.go lox/stmt.go lox/environment.go

.PHONY: all
all: tags jlox

jlox: $(SRCS)
#	go build -o $@ ./lox

#	Disable optimizations and inlining; makes it easier to step through
#	with the debugger:
	go build -gcflags "all=-N -l" -o $@ ./lox

tags: $(SRCS)
	gotags -f $@ -R ./lox

.PHONY: clean
clean:
	rm ./jlox
	rm tags

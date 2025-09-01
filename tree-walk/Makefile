SRCS=lox/ast_printer.go lox/environment.go lox/expr.go lox/interpreter.go lox/lox_callable.go lox/lox_class.go lox/lox_function.go lox/lox.go lox/lox_instance.go lox/parser.go lox/resolver.go lox/scanner.go lox/stmt.go lox/token.go lox/token_type.go

.PHONY: all
all: tags jlox

jlox: $(SRCS)
#	go build -o $@ ./lox

#	Disable optimizations and inlining; makes it easier to step through
#	with the debugger:
	go build -gcflags "all=-N -l" -o $@ ./lox

tags: $(SRCS)
	ctags -R ./lox
# For some reason gotags hardcodes the line numbers in the tags file, so the
# tags file is very brittle (e.g. if you add some text before the line that you
# want to jump to, you just jump to the line number where the tag *used* to be,
# rather than where it is now), which incentivizes frequent regeneration. I've
# switched to ctags for now, which supports Go and uses regex-based tags which
# are more stable.
#	gotags -f $@ -R ./lox

.PHONY: clean
clean:
	rm ./jlox
	rm tags

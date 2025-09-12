SRCS=src/chunk.rs src/compiler.rs src/debug.rs src/main.rs src/scanner.rs src/value.rs src/vm.rs

.PHONY: all
all: tags main

tags: $(SRCS)
	ctags -R src

main: $(SRCS)
	rustc src/main.rs

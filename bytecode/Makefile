SRCS=src/chunk.rs src/debug.rs src/main.rs src/memory.rs

.PHONY: all
all: tags main

tags: $(SRCS)
	ctags -R .

main: $(SRCS)
	rustc src/main.rs

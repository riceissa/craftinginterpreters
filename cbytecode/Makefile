H_FILES=src/chunk.h src/common.h src/compiler.h src/debug.h src/memory.h src/scanner.h src/value.h src/vm.h

C_FILES=src/chunk.c src/compiler.c src/debug.c src/main.c src/memory.c src/scanner.c src/value.c src/vm.c

.PHONY: all
all: tags main

tags: $(H_FILES) $(C_FILES)
	ctags -R src

main: $(H_FILES) $(C_FILES)
	gcc -Wall -Wextra -g -o $@ $(C_FILES)

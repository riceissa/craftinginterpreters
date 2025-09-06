H_FILES=src/chunk.h src/common.h src/debug.h src/memory.h src/value.h

C_FILES=src/chunk.c src/debug.c src/main.c src/memory.c src/value.c

.PHONY: all
all: tags main

tags: $(H_FILES) $(C_FILES)
	ctags -R src

main: $(H_FILES) $(C_FILES)
	gcc -Wall -Wextra -o $@ $(C_FILES)

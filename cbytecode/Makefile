H_FILES=$(wildcard src/*.h)
C_FILES=$(wildcard src/*.c)

.PHONY: all
all: tags main

tags: $(H_FILES) $(C_FILES)
	ctags -R src

main: $(H_FILES) $(C_FILES)
	gcc -Wall -Wextra -g -o $@ $(C_FILES)

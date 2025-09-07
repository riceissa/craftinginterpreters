# Crafting Interpreters

You can read the book at https://craftinginterpreters.com/

This repo contains my implementations of the Lox language from the book.
There are two implementations: one is a tree-walk interpreter from Part II of
the book, and the other is a bytecode virtual machine from Part III of
the book. (The latter does not actually exist yet.)

## Tree-walk interpreter

The tree-walk interpreter is in the `tree-walk/` directory of the repo.

I mostly just straightforwardly translate the interpreter from Java (which is
what the book uses, for this part) to Go. The biggest differences are:

1. Go doesn't have exceptions, so all of the error-handling is done in the
   Go-style of multiple return values.
2. I don't understand/like the visitor pattern
   ([I am not the only one](https://grugbrain.dev/#grug-on-visitor-pattern)),
   so instead I use switch cases to decide which function to call.
3. I didn't do the meta-programming thing from GenerateAst.java, since Go
   doesn't have as much boilerplate when defining structs.

Everything from Part II of the book has been implemented. (It's possible
that I didn't implement everything _correctly_, however, so there may
be some differences from the reference implementation of the book,
especially in the error-handling.)

## Bytecode virtual machine

There are two versions of the bytecode virtual machine.

* One is the C version from the book, which I am copying as I read along. This
version lives in the `cbytecode/` directory of the repo.

* The other version is written in Rust, and lives in the `bytecode/` directory of
the repo.
The Rust version is cheating a little, by using standard data structures
available in the language such as the `Vec` type for dynamic arrays.
(I originally intended to implement my own version of `Vec` and other
data structures, but it clearly became apparent that this was going to be
way too much work because the compiler really doesn't want you to do "unsafe"
things, and I was getting very tricky errors and couldn't even compile
my code. This is my first time using Rust, so I am hopeful I can come back
and maybe add in my own implementations later on.)
This means
that I don't really do any memory management, instead relying on Rust's RAII
for memory management. (To be clear, I am skeptical of RAII, but I am trying my
best to go along with the language's idioms since this is my first time using
Rust.)

Everything up to and including Chapter 15 is done.

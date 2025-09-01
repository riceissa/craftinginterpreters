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
be some differences, especially in the error-handling, from the reference
implementation of the book.)

## Bytecode virtual machine

The bytecode virtual machine lives in the `bytecode/` directory of the repo.

Work on the bytecode virtual machine has not begun yet.

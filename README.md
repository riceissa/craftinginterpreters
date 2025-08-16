# Crafting Interpreters

You can read the book at https://craftinginterpreters.com/

I mostly just straightforwardly translate the interpreter from Java (which is
what the book uses, for the first part) to Go. The biggest differences are:

1. Go doesn't have exceptions, so all of the error-handling is done in the
   Go-style of multiple return values.
2. I don't understand/like the visitor pattern, so instead I use switch cases
   to decide which function to call.

So far, everything up to and including Chapter 9 is done.

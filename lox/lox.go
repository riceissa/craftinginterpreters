package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "runtime"
)

var hadError = false
var hadRuntimeError = false

type RuntimeError struct {
	token   Token
	message string
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		fmt.Printf("running script %v\n", os.Args[1])
		runFile(os.Args[1])
	} else {
		fmt.Println("doing runPrompt()")
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	run(string(bytes))

	if hadError {
		os.Exit(65)
	}

	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			line := scanner.Text()
			run(line)
			hadError = false
		} else {
			// encountered EOF?
			fmt.Println("EOF detected??")
			break
		}
		// not sure if the following is needed...
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func run(source string) {
	// fmt.Println(source)
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	// for index, tok := range tokens {
	// 	fmt.Println(index, tok)
	// }
	// os.Exit(0)

	parser := Parser{tokens: tokens}
	statements := parser.parse()

	fmt.Print(print_statements(statements, 0))

	if hadError { return }

	// runtime.Breakpoint()
	resolver := NewResolver(&interpreter)
	resolver.resolveStatements(statements)

	// Stop if there was a resolution error.
	if hadError { return }

	// runtime.Breakpoint()
	interpreter.interpret(statements)
	// fmt.Println(print_expr(expression))

}

func log_error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %v: %v\n", line, where, message)
	hadError = true
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %v] %v", e.token.line, e.message)
}

func runtimeError(e RuntimeError) {
	fmt.Fprintln(os.Stderr, e.Error())
	hadRuntimeError = true
}

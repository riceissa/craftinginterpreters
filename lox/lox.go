package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError = false
var hadRuntimeError = false

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
	parser := Parser{tokens: tokens}
	expression := parser.parse()

	if hadError {
		return
	}

	interpret(expression)
	// fmt.Println(print_expr(expression))

	// for index, tok := range tokens {
	// 	fmt.Println(index, tok)
	// }
}

func log_error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %v: %v\n", line, where, message)
	hadError = true
}

func runtimeError(e error) {
	// TODO: our error type needs to have the token/line number to be able to report it.
	fmt.Printf("%v\n[line ??]\n", e)
	hadRuntimeError = true
}

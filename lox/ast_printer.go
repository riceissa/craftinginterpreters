package main

import (
	"fmt"
	"strings"
)

func print_expr(expr Expr) string {
	switch v := expr.(type) {
	case Binary:
		return parenthesize(v.operator.lexeme, v.left, v.right)
	case Grouping:
		return parenthesize("group", v.expression)
	case Literal:
		if v.value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", v.value)
	case Unary:
		return parenthesize(v.operator.lexeme, v.right)
	default:
		return "inside ast_printer.go:print_expr ... don't know how to stringify this guy"
	}
}

func parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(print_expr(expr))
	}
	builder.WriteString(")")
	return builder.String()
}

func test_ast_printer() {
	var expression Expr = Binary{
		Unary{
			Token{MINUS, "-", nil, 1},
			Literal{123},
		},
		Token{STAR, "*", nil, 1},
		Grouping{Literal{45.67}},
	}
	fmt.Println(print_expr(expression))
}

package main

import (
	"fmt"
	"strings"
)

func print_statements(statements []Stmt, indent int) string {
	var builder strings.Builder
	for _, stmt := range statements {
		builder.WriteString(print_stmt(stmt, indent))
		// builder.WriteString("; ")
	}
	return builder.String()
}

func print_indent(indent int) string {
	var builder strings.Builder
	for _ = range indent {
		builder.WriteString(" ")
	}
	return builder.String()

}

func print_stmt(stmt Stmt, indent int) string {
	switch v := stmt.(type) {
	case If:
		return print_indent(indent) + "<if>"
	case Print:
		return print_indent(indent) + "<print>"
	case Expression:
		return fmt.Sprintf("%v<expr: %v>", print_indent(indent), print_expr(v.expression))
	case Var:
		return fmt.Sprintf("%v<var: %v = %v>\n", print_indent(indent), v.name, print_expr(v.initializer))
	case Block:
		return fmt.Sprintf("%v<block: \n%v\n%v>\n", print_indent(indent), print_statements(v.statements, indent + 4), print_indent(indent))
	case While:
		return fmt.Sprintf("%v<while: (%v)\n%v\n%v>\n", print_indent(indent), print_expr(v.condition), print_stmt(v.body, indent + 4), print_indent(indent))
	case Function:
		return print_indent(indent) + "<function>"
	case Return:
		return print_indent(indent) + "<return>"
	default:
		panic(fmt.Sprintf("Unreachable. stmt has value %v; its type is %T which we don't know how to handle.", stmt, stmt))
	}
}

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
	case Variable:
		return "(variable " + v.name.lexeme + " )"
	case Assign:
		return parenthesize(v.name.lexeme, v.value)
	default:
		panic(fmt.Sprintf("Unreachable. expr has value %v; its type is %T which we don't know how to handle.", expr, expr))
	}
}

/*
func blockify(name string, stmts ...Stmt) string {
	var builder strings.Builder
	builder.WriteString("{")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(print_expr(expr))
	}
	builder.WriteString("}")
	return builder.String()
}
*/

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

package main

import (
	"fmt"
	"strings"
)

var environment = NewEnvironment()

func interpret(statements []Stmt) {
	for _, statement := range statements {
		err := execute(statement)
		if err != nil {
			runtimeError(err)
			// fmt.Printf("We have a problem when executing %v:\n%v\n", statement, err)
		}
	}
}

func execute(stmt Stmt) error {
	switch v := stmt.(type) {
	case Print:
		interpret_print_stmt(v)
		return nil
	case Expression:
		interpret_expression_stmt(v)
		return nil
	case Var:
		interpret_var_stmt(v)
		return nil
	default:
		panic(fmt.Sprintf("Unreachable. stmt has value %v; its type is %T which we don't know how to handle.", stmt, stmt))
	}
}

func interpret_var_stmt(stmt Var) error {
	var value any
	var err error
	if stmt.initializer != nil {
		value, err = evaluate(stmt.initializer)
		if err != nil {
			return err
		}
	}

	environment.define(stmt.name.lexeme, value)
	return nil
}

func interpret_variable_expr(expr Variable) (any, error) {
	return environment.get(expr.name)
}

func evaluate(expr Expr) (any, error) {
	switch v := expr.(type) {
	case Binary:
		return interpret_binary_expr(v)
	case Grouping:
		return interpret_grouping_expr(v)
	case Literal:
		return interpret_literal_expr(v)
	case Unary:
		return interpret_unary_expr(v)
	case Variable:
		return interpret_variable_expr(v)
	default:
		panic(fmt.Sprintf("Unreachable. expr has value %v; its type is %T which we don't know how to handle.", expr, expr))
	}
}

func interpret_binary_expr(expr Binary) (any, error) {
	left, _ := evaluate(expr.left)
	right, _ := evaluate(expr.right)

	switch expr.operator.token_type {
	case GREATER:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64), err
	case GREATER_EQUAL:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64), err
	case LESS:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) < right.(float64), err
	case LESS_EQUAL:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) <= right.(float64), err
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case MINUS:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) - right.(float64), err
	case PLUS:
		leftFloat, leftIsFloat := left.(float64)
		rightFloat, rightIsFloat := right.(float64)
		if leftIsFloat && rightIsFloat {
			return leftFloat + rightFloat, nil
		}

		leftString, leftIsString := left.(string)
		rightString, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return leftString + rightString, nil
		}

		return nil, fmt.Errorf("Operands must be two numbers or two strings.", expr.operator)
	case SLASH:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) / right.(float64), err
	case STAR:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) * right.(float64), err
	}

	return nil, fmt.Errorf("Reached the unreachable.")
}

func interpret_grouping_expr(expr Grouping) (any, error) {
	result, err := evaluate(expr.expression)
	return result, err
}

func interpret_literal_expr(expr Literal) (any, error) {
	return expr.value, nil
}

func interpret_unary_expr(expr Unary) (any, error) {
	right, err := evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.token_type {
	case BANG:
		return !isTruthy(right), nil
	case MINUS:
		err := checkNumberOperand(expr.operator, right)
		return -(right.(float64)), err
	}

	return nil, fmt.Errorf("Reached the unreachable.")
}

func interpret_expression_stmt(stmt Expression) {
	evaluate(stmt.expression)
}

func interpret_print_stmt(stmt Print) error {
	value, err := evaluate(stmt.expression)
	if err != nil {
		return err
	}
	fmt.Println(stringify(value))
	return nil
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	switch object.(type) {
	case float64:
		text := fmt.Sprintf("%v", object)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	default:
		return fmt.Sprintf("%v", object)
	}
}

func checkNumberOperand(operator Token, operand any) error {
	switch operand.(type) {
	case float64:
		return nil
	default:
		return fmt.Errorf("Operand %q must be a number.", operand)
	}
}

func checkNumberOperands(operator Token, left any, right any) error {
	_, leftIsFloat := left.(float64)
	_, rightIsFloat := right.(float64)
	if leftIsFloat && rightIsFloat {
		return nil
	}
	return fmt.Errorf("Operands %q must be numbers.", operator)
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	switch object.(type) {
	case bool:
		return object.(bool)
	default:
		return true
	}
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	// TODO: might need to use something like reflect.DeepEqual
	// or whatever; I don't know if Go's native == operator
	// is what we want here.
	return a == b
}

package main

import (
	"log"
	"fmt"
	"strings"
	"errors"
)

func interpret(expression Expr) {
	value, err := evaluate(expression)
	if err != nil {
		log.Fatal(fmt.Sprintf("we have a problem when interpreting %v", expression))
	} else {
		fmt.Println(stringify(value))
	}
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	switch object.(type) {
	case float64:
		text := fmt.Sprintf("%v", object)
		if strings.HasSuffix(text, ".0") {
			text = text[0:len(text) - 2]
		}
		return text
	default:
		return fmt.Sprintf("%v", object)
	}
}

func interpret_literal_expr(expr Literal) any {
	return expr.value
}

func interpret_unary_expr(expr Unary) any {
	right, _ := evaluate(expr.right)

	switch expr.operator.token_type {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		checkNumberOperand(expr.operator, right)
		return -(right.(float64))
	}

	// Unreachable
	return nil
}

func checkNumberOperand(operator Token, operand any) error {
	switch operand.(type) {
	case float64:
		return nil
	default:
		return errors.New(fmt.Sprintf("Operand %v must be a number.", operand))
	}
}

func checkNumberOperands(operator Token, left any, right any) error {
	_, leftIsFloat := left.(float64)
	_, rightIsFloat := right.(float64)
	if leftIsFloat && rightIsFloat {
		return nil
	}
	return errors.New(fmt.Sprintf("Operands %v must be numbers.", operator))
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
	if (a == nil && b == nil) {
		return true
	}
	if (a == nil) {
		return false
	}
	// TODO: might need to use something like reflect.DeepEqual
	// or whatever; I don't know if Go's native == operator
	// is what we want here.
	return a == b
}

func interpret_grouping_expr(expr Grouping) any {
	result, _ := evaluate(expr.expression)
	return result
}

func evaluate(expr Expr) (any, error) {
	return expr, nil // ?? idk what to do here yet
}

func interpret_binary_expr(expr Binary) any {
	left, _ := evaluate(expr.left)
	right, _ := evaluate(expr.right)

	switch expr.operator.token_type {
	case GREATER:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case MINUS:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) - right.(float64)
	case PLUS:
		leftFloat, leftIsFloat := left.(float64)
		rightFloat, rightIsFloat := right.(float64)
		if leftIsFloat && rightIsFloat {
			return leftFloat + rightFloat
		}

		leftString, leftIsString := left.(string)
		rightString, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return leftString + rightString
		}

		return fmt.Errorf("Operands must be two numbers or two strings.", expr.operator)
	case SLASH:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) * right.(float64)
	}

	// Unreachable.
	return nil
}

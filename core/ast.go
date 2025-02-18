package core

import (
	"fmt"
)

type Expr interface {
	fmt.Stringer
	Interpret(environment map[string]interface{}) interface{}
}

type BinaryExpr struct {
	Expr
	rightExpr Expr
	operator  Token
	leftExpr  Expr
}

type UnaryExpr struct {
	Expr
	operator Token
	expr     Expr
}

type GroupingExpr struct {
	Expr
	expr Expr
}

type LiteralExpr struct {
	Expr
	value Token
}

type AssignExpr struct {
	Expr
	name Token
	expr Expr
}

type InvalidExpr struct {
	Expr
}

func (expr AssignExpr) String() string {
	return fmt.Sprintf("(%v = %v)", expr.name.String(), expr.expr.String())
}
func (expr BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", expr.leftExpr.String(), expr.operator.String(), expr.rightExpr.String())
}
func (expr UnaryExpr) String() string {
	return fmt.Sprintf("(%v%v)", expr.operator.String(), expr.expr.String())
}
func (expr GroupingExpr) String() string {
	return fmt.Sprintf("(%v)", expr.expr.String())
}
func (expr LiteralExpr) String() string {
	return expr.value.String()
}
func (expr InvalidExpr) String() string {
	return "Invalid Expression"
}

func (expr AssignExpr) Interpret(environment map[string]interface{}) interface{} {
	data := expr.expr.Interpret(environment)
	environment[expr.name.String()] = data
	return data
}

func (expr BinaryExpr) Interpret(environment map[string]interface{}) interface{} {
	left := expr.leftExpr.Interpret(environment)
	right := expr.rightExpr.Interpret(environment)
	l, _ := left.(float64)
	r, _ := right.(float64)

	switch expr.operator.Type {
	case PLUS:
		return l + r
	case MINUS:
		return l - r
	case STAR:
		return l * r
	case SLASH:
		return l / r
	case EQUAL_EQ:
		return left == right
	case BANG_EQ:
		return left != right
	case GREATER:
		return l > r
	case GREATER_EQ:
		return l >= r
	case LESS:
		return l < r
	case LESS_EQ:
		return l <= r
	default:
		return nil
	}
}

func (expr UnaryExpr) Interpret(environment map[string]interface{}) interface{} {
	res := expr.expr.Interpret(environment)
	switch expr.operator.Type {
	case BANG:
		return !res.(bool)
	case MINUS:
		return -res.(float64)
	default:
		return nil
	}
}

func (expr GroupingExpr) Interpret(environment map[string]interface{}) interface{} {
	return expr.expr.Interpret(environment)
}

func (expr LiteralExpr) Interpret(environment map[string]interface{}) interface{} {
	switch expr.value.Type {
	case TRUE:
		return true
	case FALSE:
		return false
	case IDENTIFIER:
		return environment[expr.value.String()]
	default:
		return expr.value.Value
	}
}

func (expr InvalidExpr) Interpret(environment map[string]interface{}) interface{} {
	return nil
}

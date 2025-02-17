package core

import (
	"fmt"
)

type Expr interface {
	fmt.Stringer
	Interpret() interface{}
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

func (expr AssignExpr) Interpret() interface{} {
	return expr.expr.Interpret()
}

func (expr BinaryExpr) Interpret() interface{} {
	left := expr.leftExpr.Interpret()
	right := expr.rightExpr.Interpret()
	switch expr.operator.Type {
	case PLUS:
		return left.(float64) + right.(float64)
	case MINUS:
		return left.(float64) - right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case EQUAL_EQ:
		return left == right
	case BANG_EQ:
		return left != right
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQ:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQ:
		return left.(float64) <= right.(float64)
	default:
		return nil
	}
}

func (expr UnaryExpr) Interpret() interface{} {
	res := expr.expr.Interpret()
	switch expr.operator.Type {
	case BANG:
		return !res.(bool)
	case MINUS:
		return -res.(float64)
	default:
		return nil
	}
}

func (expr GroupingExpr) Interpret() interface{} {
	return expr.expr.Interpret()
}

func (expr LiteralExpr) Interpret() interface{} {
	switch expr.value.Type {
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return expr.value.Value
	}
}

func (expr InvalidExpr) Interpret() interface{} {
	return nil
}

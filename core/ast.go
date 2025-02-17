package core

import (
	"fmt"
)

type Expr interface {
	fmt.Stringer
}

type BinaryExpr struct {
	rightExpr Expr
	operator  Token
	leftExpr  Expr
}

type UnaryExpr struct {
	operator Token
	expr     Expr
}

type GroupingExpr struct {
	expr Expr
}

type LiteralExpr struct {
	value Token
}

type InvalidExpr struct {
}

func (expr BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", expr.leftExpr.String(), expr.operator.String(), expr.rightExpr.String())
}
func (expr UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", expr.operator.String(), expr.expr.String())
}
func (expr GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", expr.expr.String())
}
func (expr LiteralExpr) String() string {
	return expr.value.String()
}
func (expr InvalidExpr) String() string {
	return "Invalid Expression"
}

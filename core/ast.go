package core

import (
	"fmt"
)

type Expr interface {
	fmt.Stringer
	Pos() int
}

type BinaryExpr struct {
	pos       int
	rightExpr Expr
	operator  Token
	leftExpr  Expr
}

type UnaryExpr struct {
	pos      int
	operator Token
	expr     Expr
}

type GroupingExpr struct {
	pos  int
	expr Expr
}

type LiteralExpr struct {
	pos   int
	value Token
}

func (expr BinaryExpr) Pos() int   { return expr.pos }
func (expr UnaryExpr) Pos() int    { return expr.pos }
func (expr GroupingExpr) Pos() int { return expr.pos }
func (expr LiteralExpr) Pos() int  { return expr.pos }

func (expr BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", expr.leftExpr.String(), expr.operator.String(), expr.rightExpr.String())
}
func (expr UnaryExpr) String() string {
	return fmt.Sprintf("%s%s", expr.operator.String(), expr.expr.String())
}
func (expr GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", expr.expr.String())
}
func (expr LiteralExpr) String() string {
	return expr.value.String()
}

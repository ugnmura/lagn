package core

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func CreateParser(tokens []Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (parser *Parser) Parse() Expr {
	return parser.expression()
}

func (parser *Parser) expression() Expr {
	return parser.assignment()
}

func (parser *Parser) assignment() Expr {
	if parser.match(IDENTIFIER) {
		name := parser.tokens[parser.current-1]
		for parser.match(EQUAL) {
			expr := parser.assignment()
			return AssignExpr{
				name: name,
				expr: expr,
			}
		}
	}
	return parser.equality()
}

func (parser *Parser) equality() Expr {
	expr := parser.comparison()

	for parser.match(EQUAL_EQ, BANG_EQ) {
		operator := parser.tokens[parser.current-1]
		rightExpr := parser.comparison()
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr
}

func (parser *Parser) comparison() Expr {
	expr := parser.term()

	for parser.match(GREATER, GREATER_EQ, LESS, LESS_EQ) {
		operator := parser.tokens[parser.current-1]
		rightExpr := parser.comparison()
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr
}

func (parser *Parser) term() Expr {
	expr := parser.factor()

	for parser.match(PLUS, MINUS) {
		operator := parser.tokens[parser.current-1]
		rightExpr := parser.factor()
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr
}

func (parser *Parser) factor() Expr {
	expr := parser.unary()

	for parser.match(STAR, SLASH) {
		operator := parser.tokens[parser.current-1]
		rightExpr := parser.unary()
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr
}

func (parser *Parser) unary() Expr {
	if parser.match(BANG, MINUS) {
		expr := parser.unary()
		return UnaryExpr{
			operator: parser.tokens[parser.current-1],
			expr:     expr,
		}
	}

	return parser.primary()
}

func (parser *Parser) primary() Expr {
	if parser.match(IDENTIFIER, NUMBER, STRING, TRUE, FALSE) {
		return LiteralExpr{
			value: parser.tokens[parser.current-1],
		}
	}

	if parser.match(LEFT_PAREN) {
		expr := parser.expression()
		parser.consume(RIGHT_PAREN, "Expected ')' after expression")
		return GroupingExpr{
			expr: expr,
		}
	}

	fmt.Printf("[ERROR] Syntax Error at Line %d\n", parser.tokens[parser.current].Line)
	return InvalidExpr{}
}

func (parser *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if parser.tokens[parser.current].Type == tokenType {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser *Parser) advance() Token {
	result := parser.tokens[parser.current]

	if !parser.isAtEnd() {
		parser.current++
	}

	return result
}

func (parser *Parser) consume(tokenType TokenType, message string) Token {
	if parser.check(tokenType) {
		return parser.advance()
	}

	panic(fmt.Sprintf("[ERROR] %s at Line %d", message, parser.tokens[parser.current].Line))
}

func (parser Parser) check(tokenType TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.tokens[parser.current].Type == tokenType
}

func (parser Parser) isAtEnd() bool {
	return parser.current == len(parser.tokens)-1
}

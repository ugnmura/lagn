package core

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE

	COMMA
	DOT

	MINUS
	MINUS_MINUS
	MINUS_EQ
	PLUS
	PLUS_PLUS
	PLUS_EQ
	SEMI
	SLASH
	SLASH_EQ
	STAR
	STAR_EQ

	BIN_AND
	BIN_AND_EQ
	BIN_OR
	BIN_OR_EQ
	BIN_XOR
	BIN_XOR_EQ

	BANG
	BANG_EQ
	EQUAL
	EQUAL_EQ
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ

	IDENTIFIER
	STRING
	NUMBER

	AND
	OR
	FOR
	WHILE
	IF
	ELSE
	RETURN

	EOF
)

type Token struct {
	Type  TokenType
	Value []rune
	Line  int
}

func (token Token) String() string {
	return string(token.Value)
}

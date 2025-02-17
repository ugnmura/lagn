package core

type Parser struct {
	tokens  []Token
	current int
}

func (parser Parser) Parse() {
	parser.Equality()
}

func (parser Parser) Equality() {
	for parser.match(EQUAL_EQ, BANG_EQ) {
	}
}

func (parser Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if parser.tokens[parser.current].Type == tokenType {
			return true
		}
	}

	return false
}

func (parser Parser) advance() Token {
	if parser.isAtEnd() {
	}

	result := parser.tokens[parser.current]
	parser.current++
	return result
}

func (parser Parser) isAtEnd() bool {
	return parser.current == len(parser.tokens)-1
}

package core

import (
	"fmt"
	"unicode"
)

type Scanner struct {
	Source  []rune
	Tokens  []Token
	Start   int
	Current int
	Line    int
	Column  int
}

func CreateScanner(source string) Scanner {
	return Scanner{
		Source:  []rune(source),
		Start:   0,
		Current: 0,
		Line:    0,
	}
}

func (scanner *Scanner) AddToken(token TokenType) {
	scanner.Tokens = append(scanner.Tokens, Token{
		Type: token,
		Line: scanner.Line,
	})
}

func (scanner *Scanner) AddTokenWithValue(token TokenType, value []rune) {
	scanner.Tokens = append(scanner.Tokens, Token{
		Type:  token,
		Line:  scanner.Line,
		Value: value,
	})
}

func (scanner *Scanner) Advance() rune {
	result := scanner.Source[scanner.Current]
	scanner.Current++
	scanner.Column++
	return result
}

func (scanner *Scanner) ScanToken() {
	c := scanner.Advance()
	switch c {
	case rune('('):
		scanner.AddToken(LEFT_PAREN)
	case rune(')'):
		scanner.AddToken(RIGHT_PAREN)
	case rune('{'):
		scanner.AddToken(LEFT_PAREN)
	case rune('}'):
		scanner.AddToken(RIGHT_BRACE)
	case rune('+'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(PLUS_EQ)
		} else if scanner.PeekCurrent() == rune('+') {
			scanner.Advance()
			scanner.AddToken(PLUS_PLUS)
		} else {
			scanner.AddToken(PLUS)
		}
	case rune('-'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(MINUS_EQ)
		} else if scanner.PeekCurrent() == rune('-') {
			scanner.Advance()
			scanner.AddToken(MINUS_MINUS)
		} else {
			scanner.AddToken(MINUS)
		}
	case rune('*'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(STAR_EQ)
		} else {
			scanner.AddToken(STAR)
		}
	case rune('/'):
		if scanner.PeekCurrent() == rune('/') {
			scanner.Advance()
			for !scanner.CurrentAtEnd() && scanner.PeekCurrent() != '\n' {
				scanner.Advance()
			}
		} else if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(SLASH_EQ)
		} else {
			scanner.AddToken(SLASH)
		}
	case rune('.'):
		scanner.AddToken(DOT)
	case rune(','):
		scanner.AddToken(COMMA)
	case rune(';'):
		scanner.AddToken(SEMI)
	case rune('"'):
		scanner.ScanString()
	case rune(' '):
	case rune('\t'):
	case rune('\r'):
	case rune('\n'):
		scanner.Line++
		scanner.Column = 0
	default:
		if unicode.IsDigit(c) {
			scanner.ScanNumber()
		} else if unicode.IsLetter(c) {
			scanner.ScanIdentifier()
		} else {
			fmt.Printf("Unexpected Character at [%d|%d]\n", scanner.Line, scanner.Column)
		}
	}
}

func (scanner *Scanner) ScanNumber() {
	for unicode.IsDigit(scanner.PeekCurrent()) {
		scanner.Advance()
	}
	if scanner.PeekCurrent() == rune('.') && unicode.IsDigit(scanner.Peek(scanner.Current+1)) {
		scanner.Advance()
		for unicode.IsDigit(scanner.PeekCurrent()) {
			scanner.Advance()
		}
	}

	scanner.AddTokenWithValue(NUMBER, scanner.Source[scanner.Start:scanner.Current])
}

func (scanner *Scanner) ScanString() {
	for !scanner.CurrentAtEnd() && scanner.PeekCurrent() != rune('"') {
		scanner.Advance()
	}

	if scanner.CurrentAtEnd() {
		fmt.Printf("Unterminated String at [%d|%d]\n", scanner.Line, scanner.Column)
		return
	}

	scanner.Advance()
	scanner.AddTokenWithValue(STRING, scanner.Source[scanner.Start+1:scanner.Current-1])
}

func (scanner *Scanner) ScanIdentifier() {
	for !scanner.CurrentAtEnd() && (unicode.IsLetter(scanner.PeekCurrent()) || unicode.IsDigit(scanner.PeekCurrent())) {
		scanner.Advance()
	}

	scanner.AddTokenWithValue(IDENTIFIER, scanner.Source[scanner.Start:scanner.Current])
}

func (scanner *Scanner) PeekCurrent() rune {
	return scanner.Peek(scanner.Current)
}

func (scanner *Scanner) Peek(toPeek int) rune {
	if scanner.IsAtEnd(toPeek) {
		return rune(0)
	}
	return scanner.Source[toPeek]
}

func (scanner Scanner) CurrentAtEnd() bool {
	return scanner.IsAtEnd(scanner.Current)
}

func (scanner Scanner) IsAtEnd(toCheck int) bool {
	return toCheck >= len(scanner.Source)
}

func (scanner *Scanner) ScanTokens() {
	for !scanner.CurrentAtEnd() {
		scanner.Start = scanner.Current
		scanner.ScanToken()
	}

	scanner.Tokens = append(scanner.Tokens, Token{
		Type: EOF,
		Line: scanner.Line,
	})
}

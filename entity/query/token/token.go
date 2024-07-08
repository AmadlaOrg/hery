package token

import (
	"fmt"
	"unicode"
)

type Type int

const (
	Identifier Type = iota
	String
	OpenParen
	CloseParen
	Dot
)

type Token struct {
	Type  Type
	Value string
}

// Lex
func Lex(query string) ([]Token, error) {
	var tokens []Token
	runes := []rune(query)
	i := 0
	for i < len(runes) {
		switch {
		case unicode.IsLetter(runes[i]) || runes[i] == '_':
			start := i
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_' || runes[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{Type: Identifier, Value: string(runes[start:i])})
		case runes[i] == '"':
			i++
			start := i
			for i < len(runes) && runes[i] != '"' {
				i++
			}
			if i >= len(runes) {
				return nil, fmt.Errorf("unclosed string")
			}
			tokens = append(tokens, Token{Type: String, Value: string(runes[start:i])})
			i++
		case runes[i] == '(':
			tokens = append(tokens, Token{Type: OpenParen, Value: "("})
			i++
		case runes[i] == ')':
			tokens = append(tokens, Token{Type: CloseParen, Value: ")"})
			i++
		case runes[i] == '.':
			tokens = append(tokens, Token{Type: Dot, Value: "."})
			i++
		case unicode.IsSpace(runes[i]):
			i++
		default:
			return nil, fmt.Errorf("unexpected character: %c", runes[i])
		}
	}
	return tokens, nil
}

// Parse
func Parse(tokens []Token) (map[string]string, []string, error) {
	queries := make(map[string]string)
	order := []string{}
	i := 0
	for i < len(tokens) {
		if tokens[i].Type == Identifier {
			key := tokens[i].Value
			if i+3 < len(tokens) && tokens[i+1].Type == OpenParen && tokens[i+2].Type == String && tokens[i+3].Type == CloseParen {
				queries[key] = tokens[i+2].Value
				order = append(order, key)
				i += 4
			} else {
				i++
			}
		} else {
			i++
		}
	}
	return queries, order, nil
}

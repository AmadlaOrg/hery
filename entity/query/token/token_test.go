package token

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		query string
		want  []Token
	}{
		{
			query: `entities._entity("github.com/AmadlaOrg/Entity@v1*").name("WebServer")`,
			want: []Token{
				{Type: Identifier, Value: "entities._entity"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "github.com/AmadlaOrg/Entity@v1*"},
				{Type: CloseParen, Value: ")"},
				{Type: Dot, Value: "."},
				{Type: Identifier, Value: "name"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "WebServer"},
				{Type: CloseParen, Value: ")"},
			},
		},
		{
			query: `entities._entity("example")`,
			want: []Token{
				{Type: Identifier, Value: "entities._entity"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "example"},
				{Type: CloseParen, Value: ")"},
			},
		},
		{
			query: `entities`,
			want: []Token{
				{Type: Identifier, Value: "entities"},
			},
		},
	}

	for _, tt := range tests {
		got, err := Lex(tt.query)
		if err != nil {
			t.Fatalf("Lex() error = %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Lex() = %v, want %v", got, tt.want)
		}
	}
}

func TestLexErrors(t *testing.T) {
	tests := []struct {
		query string
	}{
		{query: `entities._entity("unclosed string`},
		{query: `entities._entity("another unclosed string`},
	}

	for _, tt := range tests {
		_, err := Lex(tt.query)
		if err == nil {
			t.Errorf("Lex() expected error for query = %v", tt.query)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		tokens []Token
		want   map[string]string
		order  []string
	}{
		{
			tokens: []Token{
				{Type: Identifier, Value: "entities._entity"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "github.com/AmadlaOrg/Entity@v1*"},
				{Type: CloseParen, Value: ")"},
				{Type: Dot, Value: "."},
				{Type: Identifier, Value: "name"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "WebServer"},
				{Type: CloseParen, Value: ")"},
				{Type: Dot, Value: "."},
				{Type: Identifier, Value: "category"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "Application"},
				{Type: CloseParen, Value: ")"},
			},
			want: map[string]string{
				"entities._entity": "github.com/AmadlaOrg/Entity@v1*",
				"name":             "WebServer",
				"category":         "Application",
			},
			order: []string{
				"entities._entity",
				"name",
				"category",
			},
		},
		{
			tokens: []Token{
				{Type: Identifier, Value: "entities._entity"},
				{Type: OpenParen, Value: "("},
				{Type: String, Value: "example"},
				{Type: CloseParen, Value: ")"},
			},
			want: map[string]string{
				"entities._entity": "example",
			},
			order: []string{
				"entities._entity",
			},
		},
		{
			tokens: []Token{
				{Type: Identifier, Value: "entities"},
			},
			want:  map[string]string{},
			order: []string{},
		},
	}

	for _, tt := range tests {
		got, order, err := Parse(tt.tokens)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Parse() got = %v, want %v", got, tt.want)
		}
		if !reflect.DeepEqual(order, tt.order) {
			t.Errorf("Parse() order = %v, want %v", order, tt.order)
		}
	}
}

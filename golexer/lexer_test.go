package golexer

import (
	"testing"
)

func TestLexerFullCoverage(t *testing.T) {
	input := `
// Keywords
let const fn if else while for return break continue true false null

// Identifiers
valid_identifier _underscore CamelCase variable123

// Numbers
42 3.14 1e10 2.5e-3 1E+5

// Strings
"hello" "world with spaces" "escaped\"quote" "newline\ntest"

// Chars
'a' '\n' '\t' '\\' '\''

// Comments
let x = 5; // line comment
/* block comment */ let y = 10;

// Operators
= + - * / == != < <= > >= += -= *= /= && || !

// Errors
123abc
invalid#name
'
'\x'
"unterminated
"escape\q"
&
|
`

	lexer := NewLexer(input)

	for {
		tok := lexer.NextToken()
		if tok.Type == EOF {
			break
		}
		if tok.Type == ILLEGAL {
			t.Logf("Found illegal token: %q", tok.Literal)
		} else {
			t.Logf("%s -> %q", tok.Type, tok.Literal)
		}
	}

	errors := lexer.GetErrors()
	if len(errors) > 0 {
		for i, e := range errors {
			t.Logf("Lexer error[%d]: %s", i, e.Message)
		}
	}
}

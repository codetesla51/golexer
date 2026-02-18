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

// Test hex, binary, and octal numbers
func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
		literal  string
	}{
		{"0x10", NUMBER, "0x10"},
		{"0xFF", NUMBER, "0xFF"},
		{"0b1010", NUMBER, "0b1010"},
		{"0o755", NUMBER, "0o755"},
		{"0777", NUMBER, "0777"},
		{"0", NUMBER, "0"},
		{"999", NUMBER, "999"},
		{"3.14", NUMBER, "3.14"},
		{"1e10", NUMBER, "1e10"},
		{"1E10", NUMBER, "1E10"},
		{"2.5e-3", NUMBER, "2.5e-3"},
		{"1.0e+5", NUMBER, "1.0e+5"},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != tt.expected {
			t.Errorf("Input %q: expected type %s, got %s", tt.input, tt.expected, tok.Type)
		}
		if tok.Literal != tt.literal {
			t.Errorf("Input %q: expected literal %q, got %q", tt.input, tt.literal, tok.Literal)
		}
	}
}

// Test invalid numbers
func TestInvalidNumbers(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"0x", true},     // hex without digits
		{"0b", true},     // binary without digits
		{"0o", true},     // octal without digits
		{"0xG", true},    // invalid hex digit
		{"0b2", true},    // binary with invalid digit
		{"0o8", true},    // octal with invalid digit
		{"1e", true},     // exponent without digits
		{"1e+", true},    // exponent with sign but no digits
		{"123abc", true}, // number followed by letters
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		_ = lexer.NextToken()
		hasError := lexer.HasErrors()

		if hasError != tt.hasError {
			t.Errorf("Input %q: expected hasError=%v, got %v", tt.input, tt.hasError, hasError)
			if lexer.HasErrors() {
				for _, err := range lexer.GetErrors() {
					t.Logf("  Error: %s", err.Message)
				}
			}
		}
	}
}

// Test string escapes
func TestStringEscapes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello\n"`, "hello\n"},
		{`"hello\\world"`, "hello\\world"},
		{`"quote\"inside"`, "quote\"inside"},
		{`"\t\n\r"`, "\t\n\r"},
		{`"\a\b\f\v"`, "\a\b\f\v"},
		{`"\000"`, "\000"},
		{`"\x41"`, "A"},
		{`"\xFF"`, string(rune(0xFF))},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != STRING {
			t.Errorf("Input %q: expected STRING, got %s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expected {
			t.Errorf("Input %q: expected %q, got %q", tt.input, tt.expected, tok.Literal)
		}
	}
}

// Test character literals
func TestCharLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`'a'`, "a"},
		{`'\n'`, "\n"},
		{`'\t'`, "\t"},
		{`'\\'`, "\\"},
		{`'\''`, "'"},
		{`'"'`, "\""},
		{`'\x41'`, "A"},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != CHAR {
			t.Errorf("Input %q: expected CHAR, got %s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expected {
			t.Errorf("Input %q: expected %q, got %q", tt.input, tt.expected, tok.Literal)
		}
	}
}

// Test invalid character literals
func TestInvalidCharLiterals(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{`'`, true},    // unterminated
		{`'\x'`, true}, // invalid hex escape
		{`'ab'`, true}, // more than one char
		{`'`, true},    // single quote at EOF
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		_ = lexer.NextToken()
		hasError := lexer.HasErrors()

		if hasError != tt.hasError {
			t.Errorf("Input %q: expected hasError=%v, got %v", tt.input, tt.hasError, hasError)
		}
	}
}

// Test operators
func TestOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"+ - * / %", []TokenType{PLUS, MINUS, MULTIPLY, DIVIDE, MODULUS}},
		{"= == !=", []TokenType{ASSIGN, EQL, NOT_EQL}},
		{"< <= > >=", []TokenType{LESS_THAN, LESS_THAN_EQL, GREATER_THAN, GREATER_THAN_EQL}},
		{"+= -= *= /= %=", []TokenType{PLUS_ASSIGN, MINUS_ASSIGN, MULTIPLY_ASSIGN, DIVIDE_ASSIGN, MODULUS_ASSIGN}},
		{"&& ||", []TokenType{AND, OR}},
		{"!", []TokenType{BANG}},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		for i, expectedType := range tt.expected {
			tok := lexer.NextToken()
			if tok.Type != expectedType {
				t.Errorf("Input %q[%d]: expected %s, got %s", tt.input, i, expectedType, tok.Type)
			}
		}
	}
}

// Test invalid operators
func TestInvalidOperators(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"&", true},   // single & is invalid
		{"|", true},   // single | is invalid
		{"&&", false}, // && is valid
		{"||", false}, // || is valid
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		_ = lexer.NextToken()
		hasError := lexer.HasErrors()

		if hasError != tt.hasError {
			t.Errorf("Input %q: expected hasError=%v, got %v", tt.input, tt.hasError, hasError)
		}
	}
}

// Test comments
func TestComments(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []TokenType
	}{
		{"// comment\n42", []TokenType{NUMBER, EOF}},
		{"x // comment", []TokenType{IDENT, EOF}},
		{"/* comment */ 42", []TokenType{NUMBER, EOF}},
		{"/* nested /* comment */ 42", []TokenType{NUMBER, EOF}},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		for i, expectedType := range tt.expectedTokens {
			tok := lexer.NextToken()
			if tok.Type != expectedType {
				t.Errorf("Input %q[%d]: expected %s, got %s (literal: %q)", tt.input, i, expectedType, tok.Type, tok.Literal)
			}
		}
	}
}

// Test unterminated comment
func TestUnterminatedBlockComment(t *testing.T) {
	lexer := NewLexer("/* unterminated")
	_ = lexer.NextToken()

	if !lexer.HasErrors() {
		t.Errorf("Expected error for unterminated block comment")
	}
}

// Test whitespace handling
func TestWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"  42", NUMBER},
		{"\t42", NUMBER},
		{"\n42", NUMBER},
		{"\r\n42", NUMBER},
		{"   \t\n\r  42", NUMBER},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != tt.expected {
			t.Errorf("Input with whitespace: expected %s, got %s", tt.expected, tok.Type)
		}
	}
}

// Test keywords
func TestKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"let", LET},
		{"const", CONST},
		{"fn", FN},
		{"if", IF},
		{"else", ELSE},
		{"while", WHILE},
		{"for", FOR},
		{"return", RETURN},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"true", TRUE},
		{"false", FALSE},
		{"null", NULL},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != tt.expected {
			t.Errorf("Input %q: expected %s, got %s", tt.input, tt.expected, tok.Type)
		}
	}
}

// Test identifiers vs keywords
func TestIdentifiersVsKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"x", IDENT},
		{"variable", IDENT},
		{"_private", IDENT},
		{"CamelCase", IDENT},
		{"snake_case", IDENT},
		{"letme", IDENT}, // starts with "let" but is different
		{"letter", IDENT},
		{"let", LET}, // exact keyword
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != tt.expected {
			t.Errorf("Input %q: expected %s, got %s", tt.input, tt.expected, tok.Type)
		}
	}
}

// Test delimiters
func TestDelimiters(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"( )", []TokenType{LPAREN, RPAREN}},
		{"{ }", []TokenType{LBRACE, RBRACE}},
		{"[ ]", []TokenType{LBRACKET, RBRACKET}},
		{", ; :", []TokenType{COMMA, SEMICOLON, COLON}},
		{".", []TokenType{DOT}},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		for i, expectedType := range tt.expected {
			tok := lexer.NextToken()
			if tok.Type != expectedType {
				t.Errorf("Input %q[%d]: expected %s, got %s", tt.input, i, expectedType, tok.Type)
			}
		}
	}
}

// Test backtick strings
func TestBacktickStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"`hello`", "hello"},
		{"`raw\nstring`", "raw\nstring"},
		{"`with 'quotes'`", "with 'quotes'"},
		{"`with \"quotes\"`", "with \"quotes\""},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tok := lexer.NextToken()

		if tok.Type != BACKTICK_STRING {
			t.Errorf("Input %q: expected BACKTICK_STRING, got %s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expected {
			t.Errorf("Input %q: expected %q, got %q", tt.input, tt.expected, tok.Literal)
		}
	}
}

// Test position tracking
func TestPositionTracking(t *testing.T) {
	input := "x y\nz"
	lexer := NewLexer(input)

	tokens := []struct {
		expectedType TokenType
		line         int
		column       int
	}{
		{IDENT, 1, 1},
		{IDENT, 1, 3},
		{IDENT, 2, 1},
		{EOF, 2, 2},
	}

	for i, tt := range tokens {
		tok := lexer.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("Token %d: expected type %s, got %s", i, tt.expectedType, tok.Type)
		}
		if tok.Line != tt.line {
			t.Errorf("Token %d: expected line %d, got %d", i, tt.line, tok.Line)
		}
		if tok.Column != tt.column {
			t.Errorf("Token %d: expected column %d, got %d", i, tt.column, tok.Column)
		}
	}
}

// Test complex expression
func TestComplexExpression(t *testing.T) {
	input := "let x = 10 + 20 * 30;"
	lexer := NewLexer(input)

	expected := []TokenType{LET, IDENT, ASSIGN, NUMBER, PLUS, NUMBER, MULTIPLY, NUMBER, SEMICOLON, EOF}

	for i, expectedType := range expected {
		tok := lexer.NextToken()
		if tok.Type != expectedType {
			t.Errorf("Token %d: expected %s, got %s (literal: %q)", i, expectedType, tok.Type, tok.Literal)
		}
	}
}

// Test config loading with additional keywords
func TestConfigAdditionalKeywords(t *testing.T) {
	// First test without config - "async" should be IDENT
	lexer := NewLexer("async await")
	tok := lexer.NextToken()
	if tok.Type != IDENT {
		t.Errorf("Without config: 'async' should be IDENT, got %s", tok.Type)
	}

	// Now load config and test
	lexer2 := NewLexerWithConfig("async await unless until", "../examples/config.json")

	tests := []struct {
		expected TokenType
		literal  string
	}{
		{"ASYNC", "async"},
		{"AWAIT", "await"},
		{"UNLESS", "unless"},
		{"UNTIL", "until"},
	}

	for _, tt := range tests {
		tok := lexer2.NextToken()
		if tok.Type != tt.expected {
			t.Errorf("With config: '%s' expected %s, got %s", tt.literal, tt.expected, tok.Type)
		}
		if tok.Literal != tt.literal {
			t.Errorf("With config: expected literal %q, got %q", tt.literal, tok.Literal)
		}
	}
}

// Test config loading with additional punctuation
// Note: Config merges to global state, so once loaded it persists
func TestConfigAdditionalPunctuation(t *testing.T) {
	// Load config and test punctuation tokens
	lexer := NewLexerWithConfig("@ # $", "../examples/config.json")

	tests := []struct {
		expected TokenType
		literal  string
	}{
		{"AT_SYMBOL", "@"},
		{"HASH", "#"},
		{"DOLLAR", "$"},
	}

	for _, tt := range tests {
		tok := lexer.NextToken()
		if tok.Type != tt.expected {
			t.Errorf("With config: '%s' expected %s, got %s", tt.literal, tt.expected, tok.Type)
		}
	}
}

// Test config file not found handling
func TestConfigNotFound(t *testing.T) {
	// Should not panic, just warn and use defaults
	lexer := NewLexerWithConfig("let x = 5", "nonexistent.json")

	tok := lexer.NextToken()
	if tok.Type != LET {
		t.Errorf("Should still work with missing config: expected LET, got %s", tok.Type)
	}
}

// Test TokenizeAll helper method
func TestTokenizeAll(t *testing.T) {
	lexer := NewLexer("let x = 5;")
	tokens, errors := lexer.TokenizeAll()

	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}

	expectedTypes := []TokenType{LET, IDENT, ASSIGN, NUMBER, SEMICOLON}
	if len(tokens) != len(expectedTypes) {
		t.Errorf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, tt := range expectedTypes {
		if tokens[i].Type != tt {
			t.Errorf("Token %d: expected %s, got %s", i, tt, tokens[i].Type)
		}
	}
}

// Test error collection
func TestErrorCollection(t *testing.T) {
	lexer := NewLexer("123abc @ 456def")
	lexer.TokenizeAll()

	if !lexer.HasErrors() {
		t.Errorf("Expected errors for invalid tokens")
	}

	errors := lexer.GetErrors()
	if len(errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(errors))
	}

	// Check that errors have position info
	for _, err := range errors {
		if err.Line < 1 || err.Column < 1 {
			t.Errorf("Error should have valid position: line=%d, col=%d", err.Line, err.Column)
		}
	}
}

package golexer

import (
	"testing"
	"strings"
)

func TestTokenizeBasicOperators(t *testing.T) {
	input := `= + - * / == != < <= > >= += -= *= /= && || !`
	
	expected := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{ASSIGN, "="},
		{PLUS, "+"},
		{MINUS, "-"},
		{MULTIPLY, "*"},
		{DIVIDE, "/"},
		{EQL, "=="},
		{NOT_EQL, "!="},
		{LESS_THAN, "<"},
		{LESS_THAN_EQL, "<="},
		{GREATER_THAN, ">"},
		{GREATER_THAN_EQL, ">="},
		{PLUS_ASSIGN, "+="},
		{MINUS_ASSIGN, "-="},
		{MULTIPLY_ASSIGN, "*="},
		{DIVIDE_ASSIGN, "/="},
		{AND, "&&"},
		{OR, "||"},
		{BANG, "!"},
		{EOF, ""},
	}

	lexer := NewLexer(input)
	
	for i, tt := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenizeKeywords(t *testing.T) {
	input := `let const fn if else while for return break continue true false null`
	
	expected := []TokenType{
		LET, CONST, FN, IF, ELSE, WHILE, FOR, RETURN, BREAK, CONTINUE, TRUE, FALSE, NULL, EOF,
	}

	lexer := NewLexer(input)
	
	for i, expectedType := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expectedType, tok.Type)
		}
	}
}

func TestTokenizeIdentifiers(t *testing.T) {
	input := `valid_identifier _underscore CamelCase variable123`
	
	expected := []string{"valid_identifier", "_underscore", "CamelCase", "variable123"}

	lexer := NewLexer(input)
	
	for i, expectedLiteral := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != IDENT {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, IDENT, tok.Type)
		}
		
		if tok.Literal != expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenizeNumbers(t *testing.T) {
	input := `42 3.14 1e10 2.5e-3 1E+5`
	
	expected := []string{"42", "3.14", "1e10", "2.5e-3", "1E+5"}

	lexer := NewLexer(input)
	
	for i, expectedLiteral := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != NUMBER {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, NUMBER, tok.Type)
		}
		
		if tok.Literal != expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenizeStrings(t *testing.T) {
	input := `"hello" "world with spaces" "escaped\"quote" "newline\ntest"`
	
	expected := []string{"hello", "world with spaces", "escaped\"quote", "newline\ntest"}

	lexer := NewLexer(input)
	
	for i, expectedLiteral := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != STRING {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, STRING, tok.Type)
		}
		
		if tok.Literal != expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenizeCharacters(t *testing.T) {
	input := `'a' '\n' '\t' '\\' '\''`
	
	expected := []string{"a", "\n", "\t", "\\", "'"}

	lexer := NewLexer(input)
	
	for i, expectedLiteral := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != CHAR {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, CHAR, tok.Type)
		}
		
		if tok.Literal != expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expectedLiteral, tok.Literal)
		}
	}
}

func TestComments(t *testing.T) {
	input := `let x = 5; // line comment
	/* block comment */ let y = 10;`
	
	expected := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{NUMBER, "5"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{NUMBER, "10"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	lexer := NewLexer(input)
	
	for i, tt := range expected {
		tok := lexer.NextToken()
		
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"123abc", "identifier must start with letter or underscore"},
		{"invalid#name", "unexpected character: #"},
		{"'", "unterminated char literal"},
		{"'\\x'", "unknown escape sequence '\\x' in char literal"},
		{`"unterminated`, "unterminated string literal"},
		{`"escape\q"`, "unknown escape sequence '\\q' in string literal"},
		{"&", "unexpected character: &"},
		{"|", "unexpected character: |"},
	}

	for i, tt := range tests {
		lexer := NewLexer(tt.input)
		
		// Consume all tokens to trigger errors
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
		
		errors := lexer.GetErrors()
		if len(errors) == 0 {
			t.Fatalf("test[%d] - expected error but got none for input: %q", i, tt.input)
		}
		
		found := false
		for _, err := range errors {
			if strings.Contains(err.Message, tt.expectedError) {
				found = true
				break
			}
		}
		
		if !found {
			t.Fatalf("test[%d] - expected error containing %q, got: %v", i, tt.expectedError, errors[0].Message)
		}
	}
}

func TestUnicodeIdentifiers(t *testing.T) {
	input := `变量 переменная 変数 _مؤشر`
	
	lexer := NewLexer(input)
	
	for i := 0; i < 4; i++ {
		tok := lexer.NextToken()
		
		if tok.Type != IDENT {
			t.Fatalf("test[%d] - expected IDENT, got %q", i, tok.Type)
		}
		
		if tok.Literal == "" {
			t.Fatalf("test[%d] - empty literal for Unicode identifier", i)
		}
	}
	
	// Should not have any errors
	if lexer.HasErrors() {
		t.Fatalf("unexpected errors for Unicode identifiers: %v", lexer.GetErrors())
	}
}

func TestTokenizeAll(t *testing.T) {
	input := `let x = 42; // simple assignment`
	
	tokens, errors := NewLexer(input).TokenizeAll()
	
	if len(errors) != 0 {
		t.Fatalf("unexpected errors: %v", errors)
	}
	
	expectedTypes := []TokenType{LET, IDENT, ASSIGN, NUMBER, SEMICOLON}
	
	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}
	
	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Fatalf("token[%d] - expected %q, got %q", i, expectedType, tokens[i].Type)
		}
	}
}

func TestPositionTracking(t *testing.T) {
	input := `let x = 5;
fn main() {
    return x;
}`
	
	lexer := NewLexer(input)
	
	// Test first token position
	tok := lexer.NextToken() // "let"
	if tok.Line != 1 || tok.Column != 1 {
		t.Fatalf("expected position (1,1), got (%d,%d)", tok.Line, tok.Column)
	}
	
	// Skip to second line
	for tok.Type != FN {
		tok = lexer.NextToken()
	}
	
	if tok.Line != 2 {
		t.Fatalf("expected line 2, got %d", tok.Line)
	}
}

func BenchmarkLexer(b *testing.B) {
	input := `
	let x = 42;
	const PI = 3.14159;
	fn fibonacci(n int) int {
		if n <= 1 {
			return n;
		}
		return fibonacci(n-1) + fibonacci(n-2);
	}
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}
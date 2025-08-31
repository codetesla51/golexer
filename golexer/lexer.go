/*
GoLexer - A Comprehensive Lexical Analyzer for Go
Author: Uthman Dev
GitHub: https://github.com/codetesla51/golexer
License: MIT

Core Lexical Analyzer Implementation
This is the heart of the GoLexer library, providing comprehensive
tokenization capabilities for modern programming languages.

Features:
- Multi-format number support (decimal, hex, binary, octal, scientific)
- String literals with escape sequences and raw strings
- Character literals with full escape sequence support
- Complete operator set including compound assignments
- Unicode identifier support
- Comprehensive error recovery and reporting
- Line/column position tracking
- Comment handling (line and block comments)

The lexer processes UTF-8 encoded source code and maintains precise
position information for all tokens and errors, making it suitable
for building IDEs, compilers, and language tools.
*/

package golexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Operator defines a single or compound operator
type Operator struct {
	Single       string
	SingleType   TokenType
	Compound     string
	CompoundType TokenType
}

// operators defines all operators with their single and compound forms
var operators = []Operator{
	{"=", ASSIGN, "==", EQL},
	{"+", PLUS, "+=", PLUS_ASSIGN},
	{"-", MINUS, "-=", MINUS_ASSIGN},
	{"*", MULTIPLY, "*=", MULTIPLY_ASSIGN},
	{"/", DIVIDE, "/=", DIVIDE_ASSIGN},
	{"%", MODULUS, "%=", MODULUS_ASSIGN},
	{"!", BANG, "!=", NOT_EQL},
	{"<", LESS_THAN, "<=", LESS_THAN_EQL},
	{">", GREATER_THAN, ">=", GREATER_THAN_EQL},
	{"&", "", "&&", AND}, // Single & is invalid
	{"|", "", "||", OR},  // Single | is invalid
}

// singleCharTokens maps single characters to their token types
var singleCharTokens = map[rune]TokenType{
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
	',': COMMA,
	';': SEMICOLON,
	':': COLON,
	'.': DOT,
}

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
	errors       []*LexError
}

// NewLexer creates a new lexer instance with the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
		errors: make([]*LexError, 0),
	}
	l.readChar()
	return l
}

// load config
func NewLexerWithConfig(input, configFile string) *Lexer {
	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config file '%s': %v\n", configFile, err)
		fmt.Fprintf(os.Stderr, "Continuing with default configuration...\n")
	} else {
		config.MergeWithDefaults()
	}

	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
		errors: make([]*LexError, 0),
	}
	l.readChar()
	return l
}

// GetErrors returns all lexical errors encountered during tokenization
func (l *Lexer) GetErrors() []*LexError {
	return l.errors
}

// HasErrors returns true if any lexical errors were encountered
func (l *Lexer) HasErrors() bool {
	return len(l.errors) > 0
}

// TokenizeAll returns all tokens from the input along with any errors
func (l *Lexer) TokenizeAll() ([]Token, []*LexError) {
	var tokens []Token

	for {
		tok := l.NextToken()
		if tok.Type == EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens, l.errors
}

func (l *Lexer) addError(message string) {
	l.errors = append(l.errors, &LexError{
		Message: message,
		Line:    l.line,
		Column:  l.column,
	})
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
	}

	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isHexDigit(ch rune) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isBinaryDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isOctalDigit(ch rune) bool {
	return ch >= '0' && ch <= '7'
}

func (l *Lexer) readIdentifier() string {
	start := l.position

	// First character must be letter or underscore
	if !isLetter(l.ch) {
		l.addError("identifier must start with a letter or underscore")
		return ""
	}

	// Read the identifier - continue while we have letters or digits
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position

	// Check for hex, binary, or octal prefixes
	if l.ch == '0' {
		next := l.peekChar()
		if next == 'x' || next == 'X' {
			return l.readHexNumber()
		}
		if next == 'b' || next == 'B' {
			return l.readBinaryNumber()
		}
		if next == 'o' || next == 'O' {
			return l.readOctalNumber()
		}
		// Traditional octal (starts with 0)
		if isOctalDigit(next) {
			return l.readTraditionalOctal()
		}
	}

	// Regular decimal number
	for isDigit(l.ch) {
		l.readChar()
	}

	// Float with decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar() // consume 'e' or 'E'

		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}

		if !isDigit(l.ch) {
			l.addError("invalid scientific notation: exponent must contain digits")
		} else {
			for isDigit(l.ch) {
				l.readChar()
			}
		}
	}

	// Check for invalid trailing characters
	if isLetter(l.ch) && l.ch != 0 {
		l.addError("invalid number: numbers cannot be followed by letters")
		// Skip the invalid characters to avoid cascading errors
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) readHexNumber() string {
	start := l.position
	l.readChar() // skip '0'
	l.readChar() // skip 'x' or 'X'

	if !isHexDigit(l.ch) {
		l.addError("invalid hexadecimal number: must contain at least one hex digit after 0x")
		return l.input[start:l.position]
	}

	for isHexDigit(l.ch) {
		l.readChar()
	}

	// Check for invalid trailing characters
	if isLetter(l.ch) && l.ch != 0 {
		l.addError("invalid hexadecimal number: contains non-hex characters")
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) readBinaryNumber() string {
	start := l.position
	l.readChar() // skip '0'
	l.readChar() // skip 'b' or 'B'

	if !isBinaryDigit(l.ch) {
		l.addError("invalid binary number: must contain at least one binary digit after 0b")
		return l.input[start:l.position]
	}

	for isBinaryDigit(l.ch) {
		l.readChar()
	}

	// Check for invalid trailing characters
	if (isDigit(l.ch) && !isBinaryDigit(l.ch)) || isLetter(l.ch) {
		l.addError("invalid binary number: contains non-binary characters")
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) readOctalNumber() string {
	start := l.position
	l.readChar() // skip '0'
	l.readChar() // skip 'o' or 'O'

	if !isOctalDigit(l.ch) {
		l.addError("invalid octal number: must contain at least one octal digit after 0o")
		return l.input[start:l.position]
	}

	for isOctalDigit(l.ch) {
		l.readChar()
	}

	// Check for invalid trailing characters
	if (isDigit(l.ch) && !isOctalDigit(l.ch)) || isLetter(l.ch) {
		l.addError("invalid octal number: contains non-octal characters")
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) readTraditionalOctal() string {
	start := l.position

	for isOctalDigit(l.ch) {
		l.readChar()
	}

	// Check for invalid trailing characters
	if (isDigit(l.ch) && !isOctalDigit(l.ch)) || isLetter(l.ch) {
		l.addError("invalid octal number: contains non-octal characters")
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) readEscapeSequence() rune {
	l.readChar() // consume backslash
	if l.ch == 0 {
		l.addError("unterminated escape sequence")
		return 0
	}

	switch l.ch {
	case 'a':
		return '\a' // bell
	case 'b':
		return '\b' // backspace
	case 'f':
		return '\f' // form feed
	case 'n':
		return '\n' // newline
	case 'r':
		return '\r' // carriage return
	case 't':
		return '\t' // tab
	case 'v':
		return '\v' // vertical tab
	case '\\':
		return '\\'
	case '\'':
		return '\''
	case '"':
		return '"'
	case '0':
		return '\000' // null character
	case 'x':
		// Hex escape sequence \xNN
		l.readChar()
		if !isHexDigit(l.ch) {
			l.addError("invalid hex escape sequence: expected hex digit after \\x")
			return 0
		}
		first := l.ch
		l.readChar()
		if !isHexDigit(l.ch) {
			l.addError("invalid hex escape sequence: expected two hex digits after \\x")
			return 0
		}
		second := l.ch
		// Convert hex digits to rune
		var val rune
		if first >= '0' && first <= '9' {
			val = (first - '0') * 16
		} else if first >= 'a' && first <= 'f' {
			val = (first - 'a' + 10) * 16
		} else if first >= 'A' && first <= 'F' {
			val = (first - 'A' + 10) * 16
		}
		if second >= '0' && second <= '9' {
			val += second - '0'
		} else if second >= 'a' && second <= 'f' {
			val += second - 'a' + 10
		} else if second >= 'A' && second <= 'F' {
			val += second - 'A' + 10
		}
		return val
	default:
		l.addError(fmt.Sprintf("unknown escape sequence '\\%c'", l.ch))
		return l.ch
	}
}

func (l *Lexer) readCharLiteral() string {
	var result strings.Builder

	l.readChar() // consume opening '

	if l.ch == 0 {
		l.addError("unterminated character literal")
		return ""
	}

	if l.ch == '\n' {
		l.addError("character literal cannot contain newline")
		return ""
	}

	if l.ch == '\\' {
		char := l.readEscapeSequence()
		if char != 0 {
			result.WriteRune(char)
		}
	} else {
		result.WriteRune(l.ch)
	}

	l.readChar()
	if l.ch != '\'' {
		l.addError("character literal must be closed with single quote")
	} else {
		l.readChar() // consume closing '
	}

	return result.String()
}

func (l *Lexer) readString() string {
	var result strings.Builder

	for {
		l.readChar()
		if l.ch == 0 {
			l.addError("unterminated string literal")
			break
		}
		if l.ch == '"' {
			break
		}
		if l.ch == '\\' {
			char := l.readEscapeSequence()
			if char != 0 {
				result.WriteRune(char)
			}
		} else {
			result.WriteRune(l.ch)
		}
	}

	return result.String()
}

func (l *Lexer) readBacktickString() string {
	var result strings.Builder

	for {
		l.readChar()
		if l.ch == 0 {
			l.addError("unterminated backtick string literal")
			break
		}
		if l.ch == '`' {
			break
		}
		result.WriteRune(l.ch)
	}

	return result.String()
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // consume initial '*'
	for {
		if l.ch == 0 {
			l.addError("unterminated block comment")
			return
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip '*'
			l.readChar() // skip '/'
			break
		}
		l.readChar()
	}
}

// tryOperator attempts to match an operator and returns the token if found
func (l *Lexer) tryOperator(line, column int) (Token, bool) {
	for _, op := range operators {
		if l.ch == rune(op.Single[0]) {
			// Handle special cases for & and | which require compound form
			if (l.ch == '&' || l.ch == '|') && op.Compound != "" {
				if l.peekChar() == rune(op.Compound[1]) {
					ch := l.ch
					l.readChar() // consume first character
					result := Token{
						Type:    op.CompoundType,
						Literal: string(ch) + string(l.ch),
						Line:    line,
						Column:  column,
					}
					return result, true
				} else {
					// Single & or | is an error
					suggestion := op.Compound
					l.addError(fmt.Sprintf("unexpected character '%c' - did you mean '%s'?", l.ch, suggestion))
					return Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}, true
				}
			}

			// Check for compound operator
			if op.Compound != "" && len(op.Compound) > 1 && l.peekChar() == rune(op.Compound[1]) {
				ch := l.ch
				l.readChar() // consume first character
				return Token{
					Type:    op.CompoundType,
					Literal: string(ch) + string(l.ch),
					Line:    line,
					Column:  column,
				}, true
			}

			// Return single operator (if it has a valid single form)
			if op.SingleType != "" {
				return Token{
					Type:    op.SingleType,
					Literal: op.Single,
					Line:    line,
					Column:  column,
				}, true
			}
		}
	}
	return Token{}, false
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	line := l.line
	column := l.column

	// Handle comments FIRST (before operators)
	if l.ch == '/' {
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		}
		// If not a comment, fall through to operator handling
	}

	// Handle identifiers and keywords
	if isLetter(l.ch) {
		literal := l.readIdentifier()
		if literal == "" {
			return Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
		return Token{
			Type:    LookupIdent(literal),
			Literal: literal,
			Line:    line,
			Column:  column,
		}
	}

	// Handle numbers
	if isDigit(l.ch) {
		errorCountBefore := len(l.errors)
		literal := l.readNumber()

		// Check if errors were added during number parsing
		var tokType TokenType = NUMBER
		if len(l.errors) > errorCountBefore {
			tokType = ILLEGAL
		}

		return Token{
			Type:    tokType,
			Literal: literal,
			Line:    line,
			Column:  column,
		}
	}

	// Try operators
	if opTok, found := l.tryOperator(line, column); found {
		l.readChar()
		return opTok
	}

	// Handle special cases that need custom logic
	switch l.ch {
	case '\'':
		char := l.readCharLiteral()
		tok = Token{Type: CHAR, Literal: char, Line: line, Column: column}
		// readCharLiteral already consumed the closing quote
		return tok
	case '"':
		str := l.readString()
		tok = Token{
			Type:    STRING,
			Literal: str,
			Line:    line,
			Column:  column,
		}
	case '`':
		str := l.readBacktickString()
		tok = Token{
			Type:    BACKTICK_STRING,
			Literal: str,
			Line:    line,
			Column:  column,
		}
	case 0:
		tok = Token{Type: EOF, Literal: "", Line: line, Column: column}
	default:
		// Check single character tokens
		if tokenType, exists := singleCharTokens[l.ch]; exists {
			tok = Token{Type: tokenType, Literal: string(l.ch), Line: line, Column: column}
		} else {
			l.addError(fmt.Sprintf("unexpected character '%c' (Unicode: U+%04X)", l.ch, l.ch))
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	}

	l.readChar()
	return tok
}

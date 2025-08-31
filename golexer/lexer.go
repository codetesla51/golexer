package golexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

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
		l.column = 0
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

func (l *Lexer) readIdentifier() string {
	start := l.position
	
	// First character must be letter or underscore
	if !isLetter(l.ch) {
		l.addError("identifier must start with letter or underscore")
		return ""
	}
	
	// Read the identifier - continue while we have letters or digits
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.readPosition]
}

func (l *Lexer) readNumber() string {
	start := l.position
	
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		l.readChar() // consume 'e' or 'E'

		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		
		if !isDigit(l.ch) {
			l.addError("invalid number: missing digits in exponent")
		} else {
			for isDigit(l.ch) {
				l.readChar()
			}
		}
	}
	
	if isLetter(l.ch) && l.ch != 0 {
		l.addError("invalid number: number cannot be followed by letters")
		// Skip the invalid characters to avoid cascading errors
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	// Use readPosition (next char position) for the end slice position
	return l.input[start:l.readPosition]
}

func (l *Lexer) readCharLiteral() string {
	var result strings.Builder

	l.readChar() // consume opening '

	if l.ch == 0 {
		l.addError("unterminated char literal")
		return ""
	}

	if l.ch == '\\' {
		l.readChar()
		if l.ch == 0 {
			l.addError("unterminated escape sequence in char literal")
			return ""
		}

		switch l.ch {
		case 'n':
			result.WriteRune('\n')
		case 't':
			result.WriteRune('\t')
		case 'r':
			result.WriteRune('\r')
		case '\\':
			result.WriteRune('\\')
		case '\'':
			result.WriteRune('\'')
		case '"':
			result.WriteRune('"')
		default:
			l.addError(fmt.Sprintf("unknown escape sequence '\\%c' in char literal", l.ch))
			result.WriteRune(l.ch)
		}
	} else {
		result.WriteRune(l.ch)
	}

	l.readChar()
	if l.ch != '\'' {
		l.addError("unterminated char literal")
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
			l.readChar() // consume closing quote
			break
		}
		if l.ch == '\\' {
			l.readChar()
			if l.ch == 0 {
				l.addError("unterminated escape sequence in string literal")
				break
			}
			switch l.ch {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case '"':
				result.WriteRune('"')
			case '\\':
				result.WriteRune('\\')
			default:
				l.addError(fmt.Sprintf("unknown escape sequence '\\%c' in string literal", l.ch))
				result.WriteRune(l.ch)
			}
		} else {
			result.WriteRune(l.ch)
		}
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

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	line := l.line
	column := l.column

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

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    EQL,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{
				Type:    ASSIGN,
				Literal: string(l.ch),
				Line:    line,
				Column:  column,
			}
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    PLUS_ASSIGN,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{Type: PLUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    MINUS_ASSIGN,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{Type: MINUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    NOT_EQL,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{
				Type:    BANG,
				Literal: string(l.ch),
				Line:    line,
				Column:  column,
			}
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    DIVIDE_ASSIGN,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{Type: DIVIDE, Literal: string(l.ch), Line: line, Column: column}
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    MULTIPLY_ASSIGN,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{Type: MULTIPLY, Literal: string(l.ch), Line: line, Column: column}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    AND,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			l.addError(fmt.Sprintf("unexpected character: %c", l.ch))
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    OR,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			l.addError(fmt.Sprintf("unexpected character: %c", l.ch))
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    LESS_THAN_EQL,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{
				Type:    LESS_THAN,
				Literal: string(l.ch),
				Line:    line,
				Column:  column,
			}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    GREATER_THAN_EQL,
				Literal: string(ch) + string(l.ch),
				Line:    line,
				Column:  column,
			}
		} else {
			tok = Token{
				Type:    GREATER_THAN,
				Literal: string(l.ch),
				Line:    line,
				Column:  column,
			}
		}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Line: line, Column: column}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Line: line, Column: column}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '[':
		tok = Token{Type: LBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch), Line: line, Column: column}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch), Line: line, Column: column}
	case ':':
		tok = Token{Type: COLON, Literal: string(l.ch), Line: line, Column: column}
	case '\'':
		char := l.readCharLiteral()
		tok = Token{Type: CHAR, Literal: char, Line: line, Column: column}
	case '"':
		str := l.readString()
		tok = Token{
			Type:    STRING,
			Literal: str,
			Line:    line,
			Column:  column,
		}
	case 0:
		tok = Token{Type: EOF, Literal: "", Line: line, Column: column}
	default:
		l.addError(fmt.Sprintf("unexpected character: %c", l.ch))
		tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
	}

	l.readChar()
	return tok
}
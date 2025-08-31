package golexer

// TokenType represents the type of a token
type TokenType string

// Token represents a single token with its type, literal value, and position
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Token type constants
const (
	ILLEGAL  = "ILLEGAL"
	EOF      = "EOF"
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
	NUMBER   = "NUMBER"

	// Logical operators
	BANG = "!"
	AND  = "&&"
	OR   = "||"

	// Comparison operators
	NOT_EQL          = "!="
	LESS_THAN        = "<"
	LESS_THAN_EQL    = "<="
	GREATER_THAN     = ">"
	GREATER_THAN_EQL = ">="
	EQL              = "=="

	// Assignment operators
	PLUS_ASSIGN     = "+="
	MINUS_ASSIGN    = "-="
	MULTIPLY_ASSIGN = "*="
	DIVIDE_ASSIGN   = "/="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Brackets
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Identifiers and Keywords
	IDENT    = "IDENT"
	LET      = "LET"
	CONST    = "CONST"
	FN       = "FN"
	IF       = "IF"
	ELSE     = "ELSE"
	WHILE    = "WHILE"
	FOR      = "FOR"
	RETURN   = "RETURN"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	STRING   = "STRING"

	// Type tokens
	TYPE_INT    = "TYPE_INT"
	TYPE_FLOAT  = "TYPE_FLOAT"
	TYPE_STRING = "TYPE_STRING"
	TYPE_BOOL   = "TYPE_BOOL"
	TYPE_CHAR   = "TYPE_CHAR"
	CHAR        = "CHAR"
)

// keywords maps string literals to their corresponding token types
var keywords = map[string]TokenType{
	"let":      LET,
	"const":    CONST,
	"fn":       FN,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"return":   RETURN,
	"break":    BREAK,
	"continue": CONTINUE,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	// Type keywords
	"int":    TYPE_INT,
	"float":  TYPE_FLOAT,
	"string": TYPE_STRING,
	"bool":   TYPE_BOOL,
	"char":   TYPE_CHAR,
}

// LookupIdent checks if an identifier is a keyword and returns the appropriate token type
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

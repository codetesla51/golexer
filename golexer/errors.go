package golexer

import "fmt"

// LexError represents a lexical analysis error with position information
type LexError struct {
	Message string
	Line    int
	Column  int
}

// Error implements the error interface
func (e *LexError) Error() string {
	return fmt.Sprintf("lexical error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
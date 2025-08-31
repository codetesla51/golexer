/*
GoLexer - A Comprehensive Lexical Analyzer for Go
Author: Uthman Dev
GitHub: https://github.com/codetesla51/golexer
License: MIT

Lexical Error Handling
Defines error types and structures for reporting lexical analysis errors
with precise position information. The LexError type implements the
standard Go error interface while providing line and column details
for debugging and IDE integration.

Features:
- Detailed error messages with context
- Precise line and column tracking
- Standard Go error interface compliance
*/

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

# GoLexer

[![Go Reference](https://pkg.go.dev/badge/github.com/codetesla51/golexer.svg)](https://pkg.go.dev/github.com/codetesla51/golexer)
[![Go Report Card](https://goreportcard.com/badge/github.com/codetesla51/golexer)](https://goreportcard.com/report/github.com/codetesla51/golexer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive lexical analyzer (tokenizer) library for Go. Designed for building programming languages, domain-specific languages (DSLs), configuration parsers, and template engines.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Supported Tokens](#supported-tokens)
- [Error Handling](#error-handling)
- [Examples](#examples)
- [Extending the Lexer](#extending-the-lexer)
- [Testing](#testing)
- [Performance](#performance)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/codetesla51/golexer
```

Requires Go 1.21 or later.

## Quick Start

### Basic Tokenization

```go
package main

import (
    "fmt"
    "github.com/codetesla51/golexer/golexer"
)

func main() {
    source := `let x = 42 + 3.14;`
    
    lexer := golexer.NewLexer(source)
    
    for {
        token := lexer.NextToken()
        if token.Type == golexer.EOF {
            break
        }
        
        fmt.Printf("%-15s %-15s Line:%d Col:%d\n", 
            token.Type, token.Literal, token.Line, token.Column)
    }
    
    // Check for lexical errors
    if lexer.HasErrors() {
        for _, err := range lexer.GetErrors() {
            fmt.Printf("Error: %s\n", err.Error())
        }
    }
}
```

### Batch Processing

```go
lexer := golexer.NewLexer(source)
tokens, errors := lexer.TokenizeAll()

fmt.Printf("Processed %d tokens with %d errors\n", len(tokens), len(errors))

// Analyze token distribution
counts := make(map[golexer.TokenType]int)
for _, token := range tokens {
    counts[token.Type]++
}
```

### Testing the Lexer

A comprehensive test file (`test.lang`) is included that demonstrates all supported features:

```bash
# Clone the repository
git clone https://github.com/codetesla51/golexer.git
cd golexer

# Run the comprehensive test
go run cmd/main.go test.lang

# Should output: Status: ✓ PASSED with 0 lexical errors
```

The `test.lang` file contains over 1700 tokens showcasing every supported feature, making it an excellent reference for understanding the lexer's capabilities.

## API Reference

### Core Types

#### Lexer

```go
type Lexer struct {
    // Contains unexported fields for lexical analysis state
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token

// TokenizeAll processes the entire input and returns all tokens
func (l *Lexer) TokenizeAll() ([]Token, []*LexError)

// HasErrors returns true if lexical errors were encountered
func (l *Lexer) HasErrors() bool

// GetErrors returns all lexical errors with position information
func (l *Lexer) GetErrors() []*LexError
```

#### Token

```go
type Token struct {
    Type    TokenType  // Token classification
    Literal string     // Original text from source
    Line    int        // Line number (1-indexed)
    Column  int        // Column number (1-indexed)
}
```

#### Error Handling

```go
type LexError struct {
    Message string
    Line    int
    Column  int
}

func (e *LexError) Error() string
```

### Token Types

All token types are exported constants of type `TokenType`:

```go
// Control flow
golexer.LET, golexer.CONST, golexer.FN, golexer.IF, golexer.ELSE
golexer.WHILE, golexer.FOR, golexer.RETURN, golexer.BREAK, golexer.CONTINUE

// Literals
golexer.IDENT, golexer.NUMBER, golexer.STRING, golexer.BACKTICK_STRING, golexer.CHAR
golexer.TRUE, golexer.FALSE, golexer.NULL

// Operators
golexer.PLUS, golexer.MINUS, golexer.MULTIPLY, golexer.DIVIDE, golexer.MODULUS
golexer.ASSIGN, golexer.PLUS_ASSIGN, golexer.MINUS_ASSIGN, golexer.MULTIPLY_ASSIGN, golexer.DIVIDE_ASSIGN, golexer.MODULUS_ASSIGN

// Comparison
golexer.EQL, golexer.NOT_EQL, golexer.LESS_THAN, golexer.LESS_THAN_EQL
golexer.GREATER_THAN, golexer.GREATER_THAN_EQL

// Logical
golexer.AND, golexer.OR, golexer.BANG

// Delimiters
golexer.LPAREN, golexer.RPAREN, golexer.LBRACE, golexer.RBRACE
golexer.LBRACKET, golexer.RBRACKET, golexer.COMMA, golexer.SEMICOLON, golexer.COLON, golexer.DOT

// Special
golexer.EOF, golexer.ILLEGAL
```

## Supported Tokens

### Numbers

#### Decimal Numbers
- Integers: `42`, `0`, `1000`
- Floats: `3.14`, `0.5`, `42.0`
- Scientific notation: `1e10`, `2.5e-3`, `1E+5`

#### Hexadecimal Numbers
- Lowercase: `0xff`, `0x1a2b`
- Uppercase: `0xFF`, `0X1A2B`

#### Binary Numbers
- Lowercase: `0b1010`, `0b1111`
- Uppercase: `0B1010`, `0B0000`

#### Octal Numbers
- Modern syntax: `0o777`, `0O123`
- Traditional syntax: `0777`, `0123`

### String Literals

#### Regular Strings
Double-quoted strings with escape sequence processing:
```
"hello world"
"line 1\nline 2"
"tab\tseparated\tvalues"
"quote: \"hello\""
"backslash: \\"
"null char: \0"
"hex escape: \x41"  // Equals "A"
```

#### Raw Strings
Backtick-quoted strings with no escape processing:
```
`raw string with \n literal backslashes`
`file path: C:\Users\Name\file.txt`
`multi
line
string`
```

#### Character Literals
Single-quoted character literals:
```
'a', 'Z', '0', '9'
'\n', '\t', '\r', '\\'
'\x41'  // Hex escape for 'A'
```

### Escape Sequences

| Sequence | Character | Description |
|----------|-----------|-------------|
| `\a` | `\x07` | Bell (alert) |
| `\b` | `\x08` | Backspace |
| `\f` | `\x0C` | Form feed |
| `\n` | `\x0A` | Newline |
| `\r` | `\x0D` | Carriage return |
| `\t` | `\x09` | Horizontal tab |
| `\v` | `\x0B` | Vertical tab |
| `\\` | `\x5C` | Backslash |
| `\'` | `\x27` | Single quote |
| `\"` | `\x22` | Double quote |
| `\0` | `\x00` | Null character |
| `\xNN` | | Hex escape (NN = hex digits) |

### Comments

#### Line Comments
```go
let x = 5; // This is a line comment
// Full line comment
```

#### Block Comments
```go
/* Single line block comment */

/*
 * Multi-line
 * block comment
 */

let x = /* inline */ 42;
```

### Operators and Punctuation

#### Arithmetic Operators
```go
+    // Addition
-    // Subtraction  
*    // Multiplication
/    // Division
%    // Modulus
```

#### Assignment Operators
```go
=     // Assignment
+=    // Add and assign
-=    // Subtract and assign
*=    // Multiply and assign
/=    // Divide and assign
%=    // Modulus and assign
```

#### Comparison Operators
```go
==    // Equal
!=    // Not equal
<     // Less than
<=    // Less than or equal
>     // Greater than
>=    // Greater than or equal
```

#### Logical Operators
```go
&&    // Logical AND
||    // Logical OR
!     // Logical NOT
```

#### Delimiters
```go
( )   // Parentheses
{ }   // Braces
[ ]   // Brackets
,     // Comma
;     // Semicolon
:     // Colon
.     // Dot (for property access, method calls)
```

## Error Handling

The lexer implements comprehensive error handling with detailed position information:

### Error Types

1. **Invalid Numbers**: `123abc`, `0xGHI`, `0b123`, `1e`
2. **Unterminated Literals**: `"unclosed string`, `'unclosed char`, `` `unclosed backtick``
3. **Invalid Escape Sequences**: `"\q"`, `'\x'`, `'\xGG'`
4. **Unexpected Characters**: `@`, `#`, `$`, single `&` or `|`
5. **Unterminated Comments**: `/* unclosed block comment`

### Error Recovery

The lexer continues processing after encountering errors, allowing you to collect all issues in a single pass:

```go
source := `
let x = 123abc;    // Error: invalid number
let y = "valid";   // This still processes correctly
let z = 0xGHI;     // Error: invalid hex
`

lexer := golexer.NewLexer(source)
tokens, errors := lexer.TokenizeAll()

// Gets both valid tokens AND all errors
fmt.Printf("Tokens: %d, Errors: %d\n", len(tokens), len(errors))
```

### Error Information

Each error includes precise location information:

```go
for _, err := range lexer.GetErrors() {
    fmt.Printf("Line %d, Column %d: %s\n", err.Line, err.Column, err.Message)
}
// Output: Line 2, Column 9: invalid number: numbers cannot be followed by letters
```

## Examples

### Complete Test File

The repository includes a comprehensive test file (`test.lang`) that demonstrates all lexer capabilities:

```bash
# Run the comprehensive test
go run cmd/main.go test.lang

# Expected output:
# Status: ✓ PASSED
# Tokens generated: 1700+
# Unique token types: 50+
# Lexical errors: 0
```

The test file contains real-world examples of:
- All number formats (decimal, hex, binary, octal, scientific notation)
- String literals with escape sequences and raw strings
- Character literals with all supported escapes
- Unicode identifier support
- Complete operator and punctuation coverage
- Complex expressions and nested structures
- Realistic function and data structure patterns

### Building a Simple Parser

```go
package main

import (
    "fmt"
    "github.com/codetesla51/golexer/golexer"
)

// Simple expression parser using the lexer
func parseExpression(source string) {
    lexer := golexer.NewLexer(source)
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        switch tok.Type {
        case golexer.NUMBER:
            fmt.Printf("Found number: %s\n", tok.Literal)
        case golexer.IDENT:
            fmt.Printf("Found variable: %s\n", tok.Literal)
        case golexer.PLUS, golexer.MINUS, golexer.MULTIPLY, golexer.DIVIDE:
            fmt.Printf("Found operator: %s\n", tok.Literal)
        case golexer.DOT:
            fmt.Printf("Found property access\n")
        }
    }
}

func main() {
    parseExpression("object.property + 42 * y")
}
```

### Basic Token Analysis

```go
func analyzeTokens(source string) {
    lexer := golexer.NewLexer(source)
    tokens, _ := lexer.TokenizeAll()
    
    // Count different token types
    counts := make(map[golexer.TokenType]int)
    for _, token := range tokens {
        counts[token.Type]++
    }
    
    fmt.Printf("Found %d total tokens\n", len(tokens))
    for tokenType, count := range counts {
        fmt.Printf("  %s: %d\n", tokenType, count)
    }
}
```

## Extending the Lexer

The lexer architecture makes it easy to add new language features. Here's how we recently added DOT token support:

### Adding New Punctuation Tokens

**Step 1**: Add the token type constant in `golexer/token.go`:
```go
const (
    // existing constants...
    DOT = "."  // Add new token type
    // rest of constants...
)
```

**Step 2**: Add the character mapping in `golexer/lexer.go`:
```go
var singleCharTokens = map[rune]TokenType{
    // existing mappings...
    '.': DOT,  // Map character to token type
}
```

**That's it!** The lexer will now recognize dots and classify them as DOT tokens.

### Adding New Keywords

**Step 1**: Add the token type constant in `golexer/token.go`:
```go
const (
    // existing constants...
    ASYNC = "ASYNC"
    AWAIT = "AWAIT" 
    CLASS = "CLASS"
)
```

**Step 2**: Add to the keywords map in `golexer/token.go`:
```go
var keywords = map[string]TokenType{
    // existing keywords...
    "async":  ASYNC,
    "await":  AWAIT,
    "class":  CLASS,
}
```

### Adding New Operators

For compound operators, add to the operators slice in `golexer/lexer.go`:
```go
var operators = []Operator{
    // existing operators...
    {"*", MULTIPLY, "**", POWER},        // * and **
    {"?", QUESTION, "??", NULL_COALESCE}, // ? and ??
}
```

Don't forget to add the corresponding token type constants in `token.go`:
```go
const (
    // existing constants...
    POWER        = "**"
    QUESTION     = "?"
    NULL_COALESCE = "??"
)
```

### Real Example: Adding DOT Support

Here's exactly what we did to add DOT token support:

**Before**: `object.property` would generate "unexpected character '.'" errors

**After adding DOT support**:
1. Added `DOT = "."` to token constants
2. Added `'.': DOT,` to singleCharTokens map
3. Now `object.property` tokenizes as: `IDENT("object")`, `DOT(".")`, `IDENT("property")`

**Result**: Zero lexical errors, proper tokenization of property access patterns.

### Testing New Features

When adding new tokens:

1. **Add examples to test.lang** - include your new token in realistic contexts
2. **Run the test**: `go run cmd/main.go test.lang`
3. **Verify zero errors** and check token distribution
4. **Add unit tests** for edge cases

## Performance

The lexer is designed for performance with a focus on simplicity and correctness:

- **UTF-8 support** with efficient rune processing
- **Single-pass tokenization** with error recovery
- **Memory efficient** token generation
- **1700+ tokens processed** in the comprehensive test with zero errors

Run the included performance test:
```bash
go run cmd/main.go test.lang
# Processes 1700+ tokens across 50+ token types instantly
```

## Unicode Support

The lexer uses Go's UTF-8 support for processing source code. Identifiers can contain Unicode letters as defined by Go's `unicode.IsLetter()` function:

```go
// These work because they are Unicode letters
переменная = 42;  // Russian
变量 = "test";     // Chinese
```

**Note**: For maximum compatibility, the included `test.lang` uses ASCII identifiers, but Unicode support is fully implemented.

## Testing

### Running Tests

```bash
# Run the test suite
go test ./...

# Run with coverage
go test -cover ./...

# Test with the comprehensive example
go run cmd/main.go test.lang
```

### Comprehensive Test File

The `test.lang` file serves as both documentation and validation:
- **1700+ tokens** demonstrating every feature
- **50+ unique token types** 
- **Real-world code patterns** including functions, arrays, objects
- **All number formats** and string variations
- **Zero lexical errors** when properly processed

### Expected Test Results

When running `go run cmd/main.go test.lang`, you should see:
```
Status: ✓ PASSED
Tokens generated: 1700+
Unique token types: 50+
Lexical errors: 0
```

If you see lexical errors, it typically indicates:
1. Missing token type support (like we had with DOT)
2. Unicode compatibility issues
3. Malformed test input

## Error Handling

### Error Types and Messages

The lexer provides detailed error messages for common issues:

| Error Type | Example Input | Error Message |
|------------|---------------|---------------|
| Invalid number | `123abc` | `invalid number: numbers cannot be followed by letters` |
| Invalid hex | `0xGHI` | `invalid hexadecimal number: contains non-hex characters` |
| Invalid binary | `0b123` | `invalid binary number: contains non-binary characters` |
| Invalid octal | `0o89` | `invalid octal number: contains non-octal characters` |
| Unterminated string | `"hello` | `unterminated string literal` |
| Invalid escape | `"\q"` | `unknown escape sequence '\q'` |
| Unterminated comment | `/* comment` | `unterminated block comment` |
| Invalid operator | `&` | `unexpected character '&' - did you mean '&&'?` |
| Unexpected character | `@` | `unexpected character '@' (Unicode: U+0040)` |

### Error Recovery

The lexer implements error recovery strategies:

1. **Continue after errors** - doesn't stop on first error
2. **Skip invalid sequences** - prevents cascading errors
3. **Collect all errors** - single pass finds all issues
4. **Maintain position tracking** - accurate line/column even after errors

## Contributing

Contributions are welcome! Please read our contributing guidelines:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-operator`
3. **Add tests** for any new functionality
4. **Test with test.lang**: Ensure `go run cmd/main.go test.lang` still passes
5. **Ensure tests pass**: `go test ./...`
6. **Submit a pull request**

### Development Setup

```bash
git clone https://github.com/codetesla51/golexer.git
cd golexer
go mod download
go test ./...
go run cmd/main.go test.lang  # Should show 0 errors
```

### Adding New Features

When extending the lexer:

1. **Identify the token type needed** (keyword, operator, punctuation)
2. **Add token constant** in `token.go`
3. **Add recognition logic** in `lexer.go` (keywords map, operators slice, or singleCharTokens map)
4. **Add examples to test.lang**
5. **Verify zero errors**: `go run cmd/main.go test.lang`
6. **Add unit tests** for edge cases

### Code Style

- Follow `go fmt` formatting
- Use `go vet` to check for issues  
- Add tests for new features
- Update documentation for API changes
- Ensure `test.lang` validates your changes


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by the lexical analysis techniques described in "Crafting Interpreters" by Robert Nystrom and "Writing An Interpreter In Go" by Thorsten Ball.

---

**Author**: Uthman Dev  
**Repository**: https://github.com/codetesla51/golexer  
**License**: MIT
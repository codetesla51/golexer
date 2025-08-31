# GoLexer

[![Go Reference](https://pkg.go.dev/badge/github.com/codetesla51/golexer.svg)](https://pkg.go.dev/github.com/codetesla51/golexer)
[![Go Report Card](https://goreportcard.com/badge/github.com/codetesla51/golexer)](https://goreportcard.com/report/github.com/codetesla51/golexer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive lexical analyzer (tokenizer) for Go that transforms source code into structured tokens. Built for creating compilers, interpreters, DSLs, configuration parsers, and code analysis tools.

## What is GoLexer?

A lexical analyzer breaks down source code into meaningful units called tokens. For example, `let x = 42 + y;` becomes tokens: `LET`, `IDENT(x)`, `ASSIGN(=)`, `NUMBER(42)`, `PLUS(+)`, `IDENT(y)`, `SEMICOLON(;)`. This is the foundation for building programming language tools.

## Features

- **Rich Token Set**: 50+ built-in token types covering modern programming constructs
- **Multiple Number Formats**: Decimal, hex, binary, octal, scientific notation with full validation
- **Advanced String Processing**: Regular strings, raw backtick strings, character literals with complete escape sequences
- **JSON Configuration System**: Extend the lexer with custom keywords, operators, and punctuation without code changes
- **Robust Error Recovery**: Continues processing after errors, collecting all issues with precise position tracking
- **UTF-8 Unicode Support**: Full support for international identifiers and multibyte characters
- **High Performance**: Single-pass tokenization with minimal memory allocations
- **Comprehensive Testing**: Validated with 1700+ tokens across complex real-world code patterns

## Installation

```bash
go get github.com/codetesla51/golexer
```

**Requirements**: Go 1.21 or later

## Quick Start

### Basic Tokenization

```go
package main

import (
    "fmt"
    "github.com/codetesla51/golexer/golexer"
)

func main() {
    source := `let total = 42 + 3.14 * count;`
    lexer := golexer.NewLexer(source)
    
    for {
        token := lexer.NextToken()
        if token.Type == golexer.EOF {
            break
        }
        fmt.Printf("%-12s %-10s (Line %d, Col %d)\n", 
            token.Type, token.Literal, token.Line, token.Column)
    }
}
```

**Output:**
```
LET          let        (Line 1, Col 1)
IDENT        total      (Line 1, Col 5)
=            =          (Line 1, Col 11)
NUMBER       42         (Line 1, Col 13)
+            +          (Line 1, Col 16)
NUMBER       3.14       (Line 1, Col 18)
*            *          (Line 1, Col 23)
IDENT        count      (Line 1, Col 25)
;            ;          (Line 1, Col 30)
```

### Batch Processing with Error Handling

```go
lexer := golexer.NewLexer(source)
tokens, errors := lexer.TokenizeAll()

fmt.Printf("Generated %d tokens\n", len(tokens))

if len(errors) > 0 {
    fmt.Printf("Found %d errors:\n", len(errors))
    for _, err := range errors {
        fmt.Printf("  %s\n", err.Error())
    }
} else {
    fmt.Println("No lexical errors found")
}
```

## Configuration System

Extend the lexer functionality using JSON configuration files. This allows you to add domain-specific keywords and operators without modifying the source code.

### Creating Configuration

Create `config.json`:
```json
{
  "additionalKeywords": {
    "unless": "UNLESS",
    "until": "UNTIL",
    "async": "ASYNC",
    "await": "AWAIT"
  },
  "additionalOperators": {
    "**": "POWER",
    "??": "NULL_COALESCE",
    "?.": "SAFE_NAVIGATION"
  },
  "additionalPunctuation": {
    "@": "AT_SYMBOL",
    "#": "HASH",
    "$": "DOLLAR"
  }
}
```

### Using Configuration

```go
lexer := golexer.NewLexerWithConfig(source, "config.json")

// Now recognizes extended syntax:
source := `
unless error {
    result = value ** 2 ?? fallback
    data = object ?. property
    user = @currentUser
}
`
```

### Graceful Error Handling

If the config file is missing or invalid, the lexer shows a warning and continues with defaults:
```
Warning: failed to load config file 'config.json': no such file or directory
Continuing with default configuration...
```

## Supported Tokens

### Number Formats

GoLexer supports all modern number formats with comprehensive validation:

| Format | Examples | Description |
|--------|----------|-------------|
| **Decimal** | `42`, `0`, `1000` | Standard integers |
| **Float** | `3.14`, `0.5`, `42.0` | Decimal points |
| **Scientific** | `1e10`, `2.5e-3`, `1E+5` | Exponential notation |
| **Hexadecimal** | `0xFF`, `0x1a2b`, `0X1A2B` | Base-16 with 0x prefix |
| **Binary** | `0b1010`, `0B1111` | Base-2 with 0b prefix |
| **Octal Modern** | `0o777`, `0O123` | Base-8 with 0o prefix |
| **Octal Legacy** | `0755`, `0123` | Traditional format |

### String and Character Literals

#### Regular Strings
Complete escape sequence support:
```go
"Hello, World!"           // Simple string
"Line 1\nLine 2"         // Newline
"Quote: \"Hello\""       // Escaped quote
"Tab\tSeparated"         // Tab character
"Hex: \x41\x42"          // Hex escapes (AB)
```

#### Raw Strings
No escape processing - literal text including backslashes:
```go
`Raw string with \n literal backslashes`
`File path: C:\Users\Name\file.txt`
`Multi
line
string`
```

#### Character Literals
```go
'a', 'Z', '0'            // Regular characters
'\n', '\t', '\r'         // Control characters
'\'', '\\'               // Escaped quotes
'\x41'                   // Hex escape for 'A'
```

#### Escape Sequences
| Escape | Result | Description |
|--------|--------|-------------|
| `\n` | newline | Line feed |
| `\t` | tab | Horizontal tab |
| `\r` | return | Carriage return |
| `\\` | backslash | Literal backslash |
| `\"` | quote | Double quote |
| `\'` | apostrophe | Single quote |
| `\0` | null | Null character |
| `\xNN` | hex char | Character by hex code |

### Keywords and Identifiers

#### Built-in Keywords
```
let const fn if else while for return break continue true false null
int float string bool char
```

#### Valid Identifiers
- Must start with letter or underscore
- Can contain letters, digits, underscores
- Unicode support: `café`, `résumé`, `变量`
- Examples: `variable1`, `_private`, `camelCase`, `snake_case`

### Operators

#### Arithmetic
```
+    -    *    /    %     // Basic operations
+=   -=   *=   /=   %=    // Compound assignment
```

#### Comparison
```
==   !=                   // Equality
<    <=   >    >=        // Relational
```

#### Logical
```
&&   ||   !              // AND, OR, NOT
```

**Note**: Single `&` and `|` produce helpful error messages suggesting the compound forms.

### Punctuation and Delimiters

#### Grouping
```
( ) { } [ ]              // Parentheses, braces, brackets
```

#### Separators
```
, ; : .                  // Comma, semicolon, colon, dot
```

### Comments

```go
// Line comments - rest of line ignored
let x = 42; // End of line comment

/* Block comments - can span multiple lines */
let y = /* inline */ 10;
```

## API Reference

### Core Functions

```go
// Create basic lexer
func NewLexer(input string) *Lexer

// Create lexer with JSON configuration
func NewLexerWithConfig(input, configFile string) *Lexer
```

### Tokenization Methods

```go
// Get next token (streaming)
func (l *Lexer) NextToken() Token

// Get all tokens at once (batch)
func (l *Lexer) TokenizeAll() ([]Token, []*LexError)
```

### Error Handling

```go
// Check for errors
func (l *Lexer) HasErrors() bool

// Get error details
func (l *Lexer) GetErrors() []*LexError
```

### Data Structures

#### Token
```go
type Token struct {
    Type    TokenType  // Token classification
    Literal string     // Original text
    Line    int        // Line number (1-indexed)
    Column  int        // Column position (1-indexed)
}
```

#### Error
```go
type LexError struct {
    Message string     // Error description
    Line    int        // Error line
    Column  int        // Error column
}
```

## Error Handling and Recovery

The lexer provides comprehensive error detection while continuing to process input, finding all problems in a single pass.

### Common Error Types

| Input | Error Message | Explanation |
|-------|---------------|-------------|
| `123abc` | `invalid number: numbers cannot be followed by letters` | Invalid number format |
| `0xGHI` | `invalid hexadecimal number: contains non-hex characters` | Bad hex digits |
| `0b123` | `invalid binary number: contains non-binary characters` | Invalid binary digits |
| `"hello` | `unterminated string literal` | Missing closing quote |
| `"test\q"` | `unknown escape sequence '\q'` | Invalid escape sequence |
| `&` | `unexpected character '&' - did you mean '&&'?` | Helpful suggestion |

### Error Recovery Example

```go
problemSource := `
let x = 123abc;      // Error: invalid number
let y = "valid";     // Continues processing
let z = 0xGHI;       // Error: invalid hex
let a = 42;          // Still processes correctly
`

lexer := golexer.NewLexer(problemSource)
tokens, errors := lexer.TokenizeAll()

// Result: 
// - tokens contains all valid tokens including "valid" and 42
// - errors contains detailed info about both problems
// - Processing never stops due to errors
```

## Advanced Usage Examples

### Building a Compiler Frontend

```go
type Compiler struct {
    lexer *golexer.Lexer
}

func (c *Compiler) CompileFile(filename string) error {
    source, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    c.lexer = golexer.NewLexer(string(source))
    tokens, errors := c.lexer.TokenizeAll()
    
    if len(errors) > 0 {
        return fmt.Errorf("lexical errors: %v", errors)
    }
    
    // Pass tokens to parser
    return c.parse(tokens)
}
```

### Configuration File Parser

```go
func ParseConfig(configFile string) (*AppConfig, error) {
    content, err := os.ReadFile(configFile)
    if err != nil {
        return nil, err
    }
    
    lexer := golexer.NewLexerWithConfig(string(content), "config-lang.json")
    config := &AppConfig{}
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        // Parse key: value pairs
        if tok.Type == golexer.IDENT {
            key := tok.Literal
            if colon := lexer.NextToken(); colon.Type == golexer.COLON {
                value := lexer.NextToken()
                config.Set(key, parseValue(value))
            }
        }
    }
    
    return config, nil
}
```

### Code Analysis Tool

```go
func AnalyzeCode(source string) {
    lexer := golexer.NewLexer(source)
    tokens, errors := lexer.TokenizeAll()
    
    // Generate statistics
    tokenCounts := make(map[golexer.TokenType]int)
    for _, token := range tokens {
        tokenCounts[token.Type]++
    }
    
    fmt.Printf("Analysis Results:\n")
    fmt.Printf("Total tokens: %d\n", len(tokens))
    fmt.Printf("Unique types: %d\n", len(tokenCounts))
    fmt.Printf("Errors found: %d\n", len(errors))
    
    // Show token distribution
    for tokenType, count := range tokenCounts {
        percentage := float64(count) / float64(len(tokens)) * 100
        fmt.Printf("  %-15s: %4d (%5.1f%%)\n", tokenType, count, percentage)
    }
}
```

### Syntax Highlighter

```go
func GenerateHighlighting(source string) []HighlightToken {
    lexer := golexer.NewLexer(source)
    var highlights []HighlightToken
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        highlights = append(highlights, HighlightToken{
            Text:   tok.Literal,
            Type:   mapToHighlightType(tok.Type),
            Line:   tok.Line,
            Column: tok.Column,
        })
    }
    
    return highlights
}

func mapToHighlightType(tokenType golexer.TokenType) string {
    switch tokenType {
    case golexer.LET, golexer.IF, golexer.WHILE:
        return "keyword"
    case golexer.STRING, golexer.CHAR:
        return "string"
    case golexer.NUMBER:
        return "number"
    default:
        return "default"
    }
}
```

## Performance

- **Time Complexity**: O(n) linear processing
- **Memory Usage**: O(1) streaming, O(n) batch processing
- **UTF-8 Handling**: Proper multibyte character support
- **Benchmark**: Processes 1700+ tokens across 400+ lines instantly

### When to Use Each Method

**Streaming (NextToken)**: Large files, memory constraints, real-time processing
**Batch (TokenizeAll)**: Complete analysis, small to medium files, when you need all tokens upfront

## Testing and Validation

### Running Tests

```bash
# Clone and test
git clone https://github.com/codetesla51/golexer.git
cd golexer

# Unit tests
go test ./golexer

# Comprehensive integration test
go run cmd/main.go test.lang
```

### Expected Results

The test suite processes a comprehensive example file with:
- 1700+ tokens across 400+ lines
- All number formats and string types
- Complete operator and keyword coverage
- Complex nested structures
- Zero lexical errors

Expected output:
```
=== Summary ===
File: test.lang
Lines processed: 400+
Tokens generated: 1700+
Unique token types: 45+
Lexical errors: 0
Status: ✓ PASSED
```

## Command Line Interface

The included CLI demonstrates all lexer capabilities:

```bash
# Analyze any file
go run cmd/main.go yourfile.txt

# Test with comprehensive example  
go run cmd/main.go test.lang
```

The CLI provides:
1. **Token-by-token output**: Each token with position
2. **Batch statistics**: Token counts and distribution
3. **Error reporting**: Detailed error messages with locations
4. **Summary**: Overall processing results

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Add tests for new functionality
4. Ensure `go run cmd/main.go test.lang` passes with 0 errors
5. Submit pull request

### Development Guidelines

- Follow Go conventions and use `gofmt`
- Add comprehensive comments for new features
- Include both positive and negative test cases
- Update documentation for significant changes
- Test with various input types and edge cases

## License

MIT License - see [LICENSE](LICENSE) file for details.

You can use GoLexer in commercial projects, modify it, and distribute your changes. The only requirement is including the original license notice.

## Acknowledgments

This project was inspired by the **Monkey lexer** from the excellent book ["Writing An Interpreter In Go"](https://interpreterbook.com/) by Thorsten Ball.

The foundational concepts of lexical analysis and token processing from that work provided the foundation for building this more comprehensive, production-ready lexer with extended number format support, configurable extensions, robust error recovery, and full Unicode handling.

---

**Author**: Uthman Dev  
**Repository**: https://github.com/codetesla51/golexer  

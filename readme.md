# GoLexer

[![Go Reference](https://pkg.go.dev/badge/github.com/codetesla51/golexer.svg)](https://pkg.go.dev/github.com/codetesla51/golexer)
[![Go Report Card](https://goreportcard.com/badge/github.com/codetesla51/golexer)](https://goreportcard.com/report/github.com/codetesla51/golexer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive lexical analyzer (tokenizer) library for Go. Designed for building programming languages, domain-specific languages (DSLs), configuration parsers, and template engines.

## Features

- **50+ Token Types**: Keywords, operators, literals, punctuation with precise position tracking
- **Multiple Number Formats**: Decimal, hex, binary, octal, scientific notation
- **String Processing**: Regular strings, raw strings, character literals with full escape sequences
- **JSON Configuration**: Extend lexer with custom keywords, operators, and punctuation
- **Robust Error Handling**: Graceful degradation with detailed error messages and position tracking
- **Unicode Support**: Full UTF-8 identifier support
- **Performance**: Single-pass tokenization processing 1700+ tokens instantly

## Installation

```bash
go get github.com/codetesla51/golexer
```

Requires Go 1.21 or later.

## Quick Start

### Basic Usage

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
}
```

### With Custom Configuration

Create `config.json`:
```json
{
  "additionalKeywords": {
    "unless": "UNLESS",
    "async": "ASYNC"
  },
  "additionalOperators": {
    "**": "POWER",
    "??": "NULL_COALESCE"
  },
  "additionalPunctuation": {
    "@": "AT_SYMBOL",
    "#": "HASH"
  }
}
```

Use with configuration:
```go
lexer := golexer.NewLexerWithConfig(source, "config.json")

// If config file is missing or invalid, shows warning and continues with defaults:
// Warning: failed to load config file 'config.json': no such file or directory
// Continuing with default configuration...
```

### Batch Processing

```go
lexer := golexer.NewLexer(source)
tokens, errors := lexer.TokenizeAll()

fmt.Printf("Processed %d tokens with %d errors\n", len(tokens), len(errors))
```

## API Reference

### Core Functions

```go
// Basic lexer
func NewLexer(input string) *Lexer

// Lexer with JSON configuration (graceful error handling)
func NewLexerWithConfig(input, configFile string) *Lexer

// Tokenization
func (l *Lexer) NextToken() Token
func (l *Lexer) TokenizeAll() ([]Token, []*LexError)

// Error handling
func (l *Lexer) HasErrors() bool
func (l *Lexer) GetErrors() []*LexError
```

### Token Structure

```go
type Token struct {
    Type    TokenType  // Token classification
    Literal string     // Original text
    Line    int        // Line number (1-indexed)
    Column  int        // Column number (1-indexed)
}
```

## Supported Tokens

### Numbers
GoLexer supports all modern number formats with proper validation:

- **Decimal integers**: `42`, `0`, `1000`
- **Decimal floats**: `3.14`, `0.5`, `42.0`  
- **Scientific notation**: `1e10`, `2.5e-3`, `1E+5` (supports +/- exponents)
- **Hexadecimal**: `0xff`, `0xFF`, `0x1a2b`, `0X1A2B`
- **Binary**: `0b1010`, `0b1111`, `0B1010`, `0B0000`
- **Octal modern**: `0o777`, `0O123` (explicit octal prefix)
- **Octal traditional**: `0755`, `0123` (legacy format)

### String Literals
Complete string processing with escape sequence support:

#### Regular Strings (Double-quoted)
```
"hello world"
"line 1\nline 2" 
"tab\tseparated\tvalues"
"quote: \"hello\""
"backslash: \\"
"null char: \0"
"hex escape: \x41"  // Equals "A"
```

#### Raw Strings (Backtick-quoted)
No escape processing - literal text including backslashes:
```
`raw string with \n literal backslashes`
`file path: C:\Users\Name\file.txt`
`multi
line
string`
```

#### Character Literals (Single-quoted)
```
'a', 'Z', '0', '9'     // Regular characters
'\n', '\t', '\r', '\\'  // Escape sequences  
'\x41'                  // Hex escape for 'A'
```

### Escape Sequences
| Sequence | Result | Description |
|----------|--------|-------------|
| `\n` | newline | Line break |
| `\t` | tab | Horizontal tab |
| `\r` | return | Carriage return |
| `\\` | backslash | Literal backslash |
| `\"` | quote | Double quote |
| `\'` | apostrophe | Single quote |
| `\0` | null | Null character |
| `\xNN` | character | Hex escape (NN = hex digits) |

### Keywords and Identifiers
Built-in language keywords:
```
let const fn if else while for return break continue true false null
int float string bool char
```

Valid identifiers: `variable1`, `_underscore`, `CamelCase`, `snake_case`, `mixed123`

### Operators
Complete operator set with compound assignments:

- **Arithmetic**: `+` `-` `*` `/` `%`
- **Assignment**: `=` `+=` `-=` `*=` `/=` `%=`
- **Comparison**: `==` `!=` `<` `<=` `>` `>=`
- **Logical**: `&&` `||` `!`

### Punctuation and Delimiters
```
( ) { } [ ]    // Grouping and blocks
, ; : .        // Separators and access
```

### Comments
```go
// Line comments - rest of line ignored
/* Block comments - can span multiple lines */
let x = /* inline comment */ 42;
```

## Error Handling

GoLexer provides comprehensive error handling with precise position tracking and graceful recovery strategies.

### Lexical Error Types
The lexer detects and reports various syntax errors:

| Error Type | Input Example | Error Message |
|------------|---------------|---------------|
| Invalid numbers | `123abc` | `invalid number: numbers cannot be followed by letters` |
| Invalid hex | `0xGHI` | `invalid hexadecimal number: contains non-hex characters` |
| Invalid binary | `0b123` | `invalid binary number: contains non-binary characters` |
| Unterminated strings | `"hello` | `unterminated string literal` |
| Invalid escapes | `"\q"` | `unknown escape sequence '\q'` |
| Unterminated comments | `/* comment` | `unterminated block comment` |
| Unexpected chars | `@` | `unexpected character '@' (Unicode: U+0040)` |
| Invalid operators | `&` | `unexpected character '&' - did you mean '&&'?` |

### Error Recovery
The lexer continues processing after errors, collecting all issues in a single pass:

```go
source := `
let x = 123abc;    // Error: invalid number
let y = "valid";   // This processes correctly  
let z = 0xGHI;     // Error: invalid hex
`

lexer := golexer.NewLexer(source)
tokens, errors := lexer.TokenizeAll()
// Gets both valid tokens AND all errors
```

### Configuration Error Handling
**New Feature**: Graceful configuration loading with informative error messages.

The system handles configuration issues elegantly:

```bash
# Missing configuration file
Warning: failed to load config file 'config.json': no such file or directory
Continuing with default configuration...

# Invalid JSON syntax  
Warning: failed to load config file 'config.json': invalid character '}' looking for beginning of object key string
Continuing with default configuration...

# Permission denied
Warning: failed to load config file 'config.json': permission denied
Continuing with default configuration...
```

**Key Benefits:**
- **Uninterrupted development** - never crashes due to config issues
- **Clear diagnostic messages** - users know exactly what went wrong
- **Automatic fallback** - continues with working default settings  
- **Professional UX** - handles edge cases gracefully

## Examples

### Real-world Usage Patterns

#### Building a Configuration Parser
```go
func parseConfig(source string) map[string]interface{} {
    lexer := golexer.NewLexer(source)
    config := make(map[string]interface{})
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        switch tok.Type {
        case golexer.IDENT:
            key := tok.Literal
            // Expect ':' then value
            if nextTok := lexer.NextToken(); nextTok.Type == golexer.COLON {
                valueTok := lexer.NextToken()
                config[key] = parseValue(valueTok)
            }
        }
    }
    return config
}
```

#### Token Analysis and Statistics  
```go
func analyzeCode(source string) {
    lexer := golexer.NewLexer(source)
    tokens, errors := lexer.TokenizeAll()
    
    // Count token types
    counts := make(map[golexer.TokenType]int)
    for _, token := range tokens {
        counts[token.Type]++
    }
    
    fmt.Printf("Code Analysis Results:\n")
    fmt.Printf("Total tokens: %d\n", len(tokens))
    fmt.Printf("Total errors: %d\n", len(errors))
    
    // Show most common tokens
    for tokenType, count := range counts {
        fmt.Printf("  %-15s: %d\n", tokenType, count)
    }
}
```

#### Custom Language Extension
```go
// config.json for JavaScript-like syntax
{
  "additionalKeywords": {
    "class": "CLASS",
    "extends": "EXTENDS", 
    "async": "ASYNC",
    "await": "AWAIT"
  },
  "additionalOperators": {
    "**": "EXPONENT",
    "?.": "OPTIONAL_CHAIN",
    "??": "NULL_COALESCE"
  }
}

// Usage
lexer := golexer.NewLexerWithConfig(jsCode, "js-config.json")
// Now recognizes: class MyClass extends BaseClass { async method() { ... } }
```

## Advanced Features

### Unicode and Internationalization
Full UTF-8 identifier support allows international variable names:
```go
// Valid identifiers in different languages
let переменная = 42;      // Russian
const 变量 = "value";     // Chinese  
fn función() { ... }      // Spanish
```

### Performance Characteristics
- **Single-pass scanning**: Complete tokenization in one iteration
- **Memory efficient**: Minimal allocations, reuses buffers where possible
- **UTF-8 optimized**: Proper rune handling without unnecessary conversions
- **Error recovery**: Continues processing after errors without performance penalty

**Benchmark**: Processes 1700+ tokens with 50+ token types instantly

### Extending the Lexer

#### Method 1: JSON Configuration (Recommended)
Extend without touching source code:

```json
{
  "additionalKeywords": {
    "unless": "UNLESS",
    "until": "UNTIL"
  },
  "additionalOperators": {
    "**": "POWER",
    "<=>": "SPACESHIP"
  },
  "additionalPunctuation": {
    "@": "AT_SYMBOL", 
    "#": "HASH"
  }
}
```

#### Method 2: Source Code Modification
For permanent built-in support:

1. **Add token type** in `golexer/token.go`:
   ```go
   const DOT = "."
   ```

2. **Add recognition** in `golexer/lexer.go`:
   ```go
   var singleCharTokens = map[rune]TokenType{
       '.': DOT,  // Maps '.' character to DOT token
   }
   ```

3. **Test the change**:
   ```bash
   go run cmd/main.go test.lang  # Should show 0 errors
   ```

## Testing and Validation

### Comprehensive Test Suite
Run the full validation:

```bash
git clone https://github.com/codetesla51/golexer.git
cd golexer

# Basic test
go test ./...

# Comprehensive feature test  
go run cmd/main.go test.lang
```

**Expected results:**
```
Analyzing file: test.lang
File size: 15000+ bytes

=== Token-by-token processing ===
[Shows each token with position]

=== Batch processing ===
Total tokens: 1700+
Total errors: 0

Status: ✓ PASSED
```

### What the Test Covers
The `test.lang` file validates:
- ✅ All number formats (decimal, hex, binary, octal, scientific)
- ✅ String literals with all escape sequences
- ✅ Character literals and special characters
- ✅ Complete operator and punctuation coverage
- ✅ Realistic code patterns (functions, objects, arrays)
- ✅ Comment processing (line and block)
- ✅ Complex expressions and nested structures


## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/name`
3. Add tests and ensure `go run cmd/main.go test.lang` passes
4. Test error conditions (missing configs, invalid files)
5. Submit pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

**Author**: Uthman Dev | **Repository**: https://github.com/codetesla51/golexer
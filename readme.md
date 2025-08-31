# GoLexer

[![Go Reference](https://pkg.go.dev/badge/github.com/codetesla51/golexer.svg)](https://pkg.go.dev/github.com/codetesla51/golexer)
[![Go Report Card](https://goreportcard.com/badge/github.com/codetesla51/golexer)](https://goreportcard.com/report/github.com/codetesla51/golexer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance lexical analyzer (tokenizer) for Go that transforms source code into structured tokens. Perfect for building compilers, interpreters, DSLs, configuration parsers, and code analysis tools.

## What is a Lexical Analyzer?

A lexical analyzer (also called a tokenizer or lexer) is the first phase of language processing. It reads source code as a stream of characters and groups them into meaningful units called tokens. For example, the code `let x = 42 + y;` becomes tokens: `LET`, `IDENT(x)`, `ASSIGN(=)`, `NUMBER(42)`, `PLUS(+)`, `IDENT(y)`, `SEMICOLON(;)`.

This is essential for:
- **Compilers**: Converting source code to machine code
- **Interpreters**: Executing code line by line  
- **IDEs**: Syntax highlighting and error detection
- **Code analyzers**: Finding bugs and style issues
- **Configuration parsers**: Reading structured config files

## Features

- **Rich Token Set**: 50+ built-in token types covering modern programming language constructs
- **Multiple Number Formats**: Decimal, hexadecimal, binary, octal, and scientific notation with validation
- **String Processing**: Regular strings, raw backtick strings, and character literals with complete escape sequences
- **Configurable Extensions**: Add custom keywords, operators, and punctuation via JSON configuration
- **Robust Error Recovery**: Continues processing after errors with detailed position tracking
- **UTF-8 Unicode Support**: International identifiers and proper multibyte character handling  
- **Performance Optimized**: Single-pass tokenization with minimal memory allocations
- **Position Tracking**: Precise line and column information for every token and error

## Installation

```bash
go get github.com/codetesla51/golexer
```

**Requirements**: Go 1.21 or later

## Quick Start

### Basic Tokenization

A lexer converts source code text into a sequence of tokens. Each token has a type (what kind of element it is), the original text (literal), and position information.

```go
package main

import (
    "fmt"
    "github.com/codetesla51/golexer/golexer"
)

func ParseCustomConfig(configFile string) (*AppConfig, error) {
    content, err := os.ReadFile(configFile)
    if err != nil {
        return nil, err
    }
    
    // Use custom language configuration
    lexer := golexer.NewLexerWithConfig(string(content), "custom-lang.json")
    
    config := &AppConfig{}
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        // Parse custom syntax like: server => { host: "localhost", port: 8080 }
        switch tok.Type {
        case golexer.TokenType("SERVER"):
            config.Server, err = parseServerBlock(lexer)
        case golexer.TokenType("DATABASE"):
            config.Database, err = parseDatabaseBlock(lexer)
        case golexer.TokenType("CACHE"):
            config.Cache, err = parseCacheBlock(lexer)
        }
        
        if err != nil {
            return nil, err
        }
    }
    
    return config, nil
}
```

## Architecture and Performance

### Design Principles

The GoLexer follows several key design principles that make it reliable and efficient:

**Single Responsibility**: Each component has a focused purpose:
- `lexer.go`: Core tokenization logic
- `token.go`: Token type definitions and keyword mappings  
- `errors.go`: Error handling and reporting
- `config.go`: Configuration loading and merging

**Error Recovery**: When the lexer encounters an error, it:
1. Records the error with precise position information
2. Attempts to recover and continue processing
3. Collects ALL errors in a single pass
4. Never crashes or stops unexpectedly

**Memory Efficiency**: The lexer minimizes memory allocations by:
- Reading characters one at a time (streaming approach)
- Reusing buffers where possible
- Not storing the entire token stream in memory (unless using TokenizeAll)
- Using efficient string slicing instead of copying

**Unicode First**: Proper UTF-8 handling throughout:
- Correctly handles multibyte characters
- Supports international identifiers
- Maintains accurate position tracking across character boundaries

**Extensibility**: The configuration system allows customization without touching source code:
- JSON-based configuration is easy to understand and modify
- Graceful fallback when configuration fails
- Merges custom settings with built-in defaults

### Performance Characteristics

**Time Complexity**: O(n) where n is input length
- Each character is read exactly once
- No backtracking or re-scanning
- Linear performance regardless of input complexity

**Space Complexity**: O(1) for streaming, O(n) for batch processing  
- Streaming approach (NextToken) uses constant memory
- Batch approach (TokenizeAll) stores all tokens in memory
- Error collection grows with number of errors found

**UTF-8 Optimization**: 
- Proper rune handling without unnecessary string conversions
- Efficient multibyte character processing
- Accurate position tracking across character boundaries

**Benchmark Results**: 
The test suite processes 1700+ tokens across 400+ lines of code instantly with zero errors, demonstrating the lexer's efficiency with real-world code patterns.

### When to Use Each Approach

**Streaming (NextToken)**:
```go
// Use for: Large files, memory-constrained environments, real-time processing
lexer := golexer.NewLexer(source)
for {
    token := lexer.NextToken()
    if token.Type == golexer.EOF {
        break
    }
    processTokenImmediately(token)
}
```

**Batch (TokenizeAll)**:
```go
// Use for: Small to medium files, when you need all tokens upfront
lexer := golexer.NewLexer(source)
tokens, errors := lexer.TokenizeAll()
analyzeAllTokens(tokens)
```

## Command Line Interface

The included CLI tool demonstrates the lexer's capabilities and provides a useful testing interface:

### Basic Usage

```bash
# Analyze any text file
go run cmd/main.go yourfile.txt

# Test with the comprehensive example
go run cmd/main.go test.lang
```

### CLI Output Explanation

The CLI provides three levels of analysis:

**1. Token-by-token processing**: Shows each token as it's generated
```
Type: LET             Literal: 'let'           Line:  1 Column:  1
Type: IDENT           Literal: 'variable'      Line:  1 Column:  5
Type: =               Literal: '='             Line:  1 Column:  14
```

**2. Batch processing with statistics**: Overall metrics and token distribution
```
Total tokens: 1700+
Total errors: 0

Token distribution:
  IDENT          : 245
  NUMBER         : 189
  STRING         : 67
```

**3. Syntax validation**: Error reporting and overall status
```
✓ No lexical errors found - file is syntactically valid at lexical level
Status: ✓ PASSED
```

### Using the CLI for Development

The CLI is especially useful when:
- **Testing new features**: Verify that custom configurations work correctly
- **Debugging issues**: See exactly how input is tokenized
- **Performance testing**: Process large files and measure performance
- **Validation**: Ensure your syntax is correctly recognized

## Contributing

We welcome contributions to make GoLexer even better! Here's how to get involved:

### Development Setup

```bash
# Clone the repository
git clone https://github.com/codetesla51/golexer.git
cd golexer

# The project has no external dependencies beyond Go standard library
go mod tidy

# Run the test suite
go test ./golexer

# Test with comprehensive example
go run cmd/main.go test.lang
```

### Types of Contributions

**Bug Reports**: Found an issue? Please include:
- Input that causes the problem
- Expected vs actual behavior  
- Go version and operating system
- Complete error messages

**Feature Requests**: Ideas for new functionality:
- Description of the use case
- Examples of how it would work
- Consideration of backward compatibility

**Code Contributions**: Pull requests should:
- Include tests for new functionality
- Maintain the existing code style
- Update documentation as needed
- Ensure `go run cmd/main.go test.lang` passes

### Code Style Guidelines

- **Follow Go conventions**: Use `gofmt`, follow standard naming
- **Comprehensive comments**: Explain complex logic and design decisions
- **Error handling**: Always handle errors gracefully with helpful messages  
- **Test coverage**: Include both positive and negative test cases
- **Documentation**: Update README for new features or significant changes

### Testing Your Changes

Before submitting a pull request:

1. **Unit tests pass**: `go test ./golexer`
2. **Integration test passes**: `go run cmd/main.go test.lang` shows 0 errors
3. **Error conditions work**: Test with malformed input to ensure proper error handling
4. **Configuration works**: If you modified the config system, test with various JSON files
5. **Performance maintained**: Large inputs should still process quickly

### Adding New Token Types

If you need to add new built-in token types:

1. **Add the constant** in `golexer/token.go`:
   ```go
   const (
       // ... existing tokens ...
       NEW_TOKEN = "NEW_TOKEN"
   )
   ```

2. **Add recognition logic** in `golexer/lexer.go`:
   - For keywords: Add to `keywords` map
   - For operators: Add to `operators` slice  
   - For punctuation: Add to `singleCharTokens` map

3. **Add tests** in `golexer/lexer_test.go` and update `test.lang`

4. **Update documentation** in this README

## License

MIT License - see [LICENSE](LICENSE) file for details.

This means you can:
- ✅ Use GoLexer in commercial projects
- ✅ Modify the source code
- ✅ Distribute your modifications
- ✅ Include it in proprietary software

The only requirement is including the original license notice.

## Acknowledgments

This project was inspired by the **Monkey lexer** from the excellent book ["Writing An Interpreter In Go"](https://interpreterbook.com/) by Thorsten Ball. 

The foundational concepts of lexical analysis, token-by-token processing, and error handling techniques from that work provided the inspiration for building this more comprehensive, production-ready lexer.

**What started as a learning exercise** from the Monkey language has evolved into a full-featured lexical analyzer suitable for real-world applications, with extensive number format support, configurable extensions, robust error recovery, and comprehensive Unicode handling.

**Key expansions beyond the original**:
- Support for multiple number formats (hex, binary, octal, scientific notation)
- Raw string literals and complete escape sequence handling
- JSON-based configuration system for extensibility
- Comprehensive error recovery with helpful messages
- Production-ready error handling and position tracking
- Full UTF-8 Unicode support for international development
- Extensive test suite with real-world code patterns

---

**Author**: Uthman Dev  
**Repository**: https://github.com/codetesla51/golexer  
**Documentation**: https://pkg.go.dev/github.com/codetesla51/golexer

**Questions or Issues?** Please open an issue on GitHub or check the documentation for detailed API information. main() {
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

### Understanding Token Types

Each token has a `Type` that categorizes what it represents:

- **IDENT**: Identifiers (variable names, function names)
- **NUMBER**: Numeric literals (integers, floats, hex, binary, etc.)
- **STRING**: Text enclosed in quotes
- **Keywords**: Reserved words like `let`, `if`, `while`
- **Operators**: Mathematical and logical operations like `+`, `==`, `&&`
- **Punctuation**: Structural elements like `{`, `}`, `;`, `,`

### Batch Processing with Error Handling

Instead of getting tokens one by one, you can process the entire input at once:

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

This approach is useful when you need all tokens upfront (like for a compiler) rather than streaming them.

## Configuration System

The lexer can be extended without modifying its source code using JSON configuration files. This allows you to add custom keywords, operators, and punctuation for domain-specific languages.

### How Configuration Works

The configuration system merges your custom definitions with the built-in ones:

1. **Keywords**: Custom words that should be treated as special tokens instead of regular identifiers
2. **Operators**: Custom symbol combinations that perform operations
3. **Punctuation**: Custom single characters that have special meaning

### Creating a Configuration File

Create a file named `config.json`:

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

**Format Explanation**:
- **Key**: The text to recognize (e.g., "unless", "**", "@")
- **Value**: The token type name it should become (e.g., "UNLESS", "POWER", "AT_SYMBOL")

### Using Configuration

```go
// Load lexer with custom configuration
lexer := golexer.NewLexerWithConfig(source, "config.json")

// Now the lexer recognizes extended syntax:
source := `
unless error {
    result = value ** 2 ?? fallback
    data = object ?. property
    user = @currentUser
}
`

// The lexer will now generate:
// UNLESS token for "unless"
// POWER token for "**"  
// NULL_COALESCE token for "??"
// SAFE_NAVIGATION token for "?."
// AT_SYMBOL token for "@"
```

### Graceful Configuration Error Handling

If your configuration file has problems, the lexer doesn't crash. Instead, it shows a warning and continues with default settings:

```
Warning: failed to load config file 'config.json': no such file or directory
Continuing with default configuration...
```

This happens when:
- Configuration file doesn't exist
- File contains invalid JSON
- File permissions prevent reading
- Any other file system error

Your program keeps running with the standard token set.

## Complete Token Reference

### Numbers with Full Validation

GoLexer recognizes many number formats that are common in modern programming languages:

| Format | Examples | Description | Notes |
|--------|----------|-------------|-------|
| **Decimal** | `42`, `0`, `1000` | Standard base-10 integers | Most common format |
| **Float** | `3.14`, `0.5`, `42.0` | Numbers with decimal points | Requires digits after decimal |
| **Scientific** | `1e10`, `2.5e-3`, `1E+5` | Exponential notation | e/E followed by optional +/- and digits |
| **Hexadecimal** | `0xFF`, `0x1a2b`, `0X1A2B` | Base-16 with 0x prefix | Case insensitive for both prefix and digits |
| **Binary** | `0b1010`, `0B1111` | Base-2 with 0b prefix | Only 0 and 1 digits allowed |
| **Octal Modern** | `0o777`, `0O123` | Base-8 with 0o prefix | Explicit octal format |
| **Octal Legacy** | `0755`, `0123` | Traditional octal format | Starts with 0, uses digits 0-7 |

### Understanding Scientific Notation

Scientific notation represents very large or very small numbers efficiently:
- `1e10` = 1 × 10^10 = 10,000,000,000
- `2.5e-3` = 2.5 × 10^-3 = 0.0025
- `1E+5` = 1 × 10^5 = 100,000

### Number Format Error Examples

The lexer validates number formats and provides helpful error messages:

- `123abc` → `invalid number: numbers cannot be followed by letters`
- `0xGHI` → `invalid hexadecimal number: contains non-hex characters`
- `0b123` → `invalid binary number: contains non-binary characters`
- `0o89` → `invalid octal number: contains non-octal characters`
- `1e` → `invalid scientific notation: exponent must contain digits`

### String and Character Literals

#### Regular Strings (Double-quoted)

Regular strings are enclosed in double quotes and support escape sequences:

```go
"Hello, World!"           // Simple string
"Line 1\nLine 2"         // Newline character
"Quote: \"Hello\""       // Escaped quote inside string
"Backslash: \\"          // Escaped backslash
"Tab\tSeparated"         // Tab character
"Hex: \x41\x42"          // Hex escapes (produces "AB")
```

#### Raw Strings (Backtick-quoted)

Raw strings are enclosed in backticks and contain literal text with no escape processing:

```go
`Raw string with \n literal backslashes`    // \n remains as text, not newline
`File path: C:\Users\Name\file.txt`         // Backslashes stay as-is
`Multi
line
string`                                      // Actual newlines preserved
```

**When to use raw strings**: Configuration file paths, regular expressions, multi-line text, or any content with many backslashes.

#### Character Literals (Single-quoted)

Character literals represent single characters and are enclosed in single quotes:

```go
'a', 'Z', '0', '9'       // Regular ASCII characters
'\n', '\t', '\r'         // Control characters via escape sequences
'\'', '\\'               // Escaped quotes and backslash
'\x41'                   // Hex escape for 'A'
```

#### Complete Escape Sequences

Escape sequences start with a backslash and represent special characters:

| Escape | Result | Description | Usage |
|--------|--------|-------------|-------|
| `\n` | newline | Line feed (ASCII 10) | New line in text |
| `\r` | return | Carriage return (ASCII 13) | Windows line endings |
| `\t` | tab | Horizontal tab (ASCII 9) | Alignment and indentation |
| `\v` | vtab | Vertical tab (ASCII 11) | Rare, vertical spacing |
| `\f` | form feed | Form feed (ASCII 12) | Page breaks in printing |
| `\a` | bell | Alert/bell (ASCII 7) | Audio alert |
| `\b` | backspace | Backspace (ASCII 8) | Delete previous character |
| `\\` | backslash | Literal backslash | When you need actual \ character |
| `\"` | quote | Double quote | Quote inside double-quoted string |
| `\'` | apostrophe | Single quote | Quote inside single-quoted char |
| `\0` | null | Null character (ASCII 0) | String termination marker |
| `\xNN` | hex char | Character by hex code | NN must be exactly 2 hex digits |

**Hex escape example**: `\x41` = 65 decimal = 'A', `\x0A` = 10 decimal = newline

### Keywords and Identifiers

#### Built-in Keywords

Keywords are reserved words with special meaning in the language:

```
let const fn if else while for return break continue true false null
int float string bool char
```

These cannot be used as variable names - they have special purposes in language syntax.

#### Valid Identifier Rules

Identifiers are names for variables, functions, etc. They must follow these rules:

**Valid patterns**:
- **Must start with**: Letter (a-z, A-Z) or underscore (_)
- **Can contain**: Letters, digits (0-9), underscores
- **Examples**: `variable`, `count`, `data`, `user_name`, `_private`, `__internal`, `variable1`, `temp2`, `item_123`, `camelCase`, `PascalCase`, `snake_case`

**Unicode support**: International characters are allowed:
- `café`, `résumé` (European characters)
- `变量` (Chinese characters)
- Any valid Unicode letter

**Invalid examples**:
- `123var` - starts with digit
- `my-var` - contains hyphen
- `user@domain` - contains @

### Complete Operator Set

Operators perform operations on values. GoLexer supports several categories:

#### Arithmetic Operators
```
+    -    *    /    %     // Basic math: add, subtract, multiply, divide, modulo
+=   -=   *=   /=   %=    // Compound assignment: modify and assign in one step
```

**Compound assignment explanation**: `x += 5` is equivalent to `x = x + 5`

#### Comparison Operators  
```
==   !=                   // Equality: equal, not equal
<    <=   >    >=        // Relational: less than, less/equal, greater than, greater/equal
```

#### Logical Operators
```
&&   ||   !              // Logical: AND, OR, NOT
```

**Important**: Single `&` and `|` are intentionally invalid and produce helpful error messages:
- `&` → `unexpected character '&' - did you mean '&&'?`
- `|` → `unexpected character '|' - did you mean '||'?`

This prevents common mistakes from C-style languages where single & and | have different meanings.

#### Assignment
```
=                        // Simple assignment: store value in variable
```

### Punctuation and Delimiters

These characters structure the language syntax:

#### Grouping Characters
```
( )                      // Parentheses: group expressions, function calls
{ }                      // Braces: code blocks, object literals  
[ ]                      // Brackets: array literals, indexing
```

#### Separator Characters
```
,                        // Comma: separate items in lists
;                        // Semicolon: end statements
:                        // Colon: key-value pairs, type annotations
.                        // Dot: property access, decimal points
```

### Comments

Comments are text ignored by the lexer, used for documentation:

#### Line Comments
```go
// This is a line comment - everything after // is ignored
let x = 42; // You can put comments at the end of lines
```

#### Block Comments
```go
/* This is a block comment */
/*
This is a multi-line
block comment that spans
several lines
*/
let y = /* inline comment */ 10; // Block comments can go anywhere
```

**Nested comments**: Block comments cannot be nested - `/* outer /* inner */ */` ends at the first `*/`.

## API Reference

### Core Functions

```go
// Create basic lexer with default token set
func NewLexer(input string) *Lexer

// Create lexer with JSON configuration file (handles errors gracefully)
func NewLexerWithConfig(input, configFile string) *Lexer
```

**Parameter explanation**:
- `input`: The source code text to tokenize
- `configFile`: Path to JSON configuration file (can be relative or absolute)

### Tokenization Methods

```go
// Get next token from input (streaming approach)
func (l *Lexer) NextToken() Token

// Get all tokens at once (batch approach)
func (l *Lexer) TokenizeAll() ([]Token, []*LexError)
```

**When to use each**:
- **NextToken()**: When processing large files or implementing streaming parsers
- **TokenizeAll()**: When you need all tokens upfront for analysis or when input is small

### Error Handling

```go
// Check if any errors occurred during tokenization
func (l *Lexer) HasErrors() bool

// Get detailed error information with positions
func (l *Lexer) GetErrors() []*LexError
```

### Data Structures

#### Token Structure
```go
type Token struct {
    Type    TokenType  // What kind of token (e.g., "NUMBER", "IDENT", "IF")
    Literal string     // Original text from source code
    Line    int        // Line number where token appears (starts at 1)
    Column  int        // Column position where token starts (starts at 1)
}
```

**Example**: For input `let x = 42`, the number token would be:
```go
Token{
    Type: "NUMBER",
    Literal: "42", 
    Line: 1,
    Column: 9
}
```

#### Error Structure
```go
type LexError struct {
    Message string     // Human-readable error description
    Line    int        // Line number where error occurred
    Column  int        // Column position where error occurred
}

// Implements the standard Go error interface
func (e *LexError) Error() string
```

**Example error**: For input `123abc`, you'd get:
```go
LexError{
    Message: "invalid number: numbers cannot be followed by letters",
    Line: 1,
    Column: 1
}
```

## Advanced Usage Examples

### Configuration Parser

This example shows how to build a simple configuration file parser:

```go
func parseConfigFile(filename string) (map[string]interface{}, error) {
    // Read the configuration file
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    // Create lexer to tokenize the config syntax
    lexer := golexer.NewLexer(string(content))
    config := make(map[string]interface{})
    
    // Process tokens to extract key-value pairs
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        // Look for pattern: IDENTIFIER : VALUE
        if tok.Type == golexer.IDENT {
            key := tok.Literal
            
            // Expect colon separator
            if colon := lexer.NextToken(); colon.Type == golexer.COLON {
                value := lexer.NextToken()
                config[key] = parseTokenValue(value)
            }
        }
    }
    
    return config, nil
}

func parseTokenValue(token golexer.Token) interface{} {
    switch token.Type {
    case golexer.NUMBER:
        // Convert string to number (simplified)
        if strings.Contains(token.Literal, ".") {
            if f, err := strconv.ParseFloat(token.Literal, 64); err == nil {
                return f
            }
        } else {
            if i, err := strconv.Atoi(token.Literal); err == nil {
                return i
            }
        }
    case golexer.STRING:
        return token.Literal
    case golexer.TRUE:
        return true
    case golexer.FALSE:
        return false
    }
    return token.Literal
}
```

### Code Analysis Tool

This example analyzes source code and generates statistics:

```go
func analyzeSourceCode(source string) {
    lexer := golexer.NewLexer(source)
    tokens, errors := lexer.TokenizeAll()
    
    // Count occurrences of each token type
    tokenCounts := make(map[golexer.TokenType]int)
    for _, token := range tokens {
        tokenCounts[token.Type]++
    }
    
    // Generate analysis report
    fmt.Printf("Code Analysis Report\n")
    fmt.Printf("===================\n")
    fmt.Printf("Total tokens: %d\n", len(tokens))
    fmt.Printf("Unique token types: %d\n", len(tokenCounts))
    fmt.Printf("Lexical errors: %d\n", len(errors))
    
    // Show distribution of token types
    fmt.Printf("\nToken Distribution:\n")
    for tokenType, count := range tokenCounts {
        percentage := float64(count) / float64(len(tokens)) * 100
        fmt.Printf("  %-15s: %4d (%5.1f%%)\n", tokenType, count, percentage)
    }
    
    // List any errors found
    if len(errors) > 0 {
        fmt.Printf("\nErrors Found:\n")
        for i, err := range errors {
            fmt.Printf("  [%d] %s\n", i+1, err.Error())
        }
    }
}
```

### Syntax Highlighter

This example shows how to build a basic syntax highlighter:

```go
type SyntaxToken struct {
    Text  string
    Type  string
    Start int
    End   int
}

func generateSyntaxHighlighting(source string) []SyntaxToken {
    lexer := golexer.NewLexer(source)
    var highlighted []SyntaxToken
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        // Map lexer tokens to syntax highlighting categories
        highlightType := mapTokenToHighlightType(tok.Type)
        
        highlighted = append(highlighted, SyntaxToken{
            Text:  tok.Literal,
            Type:  highlightType,
            Start: calculatePosition(tok.Line, tok.Column),
            End:   calculatePosition(tok.Line, tok.Column + len(tok.Literal)),
        })
    }
    
    return highlighted
}

func mapTokenToHighlightType(tokenType golexer.TokenType) string {
    switch tokenType {
    case golexer.LET, golexer.CONST, golexer.IF, golexer.ELSE, golexer.WHILE, golexer.FOR:
        return "keyword"
    case golexer.STRING, golexer.CHAR:
        return "string"
    case golexer.NUMBER:
        return "number"
    case golexer.IDENT:
        return "identifier"
    default:
        return "operator"
    }
}
```

## Error Handling and Recovery

GoLexer provides comprehensive error detection while continuing to process the input. This "error recovery" approach finds all problems in one pass rather than stopping at the first error.

### Error Categories

| Category | Example Input | Error Message | Explanation |
|----------|---------------|---------------|-------------|
| **Invalid Numbers** | `123abc` | `invalid number: numbers cannot be followed by letters` | Numbers can't have letters immediately after |
| **Bad Hex Numbers** | `0xGHI` | `invalid hexadecimal number: contains non-hex characters` | Hex numbers only allow 0-9, A-F |
| **Bad Binary** | `0b123` | `invalid binary number: contains non-binary characters` | Binary numbers only allow 0 and 1 |
| **Bad Octal** | `0o89` | `invalid octal number: contains non-octal characters` | Octal numbers only allow 0-7 |
| **Bad Scientific** | `1e` | `invalid scientific notation: exponent must contain digits` | Scientific notation needs digits after e/E |
| **Unterminated Strings** | `"hello` | `unterminated string literal` | String started but never closed |
| **Bad Escapes** | `"test\q"` | `unknown escape sequence '\q'` | \q is not a valid escape sequence |
| **Unterminated Comments** | `/* comment` | `unterminated block comment` | Block comment started but never closed |
| **Invalid Characters** | `@` (without config) | `unexpected character '@' (Unicode: U+0040)` | Character not recognized by lexer |
| **Helpful Suggestions** | `&` | `unexpected character '&' - did you mean '&&'?` | Common mistake with logical operators |

### Error Recovery Strategy

The lexer continues processing after encountering errors, which helps find all problems at once:

```go
problemSource := `
let x = 123abc;      // Error: invalid number
let y = "valid";     // ✅ Processes correctly after error
let z = 0xGHI;       // Error: invalid hex
let a = 42;          // ✅ Still processes correctly  
`

lexer := golexer.NewLexer(problemSource)
tokens, errors := lexer.TokenizeAll()

// Results:
// - tokens contains: LET, IDENT(x), =, ILLEGAL, ;, LET, IDENT(y), =, STRING(valid), ;, etc.
// - errors contains detailed information about both the "123abc" and "0xGHI" problems
// - Processing never stops, so you get complete analysis
```

**Why this matters**: Instead of fixing one error and re-running, you see all issues immediately. This is especially valuable in IDEs and development tools.

### Configuration Error Resilience

The configuration loading system never crashes your application:

```go
// These all handle gracefully with warning messages:
lexer1 := golexer.NewLexerWithConfig(source, "missing.json")          // File not found
lexer2 := golexer.NewLexerWithConfig(source, "invalid-syntax.json")   // Malformed JSON
lexer3 := golexer.NewLexerWithConfig(source, "no-permission.json")    // Access denied

// All show appropriate warnings and continue with built-in defaults
```

**Typical warning output**:
```
Warning: failed to load config file 'missing.json': no such file or directory
Continuing with default configuration...
```

This approach ensures your application keeps working even in deployment environments with configuration issues.

## Testing and Validation

### Running the Test Suite

```bash
# Clone the repository
git clone https://github.com/codetesla51/golexer.git
cd golexer

# Run unit tests for individual components
go test ./golexer

# Run comprehensive integration test with real-world examples
go run cmd/main.go test.lang
```

### Expected Test Results

The comprehensive test file (`test.lang`) validates every lexer feature with real-world code patterns:

```
Analyzing file: test.lang
File size: 15000+ bytes

=== Token-by-token processing ===
Type: LET             Literal: 'let'           Line:  1 Column:  1
Type: IDENT           Literal: 'a'             Line:  1 Column:  5
Type: =               Literal: '='             Line:  1 Column:  7
[... processing continues for 1700+ more tokens ...]

=== Batch processing ===
Total tokens: 1700+
Total errors: 0

Token distribution:
  IDENT          : 245    // Variable and function names
  NUMBER         : 189    // All number formats
  STRING         : 67     // String literals  
  =              : 134    // Assignment operators
  +              : 45     // Addition operators
  [... complete statistics for all token types ...]

=== Syntax validation ===
✓ No lexical errors found - file is syntactically valid at lexical level

=== Summary ===
File: test.lang
Lines processed: 400+
Tokens generated: 1700+
Unique token types: 45+
Lexical errors: 0
Status: ✓ PASSED
```

### What the Test Validates

The comprehensive test suite covers every aspect of the lexer:

**Number Format Testing**:
- ✅ All decimal formats: integers, floats, scientific notation
- ✅ All hex formats: 0xFF, 0x1a2b, case variations
- ✅ All binary formats: 0b1010, 0B1111
- ✅ All octal formats: 0o777, 0123 (legacy)
- ✅ Error conditions: invalid digits, malformed numbers

**String Processing Testing**:  
- ✅ Regular strings with all escape sequences
- ✅ Raw backtick strings with literal content
- ✅ Character literals and special characters
- ✅ Error conditions: unterminated strings, invalid escapes

**Language Feature Testing**:
- ✅ Complete operator and punctuation coverage
- ✅ All built-in keywords and identifiers
- ✅ Configuration system with custom extensions
- ✅ Comment processing (line and block comments)
- ✅ Unicode identifiers and international characters

**Real-world Pattern Testing**:
- ✅ Complex nested structures (arrays, objects, functions)
- ✅ Mixed expressions with multiple operators
- ✅ Realistic code patterns found in actual programs
- ✅ Edge cases and boundary conditions

**Error Handling Testing**:
- ✅ Position tracking accuracy
- ✅ Error recovery and continued processing
- ✅ Helpful error messages with suggestions

## Real-world Applications

### Building a Compiler Frontend

```go
type Compiler struct {
    lexer  *golexer.Lexer
    errors []error
}

func (c *Compiler) CompileFile(filename string) error {
    // Read source code
    source, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }
    
    // Tokenize the source
    c.lexer = golexer.NewLexer(string(source))
    tokens, lexErrors := c.lexer.TokenizeAll()
    
    // Check for lexical errors
    if len(lexErrors) > 0 {
        for _, lexErr := range lexErrors {
            c.errors = append(c.errors, lexErr)
        }
        return fmt.Errorf("compilation failed with %d lexical errors", len(lexErrors))
    }
    
    // Pass tokens to parser for syntax analysis
    parser := NewParser(tokens)
    ast, parseErrors := parser.Parse()
    
    // Continue with semantic analysis, code generation, etc.
    return c.generateCode(ast)
}
```

### Building an IDE Language Server

Language servers provide IDE features like syntax highlighting, error detection, and autocompletion:

```go
func provideSyntaxHighlighting(document string) []HighlightRange {
    lexer := golexer.NewLexer(document)
    var highlights []HighlightRange
    
    for {
        tok := lexer.NextToken()
        if tok.Type == golexer.EOF {
            break
        }
        
        highlights = append(highlights, HighlightRange{
            Start: Position{Line: tok.Line - 1, Column: tok.Column - 1}, // Convert to 0-based
            End:   Position{Line: tok.Line - 1, Column: tok.Column - 1 + len(tok.Literal)},
            TokenType: mapToEditorTokenType(tok.Type),
        })
    }
    
    return highlights
}

func provideErrorDiagnostics(document string) []Diagnostic {
    lexer := golexer.NewLexer(document)
    _, errors := lexer.TokenizeAll()
    
    var diagnostics []Diagnostic
    for _, err := range errors {
        diagnostics = append(diagnostics, Diagnostic{
            Range: Range{
                Start: Position{Line: err.Line - 1, Column: err.Column - 1},
                End:   Position{Line: err.Line - 1, Column: err.Column},
            },
            Severity: Error,
            Message:  err.Message,
        })
    }
    
    return diagnostics
}
```

### Building a Custom Configuration Language

Many applications need custom configuration languages. Here's how to build one:

```go
// First, create config for your custom language syntax
// custom-lang.json:
{
  "additionalKeywords": {
    "server": "SERVER",
    "database": "DATABASE", 
    "cache": "CACHE",
    "enable": "ENABLE",
    "disable": "DISABLE"
  },
  "additionalOperators": {
    "=>": "ARROW"
  }
}

func
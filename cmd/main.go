/*
GoLexer - A Comprehensive Lexical Analyzer for Go
Author: Uthman Dev
GitHub: https://github.com/codetesla51/golexer
License: MIT

Main Command-Line Interface
This is the main entry point for the GoLexer CLI tool. It demonstrates
three key usage patterns of the lexer library:
1. Token-by-token processing with detailed output
2. Batch processing with statistical analysis
3. Basic syntax validation and error reporting

Usage: go run main.go <filename>
Example: go run main.go test.lang

The tool reads any text file and performs lexical analysis, showing
tokens, their positions, and any lexical errors encountered.
*/

package main

import (
	"fmt"
	"github.com/codetesla51/golexer/golexer"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s test.lang\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Analyzing file: %s\n", filename)
	fmt.Printf("File size: %d bytes\n\n", len(content))

	// Example 1: Token-by-token processing
	fmt.Println("=== Token-by-token processing ===")
	lexer := golexer.NewLexer(string(content))

	tokenCount := 0
	for {
		tok := lexer.NextToken()
		if tok.Type == golexer.EOF {
			break
		}

		fmt.Printf("Type: %-15s Literal: %-15s Line: %2d Column: %2d\n",
			string(tok.Type), fmt.Sprintf("'%s'", tok.Literal), tok.Line, tok.Column)
		tokenCount++
	}

	// Print any errors from first pass
	if lexer.HasErrors() {
		fmt.Println("\nLexical Errors (First Pass):")
		for _, err := range lexer.GetErrors() {
			fmt.Printf("  %s\n", err.Error())
		}
	}

	fmt.Printf("\nProcessed: %d tokens\n", tokenCount)
	fmt.Println("\n=== Batch processing ===")

	// Example 2: Batch processing
	lexer2 := golexer.NewLexer(string(content))
	tokens, errors := lexer2.TokenizeAll()

	fmt.Printf("Total tokens: %d\n", len(tokens))
	fmt.Printf("Total errors: %d\n", len(errors))

	// Count token types
	tokenCounts := make(map[golexer.TokenType]int)
	for _, token := range tokens {
		tokenCounts[token.Type]++
	}

	fmt.Println("\nToken distribution:")
	for tokenType, count := range tokenCounts {
		fmt.Printf("  %-15s: %d\n", tokenType, count)
	}

	// Example 3: Enhanced syntax validation
	fmt.Println("\n=== Syntax validation ===")
	if len(errors) == 0 {
		fmt.Println("✓ No lexical errors found - file is syntactically valid at lexical level")
	} else {
		fmt.Printf("✗ Found %d lexical error(s):\n", len(errors))
		for i, err := range errors {
			fmt.Printf("  [%d] %s\n", i+1, err.Error())
		}
	}

	// Summary statistics
	fmt.Println("\n=== Summary ===")
	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Lines processed: %d\n", getLastLine(tokens))
	fmt.Printf("Tokens generated: %d\n", len(tokens))
	fmt.Printf("Unique token types: %d\n", len(tokenCounts))
	fmt.Printf("Lexical errors: %d\n", len(errors))

	if len(errors) == 0 {
		fmt.Println("Status: ✓ PASSED")
	} else {
		fmt.Println("Status: ✗ FAILED")
		os.Exit(1) // Exit with error code if lexical errors found
	}
}

// Helper function to get the last line number from tokens
func getLastLine(tokens []golexer.Token) int {
	if len(tokens) == 0 {
		return 1
	}
	maxLine := 1
	for _, token := range tokens {
		if token.Line > maxLine {
			maxLine = token.Line
		}
	}
	return maxLine
}

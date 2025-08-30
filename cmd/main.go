package main

import (
	"fmt"
	"os"
	"github.com/codetesla51/golexer/golexer"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Example 1: Token-by-token processing
	fmt.Println("=== Token-by-token processing ===")
	lexer := golexer.NewLexer(string(content))

	for {
		tok := lexer.NextToken()
		if tok.Type == golexer.EOF {
			break
		}

		fmt.Printf("Type: %-15s Literal: %-15s Line: %2d Column: %2d\n",
			string(tok.Type), fmt.Sprintf("'%s'", tok.Literal), tok.Line, tok.Column)
	}

	// Print any errors
	if lexer.HasErrors() {
		fmt.Println("\nLexical Errors:")
		for _, err := range lexer.GetErrors() {
			fmt.Printf("  %s\n", err.Error())
		}
	}

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
		fmt.Printf("  %s: %d\n", tokenType, count)
	}
	
	// Example 3: Simple syntax validation
	fmt.Println("\n=== Syntax validation ===")
	if len(errors) == 0 {
		fmt.Println("✓ No lexical errors found")
	} else {
		fmt.Printf("✗ Found %d lexical errors\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  %s\n", err.Error())
		}
	}
}
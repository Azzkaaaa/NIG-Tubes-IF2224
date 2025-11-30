package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	iox "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/common"
	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/lexer"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/parser"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/semantic"
)

func main() {
	rules := flag.String("rules", "config/tokenizer_m3.json", "path ke DFA JSON")
	in := flag.String("input", "", "path file sumber")
	flag.Parse()

	if *in == "" {
		fmt.Fprintln(os.Stderr, "missing --input <file>")
		os.Exit(2)
	}

	// Load DFA rules
	d, err := lexer.LoadJSON(*rules)
	if err != nil {
		log.Fatal(err)
	}

	// Read source file
	rr, err := iox.NewRuneReaderFromFile(*in)
	if err != nil {
		log.Fatal(err)
	}

	// Lexical analysis
	tokens, errs := lexer.New(d, rr).ScanAll()

	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}

	if len(errs) > 0 {
		os.Exit(1)
	}

	// Filter out comments
	tokens = slices.Collect(func(yield func(dt.Token) bool) {
		for _, token := range tokens {
			if token.Type != dt.COMMENT {
				if !yield(token) {
					return
				}
			}
		}
	})

	// Syntax analysis
	parseTree, err := parser.New(tokens).Parse()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	if parseTree == nil {
		fmt.Fprintln(os.Stderr, "Parse tree is nil")
		os.Exit(1)
	}

	// Semantic analysis
	analyzer := semantic.New(parseTree)
	tab, atab, btab, strtab, dst, err := analyzer.Analyze()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Semantic error: %v\n", err)
		os.Exit(1)
	}

	// Display decorated syntax tree
	if dst != nil {
		fmt.Println(dst.StringWithSymbols(tab, atab, btab, strtab))
	}
	fmt.Println()

	// Display symbol tables
	fmt.Println("=== Symbol Table (TAB) ===")
	fmt.Println(tab.String())
	fmt.Println()

	fmt.Println("=== Array Table (ATAB) ===")
	fmt.Println(atab.String())
	fmt.Println()

	fmt.Println("=== Block Table (BTAB) ===")
	fmt.Println(btab.String())
	fmt.Println()

	fmt.Println("=== String Table (STRTAB) ===")
	fmt.Println(strtab.String())
}

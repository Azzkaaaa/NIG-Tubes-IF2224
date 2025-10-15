package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/dfa"
	iox "github.com/Azzkaaaa/NIG-Tubes-IF2224/src/io"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/lexer"
)

func main() {
	rules := flag.String("rules", "src/rules/tokenizer.json", "path to DFA json")
	in := flag.String("input", "", "path to source code")
	flag.Parse()
	if *in == "" {
		fmt.Fprintln(os.Stderr, "missing --input <file>")
		os.Exit(2)
	}

	d, err := dfa.LoadJSON(*rules)
	if err != nil {
		panic(err)
	}

	rr, err := iox.NewRuneReaderFromFile(*in)
	if err != nil {
		panic(err)
	}

	lx := lexer.New(d, rr)
	tokens, errs := lx.ScanAll()

	iox.PrintTokens(tokens)
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}
}

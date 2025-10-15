package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/dfa"
	iox "github.com/Azzkaaaa/NIG-Tubes-IF2224/src/io"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/lexer"
)

func main() {
	rules := flag.String("rules", "src/rules/tokenizer.json", "path ke DFA JSON")
	in := flag.String("input", "", "path file sumber")
	out := flag.String("out", "", "opsional: file output token")
	flag.Parse()

	if *in == "" {
		fmt.Fprintln(os.Stderr, "missing --input <file>")
		os.Exit(2)
	}

	d, err := dfa.LoadJSON(*rules)
	if err != nil {
		log.Fatal(err)
	}

	rr, err := iox.NewRuneReaderFromFile(*in)
	if err != nil {
		log.Fatal(err)
	}

	tokens, errs := lexer.New(d, rr).ScanAll()

	iox.PrintTokens(tokens)

	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}

	if *out != "" {
		if err := iox.WriteTokensToFile(*out, tokens); err != nil {
			log.Fatal(err)
		}
	}
}

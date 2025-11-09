package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	iox "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/common"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/lexer"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/parser"
)

func main() {
	rules := flag.String("rules", "config/tokenizer.json", "path ke DFA JSON")
	in := flag.String("input", "", "path file sumber")
	out := flag.String("out", "", "opsional: file output parser")
	flag.Parse()

	if *in == "" {
		fmt.Fprintln(os.Stderr, "missing --input <file>")
		os.Exit(2)
	}

	d, err := lexer.LoadJSON(*rules)
	if err != nil {
		log.Fatal(err)
	}

	rr, err := iox.NewRuneReaderFromFile(*in)
	if err != nil {
		log.Fatal(err)
	}

	tokens, errs := lexer.New(d, rr).ScanAll()

	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}

	if *out != "" {
		// if err := WriteTokensAndErrorsToFile(*out, tokens, errs); err != nil {
		// 	log.Fatal(err)
		// }
	}

	if len(errs) > 0 {
		os.Exit(1)
	}

	parseTree, err := parser.New(tokens).Parse()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if parseTree != nil {
		fmt.Println(parseTree.String())
	}

}

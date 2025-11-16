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
)

func main() {
	rules := flag.String("rules", "config/tokenizer.json", "path ke DFA JSON")
	in := flag.String("input", "", "path file sumber")
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

	if len(errs) > 0 {
		os.Exit(1)
	}

	tokens = slices.Collect(func(yield func(dt.Token) bool) {
		for _, token := range tokens {
			if token.Type != dt.COMMENT {
				if !yield(token) {
					return
				}
			}
		}
	})

	parseTree, err := parser.New(tokens).Parse()

	if err != nil {
		fmt.Printf("%v", err)
	}

	if parseTree != nil {
		fmt.Println(parseTree.String())
	}

}

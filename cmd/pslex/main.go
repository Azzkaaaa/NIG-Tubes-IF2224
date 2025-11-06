package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	iox "github.com/Azzkaaaa/NIG-Tubes-IF2224/common"
	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/datatype"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/lexer"
)

func main() {
	rules := flag.String("rules", "config/tokenizer.json", "path ke DFA JSON")
	in := flag.String("input", "", "path file sumber")
	out := flag.String("out", "", "opsional: file output token")
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

	PrintTokens(tokens)

	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}

	if *out != "" {
		if err := WriteTokensAndErrorsToFile(*out, tokens, errs); err != nil {
			log.Fatal(err)
		}
	}
}

func PrintTokens(tokens []dt.Token) {
	for _, t := range tokens {
		fmt.Printf("%s(%s)\n", t.Type.String(), t.Lexeme)
	}
}

func PrintErrorslexe(errs []error) {
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}
}

func WriteTokensAndErrorsToFile(path string, tokens []dt.Token, errs []error) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, t := range tokens {
		if _, err := fmt.Fprintf(w, "%s(%s)\n", t.Type.String(), t.Lexeme); err != nil {
			return err
		}
	}
	for _, e := range errs {
		if _, err := fmt.Fprintln(w, e); err != nil {
			return err
		}
	}
	return w.Flush()
}

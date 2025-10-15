package io

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/datatype"
)

func PrintTokens(tokens []datatype.Token) {
	for _, t := range tokens {
		fmt.Printf("%s(%s)\n", t.Type.String(), t.Lexeme)
	}
}

func PrintErrors(errs []error) {
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}
}

func WriteTokensToFile(path string, tokens []datatype.Token) error {
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
	return w.Flush()
}

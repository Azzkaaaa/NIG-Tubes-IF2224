package lexer

import (
	"fmt"
	"strings"

	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/datatype"
	"github.com/Azzkaaaa/NIG-Tubes-IF2224/src/dfa"
	iox "github.com/Azzkaaaa/NIG-Tubes-IF2224/src/io"
)

type Lexer struct {
	d *dfa.DFA
	r *iox.RuneReader
}

func New(d *dfa.DFA, r *iox.RuneReader) *Lexer {
	return &Lexer{d: d, r: r}
}

func isWS(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' || r == 0
}

func (lx *Lexer) ScanAll() ([]datatype.Token, []error) {
	var toks []datatype.Token
	var errs []error

	for !lx.r.EOF() {
		startOff := lx.r.Offset()
		startLine, startCol := lx.r.Pos()
		startSnap := lx.r.Snapshot()

		lx.d.Reset()

		lastOkOff := -1
		lastOkLabel := ""
		lastSnap := startSnap

		for !lx.r.EOF() {
			ch := lx.r.Peek()
			if !lx.d.Advance(ch) {
				break
			}
			lx.r.Read()

			if lab, ok := lx.d.IsFinal(lx.d.State()); ok {
				lastOkOff = lx.r.Offset()
				lastOkLabel = lab
				lastSnap = lx.r.Snapshot()
			}
		}

		if lastOkOff != -1 {
			lx.r.Restore(lastSnap)
			lex := lx.r.Slice(startOff, lastOkOff)

			tt := mapLabel(lastOkLabel)
			// if tt == datatype.IDENTIFIER {
			// 	if _, ok := datatype.Keywords[strings.ToLower(lex)]; ok {
			// 		tt = datatype.KEYWORD
			// 	}
			// }

			lex = strings.Trim(lex, " \t\r\n\f")

			switch tt {
			case datatype.STRING_LITERAL, datatype.CHAR_LITERAL, datatype.COMMENT:
			default:
				lex = strings.ToLower(lex)
			}

			if lex != "" || tt == datatype.STRING_LITERAL || tt == datatype.CHAR_LITERAL || tt == datatype.COMMENT {
				toks = append(toks, datatype.Token{
					Type:   tt,
					Lexeme: lex,
					Line:   startLine,
					Col:    startCol,
				})
			}
			continue
		}

		if ch, ok := lx.r.Read(); ok {
			if !isWS(ch) {
				errs = append(errs, fmt.Errorf("unrecognized %q at %d:%d", ch, startLine, startCol))
			}
		}
	}

	return toks, errs
}

func mapLabel(label string) datatype.TokenType {
	switch label {
	case "KEYWORD":
		return datatype.KEYWORD
	case "IDENTIFIER":
		return datatype.IDENTIFIER
	case "NUMBER":
		return datatype.NUMBER
	case "CHAR_LITERAL":
		return datatype.CHAR_LITERAL
	case "STRING_LITERAL":
		return datatype.STRING_LITERAL
	case "ARITHMETIC_OPERATOR":
		return datatype.ARITHMETIC_OPERATOR
	case "RELATIONAL_OPERATOR":
		return datatype.RELATIONAL_OPERATOR
	case "LOGICAL_OPERATOR":
		return datatype.LOGICAL_OPERATOR
	case "ASSIGN_OPERATOR":
		return datatype.ASSIGN_OPERATOR
	case "RANGE_OPERATOR":
		return datatype.RANGE_OPERATOR
	case "SEMICOLON":
		return datatype.SEMICOLON
	case "COMMA":
		return datatype.COMMA
	case "COLON":
		return datatype.COLON
	case "DOT":
		return datatype.DOT
	case "LPARENTHESIS":
		return datatype.LPARENTHESIS
	case "RPARENTHESIS":
		return datatype.RPARENTHESIS
	case "LBRACKET":
		return datatype.LBRACKET
	case "RBRACKET":
		return datatype.RBRACKET
	case "COMMENT":
		return datatype.COMMENT
	default:
		return datatype.IDENTIFIER
	}
}

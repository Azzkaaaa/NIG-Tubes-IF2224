package semantic

import (
	"errors"
	"math"
	"strconv"
	"strings"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeToken(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parsetree.RootType != dt.TOKEN_NODE {
		return nil, semanticType{}, errors.New("parsetree root is not token node")
	}

	switch parsetree.TokenValue.Type {
	case dt.CHAR_LITERAL:
		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_CHAR_LITERAL,
			Data:     int(parsetree.TokenValue.Lexeme[0]),
		}, semanticType{StaticType: dt.TAB_ENTRY_CHAR}, nil

	case dt.STRING_LITERAL:
		value := parsetree.TokenValue.Lexeme[1 : len(parsetree.TokenValue.Lexeme)-1]

		atabIndex, _ := a.atab.FindArray(dt.AtabEntry{
			IndexType:   dt.TAB_ENTRY_INTEGER,
			ElementType: dt.TAB_ENTRY_CHAR,
			LowBound:    0,
			HighBound:   len(value),
		})

		if atabIndex == -1 {
			atabIndex = len(a.atab)
			a.atab = append(a.atab, dt.AtabEntry{
				IndexType:   dt.TAB_ENTRY_INTEGER,
				ElementType: dt.TAB_ENTRY_CHAR,
				LowBound:    0,
				HighBound:   len(value),
				ElementSize: 1,
				TotalSize:   len(value),
			})
		}

		stridx, _ := a.strtab.FindString(value)

		if stridx == -1 {
			stridx = len(a.strtab)
			a.strtab = append(a.strtab, dt.StrTabEntry{
				Length:    len(parsetree.TokenValue.Lexeme),
				String:    parsetree.TokenValue.Lexeme,
				Reference: atabIndex,
			})
		}

		return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_STR_LITERAL,
				Data:     stridx,
			}, semanticType{
				StaticType: dt.TAB_ENTRY_ARRAY,
				Reference:  atabIndex,
			}, nil

	case dt.NUMBER:
		if strings.ContainsAny(parsetree.TokenValue.Lexeme, "e") {
			parts := strings.Split(parsetree.TokenValue.Lexeme, "e")

			if len(parts) != 2 {
				return nil, semanticType{}, errors.New("unexpected error reading float")
			}

			significant, err := strconv.ParseFloat(parts[0], 64)

			if err != nil {
				return nil, semanticType{}, errors.New("unexpected error reading float")
			}

			exponent, err := strconv.ParseFloat(parts[1], 64)

			if err != nil {
				return nil, semanticType{}, errors.New("unexpected error reading float")
			}

			rawValue := significant * math.Pow(10, exponent)
			var data int

			if strconv.IntSize == 32 {
				data = int(math.Float32bits(float32(rawValue)))
			} else {
				data = int(math.Float64bits(float64(rawValue)))
			}

			return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_REAL_LITERAL,
				Data:     data,
			}, semanticType{StaticType: dt.TAB_ENTRY_REAL}, nil
		} else if strings.ContainsAny(parsetree.TokenValue.Lexeme, ".") {
			val, err := strconv.ParseFloat(parsetree.TokenValue.Lexeme, strconv.IntSize)

			if err != nil {
				return nil, semanticType{}, errors.New("unexpected error reading float")
			}

			var data int

			if strconv.IntSize == 32 {
				data = int(math.Float32bits(float32(val)))
			} else {
				data = int(math.Float64bits(val))
			}

			return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_REAL_LITERAL,
				Data:     data,
			}, semanticType{StaticType: dt.TAB_ENTRY_REAL}, nil
		} else {
			val, err := strconv.ParseInt(parsetree.TokenValue.Lexeme, 10, 64)

			if err != nil {
				return nil, semanticType{}, errors.New("unexpected error reading integer")
			}

			return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_INT_LITERAL,
				Data:     int(val),
			}, semanticType{StaticType: dt.TAB_ENTRY_INTEGER}, nil
		}

	case dt.KEYWORD:
		switch parsetree.TokenValue.Lexeme {
		case "true":
			return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_BOOL_LITERAL,
				Data:     1,
			}, semanticType{StaticType: dt.TAB_ENTRY_BOOLEAN}, nil
		case "false":
			return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_BOOL_LITERAL,
				Data:     0,
			}, semanticType{StaticType: dt.TAB_ENTRY_BOOLEAN}, nil
		default:
			return nil, semanticType{}, errors.New("unexpected keyword")
		}

	case dt.IDENTIFIER:
		index, tabEntry := a.tab.FindIdentifier(parsetree.TokenValue.Lexeme, a.root)

		if tabEntry == nil {
			return nil, semanticType{}, errors.New("undeclared variable")
		}

		if tabEntry.Object != dt.TAB_ENTRY_VAR {
			return nil, semanticType{}, errors.New("identifier does not reference variable")
		}

		return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_VARIABLE,
				Data:     index,
			}, semanticType{
				StaticType: tabEntry.Type,
				Reference:  tabEntry.Reference,
			}, nil
	}

	return nil, semanticType{}, errors.New("unexpected token type")
}

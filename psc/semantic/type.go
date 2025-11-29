package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeType(parsetree *dt.ParseTree) (int, dt.TabEntry, error) {
	if parsetree.RootType != dt.TYPE_NODE {
		return -1, dt.TabEntry{}, errors.New("expected type")
	}

	child := parsetree.Children[0]

	switch child.RootType {
	case dt.TOKEN_NODE:
		if child.TokenValue.Type == dt.IDENTIFIER {
			index, tabEntry := a.tab.FindIdentifier(child.TokenValue.Lexeme, a.root)

			if tabEntry == nil {
				return -1, dt.TabEntry{}, errors.New("undefined type")
			}

			if tabEntry.Object != dt.TAB_ENTRY_TYPE {
				return -1, dt.TabEntry{}, errors.New("identifier is not a type")
			}

			return index, *tabEntry, nil
		} else {
			switch child.TokenValue.Lexeme {
			case "integer":
				return -1, dt.TabEntry{Type: dt.TAB_ENTRY_INTEGER}, nil
			case "real":
				return -1, dt.TabEntry{Type: dt.TAB_ENTRY_REAL}, nil
			case "boolean":
				return -1, dt.TabEntry{Type: dt.TAB_ENTRY_BOOLEAN}, nil
			case "char":
				return -1, dt.TabEntry{Type: dt.TAB_ENTRY_CHAR}, nil
			default:
				return -1, dt.TabEntry{}, errors.New("unrecognized type declaration")
			}
		}
	case dt.ARRAY_TYPE_NODE:
		return a.analyzeArrayType(&child)
	default:
		return -1, dt.TabEntry{}, errors.New("unrecognized type declaration")
	}
}

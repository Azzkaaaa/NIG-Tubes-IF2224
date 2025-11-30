package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeConstDeclaration(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.CONST_DECLARATION_NODE {
		return nil, errors.New("expected const declaration")
	}

	identifier := parsetree.Children[0].TokenValue.Lexeme
	_, prev := a.tab.FindIdentifier(identifier, a.root)

	if prev != nil {
		if prev.Level == a.depth {
			return nil, errors.New("identifier redefined in the same scope")
		}
	}

	val, valtype, err := a.analyzeToken(&parsetree.Children[2])

	if err != nil {
		return nil, err
	}

	tabEntry := dt.TabEntry{
		Identifier: identifier,
		Link:       a.root,
		Object:     dt.TAB_ENTRY_CONST,
		Type:       valtype.StaticType,
		Reference:  valtype.Reference,
		Level:      a.depth,
		Data:       val.Data,
	}

	a.root = len(a.tab)
	a.tab = append(a.tab, tabEntry)

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_CONST,
		Data:     a.root,
	}, nil
}

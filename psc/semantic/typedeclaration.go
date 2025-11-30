package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeTypeDeclaration(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.TYPE_DECLARATION_NODE {
		return nil, errors.New("expected type declaration")
	}

	identifier := parsetree.Children[0].TokenValue.Lexeme
	_, prev := a.tab.FindIdentifier(identifier, a.root)

	if prev != nil {
		if prev.Level == a.depth {
			return nil, errors.New("identifier redefined in the same scope")
		}
	}

	var tabIndex int
	var tabEntry dt.TabEntry
	var err error
	if parsetree.Children[2].RootType == dt.RECORD_TYPE_NODE {
		_, _, err = a.analyzeRecordType(&parsetree.Children[2], identifier)
		if err != nil {
			return nil, err
		}

		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_TYPE,
			Data:     a.root,
		}, nil
	} else {
		tabIndex, tabEntry, err = a.analyzeType(&parsetree.Children[2])
		if tabIndex != -1 {
			tabEntry = dt.TabEntry{
				Type:      dt.TAB_ENTRY_ALIAS,
				Reference: tabIndex,
			}
		}

		tabEntry.Identifier = identifier
		tabEntry.Object = dt.TAB_ENTRY_TYPE
		tabEntry.Level = a.depth
		tabEntry.Link = a.root
		if err != nil {
			return nil, err
		}

		a.root = len(a.tab)
		a.tab = append(a.tab, tabEntry)

		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_TYPE,
			Data:     a.root,
		}, nil
	}

}

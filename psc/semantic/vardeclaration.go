package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeVarDeclaration(parsetree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.VAR_DECLARATION_NODE {
		return nil, errors.New("expected var declaration")
	}

	identifiers, err := a.analyzeIdentifierList(&parsetree.Children[0])

	if err != nil {
		return nil, err
	}

	tabIndex, tabEntry, err := a.analyzeType(&parsetree.Children[2])

	if err != nil {
		return nil, err
	}

	declarations := make([]dt.DecoratedSyntaxTree, len(identifiers))

	for i, identifier := range identifiers {
		_, check := a.tab.FindIdentifier(identifier, a.root)
		if check != nil {
			if check.Level == a.depth {
				// Get token from identifier list
				token := parsetree.Children[0].Children[i*2].TokenValue
				return nil, a.newRedeclarationError(identifier, token)
			}
		}

		if tabIndex != -1 {
			tabEntry = dt.TabEntry{
				Type:      dt.TAB_ENTRY_ALIAS,
				Reference: tabIndex,
			}
		}

		tabEntry.Identifier = identifier
		tabEntry.Link = a.root
		tabEntry.Object = dt.TAB_ENTRY_VAR
		tabEntry.Level = a.depth
		tabEntry.Data = a.stackSize

		a.stackSize += a.getTypeSize(semanticType{
			StaticType: tabEntry.Type,
			Reference:  tabEntry.Reference,
		})

		a.root = len(a.tab)
		a.tab = append(a.tab, tabEntry)

		declarations[i] = dt.DecoratedSyntaxTree{
			SelfType: dt.DST_VARIABLE,
			Data:     a.root,
		}
	}

	return declarations, nil
}

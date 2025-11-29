package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeFormalParameterList(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.FORMAL_PARAMETER_LIST_NODE {
		return nil, errors.New("expected formal parameter list")
	}

	identifierListNodes := parsetree.Children[1:len(parsetree.Children):4]
	typesList := parsetree.Children[3:len(parsetree.Children):4]

	parameters := make([]dt.DecoratedSyntaxTree, 0)

	for i, identifierListNode := range identifierListNodes {
		identifierList, err := a.analyzeIdentifierList(&identifierListNode)

		if err != nil {
			return nil, err
		}

		tabIndex, tabEntry, err := a.analyzeType(&typesList[i])

		if err != nil {
			return nil, err
		}

		for _, identifier := range identifierList {
			_, check := a.tab.FindIdentifier(identifier, a.root)
			if check != nil {
				if check.Level == a.depth {
					return nil, errors.New("cannot redeclare identifier in the same scope")
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
			tabEntry.Object = dt.TAB_ENTRY_PARAM
			tabEntry.Level = a.depth

			a.root = len(a.tab)
			a.tab = append(a.tab, tabEntry)

			parameters = append(parameters, dt.DecoratedSyntaxTree{
				Property: dt.DST_PARAMETER,
				SelfType: dt.DST_VARIABLE,
				Data:     a.root,
			})
		}
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_VARIABLE_DECLARATIONS,
		Children: parameters,
	}, nil
}

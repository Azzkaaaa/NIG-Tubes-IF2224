package semantic

import (
	"errors"
	"fmt"
	"strconv"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeFormalParameterList(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.FORMAL_PARAMETER_LIST_NODE {
		return nil, errors.New("expected formal parameter list")
	}

	parameters := make([]dt.DecoratedSyntaxTree, 0)
	isRef := false

	i := 0
	for i < len(parsetree.Children) {
		child := &parsetree.Children[i]
		if child.RootType == dt.TOKEN_NODE && child.TokenValue.Type == dt.KEYWORD && child.TokenValue.Lexeme == "variabel" {
			isRef = true
			i++
			continue
		}

		if child.RootType == dt.TOKEN_NODE {
			if child.TokenValue.Type == dt.SEMICOLON {
				isRef = false
			}
			i++
			continue
		}

		if child.RootType == dt.IDENTIFIER_LIST_NODE {
			if i+2 >= len(parsetree.Children) {
				return nil, errors.New("malformed parameter list: missing type for identifier list")
			}

			identifierListNode := child
			typeNode := &parsetree.Children[i+2]

			identifierList, err := a.analyzeIdentifierList(identifierListNode)
			if err != nil {
				return nil, err
			}

			tabIndex, tabEntry, err := a.analyzeType(typeNode)
			if err != nil {
				return nil, err
			}

			for _, identifier := range identifierList {
				_, check := a.tab.FindIdentifier(identifier, a.root)
				if check != nil && check.Level == a.depth {
					return nil, fmt.Errorf("cannot redeclare identifier '%s' in the same scope", identifier)
				}

				var entry dt.TabEntry
				if tabIndex != -1 {
					entry.Type = dt.TAB_ENTRY_ALIAS
					entry.Reference = tabIndex
				} else {
					entry.Type = tabEntry.Type
					entry.Reference = tabEntry.Reference
				}
				entry.Identifier = identifier

				entry.Link = a.root
				entry.Object = dt.TAB_ENTRY_PARAM
				entry.Level = a.depth
				entry.Normal = !isRef
				entry.Data = a.stackSize

				paramSize := 0
				if isRef {
					paramSize = strconv.IntSize
				} else {
					paramSize = a.getTypeSize(semanticType{
						StaticType: entry.Type,
						Reference:  entry.Reference,
					})
				}
				a.stackSize += paramSize
				a.root = len(a.tab)
				a.tab = append(a.tab, entry)

				parameters = append(parameters, dt.DecoratedSyntaxTree{
					Property: dt.DST_PARAMETER,
					SelfType: dt.DST_VARIABLE,
					Data:     a.root,
				})
			}

			isRef = false
			i += 3
			// If the node is not a recognized token or an identifier list, it's a syntax error.
			return nil, fmt.Errorf("unexpected node type '%s' in formal parameter list at index %d", child.RootType, i)
		}
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_VARIABLE_DECLARATIONS,
		Children: parameters,
	}, nil
}

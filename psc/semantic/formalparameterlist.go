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
	isRef := false // Flag for pass-by-reference

	i := 0
	for i < len(parsetree.Children) {
		child := &parsetree.Children[i]
		fmt.Printf("Analyzing formal parameter list child %d: %+v\n", i, child)
		// Check for 'variabel' keyword (Pascal 'var' for by-reference)
		if child.RootType == dt.TOKEN_NODE && child.TokenValue.Type == dt.KEYWORD && child.TokenValue.Lexeme == "variabel" {
			isRef = true
			i++
			continue
		}

		// Skip structural tokens, but handle semicolon specifically to reset isRef context.
		if child.RootType == dt.TOKEN_NODE {
			// A semicolon separates parameter declarations. Reset 'var' context.
			if child.TokenValue.Type == dt.SEMICOLON {
				isRef = false
			}
			i++
			continue
		}

		// Process a parameter declaration group (e.g., "a, b: integer")
		if child.RootType == dt.IDENTIFIER_LIST_NODE {
			// Expecting structure: IDENTIFIER_LIST, COLON_TOKEN, TYPE_NODE
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

			// Process each identifier in the list
			for _, identifier := range identifierList {
				// Check for redeclaration in the same scope
				_, check := a.tab.FindIdentifier(identifier, a.root)
				if check != nil && check.Level == a.depth {
					return nil, fmt.Errorf("cannot redeclare identifier '%s' in the same scope", identifier)
				}

				// Each identifier needs a distinct TabEntry to avoid aliasing bugs.
				// Initialize a new entry and copy only the necessary type info.
				var entry dt.TabEntry
				if tabIndex != -1 {
					entry.Type = dt.TAB_ENTRY_ALIAS
					entry.Reference = tabIndex
				} else {
					// Copy type information from the analyzed base type
					entry.Type = tabEntry.Type
					entry.Reference = tabEntry.Reference
				}

				// Set common properties for the new symbol table entry
				entry.Identifier = identifier

				// Link should point to the same identifier in outer scope (if exists), or -1
				entry.Link = a.root
				entry.Object = dt.TAB_ENTRY_PARAM
				entry.Level = a.depth
				entry.Normal = !isRef // True for pass-by-value, false for pass-by-reference
				entry.Data = a.stackSize

				fmt.Printf("[PARAM] Adding parameter '%s' with Link=%d, isRef=%v, StackSize=%d\n", identifier, a.root, isRef, a.stackSize)

				paramSize := 0
				if isRef {
					// Pass-by-reference stores an address (pointer size)
					paramSize = strconv.IntSize
				} else {
					// Pass-by-value stores the actual value
					paramSize = a.getTypeSize(semanticType{
						StaticType: entry.Type,
						Reference:  entry.Reference,
					})
				}
				a.stackSize += paramSize

				// Add the new entry to the symbol table
				a.root = len(a.tab)
				a.tab = append(a.tab, entry)

				// Create the corresponding node in the decorated syntax tree
				parameters = append(parameters, dt.DecoratedSyntaxTree{
					Property: dt.DST_PARAMETER,
					SelfType: dt.DST_VARIABLE,
					Data:     a.root,
				})
			}

			// Reset isRef context for the next parameter group after processing this one.
			isRef = false
			// Advance index past the IDENTIFIER_LIST, COLON, and TYPE nodes
			i += 3
		} else {
			// If the node is not a recognized token or an identifier list, it's a syntax error.
			return nil, fmt.Errorf("unexpected node type '%s' in formal parameter list at index %d", child.RootType, i)
		}
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_VARIABLE_DECLARATIONS,
		Children: parameters,
	}, nil
}

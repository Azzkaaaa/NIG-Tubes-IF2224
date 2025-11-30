package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeArrayAccess(parsetree *dt.ParseTree, prev *dt.DecoratedSyntaxTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parsetree.RootType != dt.ARRAY_ACCESS_NODE {
		return nil, semanticType{}, errors.New("parse tree node is not array access node")
	}

	root := a.root

	if prev != nil {
		switch prev.SelfType {
		case dt.DST_CONST:
			fallthrough
		case dt.DST_VARIABLE:
			fallthrough
		case dt.DST_RECORD_FIELD:
			fallthrough
		case dt.DST_FUNCTION_CALL:
			symbolIndex := prev.Data

			for a.tab[symbolIndex].Type == dt.TAB_ENTRY_ALIAS {
				symbolIndex = a.tab[symbolIndex].Reference
			}

			if a.tab[symbolIndex].Type != dt.TAB_ENTRY_RECORD {
				return nil, semanticType{}, errors.New("non record type has no fields")
			}

			root = a.btab[a.tab[symbolIndex].Reference].End
		}
	}

	token := parsetree.Children[0].TokenValue
	index, tabEntry := a.tab.FindIdentifier(token.Lexeme, root)

	if tabEntry == nil {
		return nil, semanticType{}, errors.New("undeclared identifier")
	}

	var dstType dt.DSTNodeType

	switch tabEntry.Object {
	case dt.TAB_ENTRY_CONST:
		dstType = dt.DST_CONST
	case dt.TAB_ENTRY_VAR:
		dstType = dt.DST_VARIABLE
	case dt.TAB_ENTRY_FIELD:
		dstType = dt.DST_RECORD_FIELD
	default:
		return nil, semanticType{}, errors.New("identifier does not reference a constant, variable, or field")
	}

	if tabEntry.Type != dt.TAB_ENTRY_ARRAY {
		return nil, semanticType{}, errors.New("identifier does not hold an array value")
	}

	var children []dt.DecoratedSyntaxTree
	if prev == nil {
		children = make([]dt.DecoratedSyntaxTree, 0)
	} else {
		children = []dt.DecoratedSyntaxTree{*prev}
	}

	prev = &dt.DecoratedSyntaxTree{
		Property: dt.DST_FROM,
		SelfType: dstType,
		Data:     index,
		Children: children,
	}

	// Extract nodes with step 3 (every 3rd element starting from index 2)
	recursiveNodes := make([]dt.ParseTree, 0)
	for i := 2; i < len(parsetree.Children); i += 3 {
		recursiveNodes = append(recursiveNodes, parsetree.Children[i])
	}

	return a.analyzeRecursiveArrayAccess(recursiveNodes, prev)
}

func (a *SemanticAnalyzer) analyzeRecursiveArrayAccess(nodes []dt.ParseTree, prev *dt.DecoratedSyntaxTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if nodes[0].RootType != dt.EXPRESSION_NODE {
		return nil, semanticType{}, errors.New("expected valid array index")
	}

	atabIndex := 0

	switch prev.SelfType {
	case dt.DST_VARIABLE:
		tabIndex := prev.Data

		for a.tab[tabIndex].Type == dt.TAB_ENTRY_ALIAS {
			tabIndex = a.tab[tabIndex].Reference
		}

		if a.tab[tabIndex].Type != dt.TAB_ENTRY_ARRAY {
			return nil, semanticType{}, errors.New("expected variable to point to an array")
		}

		atabIndex = a.tab[tabIndex].Reference
	case dt.DST_ARRAY_ELEMENT:
		atabIndex = prev.Data
	default:
		return nil, semanticType{}, errors.New("object cannot be indexed")
	}

	expectedType := semanticType{
		StaticType: a.atab[atabIndex].ElementType,
		Reference:  a.atab[atabIndex].ElementReference,
	}

	index, indexType, err := a.analyzeExpression(&nodes[0])

	if err != nil {
		return nil, semanticType{}, err
	}

	if indexType.StaticType != a.atab[atabIndex].IndexType {
		return nil, semanticType{}, errors.New("expression type does not match index type of array")
	}

	index.Property = dt.DST_INDEX

	var self *dt.DecoratedSyntaxTree

	if prev != nil {
		prev.Property = dt.DST_FROM
		self = &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_ARRAY_ELEMENT,
			Data:     atabIndex,
			Children: []dt.DecoratedSyntaxTree{
				*prev,
				*index,
			},
		}
	} else {
		self = &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_ARRAY_ELEMENT,
			Data:     atabIndex,
			Children: []dt.DecoratedSyntaxTree{*index},
		}
	}

	if len(nodes) == 1 {
		return self, expectedType, nil
	}

	if expectedType.StaticType != dt.TAB_ENTRY_ARRAY && len(nodes) > 1 {
		return nil, semanticType{}, errors.New("cannot access non-array type as if it was an array")
	}

	return a.analyzeRecursiveArrayAccess(nodes[1:], self)
}

package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeArrayAccess(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parsetree.RootType != dt.ARRAY_ACCESS_NODE {
		return nil, semanticType{}, errors.New("parse tree node is not array access node")
	}

	token := parsetree.Children[0].TokenValue
	index, tabEntry := a.tab.FindIdentifier(token.Lexeme, a.root)

	if tabEntry == nil {
		return nil, semanticType{}, errors.New("undeclared identifier")
	}

	switch tabEntry.Object {
	case dt.TAB_ENTRY_CONST:
	case dt.TAB_ENTRY_VAR:
	default:
		return nil, semanticType{}, errors.New("identifier does not reference a constant or an array")
	}

	if tabEntry.Type != dt.TAB_ENTRY_ARRAY {
		return nil, semanticType{}, errors.New("identifier does not hold an array value")
	}

	atabEntry := a.atab[tabEntry.Reference]

	dst, typ, err := a.analyzeExpression(&parsetree.Children[2])

	if err != nil {
		return nil, typ, err
	}

	if typ.StaticType != atabEntry.IndexType {
		return nil, typ, errors.New("type of index expression does not match index type of array")
	}

	dst.Property = dt.DST_INDEX

	return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_ARRAY_ELEMENT,
			Data:     index,
			Children: []dt.DecoratedSyntaxTree{*dst},
		}, semanticType{
			StaticType: atabEntry.ElementType,
			Reference:  atabEntry.ElementReference,
		}, nil
}

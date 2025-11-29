package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeIfStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.IF_STATEMENT_NODE {
		return nil, errors.New("expected if block")
	}

	condition, typ, err := a.analyzeExpression(&parsetree.Children[1])

	if err != nil {
		return nil, err
	}

	if typ.StaticType != dt.TAB_ENTRY_BOOLEAN {
		return nil, errors.New("condition must be a boolean type")
	}

	thenBlock, err := a.analyzeStatement(&parsetree.Children[3])

	if err != nil {
		return nil, err
	}

	elseBlock, err := a.analyzeStatement(&parsetree.Children[5])

	if err != nil {
		return nil, err
	}

	condition.Property = dt.DST_CONDITION
	thenBlock.Property = dt.DST_THEN
	elseBlock.Property = dt.DST_ELSE

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_IF_BLOCK,
		Children: []dt.DecoratedSyntaxTree{
			*condition,
			*thenBlock,
			*elseBlock,
		},
	}, nil
}

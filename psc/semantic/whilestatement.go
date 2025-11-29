package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeWhileStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.WHILE_STATEMENT_NODE {
		return nil, errors.New("expected while block")
	}

	condition, typ, err := a.analyzeExpression(&parsetree.Children[1])

	if err != nil {
		return nil, err
	}

	if typ.StaticType != dt.TAB_ENTRY_BOOLEAN {
		return nil, errors.New("condition expression must be of boolean type")
	}

	block, err := a.analyzeStatement(&parsetree.Children[3])

	if err != nil {
		return nil, err
	}

	condition.Property = dt.DST_CONDITION
	block.Property = dt.DST_EXECUTE

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_WHILE_BLOCK,
		Children: []dt.DecoratedSyntaxTree{
			*condition,
			*block,
		},
	}, nil
}

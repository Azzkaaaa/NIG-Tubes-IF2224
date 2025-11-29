package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeCompoundStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.COMPOUND_STATEMENT_NODE {
		return nil, errors.New("expected compound statement section")
	}

	statements, err := a.analyzeStatementList(&parsetree.Children[1])

	if err != nil {
		return nil, err
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_BLOCK,
		Children: statements,
	}, nil
}

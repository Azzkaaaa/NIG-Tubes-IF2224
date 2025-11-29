package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeStatementList(parsetree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.STATEMENT_LIST_NODE {
		return nil, errors.New("expected statement list section")
	}

	statements := parsetree.Children[:len(parsetree.Children):1]
	decoratedStatements := make([]dt.DecoratedSyntaxTree, len(statements))

	for i, statement := range statements {
		dst, err := a.analyzeStatement(&statement)

		if err != nil {
			return nil, err
		}

		decoratedStatements[i] = *dst
	}

	return decoratedStatements, nil
}

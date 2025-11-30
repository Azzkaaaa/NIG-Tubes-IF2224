package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeStatementList(parsetree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.STATEMENT_LIST_NODE {
		return nil, errors.New("expected statement list section")
	}

	decoratedStatements := make([]dt.DecoratedSyntaxTree, 0)

	for _, statement := range parsetree.Children {
		if statement.RootType == dt.TOKEN_NODE {
			continue
		}

		dst, err := a.analyzeStatement(&statement)

		if err != nil {
			return nil, err
		}

		if dst != nil {
			decoratedStatements = append(decoratedStatements, *dst)
		}
	}

	return decoratedStatements, nil
}

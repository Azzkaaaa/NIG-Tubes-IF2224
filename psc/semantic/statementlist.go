package semantic

import (
	"errors"
	"fmt"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeStatementList(parsetree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, error) {
	fmt.Printf("[STMT_LIST] Analyzing statement list, children count: %d\n", len(parsetree.Children))
	if parsetree.RootType != dt.STATEMENT_LIST_NODE {
		return nil, errors.New("expected statement list section")
	}

	decoratedStatements := make([]dt.DecoratedSyntaxTree, 0)

	for i, statement := range parsetree.Children {
		fmt.Printf("[STMT_LIST] Processing child %d: type=%s\n", i, statement.RootType)

		// Skip token nodes (like semicolons)
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

	fmt.Printf("[STMT_LIST] Total statements: %d\n", len(decoratedStatements))
	return decoratedStatements, nil
}

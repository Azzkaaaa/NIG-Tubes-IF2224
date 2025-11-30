package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeParameterList(parseTree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, []semanticType, error) {
	if parseTree.RootType != dt.PARAMETER_LIST_NODE {
		return nil, nil, errors.New("parse tree node is not parameter list")
	}

	parameters := parseTree.Children
	decoratedParams := make([]dt.DecoratedSyntaxTree, len(parameters)/2+1)
	paramTypes := make([]semanticType, len(parameters)/2+1)

	for i, p := range parameters {
		if i%2 == 1 {
			continue
		}

		param, typ, err := a.analyzeExpression(&p)

		if err != nil {
			return nil, nil, err
		}

		decoratedParams[i/2] = *param
		paramTypes[i/2] = typ
	}

	return decoratedParams, paramTypes, nil
}

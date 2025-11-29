package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeParameterList(parseTree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, []semanticType, error) {
	if parseTree.RootType != dt.PARAMETER_LIST_NODE {
		return nil, nil, errors.New("parse tree node is not parameter list")
	}

	parameters := parseTree.Children[:len(parseTree.Children):1]
	decoratedParams := make([]dt.DecoratedSyntaxTree, len(parameters))
	paramTypes := make([]semanticType, len(parameters))

	for i, p := range parameters {
		param, typ, err := a.analyzeExpression(&p)

		if err != nil {
			return nil, nil, err
		}

		decoratedParams[i] = *param
		paramTypes[i] = typ
	}

	return decoratedParams, paramTypes, nil
}

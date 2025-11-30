package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeAssignmentStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.ASSIGNMENT_STATEMENT_NODE {
		return nil, errors.New("expected an assignment statement node")
	}

	target, targetType, err := a.analyzeStaticAccess(&parsetree.Children[0], nil)

	if err != nil {
		return nil, err
	}

	value, valueType, err := a.analyzeExpression(&parsetree.Children[2])

	if err != nil {
		return nil, err
	}

	if !a.checkTypeEquality(targetType, valueType) {
		return nil, errors.New("target type and value type do not match")
	}

	target.Property = dt.DST_TARGET
	value.Property = dt.DST_VALUE

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_ASSIGNMENT_OPERATOR,
		Data:     int(targetType.StaticType),
		Children: []dt.DecoratedSyntaxTree{
			*target,
			*value,
		},
	}, nil
}

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

	// Check type compatibility and insert implicit cast if needed
	// Cast from value type (right) to target type (left)

	if !a.checkTypeEquality(targetType, valueType) {
		if a.canCastImplicitly(valueType, targetType) {
			value, valueType = a.insertImplicitCast(value, valueType, targetType)
		} else {
			// Get token for error reporting (from assignment operator)
			assignToken := parsetree.Children[1].TokenValue
			return nil, a.newAssignmentError(
				targetType.StaticType.String(),
				valueType.StaticType.String(),
				assignToken,
			)
		}
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

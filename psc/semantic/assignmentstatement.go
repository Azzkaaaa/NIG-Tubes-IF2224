package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeAssignmentStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.ASSIGNMENT_STATEMENT_NODE {
		return nil, errors.New("expected an assignment statement node")
	}

	var target *dt.DecoratedSyntaxTree
	var targetType semanticType
	var err error

	switch parsetree.Children[0].RootType {
	case dt.ARRAY_ACCESS_NODE:
		target, targetType, err = a.analyzeArrayAccess(&parsetree.Children[0])
		if a.tab[target.Data].Object == dt.TAB_ENTRY_CONST {
			return nil, errors.New("cannot assign value to constant array")
		}
	case dt.TOKEN_NODE:
		target, targetType, err = a.analyzeToken(&parsetree.Children[0])
		if target.SelfType != dt.DST_VARIABLE {
			return nil, errors.New("expected array element or variable as assignment target")
		}
	default:
		return nil, errors.New("expected array element or variable as assignment target")
	}

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

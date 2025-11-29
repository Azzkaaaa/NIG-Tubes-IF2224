package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeExpression(parseTree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if len(parseTree.Children) == 1 {
		return a.analyzeSimpleExpression(&parseTree.Children[0])
	} else {
		optype, err := a.analyzeRelationalOperator(&parseTree.Children[0])

		if err != nil {
			return nil, semanticType{}, err
		}

		lhs, ltype, err := a.analyzeSimpleExpression(&parseTree.Children[0])

		if err != nil {
			return nil, ltype, err
		}

		rhs, rtype, err := a.analyzeSimpleExpression(&parseTree.Children[2])

		if err != nil {
			return nil, rtype, err
		}

		if !a.checkTypeEquality(ltype, rtype) {
			return nil, ltype, errors.New("operand types do not match")
		}

		return &dt.DecoratedSyntaxTree{
			SelfType: optype,
			Children: []dt.DecoratedSyntaxTree{
				*lhs,
				*rhs,
			},
		}, semanticType{StaticType: dt.TAB_ENTRY_BOOLEAN}, nil
	}
}

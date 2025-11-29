package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeRange(parsetree *dt.ParseTree) (int, int, semanticType, error) {
	if parsetree.RootType != dt.RANGE_NODE {
		return 0, 0, semanticType{}, errors.New("expected range expression")
	}

	beginExpression, beginType, err := a.analyzeExpression(&parsetree.Children[0])

	if err != nil {
		return 0, 0, semanticType{}, err
	}

	begin, err := a.staticEvaluate(beginExpression, beginType)

	if err != nil {
		return 0, 0, semanticType{}, err
	}

	endExpression, endType, err := a.analyzeExpression(&parsetree.Children[2])

	if err != nil {
		return 0, 0, semanticType{}, err
	}

	if !a.checkTypeEquality(beginType, endType) {
		return 0, 0, semanticType{}, errors.New("range expression types do not match")
	}

	end, err := a.staticEvaluate(endExpression, endType)

	if err != nil {
		return 0, 0, semanticType{}, err
	}

	return begin, end, beginType, nil
}

package semantic

import (
	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeTerm(parseTree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if len(parseTree.Children) == 1 {
		return a.analyzeFactor(&parseTree.Children[0])
	} else {
		return a.recurseTerm(parseTree.Children)
	}
}

func (a *SemanticAnalyzer) recurseTerm(nodes []dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if len(nodes) == 1 {
		return a.analyzeFactor(&nodes[0])
	}

	optype, err := a.analyzeMultiplicativeOperator(&nodes[len(nodes)-2])

	if err != nil {
		return nil, semanticType{}, err
	}

	dst := &dt.DecoratedSyntaxTree{
		SelfType: optype,
		Children: make([]dt.DecoratedSyntaxTree, 2),
	}

	lval, ltype, err := a.recurseTerm(nodes[:len(nodes)-2])

	if err != nil {
		return nil, ltype, err
	}

	rval, rtype, err := a.analyzeFactor(&nodes[len(nodes)-1])

	if err != nil {
		return nil, rtype, err
	}

	// Promote types if needed (e.g., integer * real)
	promotedLval, promotedRval, resultType, compatible := a.promoteTypes(lval, ltype, rval, rtype)

	if !compatible {
		// Get operator token
		token := nodes[len(nodes)-2].Children[0].TokenValue
		return nil, ltype, a.newOperatorTypeError(
			token.Lexeme,
			ltype.StaticType.String(),
			rtype.StaticType.String(),
			token,
		)
	}

	dst.Children[0] = *promotedLval
	dst.Children[0].Property = dt.DST_OPERAND

	dst.Children[1] = *promotedRval
	dst.Children[1].Property = dt.DST_OPERAND

	return dst, resultType, nil
}

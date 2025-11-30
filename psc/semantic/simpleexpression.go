package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeSimpleExpression(parseTree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	beginOffset := 0
	negated := false

	if parseTree.Children[0].RootType == dt.TOKEN_NODE {
		switch parseTree.Children[0].TokenValue.Lexeme {
		case "-":
			negated = true
			fallthrough
		case "+":
			beginOffset = 1
		default:
			return nil, semanticType{}, errors.New("unknown modifier at beginning of expression")
		}
	}

	var dst *dt.DecoratedSyntaxTree
	var typ semanticType
	var err error

	if len(parseTree.Children) == 1 {
		dst, typ, err = a.analyzeTerm(&parseTree.Children[beginOffset])
	} else {
		dst, typ, err = a.recurseSimpleExpression(parseTree.Children[beginOffset:])
	}

	if err != nil {
		return nil, typ, err
	}

	if negated {
		switch typ.StaticType {
		case dt.TAB_ENTRY_INTEGER:
		case dt.TAB_ENTRY_REAL:
		default:
			return nil, typ, errors.New("cannot negate non numeric expression")
		}

		dst.Property = dt.DST_OPERAND
		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_NEG_OPERATOR,
			Children: []dt.DecoratedSyntaxTree{*dst},
		}, typ, nil
	}

	return dst, typ, nil
}

func (a *SemanticAnalyzer) recurseSimpleExpression(nodes []dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if len(nodes) == 1 {
		return a.analyzeTerm(&nodes[0])
	}

	optype, err := a.analyzeAdditiveOperator(&nodes[len(nodes)-2])

	if err != nil {
		return nil, semanticType{}, err
	}

	dst := &dt.DecoratedSyntaxTree{
		SelfType: optype,
		Children: make([]dt.DecoratedSyntaxTree, 2),
	}

	lval, ltype, err := a.recurseSimpleExpression(nodes[:len(nodes)-2])

	if err != nil {
		return nil, ltype, err
	}

	rval, rtype, err := a.analyzeTerm(&nodes[len(nodes)-1])

	if err != nil {
		return nil, rtype, err
	}

	if !a.checkTypeEquality(ltype, rtype) {
		// Get operator token
		token := nodes[len(nodes)-2].Children[0].TokenValue
		return nil, ltype, a.newOperatorTypeError(
			token.Lexeme,
			ltype.StaticType.String(),
			rtype.StaticType.String(),
			token,
		)
	}

	dst.Children[0] = *lval
	dst.Children[0].Property = dt.DST_OPERAND

	dst.Children[1] = *rval
	dst.Children[1].Property = dt.DST_OPERAND

	return dst, ltype, nil
}

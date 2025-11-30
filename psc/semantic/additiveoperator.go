package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeAdditiveOperator(parsetree *dt.ParseTree) (dt.DSTNodeType, error) {
	if parsetree.RootType != dt.ADDITIVE_OPERATOR_NODE {
		return dt.DST_ADD_OPERATOR, errors.New("operator is not additive")
	}

	switch parsetree.Children[0].TokenValue.Lexeme {
	case "+":
		return dt.DST_ADD_OPERATOR, nil
	case "-":
		return dt.DST_SUB_OPERATOR, nil
	case "atau":
		return dt.DST_OR_OPERATOR, nil
	default:
		return dt.DST_ADD_OPERATOR, errors.New("unknown additive operator")
	}
}

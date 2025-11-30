package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeMultiplicativeOperator(parsetree *dt.ParseTree) (dt.DSTNodeType, error) {
	if parsetree.RootType != dt.MULTIPLICATIVE_OPERATOR_NODE {
		return dt.DST_ADD_OPERATOR, errors.New("operator is not multiplicative")
	}

	switch parsetree.Children[0].TokenValue.Lexeme {
	case "*":
		return dt.DST_MUL_OPERATOR, nil
	case "/":
		return dt.DST_DIV_OPERATOR, nil
	case "bagi":
		return dt.DST_DIV_OPERATOR, nil
	case "mod":
		return dt.DST_MOD_OPERATOR, nil
	case "dan":
		return dt.DST_AND_OPERATOR, nil
	default:
		return dt.DST_ADD_OPERATOR, errors.New("unknown multiplicative operator")
	}
}

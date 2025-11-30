package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeRelationalOperator(parsetree *dt.ParseTree) (dt.DSTNodeType, error) {
	if parsetree.RootType != dt.RELATIONAL_OPERATOR_NODE {
		return dt.DST_ADD_OPERATOR, errors.New("operator is not relational")
	}

	switch parsetree.Children[0].TokenValue.Lexeme {
	case ">":
		return dt.DST_GT_OPERATOR, nil
	case "<":
		return dt.DST_LT_OPERATOR, nil
	case ">=":
		return dt.DST_GE_OPERATOR, nil
	case "<=":
		return dt.DST_LE_OPERATOR, nil
	case "=":
		return dt.DST_EQ_OPERATOR, nil
	case "!=":
		return dt.DST_NE_OPERATOR, nil
	default:
		return dt.DST_ADD_OPERATOR, errors.New("unknown relational operator")
	}
}

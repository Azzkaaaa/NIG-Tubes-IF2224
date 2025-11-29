package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

// Currently allows expressions from literals and previously defined constants

func (a *SemanticAnalyzer) staticEvaluate(dst *dt.DecoratedSyntaxTree, typ semanticType) (int, error) {
	switch typ.StaticType {
	case dt.TAB_ENTRY_INTEGER:
		return a.staticEvaluateInt(dst)
	case dt.TAB_ENTRY_REAL:
		return a.staticEvaluateReal(dst)
	case dt.TAB_ENTRY_CHAR:
		return a.staticEvaluateChar(dst)
	case dt.TAB_ENTRY_BOOLEAN:
		return a.staticEvaluateBool(dst)
	default:
		return -1, errors.New("cannot evaluate expression statically")
	}
}

func (a *SemanticAnalyzer) staticEvaluateInt(dst *dt.DecoratedSyntaxTree) (int, error) {
	return 0, nil
}

func (a *SemanticAnalyzer) staticEvaluateReal(dst *dt.DecoratedSyntaxTree) (int, error) {
	return 0, nil
}

func (a *SemanticAnalyzer) staticEvaluateChar(dst *dt.DecoratedSyntaxTree) (int, error) {
	return 0, nil
}

func (a *SemanticAnalyzer) staticEvaluateBool(dst *dt.DecoratedSyntaxTree) (int, error) {
	return 0, nil
}

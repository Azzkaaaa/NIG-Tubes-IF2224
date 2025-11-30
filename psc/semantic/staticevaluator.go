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
	switch dst.SelfType {
	case dt.DST_INT_LITERAL:
		return dst.Data, nil
	case dt.DST_CONST:
		// Reference to a constant - look it up in the symbol table
		constEntry := a.tab[dst.Data]
		if constEntry.Type != dt.TAB_ENTRY_INTEGER {
			return 0, errors.New("constant is not an integer")
		}
		return constEntry.Data, nil
	case dt.DST_ADD_OPERATOR:
		left, err := a.staticEvaluateInt(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateInt(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		return left + right, nil
	case dt.DST_SUB_OPERATOR:
		left, err := a.staticEvaluateInt(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateInt(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		return left - right, nil
	case dt.DST_MUL_OPERATOR:
		left, err := a.staticEvaluateInt(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateInt(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		return left * right, nil
	case dt.DST_DIV_OPERATOR:
		left, err := a.staticEvaluateInt(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateInt(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		if right == 0 {
			return 0, errors.New("division by zero in static evaluation")
		}
		return left / right, nil
	case dt.DST_MOD_OPERATOR:
		left, err := a.staticEvaluateInt(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateInt(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		if right == 0 {
			return 0, errors.New("modulo by zero in static evaluation")
		}
		return left % right, nil
	case dt.DST_NEG_OPERATOR:
		val, err := a.staticEvaluateInt(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		return -val, nil
	default:
		return 0, errors.New("cannot statically evaluate integer expression")
	}
}

func (a *SemanticAnalyzer) staticEvaluateReal(dst *dt.DecoratedSyntaxTree) (int, error) {
	switch dst.SelfType {
	case dt.DST_REAL_LITERAL:
		return dst.Data, nil
	case dt.DST_CONST:
		// Reference to a constant - look it up in the symbol table
		constEntry := a.tab[dst.Data]
		if constEntry.Type != dt.TAB_ENTRY_REAL {
			return 0, errors.New("constant is not a real")
		}
		return constEntry.Data, nil
	default:
		return 0, errors.New("cannot statically evaluate real expression")
	}
}

func (a *SemanticAnalyzer) staticEvaluateChar(dst *dt.DecoratedSyntaxTree) (int, error) {
	switch dst.SelfType {
	case dt.DST_CHAR_LITERAL:
		return dst.Data, nil
	case dt.DST_CONST:
		// Reference to a constant - look it up in the symbol table
		constEntry := a.tab[dst.Data]
		if constEntry.Type != dt.TAB_ENTRY_CHAR {
			return 0, errors.New("constant is not a char")
		}
		return constEntry.Data, nil
	default:
		return 0, errors.New("cannot statically evaluate char expression")
	}
}

func (a *SemanticAnalyzer) staticEvaluateBool(dst *dt.DecoratedSyntaxTree) (int, error) {
	switch dst.SelfType {
	case dt.DST_BOOL_LITERAL:
		return dst.Data, nil
	case dt.DST_CONST:
		// Reference to a constant - look it up in the symbol table
		constEntry := a.tab[dst.Data]
		if constEntry.Type != dt.TAB_ENTRY_BOOLEAN {
			return 0, errors.New("constant is not a boolean")
		}
		return constEntry.Data, nil
	case dt.DST_AND_OPERATOR:
		left, err := a.staticEvaluateBool(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateBool(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		if left != 0 && right != 0 {
			return 1, nil
		}
		return 0, nil
	case dt.DST_OR_OPERATOR:
		left, err := a.staticEvaluateBool(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		right, err := a.staticEvaluateBool(&dst.Children[1])
		if err != nil {
			return 0, err
		}
		if left != 0 || right != 0 {
			return 1, nil
		}
		return 0, nil
	case dt.DST_NOT_OPERATOR:
		val, err := a.staticEvaluateBool(&dst.Children[0])
		if err != nil {
			return 0, err
		}
		if val == 0 {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, errors.New("cannot statically evaluate boolean expression")
	}
}

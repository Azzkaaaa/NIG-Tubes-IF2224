package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeSubprogramCall(parseTree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parseTree.RootType != dt.SUBPROGRAM_CALL_NODE {
		return nil, semanticType{}, errors.New("parse tree node is not subprogram call node")
	}

	subprogramIdentifier := parseTree.Children[0].TokenValue.Lexeme
	index, tabEntry := a.tab.FindIdentifier(subprogramIdentifier, a.root)

	if tabEntry == nil {
		return nil, semanticType{}, errors.New("identifier is not defined")
	}

	var callType dt.DSTNodeType

	switch tabEntry.Object {
	case dt.TAB_ENTRY_FUNC:
		callType = dt.DST_FUNCTION_CALL
	case dt.TAB_ENTRY_PROC:
		callType = dt.DST_PROCEDURE_CALL
	default:
		return nil, semanticType{}, errors.New("identifier does not reference a function or procedure")
	}

	btabEntry := a.btab[tabEntry.Data]
	paramStart := btabEntry.Start
	paramEnd := btabEntry.ParamEnd

	callParams, callTypes, err := a.analyzeParameterList(&parseTree.Children[2])

	if err != nil {
		return nil, semanticType{}, err
	}

	if len(callParams) != (paramEnd - paramStart + 1) {
		return nil, semanticType{}, errors.New("parameter count mismatch")
	}

	dst := &dt.DecoratedSyntaxTree{
		SelfType: callType,
		Data:     index,
		Children: make([]dt.DecoratedSyntaxTree, len(callParams)),
	}

	for i := paramStart; i < paramEnd; i++ {
		declaredType := semanticType{
			StaticType: a.tab[i].Type,
			Reference:  a.tab[i].Reference,
		}

		if declaredType != callTypes[i] {
			return nil, semanticType{}, errors.New("parameter type mismatch")
		}

		dst.Children[i] = callParams[i]
	}

	return dst, semanticType{
		StaticType: tabEntry.Type,
		Reference:  tabEntry.Reference,
	}, nil
}

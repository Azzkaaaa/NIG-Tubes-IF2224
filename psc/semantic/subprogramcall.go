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
		token := parseTree.Children[0].TokenValue
		return nil, semanticType{}, a.newUndeclaredIdentError(subprogramIdentifier, token)
	}

	var callType dt.DSTNodeType

	switch tabEntry.Object {
	case dt.TAB_ENTRY_FUNC:
		callType = dt.DST_FUNCTION_CALL
	case dt.TAB_ENTRY_PROC:
		callType = dt.DST_PROCEDURE_CALL
	default:
		token := parseTree.Children[0].TokenValue
		return nil, semanticType{}, a.newNotCallableError(
			subprogramIdentifier,
			tabEntry.Object.String(),
			token,
		)
	}

	btabEntry := a.btab[tabEntry.Data]
	paramStart := btabEntry.Start
	paramEnd := btabEntry.ParamEnd

	callParams, callTypes, err := a.analyzeParameterList(&parseTree.Children[2])

	if err != nil {
		return nil, semanticType{}, err
	}

	if len(callParams) != (paramEnd - paramStart + 1) {
		token := parseTree.Children[0].TokenValue
		return nil, semanticType{}, a.newParameterCountError(
			paramEnd-paramStart+1,
			len(callParams),
			subprogramIdentifier,
			token,
		)
	}

	dst := &dt.DecoratedSyntaxTree{
		SelfType: callType,
		Data:     index,
		Children: make([]dt.DecoratedSyntaxTree, len(callParams)),
	}

	for i := paramStart; i < paramEnd+1; i++ {
		declaredType := semanticType{
			StaticType: a.tab[i].Type,
			Reference:  a.tab[i].Reference,
		}

		if !a.checkTypeEquality(declaredType, callTypes[i-paramStart]) {
			token := parseTree.Children[0].TokenValue
			return nil, semanticType{}, a.newParameterTypeError(
				i-paramStart,
				declaredType.StaticType.String(),
				callTypes[i-paramStart].StaticType.String(),
				subprogramIdentifier,
				token,
			)
		}

		dst.Children[i-paramStart] = callParams[i-paramStart]
	}

	return dst, semanticType{
		StaticType: tabEntry.Type,
		Reference:  tabEntry.Reference,
	}, nil
}

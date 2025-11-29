package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeFactor(parseTree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	switch parseTree.Children[0].RootType {
	case dt.TOKEN_NODE:
		if parseTree.Children[0].TokenValue.Type == dt.LPARENTHESIS {
			return a.analyzeExpression(&parseTree.Children[1])
		} else if parseTree.Children[0].TokenValue.Lexeme == "tidak" {
			dst, typ, err := a.analyzeFactor(&parseTree.Children[1])

			if err != nil {
				return nil, typ, err
			}

			if typ.StaticType != dt.TAB_ENTRY_BOOLEAN {
				return nil, typ, errors.New("not operator only works on boolean expressions")
			}

			dst.Property = dt.DST_OPERAND

			return &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_NOT_OPERATOR,
				Children: []dt.DecoratedSyntaxTree{*dst},
			}, typ, nil
		} else {
			return a.analyzeToken(&parseTree.Children[0])
		}
	case dt.ARRAY_ACCESS_NODE:
		return a.analyzeArrayAccess(&parseTree.Children[0])
	case dt.SUBPROGRAM_CALL_NODE:
		return a.analyzeSubprogramCall(&parseTree.Children[0])
	}

	return nil, semanticType{}, errors.New("idk what happened")
}

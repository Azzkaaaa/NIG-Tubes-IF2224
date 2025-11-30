package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeFactor(parseTree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parseTree.Children[0].RootType == dt.TOKEN_NODE {
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
		}
	}

	return a.analyzeAccess(&parseTree.Children[0])
}

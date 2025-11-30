package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeStaticAccess(parsetree *dt.ParseTree, prev *dt.DecoratedSyntaxTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parsetree.RootType != dt.STATIC_ACCESS_NODE {
		return nil, semanticType{}, errors.New("expected static access")
	}

	nodes := parsetree.Children[:len(parsetree.Children):2]
	root := a.root

	var typ semanticType
	var err error

	if prev != nil {
		typeIndex := prev.Data

		for a.tab[typeIndex].Type == dt.TAB_ENTRY_ALIAS {
			typeIndex = a.tab[typeIndex].Reference
		}

		if a.tab[typeIndex].Type != dt.TAB_ENTRY_RECORD {
			return nil, semanticType{}, errors.New("cannot access field of non struct object")
		}

		root = a.btab[a.tab[typeIndex].Reference].End
	}

	switch nodes[0].RootType {
	case dt.TOKEN_NODE:
		tabIndex, _ := a.tab.FindIdentifier(nodes[0].TokenValue.Lexeme, root)
		if tabIndex == -1 {
			return nil, semanticType{}, errors.New("undeclared identifier")
		}
		if prev != nil {
			prev.Property = dt.DST_FROM
			prev = &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_RECORD_FIELD,
				Data:     tabIndex,
				Children: []dt.DecoratedSyntaxTree{*prev},
			}
		} else {
			var dstType dt.DSTNodeType

			switch a.tab[tabIndex].Object {
			case dt.TAB_ENTRY_CONST:
				dstType = dt.DST_CONST
			case dt.TAB_ENTRY_PARAM:
				fallthrough
			case dt.TAB_ENTRY_RETURN:
				fallthrough
			case dt.TAB_ENTRY_VAR:
				dstType = dt.DST_VARIABLE
			default:
				return nil, semanticType{}, errors.New("unexpected object")
			}

			prev = &dt.DecoratedSyntaxTree{
				SelfType: dstType,
				Data:     tabIndex,
			}
		}
	case dt.ARRAY_ACCESS_NODE:
		prev, typ, err = a.analyzeArrayAccess(&nodes[0], prev)
		if err != nil {
			return nil, semanticType{}, err
		}
	default:
		return nil, semanticType{}, errors.New("found unexpected node")
	}

	for _, node := range nodes[1:] {
		typeIndex := prev.Data

		for a.tab[typeIndex].Type == dt.TAB_ENTRY_ALIAS {
			typeIndex = a.tab[typeIndex].Reference
		}

		if a.tab[typeIndex].Type != dt.TAB_ENTRY_RECORD {
			return nil, semanticType{}, errors.New("cannot access field of non struct object")
		}

		root = a.btab[a.tab[typeIndex].Reference].End

		switch node.RootType {
		case dt.TOKEN_NODE:
			tabIndex, _ := a.tab.FindIdentifier(node.TokenValue.Lexeme, root)
			if tabIndex == -1 {
				return nil, semanticType{}, errors.New("undeclared identifier")
			}
			prev.Property = dt.DST_FROM
			prev = &dt.DecoratedSyntaxTree{
				SelfType: dt.DST_RECORD_FIELD,
				Data:     tabIndex,
				Children: []dt.DecoratedSyntaxTree{*prev},
			}
		case dt.ARRAY_ACCESS_NODE:
			prev, typ, err = a.analyzeArrayAccess(&node, prev)
			if err != nil {
				return nil, semanticType{}, err
			}
		default:
			return nil, semanticType{}, errors.New("found unexpected node")
		}
	}

	return prev, typ, nil
}

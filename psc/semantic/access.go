package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeAccess(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, semanticType, error) {
	if parsetree.RootType != dt.ACCESS_NODE {
		return nil, semanticType{}, errors.New("expected access node")
	}

	callNode := parsetree.Children[0]
	var staticAccessNode *dt.ParseTree

	var prev *dt.DecoratedSyntaxTree
	var typ semanticType
	var err error

	if callNode.RootType == dt.SUBPROGRAM_CALL_NODE {
		if len(parsetree.Children) == 3 {
			staticAccessNode = &parsetree.Children[2]
		}

		prev, typ, err = a.analyzeSubprogramCall(&callNode)

		if err != nil {
			return nil, semanticType{}, err
		}
	} else {
		staticAccessNode = &parsetree.Children[0]
	}

	if staticAccessNode != nil {
		return a.analyzeStaticAccess(staticAccessNode, prev)
	}

	return prev, typ, nil
}

package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeVarDeclarationPart(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.VAR_DECLARATION_PART_NODE {
		return nil, errors.New("expected var declaration part")
	}

	declarationNodes := parsetree.Children[1:]
	declarations := make([]dt.DecoratedSyntaxTree, 0)

	for _, node := range declarationNodes {
		partialDeclarations, err := a.analyzeVarDeclaration(&node)

		if err != nil {
			return nil, err
		}

		declarations = append(declarations, partialDeclarations...)
	}

	for i := range declarations {
		declarations[i].Property = dt.DST_DECLARE
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_VARIABLE_DECLARATIONS,
		Children: declarations,
	}, nil
}

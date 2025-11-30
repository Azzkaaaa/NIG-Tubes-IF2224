package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeConstDeclarationPart(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.CONST_DECLARATION_PART_NODE {
		return nil, errors.New("expected const declaration part")
	}

	declarations := make([]dt.DecoratedSyntaxTree, len(parsetree.Children)-1)

	for i, constDeclaration := range parsetree.Children[1:] {
		declaration, err := a.analyzeConstDeclaration(&constDeclaration)

		if err != nil {
			return nil, err
		}

		declarations[i] = *declaration
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_CONSTANT_DECLARATIONS,
		Children: declarations,
	}, nil
}

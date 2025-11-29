package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeTypeDeclarationPart(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.TYPE_DECLARATION_PART_NODE {
		return nil, errors.New("expected type declaration part")
	}

	declarations := make([]dt.DecoratedSyntaxTree, len(parsetree.Children)-1)

	for i, typeDeclaration := range parsetree.Children[1:] {
		declaration, err := a.analyzeTypeDeclaration(&typeDeclaration)

		if err != nil {
			return nil, err
		}

		declarations[i] = *declaration
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_TYPE_DECLARATIONS,
		Children: declarations,
	}, nil
}

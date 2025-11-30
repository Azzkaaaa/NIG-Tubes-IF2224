package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeDeclarationPart(parsetree *dt.ParseTree) ([]dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.DECLARATION_PART_NODE {
		return nil, errors.New("expected declaration part")
	}

	declarations := make([]dt.DecoratedSyntaxTree, 0)

	for _, child := range parsetree.Children {
		var declaration *dt.DecoratedSyntaxTree
		var err error

		switch child.RootType {
		case dt.CONST_DECLARATION_PART_NODE:
			declaration, err = a.analyzeConstDeclarationPart(&child)
		case dt.TYPE_DECLARATION_PART_NODE:
			declaration, err = a.analyzeTypeDeclarationPart(&child)
		case dt.VAR_DECLARATION_PART_NODE:
			declaration, err = a.analyzeVarDeclarationPart(&child)
		case dt.SUBPROGRAM_DECLARATION_NODE:
			declaration, err = a.analyzeSubprogramDeclaration(&child)
		default:
			return nil, errors.New("unknown declaration section")
		}

		if err != nil {
			return nil, err
		}

		declarations = append(declarations, *declaration)
	}

	return declarations, nil
}

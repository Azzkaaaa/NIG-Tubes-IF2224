package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeSubprogramDeclaration(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.SUBPROGRAM_DECLARATION_NODE {
		return nil, errors.New("subprogram declaration expected")
	}

	child := parsetree.Children[0]

	switch child.RootType {
	case dt.PROCEDURE_DECLARATION_NODE:
		return a.analyzeProcedureDeclaration(&child)
	case dt.FUNCTION_DECLARATION_NODE:
		return a.analyzeFunctionDeclaration(&child)
	default:
		return nil, errors.New("expected procedure or function declaration")
	}
}

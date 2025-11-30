package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	switch parsetree.RootType {
	case dt.COMPOUND_STATEMENT_NODE:
		return a.analyzeCompoundStatement(parsetree)
	case dt.IF_STATEMENT_NODE:
		return a.analyzeIfStatement(parsetree)
	case dt.WHILE_STATEMENT_NODE:
		return a.analyzeWhileStatement(parsetree)
	case dt.FOR_STATEMENT_NODE:
		return a.analyzeForStatement(parsetree)
	case dt.ASSIGNMENT_STATEMENT_NODE:
		return a.analyzeAssignmentStatement(parsetree)
	case dt.SUBPROGRAM_CALL_NODE:
		dst, _, err := a.analyzeSubprogramCall(parsetree)
		return dst, err
	default:
		return nil, errors.New("unrecognized statement")
	}
}

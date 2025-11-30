package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeIdentifierList(parsetree *dt.ParseTree) ([]string, error) {
	if parsetree.RootType != dt.IDENTIFIER_LIST_NODE {
		return nil, errors.New("expected identifier list")
	}

	identifiers := make([]string, 0)

	for i := 0; i < len(parsetree.Children); i += 2 {
		if parsetree.Children[i].TokenValue != nil {
			identifiers = append(identifiers, parsetree.Children[i].TokenValue.Lexeme)
		}
	}

	return identifiers, nil
}

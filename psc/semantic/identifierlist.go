package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeIdentifierList(parsetree *dt.ParseTree) ([]string, error) {
	if parsetree.RootType != dt.IDENTIFIER_LIST_NODE {
		return nil, errors.New("expected identifier list")
	}

	identifierNodes := parsetree.Children[:len(parsetree.Children):2]
	identifiers := make([]string, len(identifierNodes))

	for i, v := range identifierNodes {
		identifiers[i] = v.TokenValue.Lexeme
	}

	return identifiers, nil
}

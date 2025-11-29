package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeForStatement(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.FOR_STATEMENT_NODE {
		return nil, errors.New("expected for block")
	}

	target, targetType, err := a.analyzeToken(&parsetree.Children[1])

	if err != nil {
		return nil, err
	}

	if target.SelfType != dt.DST_VARIABLE {
		return nil, errors.New("expected variable in for loop assignment")
	}

	initial, initialType, err := a.analyzeExpression(&parsetree.Children[3])

	if err != nil {
		return nil, err
	}

	if !a.checkTypeEquality(targetType, initialType) {
		return nil, errors.New("assigned expression type does not match variable type")
	}

	final, finalType, err := a.analyzeExpression(&parsetree.Children[5])

	if err != nil {
		return nil, err
	}

	if !a.checkTypeEquality(targetType, finalType) {
		return nil, errors.New("final expression type does not match variable type")
	}

	block, err := a.analyzeStatement(&parsetree.Children[7])

	if err != nil {
		return nil, err
	}

	target.Property = dt.DST_TARGET
	initial.Property = dt.DST_VALUE
	block.Property = dt.DST_EXECUTE

	switch parsetree.Children[4].TokenValue.Lexeme {
	case "ke":
		final.Property = dt.DST_UPTO
	case "turun_ke":
		final.Property = dt.DST_DOWNTO
	default:
		return nil, errors.New("expected 'ke' or 'turun_ke'")
	}

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_FOR_BLOCK,
		Children: []dt.DecoratedSyntaxTree{
			*initial,
			*target,
			*final,
			*block,
		},
	}, nil
}

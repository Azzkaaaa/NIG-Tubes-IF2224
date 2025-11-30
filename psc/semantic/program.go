package semantic

import (
	"errors"
	"fmt"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeProgram(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.PROGRAM_NODE {
		return nil, errors.New("expected program")
	}

	fmt.Println("\n[PROGRAM] Starting program analysis")
	a.printDebugState("PROGRAM_START")

	headerIndex, _, err := a.analyzeProgramHeader(&parsetree.Children[0])

	if err != nil {
		return nil, err
	}

	declarations, err := a.analyzeDeclarationPart(&parsetree.Children[1])

	if err != nil {
		return nil, err
	}

	block, err := a.analyzeCompoundStatement(&parsetree.Children[2])

	if err != nil {
		return nil, err
	}

	a.printDebugState("PROGRAM_END")

	return &dt.DecoratedSyntaxTree{
		Property: dt.DST_ROOT,
		SelfType: dt.DST_PROGRAM,
		Data:     headerIndex,
		Children: append(declarations, *block),
	}, nil
}

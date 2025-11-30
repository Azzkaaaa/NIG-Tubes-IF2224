package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeProcedureDeclaration(parsetree *dt.ParseTree) (*dt.DecoratedSyntaxTree, error) {
	if parsetree.RootType != dt.PROCEDURE_DECLARATION_NODE {
		return nil, errors.New("expected procedure declaration")
	}

	identifier := parsetree.Children[1].TokenValue.Lexeme

	_, check := a.tab.FindIdentifier(identifier, a.root)
	if check != nil {
		if check.Level == a.depth {
			return nil, errors.New("cannot redeclare identifier in the same scope")
		}
	}

	tabIndex := len(a.tab)
	a.tab = append(a.tab, dt.TabEntry{
		Identifier: identifier,
		Link:       a.root,
		Object:     dt.TAB_ENTRY_PROC,
		Level:      a.depth,
	})

	a.root = tabIndex

	root := a.root
	stackSize := a.stackSize
	a.stackSize = 0
	a.depth++

	paramSize := 0

	var parameters, block *dt.DecoratedSyntaxTree
	var declarations []dt.DecoratedSyntaxTree
	var err error

	for _, child := range parsetree.Children[2:] {
		switch child.RootType {
		case dt.FORMAL_PARAMETER_LIST_NODE:
			parameters, err = a.analyzeFormalParameterList(&child)
			paramSize = a.stackSize
		case dt.DECLARATION_PART_NODE:
			declarations, err = a.analyzeDeclarationPart(&child)
		case dt.COMPOUND_STATEMENT_NODE:
			block, err = a.analyzeCompoundStatement(&child)
		}

		if err != nil {
			return nil, err
		}
	}

	varSize := a.stackSize - paramSize

	start := 0
	paramEnd := 0
	varEnd := 0

	for _, part := range declarations {
		if part.SelfType == dt.DST_VARIABLE_DECLARATIONS && len(part.Children) != 0 {
			start = part.Children[0].Data
			varEnd = part.Children[len(part.Children)-1].Data
		}
	}

	if parameters != nil {
		start = parameters.Children[0].Data
		paramEnd = parameters.Children[len(parameters.Children)-1].Data
	}

	btabIndex := len(a.btab)
	a.btab = append(a.btab, dt.BtabEntry{
		Start:        start,
		End:          varEnd,
		ParamEnd:     paramEnd,
		ParamSize:    paramSize,
		VariableSize: varSize,
	})

	a.tab[tabIndex].Data = btabIndex

	a.depth--
	a.stackSize = stackSize
	a.root = root

	var children []dt.DecoratedSyntaxTree

	if parameters != nil {
		children = []dt.DecoratedSyntaxTree{*parameters}
	} else {
		children = []dt.DecoratedSyntaxTree{}
	}

	children = append(children, declarations...)
	children = append(children, *block)

	return &dt.DecoratedSyntaxTree{
		SelfType: dt.DST_PROCEDURE,
		Data:     tabIndex,
		Children: children,
	}, nil
}

package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeProgramHeader(parsetree *dt.ParseTree) (int, dt.TabEntry, error) {
	if parsetree.RootType != dt.PROGRAM_HEADER_NODE {
		return -1, dt.TabEntry{}, errors.New("expected program header")
	}

	identifier := parsetree.Children[1].TokenValue.Lexeme

	_, check := a.tab.FindIdentifier(identifier, a.root)
	if check != nil {
		if check.Level == a.depth {
			return -1, dt.TabEntry{}, errors.New("program name is reserved")
		}
	}

	tabIndex := len(a.tab)

	// Search for outer scope occurrence of same identifier
	linkIndex := -1
	if a.root != -1 && a.root < len(a.tab) {
		outerRoot := a.tab[a.root].Link
		if outerRoot != -1 {
			linkIndex, _ = a.tab.FindIdentifier(identifier, outerRoot)
		}
	}

	tabEntry := dt.TabEntry{
		Identifier: identifier,
		Link:       linkIndex,
		Object:     dt.TAB_ENTRY_PROGRAM,
		Type:       dt.TAB_ENTRY_NONE,
		Level:      a.depth,
	}

	a.tab = append(a.tab, tabEntry)

	a.root++

	return tabIndex, tabEntry, nil
}

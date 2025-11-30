package semantic

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeRecordType(parsetree *dt.ParseTree, identifier string) (int, dt.TabEntry, error) {
	if parsetree.RootType != dt.RECORD_TYPE_NODE {
		return -1, dt.TabEntry{}, errors.New("expected record type")
	}

	btabIndex := len(a.btab)

	entry := dt.TabEntry{
		Identifier: identifier,
		Link:       a.root,
		Object:     dt.TAB_ENTRY_TYPE,
		Type:       dt.TAB_ENTRY_RECORD,
		Reference:  btabIndex,
		Normal:     false,
		Level:      a.depth,
	}
	a.tab = append(a.tab, entry)
	a.root = len(a.tab) - 1

	btabEntry := dt.BtabEntry{
		Start:        len(a.tab),
		ParamEnd:     0,
		ReturnEnd:    0,
		End:          0,
		ParamSize:    0,
		ReturnSize:   0,
		VariableSize: 0,
	}

	oldRoot := a.root
	oldDepth := a.depth

	a.depth++

	for _, child := range parsetree.Children {
		if child.RootType != dt.VAR_DECLARATION_NODE {
			if child.RootType == dt.TOKEN_NODE {
				continue
			}
			return -1, dt.TabEntry{}, errors.New("expected a var declaration node for record field")
		}

		oldStackSize := a.stackSize
		a.stackSize = btabEntry.VariableSize

		_, err := a.analyzeVarDeclaration(&child)
		if err != nil {
			return -1, dt.TabEntry{}, err
		}

		for j := btabEntry.Start; j < len(a.tab); j++ {
			if a.tab[j].Object == dt.TAB_ENTRY_VAR {
				a.tab[j].Object = dt.TAB_ENTRY_FIELD
			}
		}
		btabEntry.VariableSize = a.stackSize
		a.stackSize = oldStackSize
	}

	btabEntry.End = len(a.tab) - 1

	totalSize := 0
	for i := btabEntry.Start; i < btabEntry.End; i++ {
		fieldType := semanticType{
			StaticType: a.tab[i].Type,
			Reference:  a.tab[i].Reference,
		}
		totalSize += a.getTypeSize(fieldType)
	}
	btabEntry.VariableSize = totalSize

	a.root = oldRoot
	a.depth = oldDepth

	a.btab = append(a.btab, btabEntry)

	return a.root, entry, nil
}

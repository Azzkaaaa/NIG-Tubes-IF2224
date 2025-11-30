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

	// Create a new block table entry for this record
	btabEntry := dt.BtabEntry{
		Start:        len(a.tab),
		ParamEnd:     0,
		ReturnEnd:    0,
		End:          0,
		ParamSize:    0,
		ReturnSize:   0,
		VariableSize: 0,
	}

	// Save current root and depth
	oldRoot := a.root
	oldDepth := a.depth

	// Set new scope for record fields
	// a.root = 0
	a.depth++

	// Process field declarations
	for _, child := range parsetree.Children {
		// Process variable declaration (field)
		if child.RootType != dt.VAR_DECLARATION_NODE {
			// Skip tokens like 'rekaman' and 'selesai'
			if child.RootType == dt.TOKEN_NODE {
				continue
			}
			return -1, dt.TabEntry{}, errors.New("expected a var declaration node for record field")
		}

		// Change object type to FIELD for record fields
		oldStackSize := a.stackSize
		a.stackSize = btabEntry.VariableSize

		_, err := a.analyzeVarDeclaration(&child)
		if err != nil {
			return -1, dt.TabEntry{}, err
		}

		// Update the last added entries to be fields instead of variables
		for j := btabEntry.Start; j < len(a.tab); j++ {
			if a.tab[j].Object == dt.TAB_ENTRY_VAR {
				a.tab[j].Object = dt.TAB_ENTRY_FIELD
			}
		}
		btabEntry.VariableSize = a.stackSize
		a.stackSize = oldStackSize
	}

	// Update btab entry with final values
	btabEntry.End = len(a.tab) - 1

	// Calculate total size of all fields
	totalSize := 0
	for i := btabEntry.Start; i < btabEntry.End; i++ {
		fieldType := semanticType{
			StaticType: a.tab[i].Type,
			Reference:  a.tab[i].Reference,
		}
		totalSize += a.getTypeSize(fieldType)
	}
	btabEntry.VariableSize = totalSize

	// Restore previous scope
	a.root = oldRoot
	a.depth = oldDepth

	// Add to btab
	a.btab = append(a.btab, btabEntry)

	// Return the record type entry
	return a.root, entry, nil
}

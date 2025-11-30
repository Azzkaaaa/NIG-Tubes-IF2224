package semantic

import (
	"strconv"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

type SemanticAnalyzer struct {
	parseTree *dt.ParseTree
	tab       dt.Tab
	atab      dt.Atab
	btab      dt.Btab
	strtab    dt.StrTab
	root      int
	depth     int
	stackSize int
}

type semanticType struct {
	StaticType dt.TabEntryType
	Reference  int
}

func New(parseTree *dt.ParseTree) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		parseTree: parseTree,
		tab: dt.Tab{
			dt.TabEntry{
				Identifier: "string",
				Link:       -1,
				Object:     dt.TAB_ENTRY_TYPE,
				Type:       dt.TAB_ENTRY_ARRAY,
				Reference:  0,
				Normal:     false,
				Level:      0,
				Data:       0,
			},
			dt.TabEntry{
				Identifier: "write",
				Link:       0,
				Object:     dt.TAB_ENTRY_PROC,
				Type:       dt.TAB_ENTRY_NONE,
				Reference:  0,
				Normal:     false,
				Level:      0,
				Data:       0,
			},
			dt.TabEntry{
				Identifier: "writeparam1",
				Link:       1,
				Object:     dt.TAB_ENTRY_PARAM,
				Type:       dt.TAB_ENTRY_ALIAS,
				Reference:  0,
				Normal:     false,
				Level:      1,
				Data:       0,
			},
			dt.TabEntry{
				Identifier: "writeparam2",
				Link:       2,
				Object:     dt.TAB_ENTRY_PARAM,
				Type:       dt.TAB_ENTRY_ALIAS,
				Reference:  0,
				Normal:     false,
				Level:      1,
				Data:       0,
			},
		},
		atab: dt.Atab{
			dt.AtabEntry{
				IndexType:        dt.TAB_ENTRY_INTEGER,
				ElementType:      dt.TAB_ENTRY_CHAR,
				ElementReference: 0,
				LowBound:         0,
				HighBound:        255,
				ElementSize:      1,
				TotalSize:        256,
			},
		},
		btab: dt.Btab{
			dt.BtabEntry{
				Start:        1,
				ParamEnd:     3,
				ReturnEnd:    3,
				End:          3,
				ParamSize:    512,
				ReturnSize:   0,
				VariableSize: 0,
			},
		},
		strtab:    make(dt.StrTab, 0),
		root:      3,
		depth:     0,
		stackSize: 0,
	}
}

func (a *SemanticAnalyzer) GetSymbols() (dt.Tab, dt.Atab, dt.Btab, dt.StrTab) {
	return a.tab, a.atab, a.btab, a.strtab
}

func (a *SemanticAnalyzer) Analyze() (dt.Tab, dt.Atab, dt.Btab, dt.StrTab, *dt.DecoratedSyntaxTree, error) { // ilangin switchcase
	dst, err := a.analyzeProgram(a.parseTree)
	tab, atab, btab, strtab := a.GetSymbols()
	return tab, atab, btab, strtab, dst, err
}

// resolveAliasType resolves a type alias to its underlying base type
// If the type is not an alias, it returns the type unchanged
func (a *SemanticAnalyzer) resolveAliasType(t semanticType) semanticType {
	// Keep resolving until we find a non-alias type
	for t.StaticType == dt.TAB_ENTRY_ALIAS {
		// Follow the reference chain
		t = semanticType{
			StaticType: a.tab[t.Reference].Type,
			Reference:  a.tab[t.Reference].Reference,
		}
	}
	return t
}

func (a *SemanticAnalyzer) checkTypeEquality(t1 semanticType, t2 semanticType) bool {
	// Resolve aliases to their base types before comparison
	resolved1 := a.resolveAliasType(t1)
	resolved2 := a.resolveAliasType(t2)

	if resolved1.StaticType == resolved2.StaticType {
		switch resolved1.StaticType {
		case dt.TAB_ENTRY_NONE:
			fallthrough
		case dt.TAB_ENTRY_INTEGER:
			fallthrough
		case dt.TAB_ENTRY_REAL:
			fallthrough
		case dt.TAB_ENTRY_BOOLEAN:
			fallthrough
		case dt.TAB_ENTRY_CHAR:
			return true
		case dt.TAB_ENTRY_RECORD:
			fallthrough
		case dt.TAB_ENTRY_ARRAY:
			return resolved1.Reference == resolved2.Reference
		}
	}

	return false
}

func (a *SemanticAnalyzer) getTypeSize(t semanticType) int {
	switch t.StaticType {
	case dt.TAB_ENTRY_INTEGER:
		return strconv.IntSize
	case dt.TAB_ENTRY_REAL:
		return strconv.IntSize
	case dt.TAB_ENTRY_BOOLEAN:
		return 1
	case dt.TAB_ENTRY_CHAR:
		return 1
	case dt.TAB_ENTRY_RECORD:
		return a.btab[t.Reference].VariableSize
	case dt.TAB_ENTRY_ARRAY:
		return a.atab[t.Reference].TotalSize
	case dt.TAB_ENTRY_ALIAS:
		return a.getTypeSize(semanticType{
			StaticType: a.tab[t.Reference].Type,
			Reference:  a.tab[t.Reference].Reference,
		})
	default:
		return 0
	}
}

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
		tab:       make(dt.Tab, 0),
		atab:      make(dt.Atab, 0),
		btab:      make(dt.Btab, 0),
		strtab:    make(dt.StrTab, 0),
		root:      0,
		depth:     0,
		stackSize: 0,
	}
}

func (a *SemanticAnalyzer) GetSymbols() (dt.Tab, dt.Atab, dt.Btab, dt.StrTab) {
	return a.tab, a.atab, a.btab, a.strtab
}

func (a *SemanticAnalyzer) Analyze() (dt.Tab, dt.Atab, dt.Btab, dt.StrTab, *dt.DecoratedSyntaxTree, error) {
	dst, err := a.analyzeProgram(a.parseTree)
	tab, atab, btab, strtab := a.GetSymbols()
	return tab, atab, btab, strtab, dst, err
}

func (a *SemanticAnalyzer) checkTypeEquality(t1 semanticType, t2 semanticType) bool {
	if t1.StaticType == t2.StaticType {
		switch t1.StaticType {
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
			fallthrough
		case dt.TAB_ENTRY_ALIAS:
			return t1.Reference == t2.Reference
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

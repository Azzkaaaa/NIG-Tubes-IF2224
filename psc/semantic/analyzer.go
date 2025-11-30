package semantic

import (
	"fmt"
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
	fmt.Println("Starting semantic analysis...")
	dst, err := a.analyzeProgram(a.parseTree)
	tab, atab, btab, strtab := a.GetSymbols()
	return tab, atab, btab, strtab, dst, err
}

func (a *SemanticAnalyzer) printDebugState(label string) {
	fmt.Printf("\n=== DEBUG STATE: %s (Depth=%d, Root=%d, StackSize=%d) ===\n", label, a.depth, a.root, a.stackSize)

	// Print Tab (Symbol Table)
	fmt.Println("\n[TAB] Symbol Table:")
	if len(a.tab) == 0 {
		fmt.Println("  <empty>")
	} else {
		fmt.Println("  Idx | Identifier      | Link | Object   | Type     | Ref | Level | Data")
		fmt.Println("  ----|-----------------|------|----------|----------|-----|-------|-----")
		for i, entry := range a.tab {
			fmt.Printf("  %3d | %-15s | %4d | %-8s | %-8s | %3d | %5d | %4d\n",
				i, entry.Identifier, entry.Link, entry.Object, entry.Type, entry.Reference, entry.Level, entry.Data)
		}
	}

	// Print Atab (Array Table)
	fmt.Println("\n[ATAB] Array Table:")
	if len(a.atab) == 0 {
		fmt.Println("  <empty>")
	} else {
		fmt.Println("  Idx | IndexType | ElementType | LowBound | HighBound | ElementSize | TotalSize")
		fmt.Println("  ----|-----------|-------------|----------|-----------|-------------|----------")
		for i, entry := range a.atab {
			fmt.Printf("  %3d | %9d | %11d | %8d | %9d | %11d | %9d\n",
				i, entry.IndexType, entry.ElementType, entry.LowBound, entry.HighBound, entry.ElementSize, entry.TotalSize)
		}
	}

	// Print Btab (Block Table)
	fmt.Println("\n[BTAB] Block Table:")
	if len(a.btab) == 0 {
		fmt.Println("  <empty>")
	} else {
		fmt.Println("  Idx | Start | End  | VariableSize")
		fmt.Println("  ----|-------|------|-------------")
		for i, entry := range a.btab {
			fmt.Printf("  %3d | %5d | %4d | %12d\n", i, entry.Start, entry.End, entry.VariableSize)
		}
	}

	// Print StrTab (String Table)
	fmt.Println("\n[STRTAB] String Table:")
	if len(a.strtab) == 0 {
		fmt.Println("  <empty>")
	} else {
		fmt.Println("  Idx | Value")
		fmt.Println("  ----|-----")
		for i, entry := range a.strtab {
			fmt.Printf("  %3d | %s\n", i, entry)
		}
	}
	fmt.Println("=== END DEBUG STATE ===\n")
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

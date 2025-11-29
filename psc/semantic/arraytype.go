package semantic

import (
	"errors"
	"strconv"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) analyzeArrayType(parsetree *dt.ParseTree) (int, dt.TabEntry, error) {
	if parsetree.RootType != dt.ARRAY_TYPE_NODE {
		return -1, dt.TabEntry{}, errors.New("expected array type")
	}

	begin, end, indexType, err := a.analyzeRange(&parsetree.Children[2])

	if err != nil {
		return -1, dt.TabEntry{}, err
	}

	if indexType.StaticType == dt.TAB_ENTRY_REAL {
		return -1, dt.TabEntry{}, errors.New("cannot use real value as array index")
	}

	_, tabEntry, err := a.analyzeType(&parsetree.Children[5])

	if err != nil {
		return -1, dt.TabEntry{}, err
	}

	atabEntry := dt.AtabEntry{
		IndexType:        indexType.StaticType,
		ElementType:      tabEntry.Type,
		ElementReference: tabEntry.Reference,
		LowBound:         begin,
		HighBound:        end,
	}

	atabIdx, _ := a.atab.FindArray(atabEntry)

	if atabIdx == -1 {
		atabIdx = len(a.atab)

		switch tabEntry.Type {
		case dt.TAB_ENTRY_ARRAY:
			atabEntry.ElementSize = a.atab[tabEntry.Reference].TotalSize
		case dt.TAB_ENTRY_RECORD:
			atabEntry.ElementSize = a.btab[tabEntry.Reference].VariableSize
		case dt.TAB_ENTRY_BOOLEAN:
			atabEntry.ElementSize = 1
		case dt.TAB_ENTRY_CHAR:
			atabEntry.ElementSize = 1
		case dt.TAB_ENTRY_INTEGER:
			atabEntry.ElementSize = strconv.IntSize
		case dt.TAB_ENTRY_REAL:
			atabEntry.ElementSize = strconv.IntSize
		default:
			return -1, dt.TabEntry{}, errors.New("unknown type")
		}

		atabEntry.TotalSize = atabEntry.ElementSize * (end - begin + 1)

		a.atab = append(a.atab, atabEntry)
	}

	return -1, dt.TabEntry{
		Type:      dt.TAB_ENTRY_ARRAY,
		Reference: atabIdx,
	}, nil
}

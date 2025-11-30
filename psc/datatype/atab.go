package datatype

import (
	"fmt"
)

type AtabEntry struct {
	IndexType        TabEntryType
	ElementType      TabEntryType
	ElementReference int
	LowBound         int
	HighBound        int
	ElementSize      int
	TotalSize        int
}

type Atab []AtabEntry

func (t Atab) FindArray(entry AtabEntry) (int, *AtabEntry) {
	for i, v := range t {
		if v.IndexType != entry.IndexType {
			continue
		}

		if v.ElementType != entry.ElementType {
			continue
		}

		if v.ElementReference != entry.ElementReference {
			switch v.ElementType {
			case TAB_ENTRY_ARRAY:
				fallthrough
			case TAB_ENTRY_RECORD:
				continue
			}
		}

		if v.LowBound != entry.LowBound {
			continue
		}

		if v.HighBound != entry.HighBound {
			continue
		}

		return i, &v
	}

	return -1, nil
}

func (t Atab) String() string {
	if len(t) == 0 {
		return "<empty array table>"
	}

	out := "Idx  IdxType      ElemType     ElemRef  Low   High  ElemSize  TotalSize\n"
	out += "---- ------------ ------------ -------- ----- ----- --------- ----------\n"

	for i, e := range t {
		out += fmt.Sprintf(
			"%-4d %-12s %-12s %-8d %-5d %-5d %-9d %-10d\n",
			i,
			e.IndexType.String(),
			e.ElementType.String(),
			e.ElementReference,
			e.LowBound,
			e.HighBound,
			e.ElementSize,
			e.TotalSize,
		)
	}

	return out
}

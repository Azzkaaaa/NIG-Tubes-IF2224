package datatype

import (
	"fmt"
)

type AtabEntryType int

const (
	ATAB_ENTRY_INTEGER AtabEntryType = iota
	ATAB_ENTRY_REAL
	ATAB_ENTRY_BOOLEAN
	ATAB_ENTRY_CHAR
	ATAB_ENTRY_ARRAY
	ATAB_ENTRY_RECORD
)

func (o AtabEntryType) String() string {
	names := []string{
		"integer",
		"real",
		"boolean",
		"char",
		"array",
		"record",
	}

	if int(o) < 0 || int(o) > len(names) {
		return "unknown"
	}

	return names[o]
}

type AtabEntry struct {
	Index            int
	IndexType        AtabEntryType
	ElementType      AtabEntryType
	ElementReference int
	LowBound         int
	HighBound        int
	ElementSize      int
	TotalSize        int
}

type Atab []AtabEntry

func (t Atab) String() string {
	if len(t) == 0 {
		return "<empty array table>"
	}

	out := "Idx  IdxType      ElemType     ElemRef  Low   High  ElemSize  TotalSize\n"
	out += "---- ------------ ------------ -------- ----- ----- --------- ----------\n"

	for _, e := range t {
		out += fmt.Sprintf(
			"%-4d %-12s %-12s %-8d %-5d %-5d %-9d %-10d\n",
			e.Index,
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

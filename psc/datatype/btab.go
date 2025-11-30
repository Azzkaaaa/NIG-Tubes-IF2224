package datatype

import (
	"fmt"
)

type BtabEntry struct {
	Start        int
	ParamEnd     int
	ReturnEnd    int
	End          int
	ParamSize    int
	ReturnSize   int
	VariableSize int
}

type Btab []BtabEntry

func (t Btab) String() string {
	if len(t) == 0 {
		return "<empty block table>"
	}

	out := "Idx  Start  End    ParamEnd  ReturnEnd  ParamSize  ReturnSize  VarSize\n"
	out += "---- ------ ------ --------- ---------- ---------- ----------- --------\n"

	for i, e := range t {
		out += fmt.Sprintf(
			"%-4d %-6d %-6d %-9d %-10d %-10d %-11d %-8d\n",
			i,
			e.Start,
			e.End,
			e.ParamEnd,
			e.ReturnEnd,
			e.ParamSize,
			e.ReturnSize,
			e.VariableSize,
		)
	}

	return out
}

package datatype

import (
	"errors"
	"fmt"
)

type TabEntryObject int

const (
	TAB_ENTRY_PROGRAM TabEntryObject = iota
	TAB_ENTRY_CONST
	TAB_ENTRY_VAR
	TAB_ENTRY_TYPE
	TAB_ENTRY_PROC
	TAB_ENTRY_FUNC
)

func (o TabEntryObject) String() string {
	names := []string{
		"program",
		"constant",
		"variable",
		"type",
		"procedure",
		"function",
	}

	if int(o) < 0 || int(o) > len(names) {
		return "unknown"
	}

	return names[o]
}

type TabEntryType int

const (
	TAB_ENTRY_INTEGER TabEntryType = iota
	TAB_ENTRY_REAL
	TAB_ENTRY_BOOLEAN
	TAB_ENTRY_CHAR
	TAB_ENTRY_ARRAY
	TAB_ENTRY_RECORD
)

func (o TabEntryType) String() string {
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

type TabEntry struct {
	Index      int
	Identifier string
	Link       int
	Object     TabEntryObject
	Type       TabEntryType
	Reference  int
	Normal     bool
	Level      int
	Data       int
}

type Tab []TabEntry

func (t Tab) GetByIdentifier(identifier string) (TabEntry, error) {
	for i := 0; i < len(t); i++ {
		if t[i].Identifier == identifier {
			return t[i], nil
		}
	}

	return TabEntry{}, errors.New("index out of range")
}

func (t Tab) String() string {
	if len(t) == 0 {
		return "<empty symbol table>"
	}

	// Table header
	out := "Idx  Identifier       Link  Object        Type          Ref   Norm  Level  Data\n"
	out += "---- ---------------- ----- ------------- ------------- ----- ----- ------ -----\n"

	for _, e := range t {
		identifier := e.Identifier

		if len(identifier) > 16 {
			identifier = identifier[:13] + "..."
		}

		out += fmt.Sprintf(
			"%-4d %-16s %-5d %-13s %-13s %-5d %-5t %-6d %-5d\n",
			e.Index,
			identifier,
			e.Link,
			e.Object.String(),
			e.Type.String(),
			e.Reference,
			e.Normal,
			e.Level,
			e.Data,
		)
	}

	return out
}

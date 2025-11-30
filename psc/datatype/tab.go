package datatype

import (
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
	TAB_ENTRY_PARAM
	TAB_ENTRY_FIELD
	TAB_ENTRY_RETURN
)

func (o TabEntryObject) String() string {
	names := []string{
		"program",
		"constant",
		"variable",
		"type",
		"procedure",
		"function",
		"parameter",
		"field",
		"return",
	}

	if int(o) < 0 || int(o) > len(names) {
		return "unknown"
	}

	return names[o]
}

type TabEntryType int

const (
	TAB_ENTRY_NONE TabEntryType = iota
	TAB_ENTRY_INTEGER
	TAB_ENTRY_REAL
	TAB_ENTRY_BOOLEAN
	TAB_ENTRY_CHAR
	TAB_ENTRY_ARRAY
	TAB_ENTRY_RECORD
	TAB_ENTRY_ALIAS
)

func (o TabEntryType) String() string {
	names := []string{
		"none",
		"integer",
		"real",
		"boolean",
		"char",
		"array",
		"record",
		"alias",
	}

	if int(o) < 0 || int(o) > len(names) {
		return "unknown"
	}

	return names[o]
}

type TabEntry struct {
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

func (t *Tab) FindIdentifier(id string, start int) (int, *TabEntry) {
	current := start

	if current == -1 || len(*t) == 0 {
		return -1, nil
	}

	for i := 0; current != -1 && i < 100; i++ {
		if current >= len(*t) {
			return -1, nil
		}

		if (*t)[current].Identifier == id {
			return current, &(*t)[current]
		}
		current = (*t)[current].Link
	}

	return -1, nil
}

func (t Tab) String() string {
	if len(t) == 0 {
		return "<empty symbol table>"
	}

	out := "Idx  Identifier       Link  Object        Type          Ref   Norm  Level  Data\n"
	out += "---- ---------------- ----- ------------- ------------- ----- ----- ------ -----\n"

	for i, e := range t {
		identifier := e.Identifier

		if len(identifier) > 16 {
			identifier = identifier[:13] + "..."
		}

		out += fmt.Sprintf(
			"%-4d %-16s %-5d %-13s %-13s %-5d %-5t %-6d %-5d\n",
			i,
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

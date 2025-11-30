package datatype

import (
	"fmt"
	"strings"
)

// String returns a string representation of the DST
func (t DecoratedSyntaxTree) String() string {
	var sb strings.Builder
	t.writeString(&sb, "", true, true, nil, nil, nil, nil)
	return sb.String()
}

// StringWithSymbols returns a string representation of the DST with symbol table information
func (t DecoratedSyntaxTree) StringWithSymbols(tab Tab, atab Atab, btab Btab, strtab StrTab) string {
	var sb strings.Builder
	t.writeString(&sb, "", true, true, &tab, &atab, &btab, &strtab)
	return sb.String()
}

// formatDataWithSymbols formats the Data field with symbol table information
func formatDataWithSymbols(nodeType DSTNodeType, data int, tab *Tab, atab *Atab, btab *Btab, strtab *StrTab) string {
	if tab == nil {
		return fmt.Sprintf(" (%d)", data)
	}

	switch nodeType {
	case DST_CHAR_LITERAL:
		// Display as character
		return fmt.Sprintf(": '%c'", rune(data))

	case DST_STR_LITERAL:
		// Display string from strtab
		if strtab != nil && data >= 0 && data < len(*strtab) {
			return fmt.Sprintf(": \"%s\"", (*strtab)[data].String)
		}
		return fmt.Sprintf(" (strtab[%d])", data)

	case DST_INT_LITERAL:
		// Display integer value
		return fmt.Sprintf(": %d", data)

	case DST_BOOL_LITERAL:
		// Display boolean value
		if data == 1 {
			return ": true"
		}
		return ": false"

	case DST_CAST_OPERATOR:
		// Display target type
		typeName := TabEntryType(data).String()
		return fmt.Sprintf(": to %s", typeName)

	case DST_VARIABLE, DST_CONST, DST_TYPE:
		// Display identifier name from symbol table
		if data >= 0 && data < len(*tab) {
			return fmt.Sprintf(": %s (tab[%d])", (*tab)[data].Identifier, data)
		}
		return fmt.Sprintf(" (tab[%d])", data)

	case DST_FUNCTION_CALL, DST_PROCEDURE_CALL:
		// Display subprogram name
		if data >= 0 && data < len(*tab) {
			return fmt.Sprintf(": %s (tab[%d])", (*tab)[data].Identifier, data)
		}
		return fmt.Sprintf(" (tab[%d])", data)

	case DST_ARRAY_ELEMENT, DST_RECORD_FIELD:
		// Display identifier name
		if data >= 0 && data < len(*tab) {
			return fmt.Sprintf(": %s (tab[%d])", (*tab)[data].Identifier, data)
		}
		return fmt.Sprintf(" (tab[%d])", data)

	case DST_FUNCTION, DST_PROCEDURE, DST_PROGRAM:
		// Display identifier name
		if data >= 0 && data < len(*tab) {
			return fmt.Sprintf(": %s (tab[%d])", (*tab)[data].Identifier, data)
		}
		return fmt.Sprintf(" (tab[%d])", data)

	default:
		// For other node types, just show the data value
		if data != 0 {
			return fmt.Sprintf(" (%d)", data)
		}
		return ""
	}
}

// writeString recursively writes the tree structure to a string builder
func (t DecoratedSyntaxTree) writeString(sb *strings.Builder, prefix string, isLast bool, isRoot bool, tab *Tab, atab *Atab, btab *Btab, strtab *StrTab) {
	if prefix != "" {
		if isLast {
			sb.WriteString(prefix + "└─")
		} else {
			sb.WriteString(prefix + "├─")
		}
	}

	// Write in format "property: type"
	if t.Property != DST_ROOT {
		sb.WriteString(t.Property.String())
		sb.WriteString(": ")
	}
	sb.WriteString(t.SelfType.String())

	// Write data with symbol information if available
	dataStr := formatDataWithSymbols(t.SelfType, t.Data, tab, atab, btab, strtab)
	sb.WriteString(dataStr)

	sb.WriteString("\n")

	// Write children
	childPrefix := prefix
	if isRoot || prefix != "" {
		if isLast {
			childPrefix += "  "
		} else {
			childPrefix += "│ "
		}
	}

	for i, child := range t.Children {
		child.writeString(sb, childPrefix, i == len(t.Children)-1, false, tab, atab, btab, strtab)
	}
}

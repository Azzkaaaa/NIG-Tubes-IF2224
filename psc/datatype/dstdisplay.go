package datatype

import (
	"fmt"
	"strings"
)

// String returns a string representation of the DST
func (t DecoratedSyntaxTree) String() string {
	var sb strings.Builder
	t.writeString(&sb, "", true, true)
	return sb.String()
}

// writeString recursively writes the tree structure to a string builder
func (t DecoratedSyntaxTree) writeString(sb *strings.Builder, prefix string, isLast bool, isRoot bool) {
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

	// Write data if non-zero
	if t.Data != 0 {
		sb.WriteString(fmt.Sprintf(" (%d)", t.Data))
	}

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
		child.writeString(sb, childPrefix, i == len(t.Children)-1, false)
	}
}

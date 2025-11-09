package datatype

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	PROGRAM_NODE NodeType = iota
	PROGRAM_HEADER_NODE
	DECLARATION_PART_NODE
	CONST_DECLARATION_NODE
	TYPE_DECLARATION_NODE
	VAR_DECLARATION_NODE
	IDENTIFIER_LIST_NODE
	TYPE_NODE
	ARRAY_TYPE_NODE
	RANGE_NODE
	SUBPROGRAM_DECLARATION_NODE
	PROCEDURE_DECLARATION_NODE
	FUNCTION_DECLARATION_NODE
	FORMAL_PARAMETER_LIST_NODE
	COMPOUND_STATEMENT_NODE
	STATEMENT_LIST_NODE
	ASSIGNMENT_STATEMENT_NODE
	IF_STATEMENT_NODE
	WHILE_STATEMENT_NODE
	FOR_STATEMENT_NODE
	SUBPROGRAM_CALL_NODE
	PARAMETER_LIST_NODE
	EXPRESSION_NODE
	SIMPLE_EXPRESSION_NODE
	TERM_NODE
	FACTOR_NODE
	RELATIONAL_OPERATOR_NODE
	ADDITIVE_OPERATOR_NODE
	MULTIPLICATIVE_OPERATOR_NODE
	TOKEN_NODE
)

type ParseTree struct {
	RootType   NodeType
	TokenValue *Token
	Children   []ParseTree
}

func (t NodeType) String() string {
	names := [...]string{
		"<program>",
		"<program-header>",
		"<declaration-part>",
		"<const-declaration>",
		"<type-declaration>",
		"<var-declaration>",
		"<identifier-list>",
		"<type>",
		"<array-type>",
		"<range>",
		"<subprogram-declaration>",
		"<procedure-declaration>",
		"<function-declaration>",
		"<formal-parameter-list>",
		"<compound-statement>",
		"<statement-list>",
		"<assignment-statement>",
		"<if-statement>",
		"<while-statement>",
		"<for-statement>",
		"<procedure/function-call>",
		"<parameter-list>",
		"<expression>",
		"<simple-expression>",
		"<term>",
		"<factor>",
		"<relational-operator>",
		"<additive-operator>",
		"<multiplicative-operator>",
	}

	if int(t) < 0 || int(t) >= len(names) {
		return "UNKNOWN"
	}

	return names[t]
}

func (t ParseTree) String() string {
	var sb strings.Builder
	t.writeString(&sb, 0)
	return sb.String()
}

func (t ParseTree) writeString(sb *strings.Builder, indent int) {
	prefix := strings.Repeat("  ", indent)
	sb.WriteString(prefix)
	if t.RootType != TOKEN_NODE {
		sb.WriteString(t.RootType.String())
	}

	if t.TokenValue != nil {
		sb.WriteString(t.TokenValue.Type.String())
		sb.WriteString(fmt.Sprintf(" (%s)", t.TokenValue.Lexeme))
	}

	sb.WriteString("\n")

	for _, child := range t.Children {
		child.writeString(sb, indent+1)
	}
}

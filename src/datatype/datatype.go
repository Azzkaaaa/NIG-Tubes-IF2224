package datatype

type TokenType int

const (
	KEYWORD TokenType = iota
	IDENTIFIER
	ARITHMETIC_OPERATOR
	RELATIONAL_OPERATOR
	LOGICAL_OPERATOR
	ASSIGN_OPERATOR
	NUMBER
	CHAR_LITERAL
	STRING_LITERAL
	SEMICOLON
	COMMA
	COLON
	DOT
	LPARENTHESIS
	RPARENTHESIS
	LBRACKET
	RBRACKET
	RANGE_OPERATOR
	COMMENT
)

type Token struct {
	Type   TokenType
	Lexeme string // substring aslinya
	Line   int    // posisi awal token (untuk error/report)
	Col    int
}

func (t TokenType) String() string {
	names := [...]string{
		"KEYWORD", "IDENTIFIER", "ARITHMETIC_OPERATOR", "RELATIONAL_OPERATOR", "LOGICAL_OPERATOR",
		"ASSIGN_OPERATOR", "NUMBER", "CHAR_LITERAL", "STRING_LITERAL", "SEMICOLON", "COMMA", "COLON",
		"DOT", "LPARENTHESIS", "RPARENTHESIS", "LBRACKET", "RBRACKET", "RANGE_OPERATOR",
		"COMMENT",
	}
	if int(t) < 0 || int(t) >= len(names) {
		return "UNKNOWN"
	}
	return names[t]
}

var Keywords = map[string]struct{}{
	"program": {}, "var": {}, "begin": {}, "end": {}, "if": {}, "then": {}, "else": {},
	"while": {}, "do": {}, "for": {}, "to": {}, "downto": {},
	"integer": {}, "real": {}, "boolean": {}, "char": {},
	"array": {}, "of": {}, "procedure": {}, "function": {}, "const": {}, "type": {},
}

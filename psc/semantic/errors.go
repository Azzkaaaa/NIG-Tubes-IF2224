package semantic

import (
	"fmt"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

// SemanticError represents a semantic analysis error with location information
type SemanticError struct {
	Message string
	Line    int
	Column  int
	Context string
	Token   *dt.Token
}

// Error implements the error interface
func (e *SemanticError) Error() string {
	if e.Token != nil {
		return fmt.Sprintf("Semantic error at line %d, column %d: %s\nContext: %s\nNear: '%s'",
			e.Line, e.Column, e.Message, e.Context, e.Token.Lexeme)
	}
	return fmt.Sprintf("Semantic error at line %d, column %d: %s\nContext: %s",
		e.Line, e.Column, e.Message, e.Context)
}

// NewSemanticError creates a new semantic error with location information
func NewSemanticError(message string, token *dt.Token, context string) *SemanticError {
	line := 0
	column := 0
	if token != nil {
		line = token.Line
		column = token.Col
	}
	return &SemanticError{
		Message: message,
		Line:    line,
		Column:  column,
		Context: context,
		Token:   token,
	}
}

// Error type constants for common semantic errors
const (
	ErrRedeclaration      = "identifier already declared in this scope"
	ErrUndeclaredIdent    = "undeclared identifier"
	ErrTypeMismatch       = "type mismatch"
	ErrIncompatibleTypes  = "incompatible types"
	ErrInvalidOperation   = "invalid operation"
	ErrConstantExpected   = "constant expression expected"
	ErrArrayBoundsInvalid = "invalid array bounds"
	ErrDivisionByZero     = "division by zero"
	ErrNotAFunction       = "identifier is not a function"
	ErrNotAProcedure      = "identifier is not a procedure"
	ErrWrongArgCount      = "wrong number of arguments"
	ErrCannotAssign       = "cannot assign to this expression"
)

// Helper functions to create specific error types

func (a *SemanticAnalyzer) newRedeclarationError(identifier string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("%s: '%s'", ErrRedeclaration, identifier),
		token,
		"identifier declaration",
	)
}

func (a *SemanticAnalyzer) newUndeclaredIdentError(identifier string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("%s: '%s'", ErrUndeclaredIdent, identifier),
		token,
		"identifier reference",
	)
}

func (a *SemanticAnalyzer) newTypeMismatchError(expected string, got string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("%s: expected %s, got %s", ErrTypeMismatch, expected, got),
		token,
		"type checking",
	)
}

func (a *SemanticAnalyzer) newIncompatibleTypesError(type1 string, type2 string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("%s: %s and %s", ErrIncompatibleTypes, type1, type2),
		token,
		"type compatibility check",
	)
}

func (a *SemanticAnalyzer) newConstantExpectedError(token *dt.Token) error {
	return NewSemanticError(
		ErrConstantExpected,
		token,
		"static evaluation",
	)
}

func (a *SemanticAnalyzer) newArrayBoundsError(token *dt.Token) error {
	return NewSemanticError(
		ErrArrayBoundsInvalid,
		token,
		"array type declaration",
	)
}

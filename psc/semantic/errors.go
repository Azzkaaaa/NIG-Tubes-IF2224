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

func (a *SemanticAnalyzer) newParameterCountError(expected int, got int, funcName string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("parameter count mismatch for '%s': expected %d, got %d", funcName, expected, got),
		token,
		"subprogram call",
	)
}

func (a *SemanticAnalyzer) newParameterTypeError(paramIndex int, expected string, got string, funcName string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("parameter %d type mismatch for '%s': expected %s, got %s", paramIndex+1, funcName, expected, got),
		token,
		"subprogram call",
	)
}

func (a *SemanticAnalyzer) newNotCallableError(identifier string, actualType string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("'%s' is not callable (it is a %s)", identifier, actualType),
		token,
		"subprogram call",
	)
}

func (a *SemanticAnalyzer) newInvalidArrayAccessError(identifier string, actualType string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("cannot index '%s': not an array (type is %s)", identifier, actualType),
		token,
		"array access",
	)
}

func (a *SemanticAnalyzer) newInvalidRecordAccessError(identifier string, actualType string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("cannot access field of '%s': not a record (type is %s)", identifier, actualType),
		token,
		"record field access",
	)
}

func (a *SemanticAnalyzer) newUndeclaredFieldError(fieldName string, recordName string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("record '%s' has no field named '%s'", recordName, fieldName),
		token,
		"record field access",
	)
}

func (a *SemanticAnalyzer) newOperatorTypeError(operator string, leftType string, rightType string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("operator '%s' cannot be applied to types %s and %s", operator, leftType, rightType),
		token,
		"operator type checking",
	)
}

func (a *SemanticAnalyzer) newInvalidTypeError(identifier string, expectedKind string, actualKind string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("'%s' is not a %s (it is a %s)", identifier, expectedKind, actualKind),
		token,
		"type checking",
	)
}

func (a *SemanticAnalyzer) newConditionTypeError(actualType string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("condition must be boolean, got %s", actualType),
		token,
		"control flow statement",
	)
}

func (a *SemanticAnalyzer) newAssignmentError(targetType string, valueType string, token *dt.Token) error {
	return NewSemanticError(
		fmt.Sprintf("cannot assign %s to %s", valueType, targetType),
		token,
		"assignment statement",
	)
}

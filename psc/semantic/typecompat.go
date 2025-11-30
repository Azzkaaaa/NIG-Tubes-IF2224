package semantic

import (
	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

// checkTypeCompatibility checks if two types are compatible for operations
// Returns true if types match exactly or if an implicit cast is possible
func (a *SemanticAnalyzer) checkTypeCompatibility(t1 semanticType, t2 semanticType) bool {
	// Exact match
	if a.checkTypeEquality(t1, t2) {
		return true
	}

	// Integer can be implicitly cast to real
	if t1.StaticType == dt.TAB_ENTRY_REAL && t2.StaticType == dt.TAB_ENTRY_INTEGER {
		return true
	}
	if t1.StaticType == dt.TAB_ENTRY_INTEGER && t2.StaticType == dt.TAB_ENTRY_REAL {
		return true
	}

	// Char can be compared with char
	if t1.StaticType == dt.TAB_ENTRY_CHAR && t2.StaticType == dt.TAB_ENTRY_CHAR {
		return true
	}

	// Boolean can be compared with boolean
	if t1.StaticType == dt.TAB_ENTRY_BOOLEAN && t2.StaticType == dt.TAB_ENTRY_BOOLEAN {
		return true
	}

	return false
}

// insertImplicitCast inserts a cast node if needed
// Casts FROM right type TO left type (kanan ke kiri)
// Returns the potentially modified DST and the result type
func (a *SemanticAnalyzer) insertImplicitCast(dst *dt.DecoratedSyntaxTree, fromType semanticType, toType semanticType) (*dt.DecoratedSyntaxTree, semanticType) {
	// No cast needed if types are equal
	if a.checkTypeEquality(fromType, toType) {
		return dst, fromType
	}

	// Cast integer to real (when assigning int to real variable)
	if fromType.StaticType == dt.TAB_ENTRY_INTEGER && toType.StaticType == dt.TAB_ENTRY_REAL {
		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_CAST_OPERATOR,
			Data:     int(dt.TAB_ENTRY_REAL), // Target type
			Children: []dt.DecoratedSyntaxTree{*dst},
		}, toType
	}

	// No implicit cast for other combinations
	return dst, fromType
}

// canCastImplicitly checks if an implicit cast from one type to another is allowed
func (a *SemanticAnalyzer) canCastImplicitly(fromType semanticType, toType semanticType) bool {
	// Same type - no cast needed
	if a.checkTypeEquality(fromType, toType) {
		return true
	}

	// Integer to real is allowed
	if fromType.StaticType == dt.TAB_ENTRY_INTEGER && toType.StaticType == dt.TAB_ENTRY_REAL {
		return true
	}

	return false
}

// promoteTypes promotes two types to a common type if possible
// Returns the promoted DSTs and the common type
func (a *SemanticAnalyzer) promoteTypes(dst1 *dt.DecoratedSyntaxTree, type1 semanticType, dst2 *dt.DecoratedSyntaxTree, type2 semanticType) (*dt.DecoratedSyntaxTree, *dt.DecoratedSyntaxTree, semanticType, bool) {
	// If types are equal, no promotion needed
	if a.checkTypeEquality(type1, type2) {
		return dst1, dst2, type1, true
	}

	// If one is integer and the other is real, promote integer to real
	if type1.StaticType == dt.TAB_ENTRY_INTEGER && type2.StaticType == dt.TAB_ENTRY_REAL {
		promoted1, _ := a.insertImplicitCast(dst1, type1, type2)
		return promoted1, dst2, type2, true
	}

	if type1.StaticType == dt.TAB_ENTRY_REAL && type2.StaticType == dt.TAB_ENTRY_INTEGER {
		promoted2, _ := a.insertImplicitCast(dst2, type2, type1)
		return dst1, promoted2, type1, true
	}

	// Char types are compatible
	if type1.StaticType == dt.TAB_ENTRY_CHAR && type2.StaticType == dt.TAB_ENTRY_CHAR {
		return dst1, dst2, type1, true
	}

	// Boolean types are compatible
	if type1.StaticType == dt.TAB_ENTRY_BOOLEAN && type2.StaticType == dt.TAB_ENTRY_BOOLEAN {
		return dst1, dst2, type1, true
	}

	// Types are incompatible
	return dst1, dst2, type1, false
}

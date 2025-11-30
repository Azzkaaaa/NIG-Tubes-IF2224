package semantic

import (
	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

// checkTypeCompatibility checks if two types are compatible for operations
// Returns true if types match exactly or if an implicit cast is possible
func (a *SemanticAnalyzer) checkTypeCompatibility(t1 semanticType, t2 semanticType) bool {
	// Exact match (this already resolves aliases internally)
	if a.checkTypeEquality(t1, t2) {
		return true
	}

	// Resolve aliases for compatibility checking
	resolved1 := a.resolveAliasType(t1)
	resolved2 := a.resolveAliasType(t2)

	// Integer can be implicitly cast to real
	if resolved1.StaticType == dt.TAB_ENTRY_REAL && resolved2.StaticType == dt.TAB_ENTRY_INTEGER {
		return true
	}
	if resolved1.StaticType == dt.TAB_ENTRY_INTEGER && resolved2.StaticType == dt.TAB_ENTRY_REAL {
		return true
	}

	// Char can be compared with char
	if resolved1.StaticType == dt.TAB_ENTRY_CHAR && resolved2.StaticType == dt.TAB_ENTRY_CHAR {
		return true
	}

	// Boolean can be compared with boolean
	if resolved1.StaticType == dt.TAB_ENTRY_BOOLEAN && resolved2.StaticType == dt.TAB_ENTRY_BOOLEAN {
		return true
	}

	return false
}

// insertImplicitCast inserts a cast node if needed
// Casts FROM right type TO left type (kanan ke kiri)
// Returns the potentially modified DST and the result type (toType)
func (a *SemanticAnalyzer) insertImplicitCast(dst *dt.DecoratedSyntaxTree, fromType semanticType, toType semanticType) (*dt.DecoratedSyntaxTree, semanticType) {
	// No cast needed if types are equal (already resolves aliases)
	if a.checkTypeEquality(fromType, toType) {
		return dst, toType // Return toType to preserve LHS type
	}

	// Resolve aliases to check if cast is needed
	resolvedFrom := a.resolveAliasType(fromType)
	resolvedTo := a.resolveAliasType(toType)

	// Cast integer to real (when assigning int to real variable)
	if resolvedFrom.StaticType == dt.TAB_ENTRY_INTEGER && resolvedTo.StaticType == dt.TAB_ENTRY_REAL {
		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_CAST_OPERATOR,
			Data:     int(resolvedTo.StaticType), // Target type (resolved)
			Children: []dt.DecoratedSyntaxTree{*dst},
		}, toType // Return original toType to preserve alias
	}

	// No implicit cast for other combinations
	return dst, toType
}

// canCastImplicitly checks if an implicit cast from one type to another is allowed
func (a *SemanticAnalyzer) canCastImplicitly(fromType semanticType, toType semanticType) bool {
	// Same type - no cast needed (already resolves aliases)
	if a.checkTypeEquality(fromType, toType) {
		return true
	}

	// Resolve aliases for cast checking
	resolvedFrom := a.resolveAliasType(fromType)
	resolvedTo := a.resolveAliasType(toType)

	// Integer to real is allowed
	if resolvedFrom.StaticType == dt.TAB_ENTRY_INTEGER && resolvedTo.StaticType == dt.TAB_ENTRY_REAL {
		return true
	}

	return false
}

// promoteTypes promotes two types to a common type if possible
// Returns the promoted DSTs and the result type (always type1/LHS)
func (a *SemanticAnalyzer) promoteTypes(dst1 *dt.DecoratedSyntaxTree, type1 semanticType, dst2 *dt.DecoratedSyntaxTree, type2 semanticType) (*dt.DecoratedSyntaxTree, *dt.DecoratedSyntaxTree, semanticType, bool) {
	// If types are equal, no promotion needed (already resolves aliases)
	if a.checkTypeEquality(type1, type2) {
		return dst1, dst2, type1, true
	}

	// Resolve aliases for promotion checking
	resolved1 := a.resolveAliasType(type1)
	resolved2 := a.resolveAliasType(type2)

	// If one is integer and the other is real, cast RHS to LHS type
	// Result type is always type1 (LHS)
	if resolved1.StaticType == dt.TAB_ENTRY_INTEGER && resolved2.StaticType == dt.TAB_ENTRY_REAL {
		promoted2, _ := a.insertImplicitCast(dst2, type2, type1)
		return dst1, promoted2, type1, true
	}

	if resolved1.StaticType == dt.TAB_ENTRY_REAL && resolved2.StaticType == dt.TAB_ENTRY_INTEGER {
		promoted2, _ := a.insertImplicitCast(dst2, type2, type1)
		return dst1, promoted2, type1, true
	}

	// Char types are compatible
	if resolved1.StaticType == dt.TAB_ENTRY_CHAR && resolved2.StaticType == dt.TAB_ENTRY_CHAR {
		return dst1, dst2, type1, true
	}

	// Boolean types are compatible
	if resolved1.StaticType == dt.TAB_ENTRY_BOOLEAN && resolved2.StaticType == dt.TAB_ENTRY_BOOLEAN {
		return dst1, dst2, type1, true
	}

	// Types are incompatible
	return dst1, dst2, type1, false
}

package semantic

import (
	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func (a *SemanticAnalyzer) checkTypeCompatibility(t1 semanticType, t2 semanticType) bool {
	if a.checkTypeEquality(t1, t2) {
		return true
	}

	resolved1 := a.resolveAliasType(t1)
	resolved2 := a.resolveAliasType(t2)

	if resolved1.StaticType == dt.TAB_ENTRY_REAL && resolved2.StaticType == dt.TAB_ENTRY_INTEGER {
		return true
	}
	if resolved1.StaticType == dt.TAB_ENTRY_INTEGER && resolved2.StaticType == dt.TAB_ENTRY_REAL {
		return true
	}

	if resolved1.StaticType == dt.TAB_ENTRY_CHAR && resolved2.StaticType == dt.TAB_ENTRY_CHAR {
		return true
	}

	if resolved1.StaticType == dt.TAB_ENTRY_BOOLEAN && resolved2.StaticType == dt.TAB_ENTRY_BOOLEAN {
		return true
	}

	return false
}

func (a *SemanticAnalyzer) insertImplicitCast(dst *dt.DecoratedSyntaxTree, fromType semanticType, toType semanticType) (*dt.DecoratedSyntaxTree, semanticType) {
	if a.checkTypeEquality(fromType, toType) {
		return dst, toType
	}

	resolvedFrom := a.resolveAliasType(fromType)
	resolvedTo := a.resolveAliasType(toType)

	if resolvedFrom.StaticType == dt.TAB_ENTRY_INTEGER && resolvedTo.StaticType == dt.TAB_ENTRY_REAL {
		return &dt.DecoratedSyntaxTree{
			SelfType: dt.DST_CAST_OPERATOR,
			Data:     int(resolvedTo.StaticType),
			Children: []dt.DecoratedSyntaxTree{*dst},
		}, toType
	}

	return dst, toType
}

func (a *SemanticAnalyzer) canCastImplicitly(fromType semanticType, toType semanticType) bool {
	if a.checkTypeEquality(fromType, toType) {
		return true
	}

	resolvedFrom := a.resolveAliasType(fromType)
	resolvedTo := a.resolveAliasType(toType)

	if resolvedFrom.StaticType == dt.TAB_ENTRY_INTEGER && resolvedTo.StaticType == dt.TAB_ENTRY_REAL {
		return true
	}

	return false
}

func (a *SemanticAnalyzer) promoteTypes(dst1 *dt.DecoratedSyntaxTree, type1 semanticType, dst2 *dt.DecoratedSyntaxTree, type2 semanticType) (*dt.DecoratedSyntaxTree, *dt.DecoratedSyntaxTree, semanticType, bool) {
	if a.checkTypeEquality(type1, type2) {
		return dst1, dst2, type1, true
	}

	resolved1 := a.resolveAliasType(type1)
	resolved2 := a.resolveAliasType(type2)

	if resolved1.StaticType == dt.TAB_ENTRY_INTEGER && resolved2.StaticType == dt.TAB_ENTRY_REAL {
		promoted2, _ := a.insertImplicitCast(dst2, type2, type1)
		return dst1, promoted2, type1, true
	}

	if resolved1.StaticType == dt.TAB_ENTRY_REAL && resolved2.StaticType == dt.TAB_ENTRY_INTEGER {
		promoted2, _ := a.insertImplicitCast(dst2, type2, type1)
		return dst1, promoted2, type1, true
	}

	if resolved1.StaticType == dt.TAB_ENTRY_CHAR && resolved2.StaticType == dt.TAB_ENTRY_CHAR {
		return dst1, dst2, type1, true
	}
	if resolved1.StaticType == dt.TAB_ENTRY_BOOLEAN && resolved2.StaticType == dt.TAB_ENTRY_BOOLEAN {
		return dst1, dst2, type1, true
	}

	return dst1, dst2, type1, false
}

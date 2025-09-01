package logic

import (
	"fmt"

	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func not(expr Expression) Expression {
	if expr == nil {
		return nil
	}
	if expr == TRUE {
		return FALSE
	}
	if expr == FALSE {
		return TRUE
	}
	if oc, ok := expr.(OperatorCall); ok {
		if oc.GetOperator() == "not" {
			return oc.GetArgs()[0]
		}
		if oc.GetOperator() == "and" {
			newArgs := make([]Expression, len(oc.GetArgs()))
			for i, arg := range oc.GetArgs() {
				newArgs[i] = not(arg)
			}
			return Or(newArgs...)
		}
		if oc.GetOperator() == "or" {
			newArgs := make([]Expression, len(oc.GetArgs()))
			for i, arg := range oc.GetArgs() {
				newArgs[i] = not(arg)
			}
			return And(newArgs...)
		}
		if oc.GetOperator() == "is_not_null" {
			return IsNull(oc.GetArgs()[0])
		}
		if oc.GetOperator() == "is_null" {
			return IsNotNull(oc.GetArgs()[0])
		}
		if oc.GetOperator() == "in" {
			return NotIn(oc.GetArgs()...)
		}
		if oc.GetOperator() == "not_in" {
			return In(oc.GetArgs()...)
		}
		if oc.GetOperator() == "like" {
			return NotLike(oc.GetArgs()[0], oc.GetArgs()[1], oc.GetArgs()[2])
		}
		if oc.GetOperator() == "not_like" {
			return Like(oc.GetArgs()[0], oc.GetArgs()[1], oc.GetArgs()[2])
		}
		if oc.GetOperator() == "lt" {
			return Le(oc.GetArgs()[0], oc.GetArgs()[1])
		}
		if oc.GetOperator() == "le" {
			return Lt(oc.GetArgs()[0], oc.GetArgs()[1])
		}
	}
	panic(fmt.Sprintf("not not supported for %v", expr))
}

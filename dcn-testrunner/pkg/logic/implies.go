package logic

import (
	"reflect"

	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func Implies(a, b Expression) bool {
	l := minimize(a)
	r := minimize(b)
	return implies(l, r)
}

func implies(a, b Expression) bool {
	if a == FALSE {
		return true
	}
	if b == TRUE {
		return true
	}
	ao, ok1 := a.(OperatorCall)
	bo, ok2 := b.(OperatorCall)
	if ok2 && bo.GetOperator() == "and" {
		return every(bo.GetArgs(), func(arg Expression) bool {
			return implies(a, arg)
		})
	}
	if ok1 && ao.GetOperator() == "or" {
		return every(ao.GetArgs(), func(arg Expression) bool {
			return implies(arg, b)
		})
	}

	if ok2 && bo.GetOperator() == "or" {
		return exists(bo.GetArgs(), func(arg Expression) bool {
			return implies(a, arg)
		})
	}
	if ok1 && ao.GetOperator() == "and" {
		return exists(ao.GetArgs(), func(arg Expression) bool {
			return implies(arg, b)
		})
	}

	if !ok1 {
		// so we have a reference
		if ok2 && bo.GetOperator() == "is_not_null" {
			return equal(a, bo.GetArgs()[0])
		}
	}
	if ok1 && ok2 {
		cRefs := commonRefs(ao, bo)
		if len(cRefs) == 0 {
			return false
		}
		if len(cRefs) == 1 {
			i := cRefs[0]
			if i == 0 {
				ca := ao.GetArgs()[1].(Constant)
				cb := bo.GetArgs()[1].(Constant)
				if ao.GetOperator() == "lt" && (bo.GetOperator() == "lt" || bo.GetOperator() == "le") {
					return !cb.LessThan(ca)
				}
				if ao.GetOperator() == "le" && (bo.GetOperator() == "le") {
					return !cb.LessThan(ca)
				}
				if ao.GetOperator() == "le" && (bo.GetOperator() == "lt") {
					return ca.LessThan(cb)
				}
			} else {
				ca := ao.GetArgs()[0].(Constant)
				cb := bo.GetArgs()[0].(Constant)
				if ao.GetOperator() == "lt" && (bo.GetOperator() == "lt" || bo.GetOperator() == "le") {
					return !ca.LessThan(cb)
				}
				if ao.GetOperator() == "le" && (bo.GetOperator() == "le") {
					return !ca.LessThan(cb)
				}
				if ao.GetOperator() == "le" && (bo.GetOperator() == "lt") {
					return cb.LessThan(ca)
				}
			}
		}
	}

	return equal(a, b)
}

func commonRefs(a, b OperatorCall) []int {
	aArgs := a.GetArgs()
	bArgs := b.GetArgs()
	result := make([]int, 0)
	if len(aArgs) == 0 || len(bArgs) == 0 {
		return result
	}
	for i := 0; i < len(aArgs); i++ {
		if ra, ok := aArgs[i].(Reference); ok {
			if rb, ok2 := bArgs[i].(Reference); ok2 {
				if ra.GetName() == rb.GetName() {
					result = append(result, i)
				}
			}
		}
	}
	return result

}

func equal(a, b Expression) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}

	ar, ok := a.(Reference)
	br, ok2 := b.(Reference)
	if ok && ok2 && ar.GetName() == br.GetName() {
		return true
	}

	ao, ok1 := a.(OperatorCall)
	bo, ok2 := b.(OperatorCall)
	if ok1 && ok2 && ao.GetOperator() == bo.GetOperator() {
		if ao.GetOperator() != bo.GetOperator() {
			return false
		}
		for i := 0; i < len(ao.GetArgs()); i++ {
			if !equal(ao.GetArgs()[i], bo.GetArgs()[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func every(a []Expression, f func(Expression) bool) bool {
	for _, v := range a {
		if !f(v) {
			return false
		}
	}
	return true
}
func exists(a []Expression, f func(Expression) bool) bool {
	for _, v := range a {
		if f(v) {
			return true
		}
	}
	return false
}

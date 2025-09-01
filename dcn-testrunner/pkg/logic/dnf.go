package logic

import (
	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func MinimizedDNF(expr Expression) Expression {
	e := minimize(expr)
	e = DNF(e)
	return minimizeDNF(e)

}

func minimizeDNF(expr Expression) Expression {
	op, ok := expr.(OperatorCall)
	if !ok {
		return expr
	}

	if op.GetOperator() == "or" {
		newArgs := []Expression{}
		for _, arg := range op.GetArgs() {
			minArg := minimizeDNF(arg)
			if minArg != FALSE {
				newArgs = append(newArgs, minArg)
			}
		}
		if len(newArgs) == 0 {
			return nil
		}
		if len(newArgs) == 1 {
			return newArgs[0]
		}
		return Or(newArgs...)
	}
	if op.GetOperator() == "and" {
		oldArgs := op.GetArgs()
		for i := 0; i < len(oldArgs); i++ {
			for j := i + 1; j < len(oldArgs); j++ {
				if implies(oldArgs[i], not(oldArgs[j])) {
					return FALSE
				}
				if implies(oldArgs[j], not(oldArgs[i])) {
					return FALSE
				}
			}
		}
		return And(oldArgs...)
	}
	return expr

}

func DNF(expr Expression) Expression {
	return dnf(expr, false)

}

func dnf(expr Expression, inv bool) Expression {
	oc, ok := expr.(OperatorCall)
	if ok {
		if oc.GetOperator() == "not" {
			return dnf(oc.GetArgs()[0], !inv)
		}
		if oc.GetOperator() == "or" && inv || oc.GetOperator() == "and" && !inv {
			return expand(oc.GetArgs(), inv)
		}
		if oc.GetOperator() == "or" && !inv || oc.GetOperator() == "and" && inv {
			newArgs := []Expression{}
			for _, arg := range oc.GetArgs() {
				newArg := dnf(arg, inv)
				if arg, ok := newArg.(OperatorCall); ok && arg.GetOperator() == "or" {
					newArgs = append(newArgs, arg.GetArgs()...)
				} else {
					newArgs = append(newArgs, newArg)
				}
			}
			return Or(newArgs...)
		}
	}
	if inv {
		return Not(expr)
	}

	return expr

}

func expand(args []Expression, inv bool) Expression {
	result := [][]Expression{{}}
	for _, arg := range args {
		r := dnf(arg, inv)
		oc, ok := r.(OperatorCall)
		if !ok {
			for i := 0; i < len(result); i++ {
				result[i] = append(result[i], r)
			}
			continue
		}
		if oc.GetOperator() == "or" {
			newResult := [][]Expression{}
			for _, arg := range oc.GetArgs() {
				for _, r := range result {
					if andArg, ok := arg.(OperatorCall); ok && andArg.GetOperator() == "and" {
						newResult = append(newResult, append(r, andArg.GetArgs()...))
					} else {
						newResult = append(newResult, append(r, arg))
					}
				}
			}
			result = newResult
		} else if oc.GetOperator() == "and" {
			for i := 0; i < len(result); i++ {
				result[i] = append(result[i], oc.GetArgs()...)
			}
		} else {
			for i := 0; i < len(result); i++ {
				result[i] = append(result[i], oc)
			}
		}
	}
	ands := []Expression{}
	for _, r := range result {
		ands = append(ands, And(r...))
	}
	return Or(ands...)
}

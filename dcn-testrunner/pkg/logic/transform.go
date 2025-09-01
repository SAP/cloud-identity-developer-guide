package logic

import (
	"fmt"

	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func minimize(expr Expression) Expression {
	return Visit(
		expr,
		func(s string, args []Expression) Expression {
			switch s {
			case "and":
				return And(args...)
			case "or":
				return Or(args...)
			case "eq":
				if args[0] == TRUE {
					return args[1]
				}
				if args[1] == TRUE {
					return args[0]
				}
				if args[0] == FALSE {
					return Not(args[1])
				}
				if args[1] == FALSE {
					return Not(args[0])
				}
				return And(Le(args[0], args[1]), Le(args[1], args[0]))
			case "lt":
				return Lt(args[0], args[1])
			case "le":
				return Le(args[0], args[1])
			case "is_null":
				return IsNull(args[0])
			case "is_not_null":
				return IsNotNull(args[0])
			case "in":
				if array, ok := args[1].(ArrayConstant); ok {
					newArgs := make([]Expression, 0, len(array.Elements()))
					for _, v := range array.Elements() {
						newArgs = append(newArgs, And(Le(args[0], v), Le(v, args[0])))
					}
					return Or(newArgs...)
				}
				return In(args[0], args[1])
			case "not_in":
				if array, ok := args[1].(ArrayConstant); ok {
					if array.IsEmpty() {
						return IsNotNull(args[0])
					}
					newArgs := make([]Expression, 0, len(array.Elements()))
					for _, v := range array.Elements() {
						newArgs = append(newArgs, Or(Lt(args[0], v), Lt(v, args[0])))
					}
					return And(newArgs...)
				}
				return NotIn(args[0], args[1])
			case "like":
				return Like(args...)
			case "not_like":
				return NotLike(args...)
			case "gt":
				return Lt(args[1], args[0])
			case "ge":
				return Le(args[1], args[0])
			case "ne":
				return Or(Lt(args[0], args[1]), Lt(args[1], args[0]))
			case "between":
				return And(Le(args[1], args[0]), Le(args[0], args[1]))
			case "not_between":
				return Or(Lt(args[0], args[1]), Lt(args[1], args[0]))
			}
			panic(fmt.Sprintf("unknown operator %s", s))
		},
		func(ref Reference) Expression {
			return ref
		},
		func(c Constant) Expression {
			return c
		},
	)
}

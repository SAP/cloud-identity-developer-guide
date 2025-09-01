package z3

// func TransformToZ3(ctx *z3.Context, expr expression.Expression, s Schema) (z3.Bool, error) {
// 	f := ctx.FromBool(false)
// 	switch e := expr.(type) {
// 	case expression.OperatorCall:
// 		args := e.GetArgs()
// 		var err error
// 		switch e.GetOperator() {

// 		case "and":
// 			newArgs := make([]z3.Bool, len(args))
// 			for i, arg := range args {
// 				newArgs[i], err = TransformToZ3(ctx, arg, s)
// 				if err != nil {
// 					return f, err
// 				}
// 			}
// 			if len(newArgs) == 1 {
// 				return newArgs[0], nil
// 			}
// 			return newArgs[0].And(newArgs[1:]...), nil

// 		case "or":
// 			newArgs := make([]z3.Bool, len(args))
// 			for i, arg := range args {
// 				newArgs[i], err = TransformToZ3(ctx, arg, s)
// 				if err != nil {
// 					return f, err
// 				}
// 			}
// 			if len(newArgs) == 1 {
// 				return newArgs[0], nil
// 			}
// 			return newArgs[0].Or(newArgs[1:]...), nil
// 		}
// 		if len(args) == 0 {
// 			return f, fmt.Errorf("no arguments provided for operator %s", e.GetOperator())
// 		}
// 		var z3Op func() z3.Bool
// 		switch typeOf(s, args[0]) {
// 		case STRING:
// 			z3args := make([]z3.Uninterpreted, len(args))
// 			for i, arg := range args {
// 				z3args[i], err = z3String(ctx, arg, e.GetOperator())
// 				if err != nil {
// 					return f, err
// 				}
// 			}
// 			l, err := z3String(ctx, args[0], e.GetOperator())
// 			if err != nil {
// 				return f, err
// 			}
// 			r, err := z3String(ctx, args[1], e.GetOperator())
// 			if err != nil {
// 				return f, err
// 			}
// 			z3Op = func() z3.Bool {
// 				return l.Eq(r)
// 			}
// 		case NUMBER:
// 			l, err := z3Number(ctx, args[0])
// 			if err != nil {
// 				return f, err
// 			}
// 			r, err := z3Number(ctx, args[1])
// 			if err != nil {
// 				return f, err
// 			}
// 			z3Op = func() z3.Bool {
// 				switch e.GetOperator() {
// 				case "eq":
// 					return l.Eq(r)
// 				case "lt":
// 					return l.LT(r)
// 				case "le":
// 					return l.LE(r)
// 				case "gt":
// 					return l.GT(r)
// 				case "ge":
// 					return l.GE(r)
// 				case "ne":
// 					return l.NE(r)
// 				case "between":
// 					return l.GE(r).And(l.LE(args[2]))
// 				}
// 			}
// 		case BOOLEAN:
// 			l, err := z3Bool(ctx, args[0])

// 		case "eq":
// 			switch typeOf(s, args[0]) {
// 			case STRING:
// 				l, err := z3String(ctx, args[0], "eq")
// 				if err != nil {
// 					return f, err
// 				}
// 				r, err := z3String(ctx, args[1], "eq")
// 				if err != nil {
// 					return f, err

// 				}
// 				return l.Eq(r), nil
// 			case NUMBER:
// 				l, err := z3Number(ctx, args[0])
// 				if err != nil {
// 					return f, err
// 				}
// 				r, err := z3Number(ctx, args[1])
// 				if err != nil {
// 					return f, err
// 				}
// 				return l.Eq(r), nil
// 			case BOOLEAN:
// 				l, err := z3Bool(ctx, args[0])
// 				if err != nil {
// 					return f, err
// 				}
// 				r, err := z3Bool(ctx, args[1])
// 				if err != nil {
// 					return f, err
// 				}
// 				return l.Eq(r), nil
// 			}

// 		}
// 	}
// 	return f, nil
// }
// func z3Bool(ctx *z3.Context, e expression.Expression) (z3.Bool, error) {
// 	switch e := e.(type) {
// 	case expression.Reference:
// 		return ctx.Const(e.GetName(), ctx.BoolSort()).(z3.Bool), nil
// 	case expression.Bool:
// 		return ctx.FromBool(bool(e)), nil
// 	}
// 	return ctx.Const("bool_", ctx.BoolSort()).(z3.Bool), fmt.Errorf("unsupported expression type: %T", e)
// }

// func z3Number(ctx *z3.Context, e expression.Expression) (z3.Float, error) {
// 	switch e := e.(type) {
// 	case expression.Reference:
// 		return ctx.Const(e.GetName(), ctx.FloatSort(11, 53)).(z3.Float), nil
// 	case expression.Number:
// 		return ctx.FromFloat64(float64(e), ctx.FloatSort(11, 53)), nil
// 	}
// 	return ctx.Const("num_", ctx.FloatSort(11, 53)).(z3.Float), fmt.Errorf("unsupported expression type: %T", e)
// }
// func z3String(ctx *z3.Context, e expression.Expression, op string) (z3.Uninterpreted, error) {

// 	switch e := e.(type) {
// 	case expression.Reference:
// 		return ctx.Const("ref_"+e.GetName(), ctx.UninterpretedSort("string")).(z3.Uninterpreted), nil
// 	case expression.String:
// 		return ctx.Const(op+"_", ctx.UninterpretedSort("string")).(z3.Uninterpreted), nil
// 	}
// 	return ctx.Const(op+"_", ctx.UninterpretedSort("string")).(z3.Uninterpreted), fmt.Errorf("unsupported expression type: %T", e)
// }

// func typeOf(s Schema, e expression.Expression) InputType {
// 	switch e := e.(type) {
// 	case expression.Reference:
// 		return s.GetTypeOfReference(e.GetName())
// 	case expression.Constant:
// 		switch e.(type) {
// 		case expression.String:
// 			return STRING
// 		case expression.Number:
// 			return NUMBER
// 		case expression.Bool:
// 			return BOOLEAN
// 		case expression.NumberArray:
// 			return NUMBER_ARRAY
// 		case expression.StringArray:
// 			return STRING_ARRAY
// 		case expression.BoolArray:
// 			return BOOLEAN_ARRAY
// 		}
// 	case expression.OperatorCall:
// 		return BOOLEAN
// 	}
// 	return UNDEFINED
// }

package logic

// import (
// 	"fmt"

// 	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
// )

// type expr interface {
// 	implies(expr) bool
// }

// type logicOperator struct {
// 	isAnd bool
// 	args  []expr
// }

// type predicate interface {
// 	implies(expr) bool
// 	equals(predicate) bool
// }

// type boolRef struct {
// 	name    string
// 	negated bool
// }

// type lt struct {
// 	refName       string
// 	isRefLeft     bool
// 	value         expression.Constant
// 	secondRefName string
// }
// type le struct {
// 	refName       string
// 	isRefLeft     bool
// 	value         expression.Constant
// 	secondRefName string
// }
// type is_null struct {
// 	refName string
// 	negated bool
// }

// type in struct {
// 	refName       string
// 	value         expression.Constant
// 	secondRefName string
// }

// type like struct {
// 	refName    string
// 	pattern    string
// 	escapeChar string
// }

// func (b boolRef) equals(p predicate) bool {
// 	if p, ok := p.(boolRef); ok {
// 		return b.name == p.name && b.negated == p.negated
// 	}
// 	return false
// }

// func (b boolRef) implies(r expr) bool {
// 	switch r := r.(type) {
// 	case boolRef:
// 		return b.equals(r)
// 	case logicOperator:
// 		return r.impliedBy(b)
// 	}
// 	return false
// }

// func (l logicOperator) implies(r expr) bool {
// 	rl, ok := r.(logicOperator)
// 	if ok && rl.isAnd {
// 		return rl.impliedBy(l)
// 	}
// 	if !l.isAnd {
// 		return every(l.args, func(e expr) bool {
// 			return e.implies(r)
// 		})
// 	}
// 	if ok {
// 		return rl.impliedBy(l)
// 	}
// 	return exists(l.args, func(e expr) bool {
// 		return e.implies(r)
// 	})
// }

// func (r logicOperator) impliedBy(l expr) bool {
// 	if r.isAnd {
// 		return every(r.args, func(e expr) bool {
// 			return l.implies(e)
// 		})
// 	}
// 	return exists(r.args, func(e expr) bool {
// 		return l.implies(e)
// 	})
// }

// // func (l lt) implies(r expr) bool {
// // 	switch r := r.(type) {
// // 	case lt:
// // 		return l.refName == r.refName && l.isRefLeft == r.isRefLeft && l.value.Equals(r.value)
// // 	case logicOperator:
// // 		return r.impliedBy(l)
// // 	}
// // 	return false
// // }

// func every(e []expr, f func(expr) bool) bool {
// 	for _, v := range e {
// 		if !f(v) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func exists(e []expr, f func(expr) bool) bool {
// 	for _, v := range e {
// 		if f(v) {
// 			return true
// 		}
// 	}
// 	return false
// }

// // func (l logicOperator) impliedBy(e expr) bool {

// // }

// func Implies_old(left, right expr) (bool, error) {
// 	return false, fmt.Errorf("not implemented")
// }

// func TransformOperators(e expression.Expression) (expr, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

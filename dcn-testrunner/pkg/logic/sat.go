package logic

import (
	"math"

	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type InputComponents map[string][]Constant

// FindDiff checks if two expressions are equivalent by evaluating them with all possible combinations of input values.
// it returns an input that makes the expressions different or nil if they are equivalent.
func FindDiff(a, b Expression) Input {
	ic := make(InputComponents)
	min_a := minimize(a)
	min_b := minimize(b)
	findTestValues(min_a, ic)
	findTestValues(min_b, ic)

	inputs := cartesianProduct(ic)
	for _, input := range inputs {
		if a.Evaluate(input) != b.Evaluate(input) {
			return input
		}
	}
	return nil

}

func cartesianProduct(input InputComponents) []Input {
	counters := make(map[string]int, len(input))
	result := make([]Input, 0)
	if len(input) == 0 {
		return result
	}
	for k := range input {
		counters[k] = 0
	}

	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}

	done := false
	for !done {
		// create a new input
		i := make(Input, len(input))
		for k, v := range input {
			i[k] = v[counters[k]]
		}
		result = append(result, i)
		// increment the counter

		current := 0
		currentKey := keys[current]
		counters[currentKey]++
		for counters[currentKey] >= len(input[currentKey]) {
			current++
			if current >= len(keys) {
				done = true
				break
			}
			currentKey = keys[current]
			counters[currentKey]++
			for i := 0; i < current; i++ {
				counters[keys[i]] = 0
			}
		}
	}
	return result

}

func findTestValues(expr Expression, ic InputComponents) {
	e := minimize(expr)
	findTestValuesByConst(e, ic)
	findTestValuesByOtherRefs(e, ic)
	findTestValuesForIsNull(e, ic)
}

func findTestValuesForIsNull(e Expression, ic InputComponents) {
	switch e := e.(type) {
	case Reference:
		ic.Add(e.GetName(), TRUE)
		ic.Add(e.GetName(), FALSE)
	case OperatorCall:
		op := e.GetOperator()
		if op == "and" || op == "or" {
			for _, arg := range e.GetArgs() {
				findTestValuesForIsNull(arg, ic)
			}
			return
		}
		if op == "is_not_null" || op == "is_null" {
			ref := e.GetArgs()[0].(Reference)
			if _, ok := ic[ref.GetName()]; !ok {
				ic.Add(ref.GetName(), TRUE)
			}
			return
		}
	}
}

func findTestValuesByOtherRefs(e Expression, ic InputComponents) {
	switch e := e.(type) {
	case Reference:
		return
	case OperatorCall:
		op := e.GetOperator()
		if op == "and" || op == "or" {
			for _, arg := range e.GetArgs() {
				findTestValuesByOtherRefs(arg, ic)
			}
			return
		}
		if op == "is_not_null" || op == "is_null" {
			return
		}
		if op == "le" || op == "lt" {
			refL, okL := e.GetArgs()[0].(Reference)
			refR, okR := e.GetArgs()[1].(Reference)
			if !okL || !okR {
				// we handled these already
				return
			}
			lValues, okL := ic[refL.GetName()]
			rValues, okR := ic[refR.GetName()]
			if !okL && !okR {
				// we hope that we can freely choose the type of the references now
				// lets assume we have a<b here. if we end up in this if block that means the expression contains
				// no comparison between a and a constant and b and a constant
				// we will run into problems in the future for expressions like a<b and a<c and c<"foo"
				// by then we should figure out a way to handle this ( maybe looking up the schema definition of the reference )

				// fo now we just add numbers
				ic.Add(refL.GetName(), Number(0))
				ic.Add(refR.GetName(), Number(0))
				ic.Add(refL.GetName(), Next(Number(0)))
				ic.Add(refR.GetName(), Next(Number(0)))
				return
			}
			if !okL {
				if op == "lt" {
					ic.Add(refL.GetName(), rValues[0])
					ic.Add(refL.GetName(), Prev(rValues[0]))
				}
				if op == "le" {
					ic.Add(refL.GetName(), rValues[0])
					ic.Add(refL.GetName(), Next(rValues[0]))
				}
				return
			}
			if !okR {
				if op == "lt" {
					ic.Add(refR.GetName(), lValues[0])
					ic.Add(refR.GetName(), Next(lValues[0]))
				}
				if op == "le" {
					ic.Add(refR.GetName(), lValues[0])
					ic.Add(refR.GetName(), Prev(lValues[0]))
				}
				return
			}
			var hasEqual, hasLesser, hasGreater bool
			for _, l := range lValues {
				for _, r := range rValues {
					hasLesser = hasLesser || l.LessThan(r)
					hasGreater = hasGreater || r.LessThan(l)
					hasEqual = hasEqual || l == r
				}
			}
			if op == "lt" {
				if hasLesser && (hasEqual || hasGreater) {
					return
				}
				if hasLesser {
					ic.Add(refL.GetName(), rValues[0])
					return
				}
				if hasGreater || hasEqual {
					ic.Add(refL.GetName(), Prev(rValues[0]))
					return
				}
			}
			if op == "le" {
				if (hasLesser || hasEqual) && hasGreater {
					return
				}
				if hasLesser || hasEqual {
					ic.Add(refL.GetName(), Next(rValues[0]))
					return
				}
				if hasGreater {
					ic.Add(refL.GetName(), rValues[0])
					return
				}
			}
		}
		if op == "in" || op == "not_in" {
			arrayRef := e.GetArgs()[1].(Reference)
			ref, ok := e.GetArgs()[0].(Reference)
			if !ok {
				// we handled these already
				return
			}
			arrayValues, ok := ic[arrayRef.GetName()]
			refValues, ok2 := ic[ref.GetName()]
			if !ok && !ok2 {
				// similar to the le and lt case we assume that we can freely choose the type of the references
				ic.Add(ref.GetName(), Number(0))
				ic.Add(arrayRef.GetName(), NumberArray{})
				ic.Add(arrayRef.GetName(), NumberArray{Number(0)})
				return
			}
			if !ok {
				c := refValues[0]
				switch c := c.(type) {
				case String:
					ic.Add(arrayRef.GetName(), StringArray{c})
					ic.Add(arrayRef.GetName(), StringArray{})
				case Number:
					ic.Add(arrayRef.GetName(), NumberArray{c})
					ic.Add(arrayRef.GetName(), NumberArray{})
				case Bool:
					ic.Add(arrayRef.GetName(), BoolArray{c})
					ic.Add(arrayRef.GetName(), BoolArray{})
				}
				return
			}
			if !ok2 {
				var c Constant
				for _, v := range arrayValues {
					va := v.(ArrayConstant)
					if !va.IsEmpty() {
						c = va.Elements()[0]
						break
					}
				}
				ic.Add(ref.GetName(), c)

				return
			}

		}
	}
}

func findTestValuesByConst(e Expression, ic InputComponents) {

	switch e := e.(type) {
	case Reference:
		ic.Add(e.GetName(), TRUE)
		ic.Add(e.GetName(), FALSE)
	case OperatorCall:
		op := e.GetOperator()
		if op == "and" || op == "or" {
			for _, arg := range e.GetArgs() {
				findTestValuesByConst(arg, ic)
			}
			return
		}
		if op == "is_not_null" || op == "is_null" {
			// we handle these later
			return
		}
		if op == "le" || op == "lt" {
			var ref Reference
			var c Constant
			refL, okL := e.GetArgs()[0].(Reference)
			refR, okR := e.GetArgs()[1].(Reference)
			if okL && okR {
				// we handle these later
				return
			}
			if okL {
				ref = refL
				c = e.GetArgs()[1].(Constant)
			} else if okR {
				ref = refR
				c = e.GetArgs()[0].(Constant)
			} else {
				panic(" unexpected operator call without reference")
			}

			if okL && op == "lt" || okR && op == "le" {
				ic.Add(ref.GetName(), c)
				ic.Add(ref.GetName(), Prev(c))
			} else {
				ic.Add(ref.GetName(), c)
				ic.Add(ref.GetName(), Next(c))
			}
			return
		}
		if op == "in" || op == "not_in" {
			ref := e.GetArgs()[1].(Reference)
			c, ok := e.GetArgs()[0].(Constant)
			if !ok {
				// we handle these later
				return
			}
			var arrayC ArrayConstant
			var emptyArrayC ArrayConstant
			switch c := c.(type) {
			case String:
				arrayC = StringArray{c}
				emptyArrayC = StringArray{}
			case Number:
				arrayC = NumberArray{c}
				emptyArrayC = NumberArray{}
			case Bool:
				arrayC = BoolArray{c}
				emptyArrayC = BoolArray{}
			}
			ic.Add(ref.GetName(), arrayC)
			ic.Add(ref.GetName(), emptyArrayC)
			return
		}
		if op == "like" {
			ref := e.GetArgs()[0].(Reference)
			pattern := e.GetArgs()[1].(Constant)
			ic.Add(ref.GetName(), pattern)
			ic.Add(ref.GetName(), String(""))
			return
		}

	}
}

func Next(c Constant) Constant {
	switch c := c.(type) {
	case String:
		return String(NextString((string)(c)))
	case Number:
		return Number(math.Nextafter(float64(c), math.MaxFloat64))
	case Bool:
		return Bool(!bool(c))
	default:
		panic("unknown constant type")
	}
}

func Prev(c Constant) Constant {
	switch c := c.(type) {
	case String:
		return String(PrevString((string)(c)))
	case Number:
		return Number(math.Nextafter(float64(c), -math.MaxFloat64))
	case Bool:
		return Bool(!bool(c))
	default:
		panic("unknown constant type")
	}
}

func (ic InputComponents) Merge(other InputComponents) InputComponents {
	for k, v := range other {
		if _, ok := ic[k]; !ok {
			ic[k] = []Constant{}
		}
		for _, val := range v {
			ic.Add(k, val)
		}
	}
	return ic
}
func (ic InputComponents) Add(key string, value Constant) {
	if _, ok := ic[key]; !ok {
		ic[key] = []Constant{}
	}
	for _, v := range ic[key] {
		if v == value {
			return
		}
	}
	ic[key] = append(ic[key], value)
}

func NextString(s string) string {
	b := []byte(s)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < 0xFF {
			b[i]++
			return string(b[:i+1])
		}
	}
	// If all bytes are 0xFF, append a 0x00 (smallest byte) to make it greater
	return s + "\x00"
}
func PrevString(s string) string {
	if len(s) == 0 {
		return ""
	}

	b := []byte(s)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] > 0x00 {
			b[i]--
			return string(b[:i+1])
		}
	}
	// All bytes are 0x00, return empty to indicate no smaller string
	return ""
}

package logic

import (
	"reflect"
	"testing"

	e "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestIsEquivalent(t *testing.T) {

	t.Run("x = 1 <=> x = 1", func(t *testing.T) {
		a := e.Eq(e.Ref("x"), e.Number(1))
		b := e.Eq(e.Ref("x"), e.Number(1))
		got := FindDiff(a, b)
		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

	t.Run("FALSE <=> TRUE", func(t *testing.T) {
		a := e.FALSE
		b := e.TRUE
		got := FindDiff(a, b)
		if got == nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})
	t.Run("x = 1 !<=> x = 2", func(t *testing.T) {
		a := e.Eq(e.Ref("x"), e.Number(1))
		b := e.Eq(e.Ref("x"), e.Number(2))
		got := FindDiff(a, b)
		if got == nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

	t.Run("x = 1 and x = 2 <=> false", func(t *testing.T) {
		a := e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Eq(e.Ref("x"), e.Number(2)))
		b := e.FALSE
		got := FindDiff(a, b)
		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

	t.Run("x < 1 and x > 2 <=> false", func(t *testing.T) {
		a := e.And(e.Lt(e.Ref("x"), e.Number(1)), e.Gt(e.Ref("x"), e.Number(2)))
		b := e.FALSE
		got := FindDiff(a, b)
		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

	t.Run("x = 1 and (y =1 or y=2) <=> x=1 and y=1 or x=1 and y=2", func(t *testing.T) {
		a := e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Or(e.Eq(e.Ref("y"), e.Number(1)), e.Eq(e.Ref("y"), e.Number(2))))
		b := e.Or(e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Eq(e.Ref("y"), e.Number(1))), e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Eq(e.Ref("y"), e.Number(2))))
		got := FindDiff(a, b)
		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

	t.Run("x = 1 and (y =1 or y=2) !<=> x=1 and y=1", func(t *testing.T) {
		a := e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Or(e.Eq(e.Ref("y"), e.Number(1)), e.Eq(e.Ref("y"), e.Number(2))))
		b := e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Eq(e.Ref("y"), e.Number(1)))
		got := FindDiff(a, b)
		if got == nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

	t.Run("x = 1 and y = 1 <=> x = 1 and (x = 2 or y=1)", func(t *testing.T) {
		a := e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Eq(e.Ref("y"), e.Number(1)))
		b := e.And(e.Eq(e.Ref("x"), e.Number(1)), e.Or(e.Eq(e.Ref("x"), e.Number(2)), e.Eq(e.Ref("y"), e.Number(1))))
		got := FindDiff(a, b)
		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})

}

func TestCartesionProduct(t *testing.T) {
	t.Run("two references only one value", func(t *testing.T) {
		input := InputComponents{
			"x": {e.String("a")},
			"y": {e.String("b")},
		}
		got := cartesianProduct(input)
		want := []e.Input{
			{"x": e.String("a"), "y": e.String("b")},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("two references one and two values", func(t *testing.T) {
		input := InputComponents{
			"x": {e.String("a"), e.String("b")},
			"y": {e.String("c")},
		}
		got := cartesianProduct(input)
		want := []e.Input{
			{"x": e.String("a"), "y": e.String("c")},
			{"x": e.String("b"), "y": e.String("c")},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("two references two and one values", func(t *testing.T) {
		input := InputComponents{
			"x": {e.String("a")},
			"y": {e.String("c"), e.String("d")},
		}
		got := cartesianProduct(input)
		want := []e.Input{
			{"x": e.String("a"), "y": e.String("c")},
			{"x": e.String("a"), "y": e.String("d")},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("two references two and two values", func(t *testing.T) {
		input := InputComponents{
			"x": {e.String("a"), e.String("b")},
			"y": {e.String("c"), e.String("d")},
		}
		got := cartesianProduct(input)
		want := []e.Input{
			{"x": e.String("a"), "y": e.String("c")},
			{"x": e.String("a"), "y": e.String("d")},
			{"x": e.String("b"), "y": e.String("c")},
			{"x": e.String("b"), "y": e.String("d")},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}

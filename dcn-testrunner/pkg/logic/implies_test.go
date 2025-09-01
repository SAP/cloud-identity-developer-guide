package logic

import (
	"testing"

	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestImplies(t *testing.T) {

	t.Run("a<1 => a<2", func(t *testing.T) {
		left := Lt(Ref("a"), Number(1))
		right := Lt(Ref("a"), Number(2))
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if Implies(right, left) {
			t.Errorf("Expected %v not to imply %v", right, left)
		}
	})

	t.Run("FALSE => TRUE", func(t *testing.T) {
		left := FALSE
		right := TRUE
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if Implies(right, left) {
			t.Errorf("Expected %v not to imply %v", right, left)
		}
	})

	t.Run("a = 2 => a>1", func(t *testing.T) {
		left := Eq(Ref("a"), Number(2))
		right := Gt(Ref("a"), Number(1))
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if Implies(right, left) {
			t.Errorf("Expected %v to not imply %v", right, left)
		}
	})

	t.Run("a = 1 <=> a=1", func(t *testing.T) {
		left := Eq(Ref("a"), Number(1))
		right := Eq(Ref("a"), Number(1))
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if !Implies(right, left) {
			t.Errorf("Expected %v to imply %v", right, left)
		}
	})

	t.Run("a and b => a or b", func(t *testing.T) {
		left := And(Ref("a"), Ref("b"))
		right := Or(Ref("a"), Ref("b"))
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if Implies(right, left) {
			t.Errorf("Expected %v not to imply %v", right, left)
		}
	})

	t.Run("a or b <=> a or b", func(t *testing.T) {
		left := Or(Ref("a"), Ref("b"))
		right := Or(Ref("a"), Ref("b"))
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if !Implies(right, left) {
			t.Errorf("Expected %v to imply %v", right, left)
		}
	})

	t.Run("a = 1 and b = 1 <=> a = 1 and (a = 2 or b=1)", func(t *testing.T) {
		left := And(Eq(Ref("a"), Number(1)), Eq(Ref("b"), Number(1)))
		right := And(Eq(Ref("a"), Number(1)), Or(Eq(Ref("a"), Number(2)), Eq(Ref("b"), Number(1))))
		if !Implies(left, right) {
			t.Errorf("Expected %v to imply %v", left, right)
		}
		if !Implies(right, left) {
			t.Errorf("Expected %v to imply %v", right, left)
		}
	})
}

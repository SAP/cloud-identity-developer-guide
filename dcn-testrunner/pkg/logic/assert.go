package logic

import (
	"fmt"
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func AssertEquivalence(a, b expression.Expression) error {

	if reflect.DeepEqual(a, b) {
		return nil
	}

	diff := FindDiff(a, b)
	if diff != nil {
		return fmt.Errorf("expressions are not equivalent, found difference with input: %v", diff)
	}

	dnfA := MinimizedDNF(a)
	dnfB := MinimizedDNF(b)

	if !Implies(dnfA, dnfB) {
		return fmt.Errorf("expressions are not equivalent after DNF minimization: A: %s, B: %s", dnfA, dnfB)
	}
	if !Implies(dnfB, dnfA) {
		return fmt.Errorf("expressions are not equivalent after DNF minimization: B: %s, A: %s", dnfB, dnfA)
	}

	return nil
}

package in_test

import (
	"dcn-testrunner/pkg/client"
	"dcn-testrunner/pkg/logic"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestMain(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Could not get current filename")
	}
	dcnPath := filepath.Join(filepath.Dir(filename), "dcn")

	c, err := client.NewClient(dcnPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("in with empty set: no input", func(t *testing.T) {
		req := client.EvaluatePoliciesRequest("*", "in_empty_set", "pkg.in_empty_set")

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.In(
			expression.Ref("$app.s"),
			expression.Ref("$app.s_a"),
		)
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("in with empty set: one input", func(t *testing.T) {
		req := client.EvaluatePoliciesRequest("*", "in_empty_set", "pkg.in_empty_set")
		req.Input.App["s"] = "x"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.FALSE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})
}

package features_test

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

	t.Run("LikeCheck", func(t *testing.T) {
		req := client.EvaluatePoliciesRequest("IsLike", "DUMMY", "cas.PolLike")
		req.Input.App["stringval"] = "TEST"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("LikeCheck failed: %v", err)
		}
	})

	t.Run("value 1TEST1", func(t *testing.T) {
		req := client.EvaluatePoliciesRequest("IsLike2", "DUMMY", "cas.PolLike")
		req.Input.App["stringval"] = "1TEST1"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("LikeCheck failed: %v", err)
		}
	})
	t.Run("value 1TEST", func(t *testing.T) {
		req := client.EvaluatePoliciesRequest("IsLike2", "DUMMY", "cas.PolLike")
		req.Input.App["stringval"] = "1TEST"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.FALSE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("LikeCheck failed: %v", err)
		}
	})

}

package like_test

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

	t.Run("like with underscore: no input", func(t *testing.T) {
		// GRANT * ON like_with_underscore WHERE s LIKE 'x_y';
		req := client.EvaluatePoliciesRequest("*", "like_with_underscore", "pkg.like_with_underscore")

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.Like(
			expression.Ref("$app.s"),
			expression.String("x_y"),
		)
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("like with underscore: no character", func(t *testing.T) {
		// GRANT * ON like_with_underscore WHERE s LIKE 'x_y';
		req := client.EvaluatePoliciesRequest("*", "like_with_underscore", "pkg.like_with_underscore")
		req.Input.App["s"] = "xy"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.FALSE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("like with underscore: one character", func(t *testing.T) {
		// GRANT * ON like_with_underscore WHERE s LIKE 'x_y';
		req := client.EvaluatePoliciesRequest("*", "like_with_underscore", "pkg.like_with_underscore")
		req.Input.App["s"] = "xay"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})
	t.Run("like with underscore: two characters", func(t *testing.T) {
		// GRANT * ON like_with_underscore WHERE s LIKE 'x_y';
		req := client.EvaluatePoliciesRequest("*", "like_with_underscore", "pkg.like_with_underscore")
		req.Input.App["s"] = "xaxy"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.FALSE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("like with percent: no input", func(t *testing.T) {
		// GRANT * ON like_with_percent WHERE s LIKE 'x%y';
		req := client.EvaluatePoliciesRequest("*", "like_with_percent", "pkg.like_with_percent")

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.Like(
			expression.Ref("$app.s"),
			expression.String("x%y"),
		)
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("like with percent: no character", func(t *testing.T) {
		// GRANT * ON like_with_percent WHERE s LIKE 'x%y';
		req := client.EvaluatePoliciesRequest("*", "like_with_percent", "pkg.like_with_percent")
		req.Input.App["s"] = "xy"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("like with percent: one character", func(t *testing.T) {
		// GRANT * ON like_with_percent WHERE s LIKE 'x%y';
		req := client.EvaluatePoliciesRequest("*", "like_with_percent", "pkg.like_with_percent")
		req.Input.App["s"] = "xay"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("like with percent: two characters", func(t *testing.T) {
		// GRANT * ON like_with_percent WHERE s LIKE 'x%y';
		req := client.EvaluatePoliciesRequest("*", "like_with_percent", "pkg.like_with_percent")
		req.Input.App["s"] = "xaxy"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("LikeWithExcape: no input", func(t *testing.T) {
		// GRANT * ON like_with_escape WHERE s LIKE 'xö%_ö_y' ESCAPE 'ö';
		req := client.EvaluatePoliciesRequest("*", "like_with_escape", "pkg.like_with_escape")

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.Like(
			expression.Ref("$app.s"),
			expression.String("xö%_ö_y"),
			expression.String("ö"),
		)
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("LikeWithExcape: no character", func(t *testing.T) {
		// GRANT * ON like_with_escape WHERE s LIKE 'xö%_ö_y' ESCAPE 'ö';
		req := client.EvaluatePoliciesRequest("*", "like_with_escape", "pkg.like_with_escape")
		req.Input.App["s"] = "x%_y"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.FALSE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("LikeWithExcape: one character", func(t *testing.T) {
		// GRANT * ON like_with_escape WHERE s LIKE 'xö%_ö_y' ESCAPE 'ö';
		req := client.EvaluatePoliciesRequest("*", "like_with_escape", "pkg.like_with_escape")
		req.Input.App["s"] = "x%C_y"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.TRUE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})
	t.Run("LikeWithExcape: no percentage", func(t *testing.T) {
		// GRANT * ON like_with_escape WHERE s LIKE 'xö%_ö_y' ESCAPE 'ö';
		req := client.EvaluatePoliciesRequest("*", "like_with_escape", "pkg.like_with_escape")
		req.Input.App["s"] = "xCC_y"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.FALSE
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

	t.Run("input but no resource on all policies", func(t *testing.T) {
		// GRANT * ON like_with_underscore WHERE s LIKE 'x_y';
		// GRANT * ON like_with_percent WHERE s LIKE 'x%y';
		// GRANT * ON like_with_escape WHERE s LIKE 'xö%_ö_y' ESCAPE 'ö';
		req := client.EvaluatePoliciesRequest("*", "", "pkg.like_with_underscore", "pkg.like_with_percent", "pkg.like_with_escape")
		req.Input.App["s"] = "x_y"

		got, err := c.EvaluatePolicies(req)
		if err != nil {
			t.Fatalf("Failed to evaluate policies: %v", err)
		}
		want := expression.Or(
			expression.Eq(
				expression.Ref("$dcl.resource"),
				expression.String("like_with_underscore"),
			),
			expression.Eq(
				expression.Ref("$dcl.resource"),
				expression.String("like_with_percent"),
			),
		)
		if err := logic.AssertEquivalence(got, want); err != nil {
			t.Errorf("got %v\n want %v\n %v", got, want, err)
		}
	})

}

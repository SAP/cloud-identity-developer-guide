package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type Client struct {
	baseURL string
}

type DCLInput struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type Input struct {
	DCL DCLInput               `json:"$dcl"`
	App map[string]interface{} `json:"$app"`
	Env map[string]interface{} `json:"$env"`
}

type PolicyEvaluationRequest struct {
	Policies []string `json:"policies"`
	Input    Input    `json:"input"`
}

type NullifyExceptRequest struct {
	Expression dcn.Expression `json:"expression"`
	KeepRefs   [][]string     `json:"keep_refs"`
}

type DefaultPoliciesRequest struct {
	Tenant string `json:"tenant"`
}
type DefaultPoliciesResponse struct {
	Policies []string `json:"policies"`
}

type EvaluationResponse struct {
	Expression dcn.Expression `json:"expression"`
}

const (
	PATH_LOAD_DCN             = "/v1/load_dcn"
	PATH_EVALUATE_POLICIES    = "/v1/evaluate_policies"
	PATH_NULLIFY_EXCEPT       = "/v1/nullify_except"
	PATH_GET_DEFAULT_POLICIES = "/v1/get_default_policies"
)

func NewClient(dcnPath string) (*Client, error) {
	c := &Client{
		baseURL: os.Getenv("DCN_TEST_SERVER_URL"),
	}
	if c.baseURL == "" {
		c.baseURL = "http://localhost:8085"
	}
	loader := dcn.NewLocalLoader(dcnPath)
	dcnContainer := <-loader.DCNChannel

	_, err := post[any](c, PATH_LOAD_DCN, dcnContainer)
	if err != nil {
		return nil, fmt.Errorf("failed to load DCN directory: %w", err)
	}
	return c, nil
}

func (c *Client) EvaluatePolicies(req PolicyEvaluationRequest) (expression.Expression, error) {
	resp, err := post[EvaluationResponse](c, PATH_EVALUATE_POLICIES, req)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policies: %w", err)
	}
	e, err := expression.FromDCN(resp.Expression, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to convert expression from DCN: %w", err)
	}
	return e.Expression, nil

}

func EvaluatePoliciesRequest(action, resource string, policies ...string) PolicyEvaluationRequest {
	return PolicyEvaluationRequest{
		Policies: policies,
		Input: Input{
			DCL: DCLInput{
				Action:   action,
				Resource: resource,
			},
			App: make(map[string]interface{}),
			Env: make(map[string]interface{}),
		},
	}
}

func (c *Client) NullifyExcept(req NullifyExceptRequest) (EvaluationResponse, error) {
	return post[EvaluationResponse](c, PATH_NULLIFY_EXCEPT, req)
}
func (c *Client) GetDefaultPolicies(req DefaultPoliciesRequest) (DefaultPoliciesResponse, error) {
	return post[DefaultPoliciesResponse](c, PATH_GET_DEFAULT_POLICIES, req)
}

func post[T any](c *Client, path string, body any) (T, error) {
	var result T
	rb, err := json.Marshal(body)
	if err != nil {
		return result, err
	}
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(rb))
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}
	return result, nil
}

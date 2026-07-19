package evaluator

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/rego"
)

type OPAEvaluator struct {
	policyPath string
}

func NewOPAEvaluator(policyPath string) (*OPAEvaluator, error) {
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("policy file not found at %s", policyPath)
	}
	return &OPAEvaluator{policyPath: policyPath}, nil
}

// Evaluate runs the input context against the Rego policies and returns (allowed, reason, error)
func (e *OPAEvaluator) Evaluate(ctx context.Context, input map[string]interface{}) (bool, string, error) {
	// Prepare the Rego evaluation query
	query, err := rego.New(
		rego.Query("allow = data.trustchain.allow; reason = data.trustchain.reason"),
		rego.Load([]string{e.policyPath}, nil),
	).PrepareForEval(ctx)
	if err != nil {
		return false, "", fmt.Errorf("preparing rego query: %w", err)
	}

	// Execute evaluation
	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, "", fmt.Errorf("evaluating rego: %w", err)
	}

	if len(results) == 0 {
		return false, "undefined evaluation result", nil
	}

	// Extract variables
	allow, ok := results[0].Bindings["allow"].(bool)
	if !ok {
		return false, "invalid type for allow result", nil
	}

	reasonStr := "no reason provided"
	if reason, ok := results[0].Bindings["reason"].(string); ok {
		reasonStr = reason
	}

	return allow, reasonStr, nil
}

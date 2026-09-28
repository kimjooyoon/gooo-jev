package decision

import "testing"

func policyFixture(t *testing.T, confidence float64, status Status) (DecisionPolicy, Spec, State, Result, Receipt) {
	t.Helper()
	policy, err := NewDecisionPolicy("policy-1", 0.8)
	if err != nil {
		t.Fatalf("NewDecisionPolicy() error = %v", err)
	}
	spec := Spec{
		ID:             "candidate-review",
		Question:       "Should this candidate be reviewed?",
		Kind:          KindChoice,
		AllowedChoices: []string{"accept", "review"},
		PolicyDigest:   "policy-1",
	}
	result := Result{
		SpecID:        spec.ID,
		Kind:          spec.Kind,
		Value:         Value{Choice: "review"},
		Confidence:    &confidence,
		EvidenceDigest: "evidence-1",
		Provider:       "jev",
		Model:          "jev-latest",
		Status:         status,
		ObservedAt:    fixedPolicyTime(),
	}
	state := State{Digest: "state-1"}
	receipt, err := Observe(spec, state, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	return policy, spec, state, result, receipt
}

func fixedPolicyTime() (resultTime time.Time) {
	return time.Unix(1_800_000_000, 0).UTC()
}

func TestEvaluateDecisionPolicyEligibleIsNonAuthorizing(t *testing.T) {
	policy, spec, state, result, receipt := policyFixture(t, 0.91, StatusObserved)
	outcome, err := EvaluateDecisionPolicy(policy, spec, state, result, receipt)
	if err != nil {
		t.Fatalf("EvaluateDecisionPolicy() error = %v", err)
	}
	if outcome.Disposition != PolicyEligible {
		t.Fatalf("disposition = %q, want %q", outcome.Disposition, PolicyEligible)
	}
	if !outcome.NonAuthorizing {
		t.Fatal("eligible outcome must remain non-authorizing")
	}
	if err := outcome.Validate(); err != nil {
		t.Fatalf("PolicyOutcome.Validate() error = %v", err)
	}
}

func TestEvaluateDecisionPolicyEscalatesBelowThreshold(t *testing.T) {
	policy, spec, state, result, receipt := policyFixture(t, 0.42, StatusObserved)
	outcome, err := EvaluateDecisionPolicy(policy, spec, state, result, receipt)
	if err != nil {
		t.Fatalf("EvaluateDecisionPolicy() error = %v", err)
	}
	if outcome.Disposition != PolicyEscalate {
		t.Fatalf("disposition = %q, want %q", outcome.Disposition, PolicyEscalate)
	}
	if outcome.Reason != "confidence-below-policy-threshold" {
		t.Fatalf("reason = %q", outcome.Reason)
	}
}

func TestEvaluateDecisionPolicyPreservesUnknownAndTampering(t *testing.T) {
	policy, spec, state, result, receipt := policyFixture(t, 0.91, StatusUnknown)
	unknown, err := EvaluateDecisionPolicy(policy, spec, state, result, receipt)
	if err != nil {
		t.Fatalf("unknown evaluation error = %v", err)
	}
	if unknown.Disposition != PolicyUnknown {
		t.Fatalf("unknown disposition = %q", unknown.Disposition)
	}

	tampered := receipt
	tampered.StateDigest = "state-2"
	outcome, err := EvaluateDecisionPolicy(policy, spec, state, result, tampered)
	if err != nil {
		t.Fatalf("tampered evaluation error = %v", err)
	}
	if outcome.Disposition != PolicyUnknown || outcome.Reason != "decision-receipt-binding-mismatch" {
		t.Fatalf("tampered outcome = %#v", outcome)
	}
	if err := outcome.Validate(); err != nil {
		t.Fatalf("tampered outcome validation error = %v", err)
	}
}

func TestEvaluateJSONPolicyUsesTheSameBoundary(t *testing.T) {
	policy, err := NewDecisionPolicy("policy-1", 0.8)
	if err != nil {
		t.Fatalf("NewDecisionPolicy() error = %v", err)
	}
	outcome, err := EvaluateJSONPolicy([]byte(validDecisionJSON()), policy)
	if err != nil {
		t.Fatalf("EvaluateJSONPolicy() error = %v", err)
	}
	if outcome.Disposition != PolicyEligible {
		t.Fatalf("JSON disposition = %q, want %q", outcome.Disposition, PolicyEligible)
	}
}

package decision

import (
	"testing"
	"time"
)

func TestRuleProviderReturnsConfiguredValue(t *testing.T) {
	score := 0.7
	spec := Spec{
		ID:           "risk-review",
		Question:     "Is the proposed change low risk?",
		Kind:         KindScore,
		PolicyDigest: "policy-1",
	}
	provider := RuleProvider{
		Values: map[string]Value{
			spec.ID: {Score: &score},
		},
		Clock: func() time.Time { return time.Unix(30, 0).UTC() },
	}
	result, err := provider.Observe(spec, State{Digest: "state-3"})
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if result.Status != StatusObserved {
		t.Fatalf("status = %q, want %q", result.Status, StatusObserved)
	}
	if result.Value.Score == nil || *result.Value.Score != score {
		t.Fatalf("score = %#v, want %v", result.Value.Score, score)
	}
}

func TestRuleProviderReturnsUnknownWhenValueIsMissing(t *testing.T) {
	spec := Spec{
		ID:           "freshness",
		Question:     "Is the evidence fresh enough?",
		Kind:         KindNoul,
		PolicyDigest: "policy-2",
	}
	provider := RuleProvider{
		Values: map[string]Value{},
		Clock:  func() time.Time { return time.Unix(40, 0).UTC() },
	}
	result, err := provider.Observe(spec, State{Digest: "state-4"})
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if result.Status != StatusUnknown {
		t.Fatalf("status = %q, want %q", result.Status, StatusUnknown)
	}
}

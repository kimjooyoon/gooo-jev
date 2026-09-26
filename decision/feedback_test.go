package decision

import (
	"testing"
	"time"
)

func TestObserveFeedbackBindsMetricToLedger(t *testing.T) {
	spec := Spec{
		ID:             "feedback-review",
		Question:       "Did the observed decision improve the result?",
		Kind:           KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest:   "policy-feedback",
	}
	state := State{Digest: "state-feedback"}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: "evidence-decision",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(90, 0).UTC(),
	}
	ledger, err := (Ledger{}).AppendVerified(spec, state, result)
	if err != nil {
		t.Fatalf("AppendVerified() error = %v", err)
	}
	feedback, err := ObserveFeedback(
		ledger,
		ledger.Entries[0].Receipt,
		FeedbackConfirmed,
		"utility.delta",
		0.25,
		"evidence-feedback",
		time.Unix(100, 0).UTC(),
	)
	if err != nil {
		t.Fatalf("ObserveFeedback() error = %v", err)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !feedback.NonAuthorizing || feedback.MetricValue != 0.25 {
		t.Fatalf("feedback = %#v", feedback)
	}
}

func TestObserveFeedbackRejectsReceiptOutsideLedger(t *testing.T) {
	spec := Spec{
		ID:             "missing-feedback",
		Question:       "Was the result useful?",
		Kind:           KindNoul,
		PolicyDigest:   "policy-missing-feedback",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Provider:       "rule",
		Status:         StatusUnknown,
		EvidenceDigest: "evidence-missing",
		ObservedAt:     time.Unix(110, 0).UTC(),
	}
	other, err := Observe(spec, State{Digest: "state-missing"}, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if _, err := ObserveFeedback(
		Ledger{},
		other,
		FeedbackUnknown,
		"utility.delta",
		0,
		"evidence-missing-feedback",
		time.Unix(120, 0).UTC(),
	); err == nil {
		t.Fatal("ObserveFeedback() accepted a receipt outside the ledger")
	}
}

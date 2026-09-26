package decision

import (
	"testing"
	"time"
)

func TestObserveFeedbackFromReplayRequiresMatchingSource(t *testing.T) {
	spec := Spec{
		ID:             "feedback-replay",
		Question:       "Did replayed evidence improve the result?",
		Kind:           KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest:   "policy-feedback-replay",
	}
	state := State{Digest: "state-feedback-replay"}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: "evidence-feedback-replay",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(130, 0).UTC(),
	}
	ledger, err := (Ledger{}).AppendVerified(spec, state, result)
	if err != nil {
		t.Fatalf("AppendVerified() error = %v", err)
	}
	feedback, err := ObserveFeedbackFromReplay(
		ledger,
		[]ReplayInput{{Spec: spec, State: state, Result: result}},
		0,
		FeedbackConfirmed,
		"utility.delta",
		-0.1,
		"evidence-feedback-replay",
		time.Unix(140, 0).UTC(),
	)
	if err != nil {
		t.Fatalf("ObserveFeedbackFromReplay() error = %v", err)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	result.Value.Choice = "no"
	if _, err := ObserveFeedbackFromReplay(
		ledger,
		[]ReplayInput{{Spec: spec, State: state, Result: result}},
		0,
		FeedbackRefuted,
		"utility.delta",
		0.1,
		"evidence-feedback-replay",
		time.Unix(141, 0).UTC(),
	); err == nil {
		t.Fatal("feedback accepted changed source after replay")
	}
}

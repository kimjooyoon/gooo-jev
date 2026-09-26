package decision

import (
	"testing"
	"time"
)

func TestReplayRecomputesSourceInputs(t *testing.T) {
	spec := Spec{
		ID:             "replay-review",
		Question:       "Does the replay preserve the original decision?",
		Kind:           KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest:   "policy-replay",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: "evidence-replay",
		Provider:       "offline-fixture",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(80, 0).UTC(),
	}
	state := State{Digest: "state-replay"}
	ledger, err := (Ledger{}).AppendVerified(spec, state, result)
	if err != nil {
		t.Fatalf("AppendVerified() error = %v", err)
	}
	input := ReplayInput{Spec: spec, State: state, Result: result}
	if err := ledger.Replay([]ReplayInput{input}); err != nil {
		t.Fatalf("Replay() error = %v", err)
	}
	input.Result.Value.Choice = "no"
	if err := ledger.Replay([]ReplayInput{input}); err == nil {
		t.Fatal("Replay() accepted changed source input")
	}
}

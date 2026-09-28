package decision

import (
	"reflect"
	"testing"
	"time"
)

func batchResultTime() time.Time {
	return time.Unix(1_800_000_100, 0).UTC()
}

func TestObserveBatchBindsSharedStateAndOrder(t *testing.T) {
	confident := 0.9
	state := State{Digest: "state-batch-1"}
	observations := []BatchObservation{
		{
			Spec: Spec{
				ID:             "route",
				Question:       "Which route should handle this?",
				Kind:           KindChoice,
				AllowedChoices: []string{"fast", "review"},
				PolicyDigest:   "policy-route",
			},
			Result: Result{
				SpecID:         "route",
				Kind:           KindChoice,
				Value:          Value{Choice: "review"},
				Confidence:     &confident,
				EvidenceDigest: "evidence-route",
				Provider:       "jev",
				Model:          "jev-latest",
				Status:         StatusObserved,
				ObservedAt:     batchResultTime(),
			},
		},
		{
			Spec: Spec{
				ID:           "urgent",
				Question:     "Is this urgent?",
				Kind:         KindNoul,
				PolicyDigest: "policy-urgent",
			},
			Result: Result{
				SpecID:         "urgent",
				Kind:           KindNoul,
				Value:          Value{Noul: true},
				EvidenceDigest: "evidence-urgent",
				Provider:       "jev",
				Model:          "jev-latest",
				Status:         StatusObserved,
				ObservedAt:     batchResultTime(),
			},
		},
	}

	batch, err := ObserveBatch(state, observations)
	if err != nil {
		t.Fatalf("ObserveBatch() error = %v", err)
	}
	if got, want := batch.SpecIDs, []string{"route", "urgent"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("spec ids = %#v, want %#v", got, want)
	}
	if !batch.NonAuthorizing {
		t.Fatal("batch must remain non-authorizing")
	}
	if err := batch.Validate(); err != nil {
		t.Fatalf("DecisionBatchReceipt.Validate() error = %v", err)
	}
}

func TestObserveBatchRejectsDuplicateAndTamperedState(t *testing.T) {
	confident := 0.8
	spec := Spec{
		ID:             "route",
		Question:       "Which route should handle this?",
		Kind:           KindChoice,
		AllowedChoices: []string{"fast", "review"},
		PolicyDigest:   "policy-route",
	}
	result := Result{
		SpecID:         "route",
		Kind:           KindChoice,
		Value:          Value{Choice: "fast"},
		Confidence:     &confident,
		EvidenceDigest: "evidence-route",
		Provider:       "jev",
		Model:          "jev-latest",
		Status:         StatusObserved,
		ObservedAt:     batchResultTime(),
	}
	if _, err := ObserveBatch(State{Digest: "state-1"}, []BatchObservation{
		{Spec: spec, Result: result},
		{Spec: spec, Result: result},
	}); err == nil {
		t.Fatal("duplicate spec id unexpectedly accepted")
	}

	batch, err := ObserveBatch(State{Digest: "state-1"}, []BatchObservation{{Spec: spec, Result: result}})
	if err != nil {
		t.Fatalf("ObserveBatch() error = %v", err)
	}
	batch.StateDigest = "state-2"
	if err := batch.Validate(); err == nil {
		t.Fatal("tampered batch state unexpectedly validated")
	}
}

func TestObserveBatchJSONPreservesUnknownReceipt(t *testing.T) {
	input := []byte(`{
		"state": {"digest": "state-json"},
		"observations": [{
			"spec": {
				"id": "freshness",
				"question": "Is the evidence fresh?",
				"kind": "noul",
				"policy_digest": "policy-freshness"
			},
			"result": {
				"spec_id": "freshness",
				"kind": "noul",
				"value": {},
				"evidence_digest": "evidence-freshness",
				"provider": "jev",
				"status": "unknown",
				"observed_at": "2026-09-28T00:00:00Z"
			}
		}]
	}`)
	batch, err := ObserveBatchJSON(input)
	if err != nil {
		t.Fatalf("ObserveBatchJSON() error = %v", err)
	}
	if batch.Receipts[0].Status != StatusUnknown {
		t.Fatalf("status = %q, want %q", batch.Receipts[0].Status, StatusUnknown)
	}
	if err := batch.Validate(); err != nil {
		t.Fatalf("unknown batch validation error = %v", err)
	}
}

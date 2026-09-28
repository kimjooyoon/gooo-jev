package decision

import (
	"encoding/json"
	"strings"
	"testing"
)

func validDecisionJSON() string {
	return `{
		"spec": {
			"id": "candidate-review",
			"question": "Should this candidate be reviewed?",
			"kind": "choice",
			"allowed_choices": ["accept", "review"],
			"policy_digest": "policy-1"
		},
		"state": {"digest": "state-1"},
		"result": {
			"spec_id": "candidate-review",
			"kind": "choice",
			"value": {"choice": "review"},
			"confidence": 0.91,
			"evidence_digest": "evidence-1",
			"provider": "jev",
			"model": "jev-latest",
			"status": "observed",
			"observed_at": "2026-09-28T00:00:00Z"
		}
	}`
}

func TestObserveJSONProducesNonAuthorizingReceipt(t *testing.T) {
	receipt, err := ObserveJSON([]byte(validDecisionJSON()))
	if err != nil {
		t.Fatalf(`ObserveJSON() error = %v`, err)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf(`Receipt.Validate() error = %v`, err)
	}
	if !receipt.NonAuthorizing {
		t.Fatal(`JSON receipt must remain non-authorizing`)
	}

	encoded, err := MarshalReceiptJSON(receipt)
	if err != nil {
		t.Fatalf(`MarshalReceiptJSON() error = %v`, err)
	}
	var output DecisionReceiptJSON
	if err := json.Unmarshal(encoded, &output); err != nil {
		t.Fatalf(`decode receipt JSON: %v`, err)
	}
	if output.DecisionDigest != receipt.DecisionDigest {
		t.Fatalf(`decision digest = %q, want %q`, output.DecisionDigest, receipt.DecisionDigest)
	}
	if !output.NonAuthorizing {
		t.Fatal(`serialized JSON receipt must remain non-authorizing`)
	}
}

func TestObserveJSONPreservesUnknown(t *testing.T) {
	input := strings.Replace(
		validDecisionJSON(),
		`"status": "observed"`,
		`"status": "unknown"`,
		1,
	)
	receipt, err := ObserveJSON([]byte(input))
	if err != nil {
		t.Fatalf(`ObserveJSON() error = %v`, err)
	}
	if receipt.Status != StatusUnknown {
		t.Fatalf(`status = %q, want %q`, receipt.Status, StatusUnknown)
	}
}

func TestObserveJSONRejectsUnknownFieldsAndTrailingValues(t *testing.T) {
	unknownField := strings.Replace(
		validDecisionJSON(),
		`"digest": "state-1"`,
		`"digest": "state-1", "unexpected": true`,
		1,
	)
	if _, err := ObserveJSON([]byte(unknownField)); err == nil {
		t.Fatal(`unknown JSON field unexpectedly accepted`)
	}
	if _, err := ObserveJSON([]byte(validDecisionJSON() + `{}`)); err == nil {
		t.Fatal(`trailing JSON value unexpectedly accepted`)
	}
}

func TestObserveJSONRejectsChoiceOutsideSchema(t *testing.T) {
	input := strings.Replace(
		validDecisionJSON(),
		`"choice": "review"`,
		`"choice": "defer"`,
		1,
	)
	if _, err := ObserveJSON([]byte(input)); err == nil {
		t.Fatal(`choice outside the declared schema unexpectedly accepted`)
	}
}

package decision

import "testing"

func TestRecordExecutionEnvelopeReviewOutcomeBindsExplicitOutcomes(t *testing.T) {
	for _, outcome := range []string{"accepted", "rejected", "hold"} {
		result := RecordExecutionEnvelopeReviewOutcome(ExecutionEnvelopeReviewOutcomeInput{
			HandoffStatus:         "handoff-required",
			HandoffDigest:         "handoff",
			ReviewRequired:        true,
			OutcomeStatus:         outcome,
			OutcomeEvidenceDigest: "evidence",
			NonAuthorizing:        true,
		})
		if result.Status != outcome || result.OutcomeDigest == "" || result.MissingStage != "" || !result.NonAuthorizing {
			t.Fatalf("%s outcome was not bound: %#v", outcome, result)
		}
	}
}

func TestRecordExecutionEnvelopeReviewOutcomeHoldsMissingEvidence(t *testing.T) {
	result := RecordExecutionEnvelopeReviewOutcome(ExecutionEnvelopeReviewOutcomeInput{
		HandoffStatus:  "handoff-required",
		HandoffDigest:  "handoff",
		ReviewRequired: true,
		NonAuthorizing: true,
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "review-outcome" || result.OutcomeDigest != "" || !result.NonAuthorizing {
		t.Fatalf("missing review outcome escaped UNKNOWN: %#v", result)
	}
}

func TestRecordExecutionEnvelopeReviewOutcomeHoldsInvalidHandoff(t *testing.T) {
	result := RecordExecutionEnvelopeReviewOutcome(ExecutionEnvelopeReviewOutcomeInput{
		HandoffStatus:  "observation-only",
		HandoffDigest:  "observation",
		NonAuthorizing: true,
	})
	if result.Status != "hold" || result.MissingStage != "review-handoff" || !result.NonAuthorizing {
		t.Fatalf("invalid handoff escaped hold: %#v", result)
	}
}

func TestRecordExecutionEnvelopeReviewOutcomeRejectsAuthorizationClaim(t *testing.T) {
	result := RecordExecutionEnvelopeReviewOutcome(ExecutionEnvelopeReviewOutcomeInput{
		HandoffStatus:         "handoff-required",
		HandoffDigest:         "handoff",
		ReviewRequired:        true,
		OutcomeStatus:         "accepted",
		OutcomeEvidenceDigest: "evidence",
		NonAuthorizing:        false,
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "authorization-boundary" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped review outcome: %#v", result)
	}
}

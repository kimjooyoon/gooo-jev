package decision

import "testing"

func reviewObservationDecision(status string) ExecutionEnvelopeJEVImprovementReviewGateDecision {
	return ExecutionEnvelopeJEVImprovementReviewGateDecision{
		Status: status, DecisionDigest: "decision-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func TestObserveExecutionEnvelopeJEVImprovementReview(t *testing.T) {
	for _, outcome := range []string{"confirmed", "refuted", "unknown"} {
		t.Run(outcome, func(t *testing.T) {
			got := ObserveExecutionEnvelopeJEVImprovementReview(ExecutionEnvelopeJEVImprovementReviewObservationInput{
				GateDecision:         reviewObservationDecision("review-eligible"),
				ReviewOutcome:        outcome,
				ReviewEvidenceDigest: "review-evidence-digest",
				NonAuthorizing:       true,
			})
			if got.Status != "observed" || got.ReviewOutcome != outcome || got.ObservationDigest == "" || got.GateDecisionDigest != "decision-digest" {
				t.Fatalf("got %+v", got)
			}
			if !got.NonExecuting || !got.NonAuthorizing {
				t.Fatalf("missing boundary %+v", got)
			}
		})
	}
}

func TestObserveExecutionEnvelopeJEVImprovementReviewPreservesUnknownStage(t *testing.T) {
	got := ObserveExecutionEnvelopeJEVImprovementReview(ExecutionEnvelopeJEVImprovementReviewObservationInput{
		GateDecision: reviewObservationDecision("UNKNOWN"),
		NonAuthorizing: true,
	})
	got.GateDecision.MissingStage = "feedback-aggregation"
	if got.Status != "UNKNOWN" || got.MissingStage != "review-disposition" || got.ObservationDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestObserveExecutionEnvelopeJEVImprovementReviewRequiresEvidence(t *testing.T) {
	got := ObserveExecutionEnvelopeJEVImprovementReview(ExecutionEnvelopeJEVImprovementReviewObservationInput{
		GateDecision:   reviewObservationDecision("review-hold"),
		ReviewOutcome:  "unknown",
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "review-evidence" || got.ObservationDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestObserveExecutionEnvelopeJEVImprovementReviewRejectsAuthorization(t *testing.T) {
	got := ObserveExecutionEnvelopeJEVImprovementReview(ExecutionEnvelopeJEVImprovementReviewObservationInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.ObservationDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

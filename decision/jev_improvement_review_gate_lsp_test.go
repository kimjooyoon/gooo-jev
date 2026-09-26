package decision

import "testing"

func reviewGateLSPDecision(status string) ExecutionEnvelopeJEVImprovementReviewGateDecision {
	return ExecutionEnvelopeJEVImprovementReviewGateDecision{
		Status: "review-eligible", FeedbackStatus: status, RevisionCandidateDigest: "candidate-digest", BindingDigest: "binding-digest", DecisionDigest: "decision-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReviewGateLSP(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   string
		code   string
	}{
		{name: "eligible", status: "stable-for-review", want: "review-eligible", code: "JEV_REVIEW_ELIGIBLE"},
		{name: "revision", status: "needs-revision", want: "revision-required", code: "JEV_REVISION_REQUIRED"},
		{name: "hold", status: "hold", want: "review-hold", code: "JEV_REVIEW_HOLD"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := reviewGateLSPDecision(test.status)
			if test.want != "review-eligible" {
				decision.Status = test.want
			}
			got := ProjectExecutionEnvelopeJEVImprovementReviewGateLSP(ExecutionEnvelopeJEVImprovementReviewGateLSPInput{Decision: decision, NonAuthorizing: true})
			if got.Status != test.want || got.Code != test.code || got.EvidenceDigest == "" || got.DecisionDigest != "decision-digest" {
				t.Fatalf("got %+v", got)
			}
			if !got.NonExecuting || !got.NonAuthorizing {
				t.Fatalf("missing boundary %+v", got)
			}
		})
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReviewGateLSPPreservesUnknown(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementReviewGateLSP(ExecutionEnvelopeJEVImprovementReviewGateLSPInput{
		Decision:       ExecutionEnvelopeJEVImprovementReviewGateDecision{Status: "UNKNOWN", MissingStage: "feedback-aggregation", DecisionDigest: "decision-digest", NonExecuting: true, NonAuthorizing: true},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-aggregation" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReviewGateLSPRequiresDecisionDigest(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementReviewGateLSP(ExecutionEnvelopeJEVImprovementReviewGateLSPInput{
		Decision:       ExecutionEnvelopeJEVImprovementReviewGateDecision{Status: "review-eligible", NonExecuting: true, NonAuthorizing: true},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "review-gate-decision" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReviewGateLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementReviewGateLSP(ExecutionEnvelopeJEVImprovementReviewGateLSPInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

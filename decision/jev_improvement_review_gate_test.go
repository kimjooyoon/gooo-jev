package decision

import "testing"

func reviewGateBinding(feedbackStatus string) ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding {
	return ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding{
		Status:                           "bound",
		PlanID:                           "triage-plan",
		CycleBindingDigest:               "cycle-binding-digest",
		RevisionCandidateDigest:          "candidate-digest",
		RevisionCandidateEvidenceDigest:  "candidate-evidence-digest",
		FeedbackStatus:                   feedbackStatus,
		FeedbackEvidenceDigest:           "feedback-evidence-digest",
		BindingDigest:                    "binding-digest",
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
}

func TestEvaluateExecutionEnvelopeJEVImprovementReviewGate(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   string
	}{
		{name: "stable", status: "stable-for-review", want: "review-eligible"},
		{name: "needs-revision", status: "needs-revision", want: "revision-required"},
		{name: "hold", status: "hold", want: "review-hold"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := EvaluateExecutionEnvelopeJEVImprovementReviewGate(ExecutionEnvelopeJEVImprovementReviewGateInput{
				Binding:        reviewGateBinding(test.status),
				NonAuthorizing: true,
			})
			if got.Status != test.want || got.FeedbackStatus != test.status || got.DecisionDigest == "" {
				t.Fatalf("got %+v", got)
			}
			if !got.NonExecuting || !got.NonAuthorizing || got.RevisionCandidateDigest != "candidate-digest" {
				t.Fatalf("missing boundary %+v", got)
			}
		})
	}
}

func TestEvaluateExecutionEnvelopeJEVImprovementReviewGatePreservesUnknown(t *testing.T) {
	got := EvaluateExecutionEnvelopeJEVImprovementReviewGate(ExecutionEnvelopeJEVImprovementReviewGateInput{
		Binding: ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding{
			Status: "UNKNOWN", MissingStage: "feedback-aggregation", NonExecuting: true, NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-aggregation" || got.DecisionDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluateExecutionEnvelopeJEVImprovementReviewGateRejectsIncompleteCandidate(t *testing.T) {
	binding := reviewGateBinding("stable-for-review")
	binding.RevisionCandidateEvidenceDigest = ""
	got := EvaluateExecutionEnvelopeJEVImprovementReviewGate(ExecutionEnvelopeJEVImprovementReviewGateInput{
		Binding:        binding,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-candidate" || got.DecisionDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluateExecutionEnvelopeJEVImprovementReviewGateRejectsAuthorization(t *testing.T) {
	got := EvaluateExecutionEnvelopeJEVImprovementReviewGate(ExecutionEnvelopeJEVImprovementReviewGateInput{
		Binding:        reviewGateBinding("stable-for-review"),
		NonAuthorizing: false,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.DecisionDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

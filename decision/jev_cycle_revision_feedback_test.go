package decision

import "testing"

func cycleRevisionFeedbackBinding(t *testing.T, feedback JEVImprovementFeedbackAggregation) ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding {
	t.Helper()
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:           revisionCandidateDirective(),
		RevisionSource:      "gooo://revision/source/one",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:      true,
	})
	if err := candidate.Validate(); err != nil {
		t.Fatal(err)
	}
	return BindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput{
		CycleBinding: ExecutionEnvelopeJEVPlanImprovementCycleBinding{
			Status:                 "bound",
			PlanID:                 "triage-plan",
			PlanDigest:             "plan-digest",
			LifecycleBindingDigest: "lifecycle-binding-digest",
			CycleStatus:            "stable-for-review",
			CycleEvidenceDigest:    "cycle-evidence-digest",
			BindingDigest:          "cycle-binding-digest",
			NonExecuting:           true,
			NonAuthorizing:         true,
		},
		RevisionCandidate: candidate,
		Feedback:          feedback,
		NonAuthorizing:    true,
	})
}

func TestBindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback(t *testing.T) {
	feedback := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
		Feedback: []JEVImprovementReplayFeedback{
			aggregateReplayFeedback(jevReplayFeedbackConfirmed),
			aggregateReplayFeedback(jevReplayFeedbackConfirmed),
		},
		NonAuthorizing: true,
	})
	got := cycleRevisionFeedbackBinding(t, feedback)
	if got.Status != "bound" || got.PlanID != "triage-plan" || got.FeedbackStatus != jevImprovementFeedbackStableForReview || got.BindingDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing || got.RevisionCandidateEvidenceDigest == "" || got.FeedbackEvidenceDigest == "" {
		t.Fatalf("missing evidence %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVImprovementCycleToRevisionFeedbackPreservesNeedsRevision(t *testing.T) {
	feedback := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
		Feedback: []JEVImprovementReplayFeedback{
			aggregateReplayFeedback(jevReplayFeedbackConfirmed),
			aggregateReplayFeedback(jevReplayFeedbackRefuted),
		},
		NonAuthorizing: true,
	})
	got := cycleRevisionFeedbackBinding(t, feedback)
	if got.Status != "bound" || got.FeedbackStatus != jevImprovementFeedbackNeedsRevision || got.BindingDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVImprovementCycleToRevisionFeedbackRejectsTampering(t *testing.T) {
	feedback := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
		Feedback:       []JEVImprovementReplayFeedback{aggregateReplayFeedback(jevReplayFeedbackConfirmed)},
		NonAuthorizing: true,
	})
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:           revisionCandidateDirective(),
		RevisionSource:      "gooo://revision/source/one",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:      true,
	})
	candidate.EvidenceDigest = "tampered"
	got := BindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput{
		CycleBinding: ExecutionEnvelopeJEVPlanImprovementCycleBinding{Status: "bound", PlanID: "triage-plan", BindingDigest: "cycle-binding-digest", NonExecuting: true, NonAuthorizing: true},
		RevisionCandidate: candidate,
		Feedback:          feedback,
		NonAuthorizing:    true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-candidate" || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVImprovementCycleToRevisionFeedbackPreservesMissingStage(t *testing.T) {
	feedback := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
		Feedback:       []JEVImprovementReplayFeedback{aggregateReplayFeedback(jevReplayFeedbackConfirmed)},
		NonAuthorizing: true,
	})
	got := BindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput{
		CycleBinding: ExecutionEnvelopeJEVPlanImprovementCycleBinding{Status: "UNKNOWN", MissingStage: "suspend-resume-source", NonExecuting: true, NonAuthorizing: true},
		Feedback:     feedback,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "suspend-resume-source" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVImprovementCycleToRevisionFeedbackRequiresAuthorizationBoundary(t *testing.T) {
	got := BindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput{
		NonAuthorizing: false,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

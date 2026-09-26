package decision

import "testing"

func cycleCandidateFeedbackBinding(status string) ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding {
	return ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding{
		Status: "bound", CycleStatus: "stable-for-review", CycleEvidenceDigest: "cycle-evidence-digest", CandidateStatus: "revision-candidate-ready-for-review", CandidateDigest: "candidate-digest", CandidateEvidenceDigest: "candidate-evidence-digest", BoundRevisionChangeDigest: "bound-revision-change-digest", BindingDigest: "binding-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func cycleCandidateFeedbackReverse(status string) ExecutionEnvelopeReverseObservationBinding {
	if status == "counterexample" {
		return BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{ObservedStatus: "ready", ExpectedStatus: "hold", ObservedEvidenceDigest: "counterexample-evidence", ExpectedEvidenceDigest: "expected-evidence", NonAuthorizing: true}))
	}
	return BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{ObservedStatus: "ready", ExpectedStatus: "ready", ObservedEvidenceDigest: "reproduced-evidence", ExpectedEvidenceDigest: "reproduced-evidence", NonAuthorizing: true}))
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(t *testing.T) {
	for _, status := range []string{"reproduced", "counterexample"} {
		t.Run(status, func(t *testing.T) {
			got := DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput{
				CandidateBinding: cycleCandidateFeedbackBinding(status),
				ReverseBinding:   cycleCandidateFeedbackReverse(status),
				MetricDigest:     "metric-digest",
				NonAuthorizing:   true,
			})
			want := jevReplayFeedbackConfirmed
			if status == "counterexample" {
				want = jevReplayFeedbackRefuted
			}
			if got.Status != want || got.FeedbackKind != want || got.EvidenceDigest == "" || got.MetricDigest == "metric-digest" {
				t.Fatalf("got %+v", got)
			}
			if err := got.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidatePreservesStage(t *testing.T) {
	binding := cycleCandidateFeedbackBinding("reproduced")
	binding.Status = "UNKNOWN"
	binding.MissingStage = "feedback-hold"
	got := DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput{
		CandidateBinding: binding,
		ReverseBinding:   cycleCandidateFeedbackReverse("reproduced"),
		MetricDigest:     "metric-digest",
		NonAuthorizing:   true,
	})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "feedback-hold" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidateRequiresMetric(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput{
		CandidateBinding: cycleCandidateFeedbackBinding("reproduced"),
		ReverseBinding:   cycleCandidateFeedbackReverse("reproduced"),
		NonAuthorizing:   true,
	})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "metric" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidateRejectsAuthorization(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput{NonAuthorizing: false})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

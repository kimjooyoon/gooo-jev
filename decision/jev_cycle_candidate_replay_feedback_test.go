package decision

import "testing"

func cycleCandidateReplayCycle() JEVImprovementCycleObservation {
	cycle := JEVImprovementCycleObservation{
		Status: "needs-revision",
		DeclarationDigest: "declaration-digest",
		IRDigest: "ir-digest",
		GenerationDigest: "generation-digest",
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest: "metric-digest",
		ChangePlanDigest: "change-plan-digest",
		DispositionStatus: "needs-revision",
		LSPCode: "jev.change-plan.abort",
		NonExecuting: true,
		NonAuthorizing: true,
	}
	cycle.EvidenceDigest = digestJEVImprovementCycleObservation(cycle.DeclarationDigest, cycle.IRDigest, cycle.GenerationDigest, cycle.ReverseObservationDigest, cycle.MetricDigest, cycle.ChangePlanDigest, cycle.DispositionStatus, cycle.LSPCode)
	return cycle
}

func cycleCandidateReplayCandidate(cycle JEVImprovementCycleObservation) ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding {
	return ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding{
		Status: "bound",
		CycleStatus: cycle.Status,
		CycleEvidenceDigest: cycle.EvidenceDigest,
		CandidateStatus: "revision-candidate-ready-for-review",
		CandidateDigest: "candidate-digest",
		CandidateEvidenceDigest: "candidate-evidence-digest",
		BoundRevisionChangeDigest: "bound-revision-change-digest",
		BindingDigest: "binding-digest",
		NonExecuting: true,
		NonAuthorizing: true,
	}
}

func cycleCandidateReplayReverse(status string) ExecutionEnvelopeReverseObservationBinding {
	observedEvidence := status + "-evidence"
	expectedEvidence := observedEvidence
	expectedStatus := "ready"
	if status == "counterexample" {
		expectedStatus = "hold"
		expectedEvidence = "expected-evidence"
	}
	return BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus: "ready",
		ExpectedStatus: expectedStatus,
		ObservedEvidenceDigest: observedEvidence,
		ExpectedEvidenceDigest: expectedEvidence,
		NonAuthorizing: true,
	}))
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservation(t *testing.T) {
	cycle := cycleCandidateReplayCycle()
	for _, status := range []string{"reproduced", "counterexample"} {
		t.Run(status, func(t *testing.T) {
			got := DeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservation(ExecutionEnvelopeJEVCycleCandidateReplayFeedbackInput{
				Cycle: cycle,
				Candidate: cycleCandidateReplayCandidate(cycle),
				ReverseBinding: cycleCandidateReplayReverse(status),
				MetricDigest: "metric-digest",
				NonAuthorizing: true,
			})
			want := jevReplayFeedbackConfirmed
			if status == "counterexample" {
				want = jevReplayFeedbackRefuted
			}
			if got.Status != want || got.FeedbackKind != want || got.EvidenceDigest == "" || got.MetricDigest == "" {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservationRejectsMixedCycle(t *testing.T) {
	cycle := cycleCandidateReplayCycle()
	candidate := cycleCandidateReplayCandidate(cycle)
	candidate.CycleEvidenceDigest = "other-cycle"
	got := DeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservation(ExecutionEnvelopeJEVCycleCandidateReplayFeedbackInput{
		Cycle: cycle,
		Candidate: candidate,
		ReverseBinding: cycleCandidateReplayReverse("reproduced"),
		MetricDigest: "metric-digest",
		NonAuthorizing: true,
	})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "candidate-cycle-binding" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservationRejectsAuthorization(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservation(ExecutionEnvelopeJEVCycleCandidateReplayFeedbackInput{})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

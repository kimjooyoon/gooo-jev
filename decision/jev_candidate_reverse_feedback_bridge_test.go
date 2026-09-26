package decision

import "testing"

func feedbackBridgeGateBinding(status string) ExecutionEnvelopeJEVCandidateReverseObservationGateBinding {
	return ExecutionEnvelopeJEVCandidateReverseObservationGateBinding{
		Status: "bound", CandidateDigest: "candidate-digest", CandidateEvidenceDigest: "candidate-evidence-digest", AdmissionDigest: "admission-digest", ReverseObservationStatus: status, ReverseObservationDigest: "reverse-observation-digest", ReverseEvidenceDigest: "reverse-evidence-digest", BindingDigest: "binding-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput{
		GateBinding:    feedbackBridgeGateBinding("reproduced"),
		MetricDigest:   "metric-digest",
		NonAuthorizing: true,
	})
	if got.Status != jevReplayFeedbackConfirmed || got.FeedbackKind != jevReplayFeedbackConfirmed || got.CandidateDigest != "candidate-digest" || got.ReplayObservationDigest != "reverse-observation-digest" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCandidateReverseObservationPreservesCounterexample(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput{
		GateBinding:    feedbackBridgeGateBinding("counterexample"),
		MetricDigest:   "metric-digest",
		NonAuthorizing: true,
	})
	if got.Status != jevReplayFeedbackRefuted || got.FeedbackKind != jevReplayFeedbackRefuted || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCandidateReverseObservationRequiresMetric(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput{
		GateBinding:    feedbackBridgeGateBinding("reproduced"),
		NonAuthorizing: true,
	})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "metric" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCandidateReverseObservationPreservesGateStage(t *testing.T) {
	binding := feedbackBridgeGateBinding("reproduced")
	binding.Status = "UNKNOWN"
	binding.MissingStage = "metric-evidence"
	got := DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput{
		GateBinding:    binding,
		MetricDigest:   "metric-digest",
		NonAuthorizing: true,
	})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "metric-evidence" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeriveJEVImprovementReplayFeedbackFromCandidateReverseObservationRejectsAuthorization(t *testing.T) {
	got := DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput{NonAuthorizing: false})
	if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

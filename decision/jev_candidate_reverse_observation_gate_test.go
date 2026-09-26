package decision

import "testing"

func candidateReverseGateCandidateBinding() ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding {
	return ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding{
		Status: "bound", CandidateDigest: "candidate-digest", CandidateEvidenceDigest: "candidate-evidence-digest", SummaryDigest: "summary-digest", MetricSourceDigest: "metric-source-digest", MetricValueDigest: "metric-value-digest", AdmissionDigest: "admission-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func candidateReverseGateReproduced() ExecutionEnvelopeReverseObservationBinding {
	return BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus: "ready", ExpectedStatus: "ready", ObservedEvidenceDigest: "reverse-evidence", ExpectedEvidenceDigest: "reverse-evidence", NonAuthorizing: true,
	}))
}

func candidateReverseGateCounterexample() ExecutionEnvelopeReverseObservationBinding {
	return BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus: "ready", ExpectedStatus: "hold", ObservedEvidenceDigest: "counterexample-evidence", ExpectedEvidenceDigest: "expected-evidence", NonAuthorizing: true,
	}))
}

func TestBindExecutionEnvelopeJEVCandidateToReverseObservation(t *testing.T) {
	got := BindExecutionEnvelopeJEVCandidateToReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationGateInput{
		CandidateBinding: candidateReverseGateCandidateBinding(),
		ReverseBinding:   candidateReverseGateReproduced(),
		NonAuthorizing:   true,
	})
	if got.Status != "bound" || got.ReverseObservationStatus != "reproduced" || got.BindingDigest == "" || got.ReverseObservationDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing || got.CandidateDigest != "candidate-digest" {
		t.Fatalf("missing boundary %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVCandidateToReverseObservationPreservesCounterexample(t *testing.T) {
	got := BindExecutionEnvelopeJEVCandidateToReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationGateInput{
		CandidateBinding: candidateReverseGateCandidateBinding(),
		ReverseBinding:   candidateReverseGateCounterexample(),
		NonAuthorizing:   true,
	})
	if got.Status != "bound" || got.ReverseObservationStatus != "counterexample" || got.ReverseEvidenceDigest != "counterexample-evidence" || got.BindingDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVCandidateToReverseObservationPreservesUnknown(t *testing.T) {
	unknown := BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{NonAuthorizing: true}))
	got := BindExecutionEnvelopeJEVCandidateToReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationGateInput{
		CandidateBinding: candidateReverseGateCandidateBinding(),
		ReverseBinding:   unknown,
		NonAuthorizing:   true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "status" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVCandidateToReverseObservationRejectsTamperedReverse(t *testing.T) {
	reverse := candidateReverseGateReproduced()
	reverse.ObservationDigest = "tampered"
	got := BindExecutionEnvelopeJEVCandidateToReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationGateInput{
		CandidateBinding: candidateReverseGateCandidateBinding(),
		ReverseBinding:   reverse,
		NonAuthorizing:   true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "reverse-observation" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVCandidateToReverseObservationRejectsAuthorization(t *testing.T) {
	got := BindExecutionEnvelopeJEVCandidateToReverseObservation(ExecutionEnvelopeJEVCandidateReverseObservationGateInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

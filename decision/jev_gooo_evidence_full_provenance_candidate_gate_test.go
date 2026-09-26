package decision

import "testing"

func evidenceFullProvenanceCandidateGateInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateInput {
	t.Helper()
	metricInput := evidenceFullProvenanceMetricInput(t)
	metric := BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(metricInput)
	candidateSource := "package candidate\nnamespace candidate\nentity Revision id \"gooo://candidate/revision\"\n"
	candidateDigest, err := digestGoooEvidenceFullProvenanceCandidateSource(candidateSource)
	if err != nil {
		t.Fatalf("candidate source digest failed: %v", err)
	}
	return ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateInput{
		Metric:                  metric,
		CandidateSource:         candidateSource,
		ProposedCandidateDigest: candidateDigest,
		NonAuthorizing:         true,
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateAdmitsConfirmed(t *testing.T) {
	gate := BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(evidenceFullProvenanceCandidateGateInput(t))
	if gate.Status != "ready" || gate.AdmissionStatus != "ADMIT" ||
		gate.CandidateSourceDigest == "" || gate.MetricEvidenceDigest == "" ||
		gate.EvidenceDigest == "" || !gate.NonExecuting || !gate.NonAuthorizing {
		t.Fatalf("gate = %#v, want confirmed ADMIT", gate)
	}
	if err := gate.Validate(); err != nil {
		t.Fatalf("gate should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateHoldsRefuted(t *testing.T) {
	input := evidenceFullProvenanceCandidateGateInput(t)
	metricInput := evidenceFullProvenanceMetricInput(t)
	metricInput.ObservedEvidenceDigest = "different-observation"
	input.Metric = BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(metricInput)
	gate := BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(input)
	if gate.Status != "ready" || gate.AdmissionStatus != "HOLD" {
		t.Fatalf("gate = %#v, want refuted HOLD", gate)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateRejectsDigestTampering(t *testing.T) {
	input := evidenceFullProvenanceCandidateGateInput(t)
	input.ProposedCandidateDigest = "tampered"
	gate := BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(input)
	if gate.Status != "UNKNOWN" || gate.MissingStage != "candidate-source-digest" {
		t.Fatalf("gate = %#v, want candidate-source-digest UNKNOWN", gate)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateRejectsUnknownMetric(t *testing.T) {
	input := evidenceFullProvenanceCandidateGateInput(t)
	input.Metric.Status = "UNKNOWN"
	gate := BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(input)
	if gate.Status != "UNKNOWN" || gate.MissingStage != "metric" {
		t.Fatalf("gate = %#v, want metric UNKNOWN", gate)
	}
}

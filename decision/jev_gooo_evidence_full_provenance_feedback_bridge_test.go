package decision

import "testing"

func evidenceFullProvenanceFeedbackBridgeInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeInput {
	t.Helper()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeInput{
		CandidateGate:        BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(evidenceFullProvenanceCandidateGateInput(t)),
		ReplayObservationDigest: "replay-observation-evidence",
		NonAuthorizing:       true,
	}
}

func TestAppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeConfirmed(t *testing.T) {
	ledger := AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge(evidenceFullProvenanceFeedbackBridgeInput(t))
	if ledger.Status != jevReplayFeedbackLedgerReady ||
		ledger.FeedbackStatus != jevReplayFeedbackConfirmed ||
		ledger.ConfirmedCount != 1 || ledger.RefutedCount != 0 ||
		ledger.UnknownCount != 0 {
		t.Fatalf("ledger = %#v, want confirmed append", ledger)
	}
	if err := ledger.Validate(); err != nil {
		t.Fatalf("ledger should validate: %v", err)
	}
}

func TestAppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeRefuted(t *testing.T) {
	input := evidenceFullProvenanceFeedbackBridgeInput(t)
	metricInput := evidenceFullProvenanceMetricInput(t)
	metricInput.ObservedEvidenceDigest = "different-observation"
	gateInput := evidenceFullProvenanceCandidateGateInput(t)
	gateInput.Metric = BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(metricInput)
	input.CandidateGate = BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(gateInput)
	ledger := AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge(input)
	if ledger.Status != jevReplayFeedbackLedgerReady ||
		ledger.FeedbackStatus != jevReplayFeedbackRefuted ||
		ledger.ConfirmedCount != 0 || ledger.RefutedCount != 1 {
		t.Fatalf("ledger = %#v, want refuted append", ledger)
	}
}

func TestAppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeKeepsUnknown(t *testing.T) {
	input := evidenceFullProvenanceFeedbackBridgeInput(t)
	input.CandidateGate.Status = "UNKNOWN"
	ledger := AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge(input)
	if ledger.Status != jevReplayFeedbackLedgerUnknown ||
		ledger.MissingStage != "candidate-gate-validation" {
		t.Fatalf("ledger = %#v, want candidate-gate-validation UNKNOWN", ledger)
	}
}

func TestAppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeRequiresObservation(t *testing.T) {
	input := evidenceFullProvenanceFeedbackBridgeInput(t)
	input.ReplayObservationDigest = ""
	ledger := AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge(input)
	if ledger.Status != jevReplayFeedbackLedgerUnknown ||
		ledger.MissingStage != "replay-observation" {
		t.Fatalf("ledger = %#v, want replay-observation UNKNOWN", ledger)
	}
}

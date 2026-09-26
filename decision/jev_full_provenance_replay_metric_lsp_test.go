package decision

import "testing"

func replayMetricLSPInput(metric JEVReplayLedgerMetric) ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPInput {
	binding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		ReplayMetric:   metric,
		NonAuthorizing: true,
	})
	return ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPInput{
		Binding:        binding,
		ReplayMetric:   metric,
		NonAuthorizing: true,
	}
}

func TestProjectJEVFullProvenanceReplayMetricToLSP(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	projected := ProjectJEVFullProvenanceReplayMetricToLSP(replayMetricLSPInput(metric))
	if projected.Status != "ready" || !projected.Publishable ||
		projected.Severity != "info" || projected.Code != "jev.replay-metric.stable" ||
		projected.EvidenceDigest == "" || projected.EvidencePrefixDigest == "" {
		t.Fatalf("projected = %#v, want stable metric projection", projected)
	}
}

func TestProjectJEVFullProvenanceReplayMetricToLSPShowsRegression(t *testing.T) {
	ledger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       replayLedgerCheckpointFixture(),
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackRefuted, "metric-b"),
		MetricDigest:   "metric-b",
		NonAuthorizing: true,
	})
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         ledger,
		NonAuthorizing: true,
	})
	projected := ProjectJEVFullProvenanceReplayMetricToLSP(replayMetricLSPInput(metric))
	if projected.Status != "diagnostic" || projected.Code != "jev.replay-metric.refuted" ||
		projected.RefutedPermille != 500 {
		t.Fatalf("projected = %#v, want refuted diagnostic", projected)
	}
}

func TestProjectJEVFullProvenanceReplayMetricToLSPShowsUnknown(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger: AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
			Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackUnknown, "metric-unknown"),
			MetricDigest:   "metric-unknown",
			NonAuthorizing: true,
		}),
		NonAuthorizing: true,
	})
	projected := ProjectJEVFullProvenanceReplayMetricToLSP(replayMetricLSPInput(metric))
	if projected.Status != "diagnostic" || projected.Code != "jev.replay-metric.unknown" ||
		projected.UnknownPermille != 1000 {
		t.Fatalf("projected = %#v, want unknown diagnostic", projected)
	}
}

func TestProjectJEVFullProvenanceReplayMetricToLSPRejectsTampering(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	input := replayMetricLSPInput(metric)
	input.Binding.EvidenceDigest = "tampered"
	projected := ProjectJEVFullProvenanceReplayMetricToLSP(input)
	if projected.Status != "UNKNOWN" || projected.Publishable ||
		projected.Code != "jev-replay-metric-binding-integrity" {
		t.Fatalf("projected = %#v, want integrity UNKNOWN", projected)
	}
}

func TestProjectJEVFullProvenanceReplayMetricToLSPRejectsAuthorization(t *testing.T) {
	projected := ProjectJEVFullProvenanceReplayMetricToLSP(ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPInput{
		NonAuthorizing: false,
	})
	if projected.Status != "UNKNOWN" || projected.NonAuthorizing ||
		projected.Code != "authorization-boundary" {
		t.Fatalf("projected = %#v, want authorization UNKNOWN", projected)
	}
}
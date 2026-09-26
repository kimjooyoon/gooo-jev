package decision

import "testing"

func TestBindJEVFullProvenanceReplayMetric(t *testing.T) {
	binding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		ReplayMetric:   MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{Ledger: replayLedgerCheckpointFixture(), NonAuthorizing: true}),
		NonAuthorizing: true,
	})
	if binding.Status != "complete" || binding.ProvenanceMetricDigest != "metric-a" ||
		binding.ReplayMetricDigest == "" || binding.MetricBindingDigest == "" ||
		binding.EvidenceDigest == "" {
		t.Fatalf("binding = %#v, want complete dual-metric binding", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestBindJEVFullProvenanceReplayMetricDistinguishesReplayMetrics(t *testing.T) {
	first := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	secondLedger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       replayLedgerCheckpointFixture(),
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackRefuted, "metric-b"),
		MetricDigest:   "metric-b",
		NonAuthorizing: true,
	})
	second := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         secondLedger,
		NonAuthorizing: true,
	})
	firstBinding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		ReplayMetric:   first,
		NonAuthorizing: true,
	})
	secondBinding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		ReplayMetric:   second,
		NonAuthorizing: true,
	})
	if firstBinding.Status != "complete" || secondBinding.Status != "complete" ||
		firstBinding.MetricBindingDigest == secondBinding.MetricBindingDigest {
		t.Fatalf("metric bindings did not change: first=%#v second=%#v", firstBinding, secondBinding)
	}
}

func TestBindJEVFullProvenanceReplayMetricRejectsMissingProvenance(t *testing.T) {
	provenance := fullProvenanceReplayLedgerFixture()
	provenance.ReverseObservationDigest = ""
	binding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		Provenance: provenance,
		ReplayMetric: MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{Ledger: replayLedgerCheckpointFixture(), NonAuthorizing: true}),
		NonAuthorizing: true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "reverse-observation" {
		t.Fatalf("binding = %#v, want reverse-observation UNKNOWN", binding)
	}
}

func TestBindJEVFullProvenanceReplayMetricRejectsTampering(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	metric.EvidenceDigest = "tampered"
	binding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		ReplayMetric:   metric,
		NonAuthorizing: true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "replay-metric" {
		t.Fatalf("binding = %#v, want replay-metric UNKNOWN", binding)
	}
}

func TestBindJEVFullProvenanceReplayMetricRejectsAuthorization(t *testing.T) {
	binding := BindJEVFullProvenanceReplayMetric(ExecutionEnvelopeJEVFullProvenanceReplayMetricInput{
		NonAuthorizing: false,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "authorization-boundary" || binding.NonAuthorizing {
		t.Fatalf("binding = %#v, want authorization UNKNOWN", binding)
	}
}
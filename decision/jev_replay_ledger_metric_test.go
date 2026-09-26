package decision

import "testing"

func TestMeasureJEVReplayLedgerMetric(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	if metric.Status != "ready" || metric.Total != 1 || metric.Confirmed != 1 ||
		metric.ConfirmedPermille != 1000 || metric.RefutedPermille != 0 ||
		metric.UnknownPermille != 0 || metric.InputLedgerDigest == "" {
		t.Fatalf("metric = %#v, want one confirmed result", metric)
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("metric should validate: %v", err)
	}

	ledger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       replayLedgerCheckpointFixture(),
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackRefuted, "metric-b"),
		MetricDigest:   "metric-b",
		NonAuthorizing: true,
	})
	mixed := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         ledger,
		NonAuthorizing: true,
	})
	if mixed.Status != "ready" || mixed.Total != 2 || mixed.Confirmed != 1 ||
		mixed.Refuted != 1 || mixed.ConfirmedPermille != 500 ||
		mixed.RefutedPermille != 500 {
		t.Fatalf("mixed metric = %#v, want 500/500 split", mixed)
	}
}

func TestMeasureJEVReplayLedgerMetricPreservesUnknownRatio(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger: AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
			Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackUnknown, "metric-unknown"),
			MetricDigest:   "metric-unknown",
			NonAuthorizing: true,
		}),
		NonAuthorizing: true,
	})
	if metric.Status != "ready" || metric.Unknown != 1 || metric.UnknownPermille != 1000 {
		t.Fatalf("metric = %#v, want all unknown", metric)
	}
}

func TestMeasureJEVReplayLedgerMetricRejectsTampering(t *testing.T) {
	ledger := replayLedgerCheckpointFixture()
	ledger.EvidenceDigest = "tampered"
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         ledger,
		NonAuthorizing: true,
	})
	if metric.Status != "UNKNOWN" || metric.MissingStage != "replay-feedback-ledger" ||
		metric.MetricDigest != "" || metric.EvidenceDigest != "" {
		t.Fatalf("metric = %#v, want ledger UNKNOWN", metric)
	}
}

func TestMeasureJEVReplayLedgerMetricRejectsAuthorization(t *testing.T) {
	metric := MeasureJEVReplayLedgerMetric(ExecutionEnvelopeJEVReplayLedgerMetricInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: false,
	})
	if metric.Status != "UNKNOWN" || metric.MissingStage != "authorization-boundary" ||
		metric.NonAuthorizing {
		t.Fatalf("metric = %#v, want authorization UNKNOWN", metric)
	}
}
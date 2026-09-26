package decision

import "testing"

func TestAggregateJEVReplayFeedbackLedger(t *testing.T) {
	first := replayLedgerCheckpointFixture()
	stable := AggregateJEVReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput{
		Ledger:         first,
		NonAuthorizing: true,
	})
	if stable.Status != jevImprovementFeedbackStableForReview || stable.Total != 1 ||
		stable.Confirmed != 1 || stable.Refuted != 0 || stable.Unknown != 0 {
		t.Fatalf("stable = %#v, want one confirmed feedback", stable)
	}
	if err := stable.Validate(); err != nil {
		t.Fatalf("stable aggregation should validate: %v", err)
	}

	revisionLedger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       first,
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackRefuted, "metric-b"),
		MetricDigest:   "metric-b",
		NonAuthorizing: true,
	})
	revision := AggregateJEVReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput{
		Ledger:         revisionLedger,
		NonAuthorizing: true,
	})
	if revision.Status != jevImprovementFeedbackNeedsRevision || revision.Total != 2 ||
		revision.Confirmed != 1 || revision.Refuted != 1 {
		t.Fatalf("revision = %#v, want confirmed plus refuted", revision)
	}
}

func TestAggregateJEVReplayFeedbackLedgerPreservesHold(t *testing.T) {
	holdLedger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackUnknown, "metric-unknown"),
		MetricDigest:   "metric-unknown",
		NonAuthorizing: true,
	})
	hold := AggregateJEVReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput{
		Ledger:         holdLedger,
		NonAuthorizing: true,
	})
	if hold.Status != jevImprovementFeedbackHold || hold.Unknown != 1 || hold.EvidenceDigest == "" {
		t.Fatalf("hold = %#v, want unknown hold", hold)
	}
}

func TestAggregateJEVReplayFeedbackLedgerRejectsTampering(t *testing.T) {
	ledger := replayLedgerCheckpointFixture()
	ledger.EvidenceDigest = "tampered"
	aggregation := AggregateJEVReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput{
		Ledger:         ledger,
		NonAuthorizing: true,
	})
	if aggregation.Status != jevImprovementFeedbackAggregateUnknown ||
		aggregation.MissingStage != "replay-feedback-ledger" ||
		aggregation.EvidenceDigest != "" {
		t.Fatalf("aggregation = %#v, want ledger UNKNOWN", aggregation)
	}
}

func TestAggregateJEVReplayFeedbackLedgerRejectsAuthorization(t *testing.T) {
	aggregation := AggregateJEVReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput{
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: false,
	})
	if aggregation.Status != jevImprovementFeedbackAggregateUnknown ||
		aggregation.MissingStage != "authorization-boundary" ||
		aggregation.NonAuthorizing {
		t.Fatalf("aggregation = %#v, want authorization UNKNOWN", aggregation)
	}
}
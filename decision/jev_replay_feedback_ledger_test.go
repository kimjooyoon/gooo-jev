package decision

import "testing"

func replayFeedbackLedgerFixture(status, metricDigest string) JEVImprovementReplayFeedback {
	feedback := JEVImprovementReplayFeedback{
		Status:                  status,
		FeedbackKind:             status,
		CandidateDigest:          "candidate-digest",
		ReplayObservationDigest: "replay-observation-digest",
		MetricDigest:             metricDigest,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	feedback.FeedbackEvidenceDigest = digestJEVImprovementReplayFeedbackEvidence(
		feedback.ReplayObservationDigest,
		feedback.MetricDigest,
		feedback.FeedbackKind,
	)
	feedback.EvidenceDigest = digestJEVImprovementReplayFeedback(
		feedback.Status,
		feedback.FeedbackKind,
		feedback.CandidateDigest,
		feedback.ReplayObservationDigest,
		feedback.MetricDigest,
		feedback.FeedbackEvidenceDigest,
	)
	return feedback
}

func TestAppendJEVImprovementReplayFeedbackLedger(t *testing.T) {
	first := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackConfirmed, "metric-a"),
		MetricDigest:   "metric-a",
		NonAuthorizing: true,
	})
	if first.Status != jevReplayFeedbackLedgerReady || first.ConfirmedCount != 1 ||
		first.RefutedCount != 0 || first.UnknownCount != 0 {
		t.Fatalf("unexpected first ledger: %+v", first)
	}
	if err := first.Validate(); err != nil {
		t.Fatalf("first ledger should validate: %v", err)
	}

	second := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       first,
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackRefuted, "metric-b"),
		MetricDigest:   "metric-b",
		NonAuthorizing: true,
	})
	if second.Status != jevReplayFeedbackLedgerReady || second.ConfirmedCount != 1 ||
		second.RefutedCount != 1 || second.UnknownCount != 0 {
		t.Fatalf("unexpected second ledger: %+v", second)
	}
	if second.PreviousEvidenceDigest != first.EvidenceDigest {
		t.Fatalf("ledger chain was not linked")
	}
	if err := second.Validate(); err != nil {
		t.Fatalf("second ledger should validate: %v", err)
	}
}

func TestAppendJEVImprovementReplayFeedbackLedgerKeepsUnknownCount(t *testing.T) {
	ledger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackUnknown, "metric-unknown"),
		MetricDigest:   "metric-unknown",
		NonAuthorizing: true,
	})
	if ledger.Status != jevReplayFeedbackLedgerReady || ledger.UnknownCount != 1 ||
		ledger.ConfirmedCount != 0 || ledger.RefutedCount != 0 {
		t.Fatalf("UNKNOWN feedback was not retained: %+v", ledger)
	}
}

func TestAppendJEVImprovementReplayFeedbackLedgerRejectsInvalidInput(t *testing.T) {
	invalid := replayFeedbackLedgerFixture(jevReplayFeedbackConfirmed, "metric-a")
	invalid.EvidenceDigest = "tampered"
	output := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Feedback:       invalid,
		MetricDigest:   "metric-a",
		NonAuthorizing: true,
	})
	if output.Status != jevReplayFeedbackLedgerUnknown || output.MissingStage != "feedback-validation" ||
		!output.NonExecuting || !output.NonAuthorizing {
		t.Fatalf("invalid feedback was accepted: %+v", output)
	}
}

func TestJEVImprovementReplayFeedbackLedgerRejectsTampering(t *testing.T) {
	ledger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackConfirmed, "metric-a"),
		MetricDigest:   "metric-a",
		NonAuthorizing: true,
	})
	ledger.ConfirmedCount++
	if err := ledger.Validate(); err == nil {
		t.Fatalf("tampered ledger unexpectedly validated")
	}
}
package decision

import "testing"

func ledgerReviewObservation(outcome string) ExecutionEnvelopeJEVImprovementReviewObservation {
	return ExecutionEnvelopeJEVImprovementReviewObservation{
		Status: "observed", ReviewOutcome: outcome, ObservationDigest: "observation-" + outcome, NonExecuting: true, NonAuthorizing: true,
	}
}

func TestAppendExecutionEnvelopeJEVImprovementReviewLedger(t *testing.T) {
	root := AppendExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerInput{
		Observation:    ledgerReviewObservation("confirmed"),
		EntryIndex:     0,
		NonAuthorizing: true,
	})
	if root.Status != "appended" || root.EntryIndex != 0 || root.LedgerEntryDigest == "" || root.PreviousLedgerDigest != "" {
		t.Fatalf("root %+v", root)
	}
	next := AppendExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerInput{
		Observation:          ledgerReviewObservation("refuted"),
		PreviousLedgerDigest: root.LedgerEntryDigest,
		EntryIndex:           1,
		NonAuthorizing:       true,
	})
	if next.Status != "appended" || next.EntryIndex != 1 || next.PreviousLedgerDigest != root.LedgerEntryDigest || next.LedgerEntryDigest == "" {
		t.Fatalf("next %+v", next)
	}
}

func TestAppendExecutionEnvelopeJEVImprovementReviewLedgerRequiresPredecessor(t *testing.T) {
	got := AppendExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerInput{
		Observation:    ledgerReviewObservation("unknown"),
		EntryIndex:     2,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "ledger-predecessor" || got.LedgerEntryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestAppendExecutionEnvelopeJEVImprovementReviewLedgerPreservesObservationStage(t *testing.T) {
	observation := ledgerReviewObservation("unknown")
	observation.Status = "UNKNOWN"
	observation.MissingStage = "review-evidence"
	got := AppendExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerInput{
		Observation:    observation,
		EntryIndex:     0,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "review-evidence" || got.LedgerEntryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestAppendExecutionEnvelopeJEVImprovementReviewLedgerRejectsAuthorization(t *testing.T) {
	got := AppendExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.LedgerEntryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

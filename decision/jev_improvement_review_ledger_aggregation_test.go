package decision

import "testing"

func aggregationLedgerEntry(index int, previous, outcome string) ExecutionEnvelopeJEVImprovementReviewLedgerEntry {
	entry := ExecutionEnvelopeJEVImprovementReviewLedgerEntry{
		Status: "appended", ReviewOutcome: outcome, EntryIndex: index, PreviousLedgerDigest: previous, ObservationDigest: "observation-" + outcome, NonExecuting: true, NonAuthorizing: true,
	}
	entry.LedgerEntryDigest, _ = digestExecutionEnvelopeJEVImprovementReviewLedgerEntry(entry)
	return entry
}

func TestAggregateExecutionEnvelopeJEVImprovementReviewLedger(t *testing.T) {
	root := aggregationLedgerEntry(0, "", "confirmed")
	next := aggregationLedgerEntry(1, root.LedgerEntryDigest, "refuted")
	got := AggregateExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput{
		Entries:        []ExecutionEnvelopeJEVImprovementReviewLedgerEntry{root, next},
		NonAuthorizing: true,
	})
	if got.Status != "needs-revision" || got.Total != 2 || got.Confirmed != 1 || got.Refuted != 1 || got.InputLedgerDigest == "" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAggregateExecutionEnvelopeJEVImprovementReviewLedgerStableAndHold(t *testing.T) {
	stable := AggregateExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput{
		Entries:        []ExecutionEnvelopeJEVImprovementReviewLedgerEntry{aggregationLedgerEntry(0, "", "confirmed")},
		NonAuthorizing: true,
	})
	if stable.Status != "stable-for-review" || stable.Confirmed != 1 {
		t.Fatalf("stable %+v", stable)
	}
	hold := AggregateExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput{
		Entries:        []ExecutionEnvelopeJEVImprovementReviewLedgerEntry{aggregationLedgerEntry(0, "", "unknown")},
		NonAuthorizing: true,
	})
	if hold.Status != "hold" || hold.Unknown != 1 {
		t.Fatalf("hold %+v", hold)
	}
}

func TestAggregateExecutionEnvelopeJEVImprovementReviewLedgerRejectsBrokenChain(t *testing.T) {
	root := aggregationLedgerEntry(0, "", "confirmed")
	next := aggregationLedgerEntry(1, "tampered-predecessor", "confirmed")
	got := AggregateExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput{
		Entries:        []ExecutionEnvelopeJEVImprovementReviewLedgerEntry{root, next},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "ledger[1]-predecessor" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestAggregateExecutionEnvelopeJEVImprovementReviewLedgerRequiresHistory(t *testing.T) {
	got := AggregateExecutionEnvelopeJEVImprovementReviewLedger(ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput{NonAuthorizing: true})
	if got.Status != "UNKNOWN" || got.MissingStage != "review-ledger-history" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

package decision

import "strings"

// ExecutionEnvelopeJEVImprovementReviewLedgerInput appends one observed review
// result to a provenance chain without mutating an existing entry.
type ExecutionEnvelopeJEVImprovementReviewLedgerInput struct {
	Observation          ExecutionEnvelopeJEVImprovementReviewObservation
	PreviousLedgerDigest string
	EntryIndex           int
	NonAuthorizing       bool
}

// ExecutionEnvelopeJEVImprovementReviewLedgerEntry is an immutable append
// result that can be used as the next self-improvement evidence input.
type ExecutionEnvelopeJEVImprovementReviewLedgerEntry struct {
	Status               string
	ReviewOutcome        string
	EntryIndex           int
	PreviousLedgerDigest string
	ObservationDigest    string
	LedgerEntryDigest    string
	MissingStage         string
	NonExecuting         bool
	NonAuthorizing       bool
}

func digestExecutionEnvelopeJEVImprovementReviewLedgerEntry(entry ExecutionEnvelopeJEVImprovementReviewLedgerEntry) (string, error) {
	return Digest(struct {
		Status               string
		ReviewOutcome        string
		EntryIndex           int
		PreviousLedgerDigest string
		ObservationDigest    string
		MissingStage         string
	}{
		Status:               entry.Status,
		ReviewOutcome:        entry.ReviewOutcome,
		EntryIndex:           entry.EntryIndex,
		PreviousLedgerDigest: entry.PreviousLedgerDigest,
		ObservationDigest:    entry.ObservationDigest,
		MissingStage:         entry.MissingStage,
	})
}

// AppendExecutionEnvelopeJEVImprovementReviewLedger records one review
// observation in order while preserving UNKNOWN and chain-boundary failures.
func AppendExecutionEnvelopeJEVImprovementReviewLedger(input ExecutionEnvelopeJEVImprovementReviewLedgerInput) ExecutionEnvelopeJEVImprovementReviewLedgerEntry {
	output := ExecutionEnvelopeJEVImprovementReviewLedgerEntry{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Observation.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	}
	if !input.Observation.NonExecuting {
		output.MissingStage = "execution-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	}
	if input.Observation.Status != "observed" || strings.TrimSpace(input.Observation.ObservationDigest) == "" {
		output.MissingStage = input.Observation.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-observation"
		}
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	}
	if input.EntryIndex < 0 {
		output.MissingStage = "ledger-index"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	}
	if input.EntryIndex == 0 && strings.TrimSpace(input.PreviousLedgerDigest) != "" {
		output.MissingStage = "ledger-root"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	}
	if input.EntryIndex > 0 && strings.TrimSpace(input.PreviousLedgerDigest) == "" {
		output.MissingStage = "ledger-predecessor"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	}
	output.Status = "appended"
	output.ReviewOutcome = input.Observation.ReviewOutcome
	output.EntryIndex = input.EntryIndex
	output.PreviousLedgerDigest = input.PreviousLedgerDigest
	output.ObservationDigest = input.Observation.ObservationDigest
	return finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
}

func finalizeExecutionEnvelopeJEVImprovementReviewLedgerEntry(output ExecutionEnvelopeJEVImprovementReviewLedgerEntry) ExecutionEnvelopeJEVImprovementReviewLedgerEntry {
	digest, err := digestExecutionEnvelopeJEVImprovementReviewLedgerEntry(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "review-ledger-entry-digest"
		output.LedgerEntryDigest = ""
		return output
	}
	output.LedgerEntryDigest = digest
	return output
}

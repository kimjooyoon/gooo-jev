package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput validates an
// ordered review ledger before it becomes the next feedback signal.
type ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput struct {
	Entries        []ExecutionEnvelopeJEVImprovementReviewLedgerEntry
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVImprovementReviewLedgerAggregation summarizes a verified
// ledger without selecting, applying, or authorizing a revision.
type ExecutionEnvelopeJEVImprovementReviewLedgerAggregation struct {
	Status            string
	MissingStage      string
	Total             int
	Confirmed         int
	Refuted           int
	Unknown           int
	InputLedgerDigest string
	EvidenceDigest    string
	NonExecuting      bool
	NonAuthorizing    bool
}

func (a ExecutionEnvelopeJEVImprovementReviewLedgerAggregation) Validate() error {
	if a.Status == "" || a.Total <= 0 || a.InputLedgerDigest == "" || a.EvidenceDigest == "" {
		return fmt.Errorf("incomplete JEV improvement review ledger aggregation")
	}
	if a.Confirmed < 0 || a.Refuted < 0 || a.Unknown < 0 || a.Confirmed+a.Refuted+a.Unknown != a.Total {
		return fmt.Errorf("invalid JEV improvement review ledger aggregation counts")
	}
	switch a.Status {
	case "stable-for-review", "needs-revision", "hold":
	default:
		return fmt.Errorf("invalid JEV improvement review ledger aggregation status")
	}
	if !a.NonExecuting || !a.NonAuthorizing {
		return fmt.Errorf("JEV improvement review ledger aggregation must be non-executing and non-authorizing")
	}
	expected, err := digestExecutionEnvelopeJEVImprovementReviewLedgerAggregationEvidence(a)
	if err != nil || a.EvidenceDigest != expected {
		return fmt.Errorf("JEV improvement review ledger aggregation evidence digest mismatch")
	}
	return nil
}

func digestExecutionEnvelopeJEVImprovementReviewLedgerAggregationEvidence(aggregation ExecutionEnvelopeJEVImprovementReviewLedgerAggregation) (string, error) {
	return Digest(struct {
		Status            string
		Total             int
		Confirmed         int
		Refuted           int
		Unknown           int
		InputLedgerDigest string
	}{
		Status:            aggregation.Status,
		Total:             aggregation.Total,
		Confirmed:         aggregation.Confirmed,
		Refuted:           aggregation.Refuted,
		Unknown:           aggregation.Unknown,
		InputLedgerDigest: aggregation.InputLedgerDigest,
	})
}

func digestExecutionEnvelopeJEVImprovementReviewLedgerEntries(entries []string) (string, error) {
	return Digest(struct{ EntryDigests []string }{EntryDigests: entries})
}

// AggregateExecutionEnvelopeJEVImprovementReviewLedger turns a verified
// append-only chain into a feedback status for the next improvement cycle.
func AggregateExecutionEnvelopeJEVImprovementReviewLedger(input ExecutionEnvelopeJEVImprovementReviewLedgerAggregationInput) ExecutionEnvelopeJEVImprovementReviewLedgerAggregation {
	output := ExecutionEnvelopeJEVImprovementReviewLedgerAggregation{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
	}
	if len(input.Entries) == 0 {
		output.MissingStage = "review-ledger-history"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
	}
	entryDigests := make([]string, 0, len(input.Entries))
	for index, entry := range input.Entries {
		if !entry.NonAuthorizing {
			output.MissingStage = fmt.Sprintf("ledger[%d]-authorization-boundary", index)
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		if !entry.NonExecuting {
			output.MissingStage = fmt.Sprintf("ledger[%d]-execution-boundary", index)
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		if entry.Status != "appended" || strings.TrimSpace(entry.LedgerEntryDigest) == "" {
			output.MissingStage = entry.MissingStage
			if strings.TrimSpace(output.MissingStage) == "" {
				output.MissingStage = fmt.Sprintf("ledger[%d]", index)
			}
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		if entry.EntryIndex != index {
			output.MissingStage = fmt.Sprintf("ledger[%d]-order", index)
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		if index == 0 && strings.TrimSpace(entry.PreviousLedgerDigest) != "" {
			output.MissingStage = "ledger-root"
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		if index > 0 && entry.PreviousLedgerDigest != input.Entries[index-1].LedgerEntryDigest {
			output.MissingStage = fmt.Sprintf("ledger[%d]-predecessor", index)
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		expectedDigest, err := digestExecutionEnvelopeJEVImprovementReviewLedgerEntry(entry)
		if err != nil || entry.LedgerEntryDigest != expectedDigest {
			output.MissingStage = fmt.Sprintf("ledger[%d]-digest", index)
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
		entryDigests = append(entryDigests, entry.LedgerEntryDigest)
		output.Total++
		switch entry.ReviewOutcome {
		case "confirmed":
			output.Confirmed++
		case "refuted":
			output.Refuted++
		case "unknown":
			output.Unknown++
		default:
			output.MissingStage = fmt.Sprintf("ledger[%d]-review-outcome", index)
			return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
		}
	}
	inputDigest, err := digestExecutionEnvelopeJEVImprovementReviewLedgerEntries(entryDigests)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Total = 0
		output.MissingStage = "ledger-input-digest"
		return finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output)
	}
	output.InputLedgerDigest = inputDigest
	switch {
	case output.Refuted > 0:
		output.Status = "needs-revision"
	case output.Unknown > 0:
		output.Status = "hold"
	default:
		output.Status = "stable-for-review"
	}
	output.EvidenceDigest, err = digestExecutionEnvelopeJEVImprovementReviewLedgerAggregationEvidence(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "ledger-aggregation-evidence"
		output.EvidenceDigest = ""
		return output
	}
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "ledger-aggregation-evidence"
		output.EvidenceDigest = ""
	}
	return output
}

func finalizeExecutionEnvelopeJEVImprovementReviewLedgerAggregation(output ExecutionEnvelopeJEVImprovementReviewLedgerAggregation) ExecutionEnvelopeJEVImprovementReviewLedgerAggregation {
	if output.EvidenceDigest != "" {
		return output
	}
	digest, err := digestExecutionEnvelopeJEVImprovementReviewLedgerAggregationEvidence(output)
	if err != nil {
		output.EvidenceDigest = ""
		output.MissingStage = "ledger-aggregation-digest"
		return output
	}
	output.EvidenceDigest = digest
	return output
}

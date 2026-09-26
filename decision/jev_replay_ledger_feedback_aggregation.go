package decision

import "strings"

// ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput adapts one
// validated append-only replay ledger to the existing cycle feedback shape.
type ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput struct {
	Ledger         JEVImprovementReplayFeedbackLedger
	NonAuthorizing bool
}

// AggregateJEVReplayFeedbackLedger returns the existing feedback aggregation
// contract without replaying, applying, or authorizing any candidate.
func AggregateJEVReplayFeedbackLedger(input ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput) JEVImprovementFeedbackAggregation {
	output := JEVImprovementFeedbackAggregation{
		Status:         jevImprovementFeedbackAggregateUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Ledger.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Ledger.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Ledger.Validate(); err != nil {
		output.MissingStage = "replay-feedback-ledger"
		return output
	}
	output.Total = input.Ledger.ConfirmedCount + input.Ledger.RefutedCount + input.Ledger.UnknownCount
	if output.Total <= 0 {
		output.MissingStage = "feedback-history"
		return output
	}
	output.Confirmed = input.Ledger.ConfirmedCount
	output.Refuted = input.Ledger.RefutedCount
	output.Unknown = input.Ledger.UnknownCount
	output.InputEvidenceDigest = input.Ledger.EvidenceDigest
	switch {
	case output.Refuted > 0:
		output.Status = jevImprovementFeedbackNeedsRevision
	case output.Unknown > 0:
		output.Status = jevImprovementFeedbackHold
	default:
		output.Status = jevImprovementFeedbackStableForReview
	}
	output.EvidenceDigest = digestJEVImprovementFeedbackAggregation(
		output.Status,
		output.Total,
		output.Confirmed,
		output.Refuted,
		output.Unknown,
		output.InputEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		output.Status = jevImprovementFeedbackAggregateUnknown
		output.MissingStage = "feedback-aggregation-evidence"
		output.EvidenceDigest = ""
	}
	if strings.TrimSpace(output.MissingStage) != "" {
		output.EvidenceDigest = ""
	}
	return output
}
package decision

// ExecutionEnvelopeEvidenceLedgerCounterexampleInput binds a reverse ledger
// observation to a bounded review disposition.
type ExecutionEnvelopeEvidenceLedgerCounterexampleInput struct {
	Status         string `json:"status"`
	FirstMismatch  string `json:"first_mismatch"`
	LedgerDigest   string `json:"ledger_digest"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeEvidenceLedgerCounterexample records whether ledger review
// is required without treating the disposition as authorization.
type ExecutionEnvelopeEvidenceLedgerCounterexample struct {
	Status         string `json:"status"`
	ReviewRequired bool   `json:"review_required"`
	LedgerDigest   string `json:"ledger_digest"`
	MismatchStage  string `json:"mismatch_stage"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ClassifyExecutionEnvelopeEvidenceLedgerCounterexample preserves the first
// mismatch and fails closed for incomplete reverse evidence.
func ClassifyExecutionEnvelopeEvidenceLedgerCounterexample(input ExecutionEnvelopeEvidenceLedgerCounterexampleInput) ExecutionEnvelopeEvidenceLedgerCounterexample {
	output := ExecutionEnvelopeEvidenceLedgerCounterexample{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MismatchStage = "authorization-boundary"
		return output
	}
	switch input.Status {
	case "counterexample":
		if input.FirstMismatch == "" || input.LedgerDigest == "" {
			output.MismatchStage = "ledger-counterexample-evidence"
			return output
		}
		output.Status = "review-required"
		output.ReviewRequired = true
		output.LedgerDigest = input.LedgerDigest
		output.MismatchStage = input.FirstMismatch
	case "reproduced":
		if input.LedgerDigest == "" {
			output.MismatchStage = "reproduced-ledger-evidence"
			return output
		}
		output.Status = "observation-only"
		output.LedgerDigest = input.LedgerDigest
	case "UNKNOWN":
		output.Status = "hold"
		output.MismatchStage = "ledger-reverse-observation"
	default:
		output.MismatchStage = "ledger-reverse-observation"
	}
	return output
}
